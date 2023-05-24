package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"gorm.io/gorm"

	_ "github.com/swaggo/http-swagger/example/go-chi/docs"

	"github.com/dmawardi/Go-Template/internal/auth"
	"github.com/dmawardi/Go-Template/internal/config"
	"github.com/dmawardi/Go-Template/internal/controller"
	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/repository"
	"github.com/dmawardi/Go-Template/internal/routes"
	"github.com/dmawardi/Go-Template/internal/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const portNumber = ":8080"

// Init state
var app config.AppConfig

// API Details
// @title           TBK Property API
// @version         1.0
// @description     This is the API server for TBK property.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/

// @securityDefinitions.apikey BearerToken
// @in header
// @name Authorization
func main() {

	// Build context
	ctx := context.Background()
	// Set context in app config
	app.Ctx = ctx
	// Load env variables
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Unable to load environment variables.")
	}

	// Set state in other packages
	controller.SetStateInHandlers(&app)
	auth.SetStateInAuth(&app)
	service.BuildServiceState(&app)

	// Create client using DbConnect
	client := db.DbConnect()
	// Set in state
	app.DbClient = client

	// Create api
	api := ApiSetup(client)

	// Setup enforcer
	e, err := auth.EnforcerSetup(client)
	if err != nil {
		log.Fatal("Couldn't setup RBAC Authorization Enforcer")
	}
	// Set enforcer in state
	app.RBEnforcer = e

	fmt.Printf("Starting application on port: %s\n", portNumber)

	// Server settings
	srv := &http.Server{
		Addr:    portNumber,
		Handler: api.Routes(),
	}

	// Listen and serve using server settings above
	err = srv.ListenAndServe()
	if err != nil {

		log.Fatal(err)
	}
}

func ApiSetup(client *gorm.DB) routes.Api {
	// user
	userRepo := repository.NewUserRepository(client)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)

	// property log
	propLogRepo := repository.NewPropertyLogRepository(client)
	propLogService := service.NewPropertyLogService(propLogRepo)
	propLogController := controller.NewPropertyLogController(propLogService)

	// property
	propRepo := repository.NewPropertyRepository(client)
	propService := service.NewPropertyService(propRepo)
	propController := controller.NewPropertyController(propService, propLogService)

	// feature
	featRepo := repository.NewFeatureRepository(client)
	featService := service.NewFeatureService(featRepo)
	featController := controller.NewFeatureController(featService)

	// contact
	contactRepo := repository.NewContactRepository(client)
	contactService := service.NewContactService(contactRepo)
	contactController := controller.NewContactController(contactService)

	// task logs
	taskLogRepo := repository.NewTaskLogRepository(client)
	taskLogService := service.NewTaskLogService(taskLogRepo)
	taskLogController := controller.NewTaskLogController(taskLogService)

	// task
	taskRepo := repository.NewTaskRepository(client)
	taskService := service.NewTaskService(taskRepo)
	taskController := controller.NewTaskController(taskService, taskLogService)

	// transaction
	transactionRepo := repository.NewTransactionRepository(client)
	transactionService := service.NewTransactionService(transactionRepo)
	transactionController := controller.NewTransactionController(transactionService)

	// Maintenance requests
	maintenanceRepo := repository.NewMaintenanceRequestRepository(client)
	maintenanceService := service.NewMaintenanceRequestService(maintenanceRepo)
	maintenanceController := controller.NewMaintenanceRequestController(maintenanceService)

	// Build API using controllers
	api := routes.NewApi(userController, propController, featController, propLogController, contactController, taskController, taskLogController, transactionController, maintenanceController)
	return api
}
