package routes

import (
	"net/http"

	"github.com/dmawardi/Go-Template/internal/auth"
	"github.com/dmawardi/Go-Template/internal/controller"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/swaggo/http-swagger/example/go-chi/docs"
)

type Api interface {
	Routes() http.Handler
}

type api struct {
	user               controller.UserController
	property           controller.PropertyController
	feature            controller.FeatureController
	propertyLog        controller.PropertyLogController
	contact            controller.ContactController
	task               controller.TaskController
	taskLog            controller.TaskLogController
	transaction        controller.TransactionController
	maintenanceRequest controller.MaintenanceRequestController
	workType           controller.WorkTypeController
}

func NewApi(user controller.UserController,
	property controller.PropertyController,
	feature controller.FeatureController,
	propertyLog controller.PropertyLogController,
	contact controller.ContactController,
	task controller.TaskController,
	taskLog controller.TaskLogController,
	trans controller.TransactionController,
	maintenance controller.MaintenanceRequestController,
	workType controller.WorkTypeController) Api {
	return &api{user, property, feature, propertyLog, contact, task, taskLog, trans, maintenance, workType}
}

func (a api) Routes() http.Handler {
	// Create new router
	mux := chi.NewRouter()
	// Use built in Chi middleware
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.Logger)

	// Public routes
	mux.Group(func(mux chi.Router) {
		// @tag.name Public Routes
		// @tag.description Unprotected routes
		mux.Get("/", controller.GetJobs)
		// Login
		mux.Post("/api/users/login", a.user.Login)

		// Create new user
		mux.Post("/api/users", a.user.Create)

		// Private routes
		mux.Group(func(mux chi.Router) {
			mux.Use(auth.AuthenticateJWT)

			// @tag.name Private routes
			// @tag.description Protected routes
			// users
			mux.Get("/api/users", a.user.FindAll)
			mux.Get("/api/users/{id}", a.user.Find)
			mux.Put("/api/users/{id}", a.user.Update)
			mux.Delete("/api/users/{id}", a.user.Delete)

			// My profile
			mux.Get("/api/me", a.user.GetMyUserDetails)
			mux.Post("/api/me", controller.HealthCheck)
			mux.Put("/api/me", a.user.UpdateMyProfile)

			// properties
			mux.Post("/api/properties", a.property.Create)
			mux.Get("/api/properties", a.property.FindAll)
			mux.Get("/api/properties/{id}", a.property.Find)
			mux.Put("/api/properties/{id}", a.property.Update)
			mux.Delete("/api/properties/{id}", a.property.Delete)

			// Property features
			mux.Post("/api/features", a.feature.Create)
			mux.Get("/api/features", a.feature.FindAll)
			mux.Get("/api/features/{id}", a.feature.Find)
			mux.Put("/api/features/{id}", a.feature.Update)
			mux.Delete("/api/features/{id}", a.feature.Delete)

			// Property Logs
			mux.Post("/api/property-logs", a.propertyLog.Create)
			mux.Get("/api/property-logs", a.propertyLog.FindAll)
			mux.Get("/api/property-logs/{id}", a.propertyLog.Find)
			mux.Put("/api/property-logs/{id}", a.propertyLog.Update)
			mux.Delete("/api/property-logs/{id}", a.propertyLog.Delete)

			// Contacts
			mux.Post("/api/contacts", a.contact.Create)
			mux.Get("/api/contacts", a.contact.FindAll)
			mux.Get("/api/contacts/{id}", a.contact.Find)
			mux.Put("/api/contacts/{id}", a.contact.Update)
			mux.Delete("/api/contacts/{id}", a.contact.Delete)

			// Tasks
			mux.Post("/api/tasks", a.task.Create)
			mux.Get("/api/tasks", a.task.FindAll)
			mux.Get("/api/tasks/{id}", a.task.Find)
			mux.Put("/api/tasks/{id}", a.task.Update)
			mux.Delete("/api/tasks/{id}", a.task.Delete)

			// Task Logs
			mux.Post("/api/task-logs", a.taskLog.Create)
			mux.Get("/api/task-logs", a.taskLog.FindAll)
			mux.Get("/api/task-logs/{id}", a.taskLog.Find)
			mux.Put("/api/task-logs/{id}", a.taskLog.Update)
			mux.Delete("/api/task-logs/{id}", a.taskLog.Delete)

			// Transactions
			mux.Post("/api/transactions", a.transaction.Create)
			mux.Get("/api/transactions", a.transaction.FindAll)
			mux.Get("/api/transactions/{id}", a.transaction.Find)
			mux.Put("/api/transactions/{id}", a.transaction.Update)
			mux.Delete("/api/transactions/{id}", a.transaction.Delete)

			// Maintenance requests
			mux.Post("/api/maintenance", a.maintenanceRequest.Create)
			mux.Get("/api/maintenance", a.maintenanceRequest.FindAll)
			mux.Get("/api/maintenance/{id}", a.maintenanceRequest.Find)
			mux.Put("/api/maintenance/{id}", a.maintenanceRequest.Update)
			mux.Delete("/api/maintenance/{id}", a.maintenanceRequest.Delete)

			// Work types
			mux.Post("/api/work-types", a.workType.Create)
			mux.Get("/api/work-types", a.workType.FindAll)
			mux.Get("/api/work-types/{id}", a.workType.Find)
			mux.Put("/api/work-types/{id}", a.workType.Update)
			mux.Delete("/api/work-types/{id}", a.workType.Delete)
		})

	})

	// Serve API Swagger docs
	mux.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/docs/swagger.json"), //The url pointing to API definition
	))

	// Build fileserver using static directory
	fileServer := http.FileServer(http.Dir("./static"))
	// Handle all calls to /static/* by stripping prefix and sending to file server
	mux.Handle("/static/*", http.StripPrefix("/static/", fileServer))

	return mux
}
