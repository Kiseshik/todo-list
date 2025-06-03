package main

import (
	"context"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	pet "github.com/Kiseshik/pet"
	"github.com/Kiseshik/pet/pkg/handler"
	"github.com/Kiseshik/pet/pkg/repository"
	"github.com/Kiseshik/pet/pkg/service"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

//hardcode чтобы все запускать самому и ковыряться
//docker run --name postgres -e POSTGRES_PASSWORD=qwerty -p 5432:5432 --rm -d postgres !!!
//docker exec -i postgres psql -U postgres -c "$(cat schema\000001_init.up.sql)"
// docker exec -it postgres psql -U postgres

// docker exec -it postgres psql -U postgres -c "\dt" тест на существование таблиц
// docker exec -it postgres psql -U postgres -c "\d users" тест на таблицы users
// docker exec -it postgres psql -U postgres -c "SELECT * FROM users;" посмотреть юзеров
// docker exec -it postgres psql -U postgres -c "SELECT * FROM todo_lists;"
// docker exec -it postgres psql -U postgres -c "SELECT * FROM users_lists;"

//docker exec -it postgres psql -U postgres -d postgres -c "SELECT tl.id, tl.title, tl.description FROM todo_lists tl INNER JOIN users_lists ul ON tl.id = ul.list_id WHERE ul.user_id = 13;"

//docker exec -it postgres psql -U postgres -d postgres -c "\dt"

//обычный запуск
//docker-compose up -d
//docker-compose down -v

func startPostgres() {
	checkCmd := exec.Command("docker", "inspect", "-f", "{{.State.Running}}", "postgres")
	if err := checkCmd.Run(); err == nil {
		logrus.Info("PostgreSQL started")
		return
	}

	cmd := exec.Command("docker-compose", "up", "-d")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		logrus.Fatalf("Error to start PostgreSQL: %v", err)
	}

	time.Sleep(5 * time.Second)
	logrus.Info("PostgreSQL successfully started")
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}

func main() {

	startPostgres()

	logrus.SetFormatter(new(logrus.JSONFormatter))

	if err := initConfig(); err != nil {
		logrus.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})

	if err != nil {
		logrus.Fatalf("failed to initialize DB: %s", err.Error())
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)
	srv := new(pet.Server)

	go func() {
		if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
			logrus.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	logrus.Info("App successfully started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("TodoApp Shutting Down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}

}
