package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func DbConnect() *gorm.DB {
	// Grab environment variables for connection
	var DB_USER string = os.Getenv("DB_USER")
	var DB_PASS string = os.Getenv("DB_PASS")
	var DB_HOST string = os.Getenv("DB_HOST")
	var DB_PORT string = os.Getenv("DB_PORT")
	var DB_NAME string = os.Getenv("DB_NAME")

	dbUrl := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", DB_HOST, DB_USER, DB_PASS, DB_NAME, DB_PORT)

	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Property{})
	db.AutoMigrate(&Feature{})
	db.AutoMigrate(&PropertyLog{})
	db.AutoMigrate(&Contact{})
	db.AutoMigrate(&Task{})
	db.AutoMigrate(&TaskLog{})
	db.AutoMigrate(&Transaction{})
	db.AutoMigrate(&WorkType{})
	db.AutoMigrate(&MaintenanceRequest{})
	db.AutoMigrate(&Vendor{})
	db.AutoMigrate(&PropertyAttachment{})

	// Build basic work types
	buildBasicWorkTypes(db)

	return db
}

func buildBasicWorkTypes(db *gorm.DB) {
	// Build basic work types
	workTypes := []WorkType{
		{Name: "Lighting"},
		{Name: "Plumbing"},
		{Name: "Electrical"},
		{Name: "Painting"},
		{Name: "Cleaning"},
		{Name: "Gardening"},
		{Name: "HVAC"},
		{Name: "Security"},
		{Name: "Fire safety"},
		{Name: "Energy management"},
		{Name: "Escalators/lifts"},
		{Name: "Facade"},
		{Name: "Other"},
	}
	createWorkOrderIfNotExist(workTypes, db)
}
