package db

import "gorm.io/gorm"

// Check if an item with a specific name exists
func workOrderExists(name string, db *gorm.DB) bool {
	var item WorkType
	result := db.First(&item, "name = ?", name)

	// Returns true if the item exists
	return result.Error == nil
}

// Create a list of work orders only if they don't already exist
func createWorkOrderIfNotExist(workOrders []WorkType, db *gorm.DB) {
	// Loop through the work orders
	for _, wo := range workOrders {
		// Check if the work order already exists
		if !workOrderExists(wo.Name, db) {
			// Create the work order
			result := db.Create(&wo)
			if result.Error != nil {
				panic("failed to create default work order types")
			}
		}
	}
}
