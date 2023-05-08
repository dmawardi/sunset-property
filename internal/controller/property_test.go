package controller_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/helpers"
	"github.com/dmawardi/Go-Template/internal/models"
)

func TestPropertyController_Find(t *testing.T) {
	// Test setup
	// Build features
	var createFeature1 = db.Feature{Feature_Name: "Canopy"}
	seedError := seedFeaturesDb([]db.Feature{createFeature1})
	if seedError != nil {
		t.Fatal("Failed to seed database for property find all test")
	}
	// Build test property
	propToCreate := &models.CreateProperty{
		Postcode:         14024,
		Property_Name:    "Kintamani No.2",
		Suburb:           "Kelapa Gading",
		City:             "Jakarta Utara",
		Street_Address_1: "Jl. Kintamani Raya no. 2",
		Street_Address_2: "Bukit Gading Villa",
		Bedrooms:         5,
		Bathrooms:        6,
		Land_Area:        400,
		Land_Metric:      "sqm",
		Description:      "A family home",
		Notes:            "The King slayer",
		// Use the first furnished property feature as relationship
		Features: []db.Feature{{ID: testConnection.features.created[0].ID}},
	}

	// Create property for test
	createdProp, err := testConnection.properties.serv.Create(propToCreate)
	if err != nil {
		t.Fatalf("failed to create test property for find by id user service test: %v", err)
	}

	// Create a request with an "id" URL parameter
	requestUrl := fmt.Sprintf("/api/properties/%v", createdProp.ID)
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		t.Fatal(err)
	}

	// // Add auth token to header
	req.Header.Set("Authorization", fmt.Sprintf("bearer %v", testConnection.accounts.admin.token))
	// // Create a response recorder
	rr := httptest.NewRecorder()

	// // Use handler with recorder and created request
	testConnection.router.ServeHTTP(rr, req)

	// Check the response status code is OK
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	// Extract JSON
	var body db.Property
	json.Unmarshal(rr.Body.Bytes(), &body)

	// Check details of returned property
	checkPropDetails(&body, createdProp, t, false)

	// Cleanup
	// Delete the created fixtures
	testConnection.dbClient.Delete(createdProp)
	testConnection.dbClient.Delete(createFeature1)
}

func TestPropertyController_FindAll(t *testing.T) {
	// Test setup
	// Build features
	var createFeature1 = db.Feature{Feature_Name: "Balcony"}
	var createFeature2 = db.Feature{Feature_Name: "Kitchen Island"}
	seedError := seedFeaturesDb([]db.Feature{createFeature1, createFeature2})
	if seedError != nil {
		t.Fatal("Failed to seed database for property find all test", seedError.Error())
	}

	// Create list of properties in db
	var listOfProperties = []db.Property{
		{
			Postcode:         14024,
			Property_Name:    "Kintamani No.12",
			Suburb:           "Kelapa Gading",
			City:             "Jakarta Utara",
			Street_Address_1: "Jl. Kintamani Raya no. 2",
			Street_Address_2: "Bukit Gading Villa",
			Bedrooms:         5,
			Bathrooms:        6,
			Land_Area:        400,
			Land_Metric:      "sqm",
			Description:      "A family home",
			Notes:            "The King slayer",
			// Use the first furnished property feature as relationship
			Features: []db.Feature{createFeature1, createFeature2},
		},
		{
			Postcode:         14024,
			Property_Name:    "Kintamani No.6",
			Suburb:           "Kelapa Gading",
			City:             "Jakarta Utara",
			Street_Address_1: "Jl. Kintamani Raya no. 2",
			Street_Address_2: "Bukit Gading Villa",
			Bedrooms:         5,
			Bathrooms:        6,
			Land_Area:        400,
			Land_Metric:      "sqm",
			Description:      "A family home",
			Notes:            "The King slayer",
			// Use the first furnished property feature as relationship
			Features: []db.Feature{createFeature1},
		},
	}
	createError := testConnection.dbClient.Create(listOfProperties)
	if createError.Error != nil {
		t.Fatal("Failed to seed database for Property Find All test", createError.Error)
	}
	// Add to state
	testConnection.properties.created = listOfProperties

	// Create a new request
	req, err := http.NewRequest("GET", "/api/properties?limit=10&offset=0&order=", nil)
	if err != nil {
		t.Fatal(err)
	}
	// Add auth token to header
	req.Header.Set("Authorization", fmt.Sprintf("bearer %v", testConnection.accounts.admin.token))
	// Create a response recorder
	rr := httptest.NewRecorder()

	// Use handler with recorder and created request
	testConnection.router.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Convert response JSON to struct
	var body []db.Property
	json.Unmarshal(rr.Body.Bytes(), &body)

	// Check length of property array (should be two with created furnishings)
	if len(body) != len(listOfProperties) {
		t.Errorf("Property array in findAll failed: expected %d, got %d", len(listOfProperties), len(body))
	}

	// Iterate through properties array received
	for _, item := range body {
		// If id is admin id
		if item.ID == testConnection.properties.created[0].ID {
			// Check details
			checkPropDetails(&item, &testConnection.properties.created[0], t, true)

		} else if item.ID == testConnection.properties.created[1].ID {
			// Else check property details
			checkPropDetails(&item, &testConnection.properties.created[1], t, true)
		}
	}

	// Test parameter inputs
	//
	var failParameterTests = []struct {
		limit                  string
		offset                 string
		order                  string
		expectedResponseStatus int
	}{
		// Bad order by
		{limit: "10", offset: "", order: "none", expectedResponseStatus: http.StatusBadRequest},
		// No limit should result in bad request
		{limit: "", offset: "", order: "", expectedResponseStatus: http.StatusBadRequest},
		// Check normal parameters functional with order by
		{limit: "20", offset: "1", order: "ID ASC", expectedResponseStatus: http.StatusOK},
	}
	for _, v := range failParameterTests {
		request := fmt.Sprintf("/api/properties?limit=%v&offset=%v&order=%v", v.limit, v.offset, v.order)
		// Create a new request
		req, err := http.NewRequest("GET", request, nil)
		if err != nil {
			t.Fatal(err)
		}
		// Add auth token to header
		req.Header.Set("Authorization", fmt.Sprintf("bearer %v", testConnection.accounts.admin.token))
		// Create a response recorder
		rr := httptest.NewRecorder()

		// Use handler with recorder and created request
		testConnection.router.ServeHTTP(rr, req)

		// Check the response status code
		if status := rr.Code; status != v.expectedResponseStatus {
			t.Errorf("handler returned wrong status code: got %v want %v",
				status, v.expectedResponseStatus)
		}
	}

	// Delete created features
	deleteResult := testConnection.dbClient.Delete(testConnection.features.created)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded features for find all")
	}
	// Delete created properties
	deleteResult = testConnection.dbClient.Delete(listOfProperties)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded features for find all")
	}
}

func TestPropertyController_Delete(t *testing.T) {
	// Test setup
	// Build features
	var createProperty = &db.Property{
		Postcode:         14024,
		Property_Name:    "Kintamani No.7",
		Suburb:           "Kelapa Gading",
		City:             "Jakarta Utara",
		Street_Address_1: "Jl. Kintamani Raya no. 2",
		Street_Address_2: "Bukit Gading Villa",
		Bedrooms:         5,
		Bathrooms:        6,
		Land_Area:        400,
		Land_Metric:      "sqm",
		Description:      "A family home",
		Notes:            "The King slayer",
	}
	// Create property for deletion test
	createResult := testConnection.dbClient.Create(createProperty)
	if createResult.Error != nil {
		t.Fatalf("Failed to seed database for Property DELETE test: %v", createResult.Error)
	}

	// Create a request with an "id" URL parameter
	requestUrl := fmt.Sprintf("/api/properties/%v", createProperty.ID)
	req, err := http.NewRequest("DELETE", requestUrl, nil)
	if err != nil {
		t.Fatal(err)
	}
	// Array of tests
	var deleteTests = []struct {
		tokenToUse             string
		expectedResponseStatus int
	}{
		// Tests of deletion failure
		{tokenToUse: testConnection.accounts.user.token, expectedResponseStatus: http.StatusForbidden},
		// Must be last
		// Tests of deletion success using admin priveleges
		{tokenToUse: testConnection.accounts.admin.token, expectedResponseStatus: http.StatusOK},
	}

	// Iterate through tests
	for _, test := range deleteTests {
		// Create a response recorder
		rr := httptest.NewRecorder()

		// Add auth token to header
		req.Header.Set("Authorization", fmt.Sprintf("bearer %v", test.tokenToUse))
		// Send deletion requestion to mock server
		testConnection.router.ServeHTTP(rr, req)
		// Check response is as expected
		if status := rr.Code; status != test.expectedResponseStatus {
			t.Errorf("User deletion test: got %v want %v.",
				status, test.expectedResponseStatus)
		}
	}
}

func TestPropertyController_Update(t *testing.T) {
	// Test setup
	var createProperty = &db.Property{
		Postcode:         14024,
		Property_Name:    "Kumran tip",
		Suburb:           "Kelapa Gading",
		City:             "Jakarta Utara",
		Street_Address_1: "Jl. Kintamani Raya no. 2",
		Street_Address_2: "Bukit Gading Villa",
		Bedrooms:         5,
		Bathrooms:        6,
		Land_Area:        400,
		Land_Metric:      "sqm",
		Description:      "A family home",
		Notes:            "The King slayer",
	}
	// Create property for update test
	createResult := testConnection.dbClient.Create(&createProperty)
	if createResult.Error != nil {
		t.Fatalf("Failed to seed database for property update test: %v", createResult.Error)
	}

	// Build test array
	var updateTests = []struct {
		data                   map[string]string
		tokenToUse             string
		expectedResponseStatus int
		checkDetails           bool
	}{
		{map[string]string{
			"Suburb":           "Kelapa Ganjing",
			"City":             "Jakarta Sutama",
			"Street_Address_1": "Jl. Burghul Raya no. 2",
			"Street_Address_2": "Swindletown",
		}, testConnection.accounts.user.token, http.StatusForbidden, false},
		{map[string]string{
			"Suburb":           "Kelapa Ganjing",
			"City":             "Jakarta Sutama",
			"Street_Address_1": "Jl. Burghul Raya no. 2",
			"Street_Address_2": "Swindletown",
		}, testConnection.accounts.admin.token, http.StatusOK, true},
		// Update should be disallowed due to being too short
		{map[string]string{
			"City":             "Tom",
			"Street_Address_1": "Crisp",
		}, testConnection.accounts.admin.token, http.StatusBadRequest, false},
		// User should be forbidden before validating
		{map[string]string{
			"City":             "que",
			"Street_Address_1": "solu",
		}, testConnection.accounts.user.token, http.StatusForbidden, false},
	}

	// Create a request url with an "id" URL parameter
	requestUrl := fmt.Sprintf("/api/properties/%v", createProperty.ID)

	// Iterate through update tests
	for _, v := range updateTests {
		// Make new request with property update in body
		req, err := http.NewRequest("PUT", requestUrl, buildReqBody(v.data))
		if err != nil {
			t.Fatal(err)
		}
		// Create a response recorder
		rr := httptest.NewRecorder()
		// Add auth token to header
		req.Header.Set("Authorization", fmt.Sprintf("bearer %v", v.tokenToUse))
		// Send update request to mock server
		testConnection.router.ServeHTTP(rr, req)
		// Check response expected vs received
		if status := rr.Code; status != v.expectedResponseStatus {
			t.Errorf("User update test: got %v want %v.",
				status, v.expectedResponseStatus)
		}

		// If need to check details
		if v.checkDetails == true {
			// Convert response JSON to struct
			var body db.Property
			json.Unmarshal(rr.Body.Bytes(), &body)

			// Update created property struct with the changes pushed through API
			createdProperty, err := updatePropChangesOnly(createProperty, v.data)
			if err != nil {
				t.Fatalf("Error updating property changes only in Update test: %v", err)
			}

			// Check user details using updated object
			checkPropDetails(&body, createdProperty, t, true)
		}
	}

	// Check for failure if incorrect ID parameter detected
	//
	var failUpdateTests = []struct {
		urlExtension           string
		expectedResponseStatus int
	}{
		// alpha character instead
		{urlExtension: "x", expectedResponseStatus: http.StatusForbidden},
		// Index out of bounds
		{urlExtension: "9", expectedResponseStatus: http.StatusBadRequest},
	}
	for _, v := range failUpdateTests {
		// Make new request with property update in body
		req, err := http.NewRequest("PUT", fmt.Sprint("/api/properties/"+v.urlExtension), buildReqBody(&db.Property{
			Postcode: 12341,
		}))
		if err != nil {
			t.Fatal(err)
		}
		// Create a response recorder
		rr := httptest.NewRecorder()
		// Add auth token to header
		req.Header.Set("Authorization", fmt.Sprintf("bearer %v", testConnection.accounts.admin.token))

		// Send update request to mock server
		testConnection.router.ServeHTTP(rr, req)
		// Check response is forbidden for
		if status := rr.Code; status != v.expectedResponseStatus {
			t.Errorf("Property update test: got %v want %v.",
				status, v.expectedResponseStatus)
		}
	}

	// Delete the created property
	testConnection.dbClient.Delete(createProperty)
}

func TestPropertyController_Create(t *testing.T) {
	// Setup
	//
	// Build features
	var createFeature1 = db.Feature{Feature_Name: "Open plan living"}
	var createFeature2 = db.Feature{Feature_Name: "Private Jacuzzi"}
	seedError := seedFeaturesDb([]db.Feature{createFeature1, createFeature2})
	if seedError != nil {
		t.Fatal("Failed to seed database for Property Find All test: ", seedError)
	}

	var createTests = []struct {
		data                   models.CreateProperty
		expectedResponseStatus int
		tokenToUse             string
	}{
		// Basic test
		{models.CreateProperty{
			Postcode:         14024,
			Property_Name:    "Shintamani No.2",
			Suburb:           "Kelapa Gading",
			City:             "Jakarta Utara",
			Street_Address_1: "Jl. Kintamani Raya no. 2",
			Street_Address_2: "Bukit Gading Villa",
			Bedrooms:         5,
			Bathrooms:        6,
			Land_Area:        400,
			Land_Metric:      "sqm",
			Description:      "A family home",
			Notes:            "The King slayer",
			Features:         []db.Feature{createFeature1, createFeature2},
		}, http.StatusForbidden, testConnection.accounts.user.token},
		// Should pass due to role of admin
		{models.CreateProperty{
			Postcode:         14024,
			Property_Name:    "Bazilarian",
			Suburb:           "Kelapa Gading",
			City:             "Kota Butara",
			Street_Address_1: "Jl. Gg. Sapi Kerbau",
			Street_Address_2: "Bukit Gading Villa",
			Bedrooms:         5,
			Bathrooms:        6,
			Land_Area:        400,
			Land_Metric:      "sqm",
			Description:      "A family home",
			Notes:            "The King slayer",
			Features:         []db.Feature{createFeature1, createFeature2},
		}, http.StatusCreated, testConnection.accounts.admin.token},
		// Create should be disallowed due to being too short
		{models.CreateProperty{
			Property_Name:    "go",
			Street_Address_1: "Jl.",
			Suburb:           "ga",
			City:             "Kota Butara",
		}, http.StatusBadRequest, testConnection.accounts.admin.token},
		// Should be a bad request due to duplicate property
		{models.CreateProperty{
			Property_Name:    "Bazilarian",
			Street_Address_1: "Jl. Gg. Sapi Kerbau",
		}, http.StatusBadRequest, testConnection.accounts.admin.token},
	}

	// Create a request url with an "id" URL parameter
	requestUrl := "/api/properties"

	for _, v := range createTests {
		// Make new request with property update in body
		req, err := http.NewRequest("POST", requestUrl, buildReqBody(v.data))
		if err != nil {
			t.Fatal(err)
		}
		// Create a response recorder
		rr := httptest.NewRecorder()
		// Add auth token to header
		req.Header.Set("Authorization", fmt.Sprintf("bearer %v", v.tokenToUse))

		// Send update request to mock server
		testConnection.router.ServeHTTP(rr, req)
		// Check response is as expected
		if status := rr.Code; status != v.expectedResponseStatus {
			t.Errorf("Property create test (%v): got %v want %v.", v.data.Property_Name,
				status, v.expectedResponseStatus)
		}

		// Init body for response extraction
		var body db.Property
		var foundProp db.Property
		// Grab ID from response body
		json.Unmarshal(rr.Body.Bytes(), &body)

		// Find the created property (to obtain full data with ID)
		testConnection.dbClient.Find(foundProp, uint(body.ID))

		// Compare found prop details with those found in returned body
		checkPropDetails(&body, &foundProp, t, true)

		// Delete the created Property
		delResult := testConnection.dbClient.Delete(&db.User{}, uint(body.ID))
		if delResult.Error != nil {
			t.Fatalf("Issue encountered deleting seeded assets for Prop create test (%v): %v", v.data.Property_Name, delResult.Error)
		}
	}
}

// Updates the parameter property struct with the updated values
func updatePropChangesOnly(createdProp *db.Property, updatedProp map[string]string) (*db.Property, error) {
	// Iterate through map and change struct values
	for k, v := range updatedProp {
		// Update each struct field using map
		err := helpers.UpdateStructField(createdProp, k, v)
		if err != nil {
			return nil, err
		}
	}
	return createdProp, nil
}

// Check the prop details
func checkPropDetails(actual *db.Property, expected *db.Property, t *testing.T, checkId bool) {
	// Only check ID if parameter checkId is true
	if checkId == true {
		// Verify that actual id is as expected
		if actual.ID != expected.ID {
			t.Errorf("found createdProperty has incorrect ID: expected %d, got %d", expected.ID, actual.ID)
		}
	}
	// Verify that details are expected
	if actual.Property_Name != expected.Property_Name {
		t.Errorf("found expected has incorrect Property_Name: expected %s, got %s", expected.Property_Name, actual.Property_Name)
	}
	if actual.Suburb != expected.Suburb {
		t.Errorf("found expected has incorrect Suburb: expected %s, got %s", expected.Suburb, actual.Suburb)
	}
	if actual.City != expected.City {
		t.Errorf("found expected has incorrect City: expected %s, got %s", expected.City, actual.City)
	}
	if actual.Street_Address_1 != expected.Street_Address_1 {
		t.Errorf("found expected has incorrect Street_Address_1: expected %s, got %s", expected.Street_Address_1, actual.Street_Address_1)
	}
	if actual.Street_Address_2 != expected.Street_Address_2 {
		t.Errorf("found expected has incorrect Street_Address_2: expected %s, got %s", expected.Street_Address_2, actual.Street_Address_2)
	}
	if actual.Bedrooms != expected.Bedrooms {
		t.Errorf("found expected has incorrect Bedrooms: expected %v, got %v", expected.Bedrooms, actual.Bedrooms)
	}
	if actual.Bathrooms != expected.Bathrooms {
		t.Errorf("found expected has incorrect Bathrooms: expected %v, got %v", expected.Bathrooms, actual.Bathrooms)
	}
	if actual.Land_Area != expected.Land_Area {
		t.Errorf("found expected has incorrect Land_Area: expected %v, got %v", expected.Land_Area, actual.Land_Area)
	}
	if actual.Land_Metric != expected.Land_Metric {
		t.Errorf("found expected has incorrect Land_Metric: expected %s, got %s", expected.Land_Metric, actual.Land_Metric)
	}
	if actual.Description != expected.Description {
		t.Errorf("found expected has incorrect Description: expected %s, got %s", expected.Description, actual.Description)
	}
	if actual.Notes != expected.Notes {
		t.Errorf("found expected has incorrect Notes: expected %s, got %s", expected.Notes, actual.Notes)
	}

	// Property Features
	// Iterate through properties array received
	for _, feature := range actual.Features {
		// Iterate through property features created for this test
		for _, createdProp := range testConnection.features.created {
			// If the feature matches the created one iterating through
			if feature.ID == createdProp.ID {
				// Check the details of the property feature
				if createdProp.Feature_Name != feature.Feature_Name {
					t.Fatalf("Feature name does not match")
				}
			}
		}
	}
}

// SETUP FUNCTIONS
//
// Seeds property features into database and adds details to testConnection
func seedFeaturesDb(featuresToCreate []db.Feature) error {
	// Create features
	result := testConnection.dbClient.Create(&featuresToCreate)
	if result.Error != nil {
		fmt.Printf("Encountered issue while seeding features")
		return result.Error
	}
	// If successful, store saved features
	testConnection.features.created = featuresToCreate
	return nil
}
