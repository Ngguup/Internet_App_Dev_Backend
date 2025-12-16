package app

import (
	"context"
	"lab1/internal/app/redis"

	"lab1/internal/app/config"
	"lab1/internal/app/dsn"
	"lab1/internal/app/repository"

	"fmt"
	"lab1/internal/app/handler"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"time"
)

type Application struct {
	Router  *gin.Engine
	Handler *handler.Handler
}

func NewApp(ctx context.Context) (*Application, error) {

	router := gin.Default()

	// -----------------------------
	// CORS middleware
	// -----------------------------
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // можно указать конкретно: "http://127.0.0.1:1420"
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	// -----------------------------

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
		Router:  router,
		Handler: hand,
	}, nil
}

func (a *Application) RunApp() {
	logrus.Info("Server start up")

	a.Handler.RegisterHandler(a.Router)
	a.Handler.RegisterStatic(a.Router)

	serverAddress := fmt.Sprintf("%s:%d", a.Handler.Config.ServiceHost, a.Handler.Config.ServicePort)
	if err := a.Router.RunTLS(serverAddress, "localhost+2.pem", "localhost+2-key.pem",); err != nil {
		logrus.Fatal(err)
	}
	logrus.Info("Server down")
}
