package controller_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/models"
)

func TestPropertyLogController_FindAll(t *testing.T) {
	// Test setup
	// Create property
	propToCreate1 := &db.Property{
		Postcode:         14024,
		Property_Name:    "Dalung No.10",
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
	createdProperties := []db.Property{*propToCreate1}
	createResult := testConnection.dbClient.Create(createdProperties)
	if createResult.Error != nil {
		t.Fatal("Failed to create property for property log find all test: ", createResult.Error)
	}
	// Add to state
	testConnection.properties.created = createdProperties
	// Create property logs
	propLogToCreate1 := &db.PropertyLog{
		PropertyID: createdProperties[0].ID,
		LogMessage: "This is a note",
	}
	propLogToCreate2 := &db.PropertyLog{
		PropertyID: createdProperties[0].ID,
		LogMessage: "This is a second note",
	}
	createdPropLogs := []db.PropertyLog{*propLogToCreate1, *propLogToCreate2}
	// Create property logs in db
	createResult = testConnection.dbClient.Create(createdPropLogs)
	if createResult.Error != nil {
		t.Fatal("Failed to create property logs for property log find all test: ", createResult.Error)
	}
	// Add to state
	testConnection.propertyLogs.created = createdPropLogs

	// Create a new request
	req, err := http.NewRequest("GET", "/api/property-logs?limit=10&offset=0&order=", nil)
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
	var body []db.PropertyLog
	json.Unmarshal(rr.Body.Bytes(), &body)

	// Check length of feature array (should be two with seeded assets)
	if len(body) != len(createdPropLogs) {
		t.Errorf("Features array in findAll failed: expected %d, got %d", len(createdPropLogs), len(body))
	}

	// Iterate through feature array received
	for _, actualPropLog := range body {
		// Iterate through created features to determine a match
		for _, createdPropLog := range createdPropLogs {
			// If match found
			if actualPropLog.ID == createdPropLog.ID {
				// Check the details of the feature
				checkPropertyLogDetails(&actualPropLog, &createdPropLog, t, true)
			}
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
	// Iterate through URL parameter tests
	for _, v := range failParameterTests {
		request := fmt.Sprintf("/api/property-logs?limit=%v&offset=%v&order=%v", v.limit, v.offset, v.order)
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

	// Clean up created fixtures
	deleteResult := testConnection.dbClient.Delete(createdPropLogs)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded prop logs: %v", deleteResult.Error)
	}
	deleteResult = testConnection.dbClient.Delete(createdProperties)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded props: %v", deleteResult.Error)
	}
}

func TestPropertyLogController_Find(t *testing.T) {
	// Test setup
	// Create property
	propToCreate1 := &db.Property{
		Postcode:         14024,
		Property_Name:    "Dalung No.24A",
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
	createdProperties := []db.Property{*propToCreate1}
	createResult := testConnection.dbClient.Create(createdProperties)
	if createResult.Error != nil {
		t.Fatal("Failed to create property for property log find all test: ", createResult.Error)
	}
	// Add to state
	testConnection.properties.created = createdProperties
	// Create property logs
	propLogToCreate1 := &db.PropertyLog{
		PropertyID: createdProperties[0].ID,
		LogMessage: "This is a note",
	}
	createdPropLogs := []db.PropertyLog{*propLogToCreate1}
	// Create property logs in db
	createResult = testConnection.dbClient.Create(createdPropLogs)
	if createResult.Error != nil {
		t.Fatal("Failed to create property logs for property log find all test: ", createResult.Error)
	}
	// Add to state
	testConnection.propertyLogs.created = createdPropLogs

	// Create a request with an "id" URL parameter
	requestUrl := fmt.Sprintf("/api/property-logs/%v", testConnection.propertyLogs.created[0].ID)
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		t.Fatal(err)
	}

	// // Add auth token to header
	req.Header.Set("Authorization", fmt.Sprintf("bearer %v", testConnection.accounts.admin.token))
	// // Create a response recorder
	rr := httptest.NewRecorder()

	// // Serve request using recorder and created request
	testConnection.router.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Extract the response body
	var body db.PropertyLog
	json.Unmarshal(rr.Body.Bytes(), &body)
	checkPropertyLogDetails(&body, &createdPropLogs[0], t, true)

	// Cleanup
	// Delete the created property logs
	deleteResult := testConnection.dbClient.Delete(createdPropLogs)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded prop logs: %v", deleteResult.Error)
	}
	// Delete the created properties
	deleteResult = testConnection.dbClient.Delete(createdProperties)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded props: %v", deleteResult.Error)
	}
}

func TestPropertyLogController_Delete(t *testing.T) {
	// Test setup
	// Create property
	propToCreate1 := &db.Property{
		Postcode:         14024,
		Property_Name:    "SmakDown No.24A",
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
	createdProperties := []db.Property{*propToCreate1}
	createResult := testConnection.dbClient.Create(createdProperties)
	if createResult.Error != nil {
		t.Fatal("Failed to create property for property log find all test: ", createResult.Error)
	}
	// Add to state
	testConnection.properties.created = createdProperties
	// Create property logs
	propLogToCreate1 := &db.PropertyLog{
		PropertyID: createdProperties[0].ID,
		LogMessage: "This is a note",
	}
	createdPropLogs := []db.PropertyLog{*propLogToCreate1}
	// Create property logs in db
	createResult = testConnection.dbClient.Create(createdPropLogs)
	if createResult.Error != nil {
		t.Fatal("Failed to create property logs for property log find all test: ", createResult.Error)
	}
	// Add to state
	testConnection.propertyLogs.created = createdPropLogs

	// Create a request with an "id" URL parameter
	requestUrl := fmt.Sprintf("/api/property-logs/%v", testConnection.propertyLogs.created[0].ID)
	req, err := http.NewRequest("DELETE", requestUrl, nil)
	if err != nil {
		t.Fatal(err)
	}
	// Array of tests
	var deleteTests = []struct {
		testName               string
		tokenToUse             string
		expectedResponseStatus int
	}{
		// Tests of deletion failure
		{testName: "Prop log basic user delete test", tokenToUse: testConnection.accounts.user.token, expectedResponseStatus: http.StatusForbidden},
		// Must be last
		// Tests of deletion success using admin privileges
		{testName: "Prop log admin delete test", tokenToUse: testConnection.accounts.admin.token, expectedResponseStatus: http.StatusOK},
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
			t.Errorf("Prop feature deletion test (%v): got %v want %v.", test.testName,
				status, test.expectedResponseStatus)
		}
	}
	// Cleanup
	// Delete the created property logs
	deleteResult := testConnection.dbClient.Delete(createdPropLogs)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded property logs: %v", deleteResult.Error)
	}
	// Delete the created properties
	deleteResult = testConnection.dbClient.Delete(createdProperties)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded props: %v", deleteResult.Error)
	}
}

func TestPropertyLogController_Update(t *testing.T) {
	// Test setup
	// Create property
	propToCreate1 := &db.Property{
		Postcode:         14024,
		Property_Name:    "Mutasada No.24A",
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
	createdProperties := []db.Property{*propToCreate1}
	createResult := testConnection.dbClient.Create(createdProperties)
	if createResult.Error != nil {
		t.Fatal("Failed to create property for property log find all test: ", createResult.Error)
	}
	// Add to state
	testConnection.properties.created = createdProperties
	// Create property logs
	propLogToCreate1 := &db.PropertyLog{
		PropertyID: createdProperties[0].ID,
		LogMessage: "This is a note",
	}
	createdPropLogs := []db.PropertyLog{*propLogToCreate1}
	// Create property logs in db
	createResult = testConnection.dbClient.Create(createdPropLogs)
	if createResult.Error != nil {
		t.Fatal("Failed to create property logs for property log find all test: ", createResult.Error)
	}
	// Add to state
	testConnection.propertyLogs.created = createdPropLogs

	// Build test array
	var updateTests = []struct {
		data                   db.PropertyLog
		tokenToUse             string
		expectedResponseStatus int
		checkDetails           bool
	}{
		// Test of update failure
		{db.PropertyLog{
			LogMessage: "This is an updated note",
		}, testConnection.accounts.user.token, http.StatusForbidden, false},
		// Update should be allowed
		{db.PropertyLog{
			LogMessage: "This is an updated note",
		}, testConnection.accounts.admin.token, http.StatusOK, true},
		// Update should be disallowed due to being too short
		{db.PropertyLog{
			LogMessage: "go",
		}, testConnection.accounts.admin.token, http.StatusBadRequest, false},
		// User should be forbidden before validating rather than Bad Request
		{db.PropertyLog{
			LogMessage: "go",
		}, testConnection.accounts.user.token, http.StatusForbidden, false},
	}

	// Create a request url with an "id" URL parameter
	requestUrl := fmt.Sprintf("/api/property-logs/%v", createdPropLogs[0].ID)

	// Iterate through update tests
	for _, v := range updateTests {
		// Make new request with feature update in body
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

		// Convert response JSON to struct
		var body db.PropertyLog
		json.Unmarshal(rr.Body.Bytes(), &body)
		// Check response expected vs received
		if status := rr.Code; status != v.expectedResponseStatus {
			t.Errorf("Update test: got %v want %v. \nBody: %v", status, v.expectedResponseStatus, body)

		}

		// If need to check details
		if v.checkDetails == true {

			// Update the expected object with the new details
			createdPropLogs[0].LogMessage = v.data.LogMessage
			// Check prop feature details using updated object
			checkPropertyLogDetails(&body, &createdPropLogs[0], t, false)
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
		// Make new request with property log update in body
		req, err := http.NewRequest("PUT", fmt.Sprint("/api/property-logs/"+v.urlExtension), buildReqBody(&db.PropertyLog{
			LogMessage: "Gustav",
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
		// Check response is as expected
		if status := rr.Code; status != v.expectedResponseStatus {
			t.Errorf("Fail update test: got %v want %v.",
				status, v.expectedResponseStatus)
		}
	}

	// Cleanup
	// Delete the created property logs
	deleteResult := testConnection.dbClient.Delete(createdPropLogs)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded property logs: %v", deleteResult.Error)
	}
	// Delete the created properties
	deleteResult = testConnection.dbClient.Delete(createdProperties)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded props: %v", deleteResult.Error)
	}
}

func TestPropertyLogController_Create(t *testing.T) {
	// Setup
	//
	// Create property
	propToCreate1 := &db.Property{
		Postcode:         14024,
		Property_Name:    "Teuku Kumar No.24A",
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
	createdProperties := []db.Property{*propToCreate1}
	createResult := testConnection.dbClient.Create(createdProperties)
	if createResult.Error != nil {
		t.Fatal("Failed to create property for property log find all test: ", createResult.Error)
	}
	// Add to state
	testConnection.properties.created = createdProperties

	var createTests = []struct {
		data                   models.CreatePropertyLog
		expectedResponseStatus int
		tokenToUse             string
	}{
		// Should fail due to user role status of basic
		{models.CreatePropertyLog{
			LogMessage: "This is a note",
			Property:   createdProperties[0],
		}, http.StatusForbidden, testConnection.accounts.user.token},
		// Should pass as user is admin
		{models.CreatePropertyLog{
			LogMessage: "This is a note",
			Property:   createdProperties[0],
		}, http.StatusCreated, testConnection.accounts.admin.token},
		// Create should be disallowed due to being too short
		{models.CreatePropertyLog{
			LogMessage: "Ta",
			Property:   createdProperties[0],
		}, http.StatusBadRequest, testConnection.accounts.admin.token},
		// Create should be disallowed due to being too long
		{models.CreatePropertyLog{
			LogMessage: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Suspendisse vulputate, nunc sit amet efficitur bibendum, sapien odio auctor nisi, a interdum magna nisl ac purus. Fusce condimentum malesuada mi at eleifend. Sed laoreet varius risus, id mattis libero tristique nec. Sed eget malesuada magna. Morbi feugiat sapien euismod neque commodo suscipit. Vivamus vehicula euismod dui, id imperdiet elit lacinia non. Integer hendrerit, enim ac gravida malesuada, dolor leo dictum purus, nec bibendum velit est vel nulla. Nulla sagittis nulla non elit imperdiet convallis. Sed bibendum sollicitudin nunc, vel facilisis nulla convallis a. Nunc id ex feugiat, finibus magna sit amet, ultricies lacus.",
			Property:   createdProperties[0],
		}, http.StatusBadRequest, testConnection.accounts.admin.token},
	}

	// Create a request url with an "id" URL parameter
	requestUrl := "/api/property-logs"

	for _, v := range createTests {
		// Make new request with prop feature creation in body
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
			t.Errorf("Property log create test (%v): got %v want %v. \nBody: %v\n", v.data.LogMessage,
				status, v.expectedResponseStatus, rr.Body.String())
		}

		// Init body for response extraction
		var body db.PropertyLog
		var foundPropLog db.PropertyLog
		// Grab ID from response body
		json.Unmarshal(rr.Body.Bytes(), &body)

		// Find the created property feature (to obtain full data with ID)
		testConnection.dbClient.Find(foundPropLog, uint(body.ID))

		// Compare found prop details with those found in returned body
		checkPropertyLogDetails(&body, &foundPropLog, t, true)

		// If the property log was created successfully, check that it's deleted after test
		if v.expectedResponseStatus == http.StatusCreated {
			// Cleanup
			//
			// Delete the created property logs
			deleteResult := testConnection.dbClient.Delete(&db.PropertyLog{}, uint(body.ID))
			if deleteResult.Error != nil {
				t.Fatalf("Couldn't clean up created property logs: %v", deleteResult.Error)
			}
		}
	}

	// Overall test cleanup
	// Delete the created properties
	deleteResult := testConnection.dbClient.Delete(createdProperties)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded props: %v", deleteResult.Error)
	}
}

// Check the property feature details (username, name, email and ID)
func checkPropertyLogDetails(actual *db.PropertyLog, expected *db.PropertyLog, t *testing.T, checkId bool) {

	// Only check ID if parameter checkId is true
	if checkId == true {
		// Verify that the actual prop feature id matches the created prop features'
		if actual.ID != expected.ID {
			t.Errorf("found feature has incorrect ID: expected %d, got %d", expected.ID, actual.ID)
		}
	}

	// Verify that the actual prop feature name matches the original feature name
	if actual.LogMessage != expected.LogMessage {
		t.Errorf("found prop log message has incorrect log message: expected %s, got %s", expected.LogMessage, actual.LogMessage)
	}
	if actual.PropertyID != expected.PropertyID {
		t.Errorf("found prop log message has incorrect property ID: expected %d, got %d", expected.PropertyID, actual.PropertyID)
	}
	if actual.UserID != expected.UserID {
		t.Errorf("found prop log message has incorrect user ID: expected %d, got %d", expected.UserID, actual.UserID)
	}
}
