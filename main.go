package main

import (
	"fmt"
	"log"
	"os"
	"soteria_go/pkg/config"
	routers "soteria_go/pkg/routers"
	middleware "soteria_go/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load("project_env_files/.env")
	if err != nil {
		log.Fatal("Error Loading Env File: ", err)
	}
	envi := os.Getenv("ENVIRONMENT")

	err = godotenv.Load(fmt.Sprintf("project_env_files/.env-%v", envi)) //
	if err != nil {
		log.Fatal("Error Loading Env File: ", err)
	}

	// Initialize DB here
	config.CreateConnection()
	// Declare & initialize fiber
	app := fiber.New(fiber.Config{
		UnescapePath: true,
	})

	// For GoRoutine implementation
	// appb := fiber.New(fiber.Config{
	// 	UnescapePath: true,
	// })

	// Configure application CORS
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// For GoRoutine implementation
	// appb.Use(cors.New(cors.Config{
	// 	AllowOrigins: "*",
	// 	AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	// }))

	// Declare & initialize logger
	app.Use(logger.New())

	// For GoRoutine implementation
	// appb.Use(logger.New())

	// Declare & initialize routes
	routers.SetupPublicRoutes(app)
	routers.SetupPrivateRoutes(app)

	// For GoRoutine implementation
	// routers.SetupPublicRoutesB(appb)
	// go func() {
	// 	err := appb.Listen(fmt.Sprintf(":8002"))
	// 	if err != nil {
	// 		log.Fatal(err.Error())
	// 	}
	// }()

	fmt.Println("Port: ", middleware.GetEnv("PORT"))
	// Serve the application
	if middleware.GetEnv("SSL") == "enabled" {
		log.Fatal(app.ListenTLS(
			fmt.Sprintf(":%s", middleware.GetEnv("PORT")),
			middleware.GetEnv("SSL_CERTIFICATE"),
			middleware.GetEnv("SSL_KEY"),
		))
	} else {
		err := app.Listen(fmt.Sprintf("%s:%s", middleware.GetEnv("RUN_HOST"), middleware.GetEnv("PORT")))
		if err != nil {
			log.Fatal(err.Error())
		}
	}
}
