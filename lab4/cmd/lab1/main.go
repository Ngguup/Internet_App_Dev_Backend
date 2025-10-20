// package main

// import (
// 	"context"
// 	"lab1/internal/pkg/app"
// 	"log"
// 	"os"
// )

// func main() {
// 	ctx := context.Background()
	
// 	application, err := app.NewApp(ctx)
// 	if err != nil {
// 		log.Println("cant create application", err)
// 		os.Exit(2)
// 	}

// 	log.Println("Application start!")
// 	application.RunApp()
// 	log.Println("Application terminated!")
// }

package main

import (
	"context"
	"lab1/internal/pkg/app"
	"log"
	"os"
)

// @title BITOP
// @version 1.0
// @description BMSTU Open IT Platform. API для работы с услугами (dataGrowthFactor), заявками (growthRequest) и пользователями.
// @contact.name API Support
// @contact.url https://vk.com/bmstu_schedule
// @contact.email bitop@spatecon.ru
// @license.name AS IS (NO WARRANTY)
// @host 127.0.0.1:8080
// @schemes http
// @BasePath /

func main() {
	ctx := context.Background()

	application, err := app.NewApp(ctx)
	if err != nil {
		log.Println("cant create application", err)
		os.Exit(2)
	}

	log.Println("Application start!")
	application.RunApp()
	log.Println("Application terminated!")
}




