package app

import (
	"context"
	"lab1/internal/app/redis"

	"lab1/internal/app/config"
	"lab1/internal/app/dsn"
	"lab1/internal/app/repository"

	"fmt"
	"lab1/internal/app/handler"
	// "lab1/internal/pkg"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// type Application struct {
// 	config *config.Config
// 	repo   *repository.Repository
// 	redis  *redis.Client
// }

type Application struct {
	Router  *gin.Engine
	Handler *handler.Handler
}

func NewApp(ctx context.Context) (*Application, error) {
	// cfg, err := config.NewConfig(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	// repo, err := repository.New(dsn.FromEnv())
	// if err != nil {
	// 	return nil, err
	// }

	router := gin.Default()
	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	postgresString := dsn.FromEnv()
	fmt.Println(postgresString)

	rep, errRep := repository.New(postgresString)
	if errRep != nil {
		logrus.Fatalf("error initializing repository: %v", errRep)
	}


	redisClient, err := redis.New(ctx, conf.Redis)
	if err != nil {
		return nil, err
	}

	hand := handler.NewHandler(rep, conf, redisClient)

	return &Application{
		Router: router,
		Handler: hand,
	}, nil
}

// func (a *Application) Run() error {
// 	log.Println("application start running")
// 	a.StartServer()
// 	log.Println("application shut down")

// 	return nil
// }

func (a *Application) RunApp() {
	logrus.Info("Server start up")

	a.Handler.RegisterHandler(a.Router)
	a.Handler.RegisterStatic(a.Router)

	serverAddress := fmt.Sprintf("%s:%d", a.Handler.Config.ServiceHost, a.Handler.Config.ServicePort)
	if err := a.Router.Run(serverAddress); err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("Server down")
}