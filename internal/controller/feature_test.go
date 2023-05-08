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

func TestFeatureController_FindAll(t *testing.T) {
	// Test setup
	// Build features
	var createFeature1 = db.Feature{Feature_Name: "Bohemian"}
	var createFeature2 = db.Feature{Feature_Name: "Enclosed living area"}
	var createdFeatures = []db.Feature{createFeature1, createFeature2}
	seedError := seedFeaturesDb(createdFeatures)
	if seedError != nil {
		t.Fatal("Failed to seed database for Feature Find All test", seedError.Error())
	}
	// Add to state
	testConnection.features.created = createdFeatures

	// Create a new request
	req, err := http.NewRequest("GET", "/api/features?limit=10&offset=0&order=", nil)
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
	var body []db.Feature
	json.Unmarshal(rr.Body.Bytes(), &body)

	// Check length of feature array (should be two with seeded assets)
	if len(body) != len(createdFeatures) {
		t.Errorf("Features array in findAll failed: expected %d, got %d", len(createdFeatures), len(body))
	}

	// Iterate through feature array received
	for _, actualFeat := range body {
		// Iterate through created features to determine a match
		for _, createdFeat := range createdFeatures {
			// If match found
			if actualFeat.ID == createdFeat.ID {
				// Check the details of the feature
				checkFeatureDetails(&actualFeat, &createdFeat, t, false)
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
		request := fmt.Sprintf("/api/features?limit=%v&offset=%v&order=%v", v.limit, v.offset, v.order)
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
	deleteResult := testConnection.dbClient.Delete(createdFeatures)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded features for find all: %v", deleteResult.Error)
	}
}

func TestFeatureController_Find(t *testing.T) {
	// Test setup
	// Build feature
	var createFeature = db.Feature{Feature_Name: "Open air living"}
	// Seed features and add to state
	seedError := seedFeaturesDb([]db.Feature{createFeature})
	if seedError != nil {
		t.Fatal("Failed to seed database for property feature find all test")
	}

	// Create a request with an "id" URL parameter
	requestUrl := fmt.Sprintf("/api/features/%v", testConnection.features.created[0].ID)
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

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
	var body db.Feature
	json.Unmarshal(rr.Body.Bytes(), &body)
	checkFeatureDetails(&body, &createFeature, t, false)

	// Cleanup
	// Delete the created feature
	deleteResult := testConnection.dbClient.Delete(testConnection.features.created)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded features for find all: %v", deleteResult.Error)
	}
}

func TestFeatureController_Delete(t *testing.T) {
	// Test setup
	// Build feature
	var createFeature = db.Feature{Feature_Name: "Soccer field"}
	// Seed features and add to state
	seedError := seedFeaturesDb([]db.Feature{createFeature})
	if seedError != nil {
		t.Fatal("Failed to seed database for Prop feature Find Delete test")
	}

	// Create a request with an "id" URL parameter
	requestUrl := fmt.Sprintf("/api/features/%v", testConnection.features.created[0].ID)
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
		{testName: "Prop feature basic user delete test", tokenToUse: testConnection.accounts.user.token, expectedResponseStatus: http.StatusForbidden},
		// Must be last
		// Tests of deletion success using admin priveleges
		{testName: "Prop feature admin delete test", tokenToUse: testConnection.accounts.admin.token, expectedResponseStatus: http.StatusOK},
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
}

func TestFeatureController_Update(t *testing.T) {
	// Test setup
	// Build feature
	var createFeature = db.Feature{Feature_Name: "Retractable canopy"}
	// Seed features and add to state
	seedError := seedFeaturesDb([]db.Feature{createFeature})
	if seedError != nil {
		t.Fatal("Failed to seed database for Prop feature Find Delete test")
	}

	// Build test array
	var updateTests = []struct {
		data                   db.Feature
		tokenToUse             string
		expectedResponseStatus int
		checkDetails           bool
	}{
		{db.Feature{Feature_Name: "Kanjing"}, testConnection.accounts.user.token, http.StatusForbidden, false},
		{db.Feature{Feature_Name: "Kanjing"}, testConnection.accounts.admin.token, http.StatusOK, true},
		// Update should be disallowed due to being too short
		{db.Feature{Feature_Name: "Ka"}, testConnection.accounts.admin.token, http.StatusBadRequest, false},
		// User should be forbidden before validating rather than Bad Request
		{db.Feature{Feature_Name: "Kanjing"}, testConnection.accounts.user.token, http.StatusForbidden, false},
	}

	// Create a request url with an "id" URL parameter
	requestUrl := fmt.Sprintf("/api/features/%v", testConnection.features.created[0].ID)

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
		var body db.Feature
		json.Unmarshal(rr.Body.Bytes(), &body)
		// Check response expected vs received
		if status := rr.Code; status != v.expectedResponseStatus {
			t.Errorf("Property feature update test: got %v want %v.",
				status, v.expectedResponseStatus)
		}

		// If need to check details
		if v.checkDetails == true {

			// Check prop feature details using updated object
			checkFeatureDetails(&body, &v.data, t, false)
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
		// Make new request with feature update in body
		req, err := http.NewRequest("PUT", fmt.Sprint("/api/features/"+v.urlExtension), buildReqBody(&db.Feature{
			Feature_Name: "Scrappy Kid",
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
			t.Errorf("Property feature update test: got %v want %v.",
				status, v.expectedResponseStatus)
		}
	}

	// Delete the created property feature
	deleteResult := testConnection.dbClient.Delete(testConnection.features.created)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded features for find all: %v", deleteResult.Error)
	}
}

func TestFeatureController_Create(t *testing.T) {
	// Setup
	//
	var createTests = []struct {
		data                   models.CreateFeature
		expectedResponseStatus int
		tokenToUse             string
	}{
		// Should fail due to user role status of basic
		{models.CreateFeature{
			Feature_Name: "Locked entry",
		}, http.StatusForbidden, testConnection.accounts.user.token},
		// Should pass as user is admin
		{models.CreateFeature{
			Feature_Name: "Locked entry",
		}, http.StatusCreated, testConnection.accounts.admin.token},
		// Create should be disallowed due to being too short
		{models.CreateFeature{
			Feature_Name: "Loc",
		}, http.StatusBadRequest, testConnection.accounts.admin.token},
		// Should be a bad request due to duplicate feature
		{models.CreateFeature{
			Feature_Name: "Locked entry",
		}, http.StatusBadRequest, testConnection.accounts.admin.token},
	}

	// Create a request url with an "id" URL parameter
	requestUrl := "/api/features"

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
			t.Errorf("Property create test (%v): got %v want %v.", v.data.Feature_Name,
				status, v.expectedResponseStatus)
		}

		// Init body for response extraction
		var body db.Feature
		var foundFeat db.Feature
		// Grab ID from response body
		json.Unmarshal(rr.Body.Bytes(), &body)

		// Find the created property feature (to obtain full data with ID)
		testConnection.dbClient.Find(foundFeat, uint(body.ID))

		// Compare found prop details with those found in returned body
		checkFeatureDetails(&body, &foundFeat, t, true)

		// Delete the created feature
		delResult := testConnection.dbClient.Delete(&db.Feature{}, uint(body.ID))
		if delResult.Error != nil {
			t.Fatalf("Issue encountered deleting seeded assets for Feature create test (%v): %v", v.data.Feature_Name, delResult.Error)
		}
	}
}

// Check the property feature details (username, name, email and ID)
func checkFeatureDetails(actual *db.Feature, expected *db.Feature, t *testing.T, checkId bool) {

	// Only check ID if parameter checkId is true
	if checkId == true {
		// Verify that the actual prop feature id matches the created prop features'
		if actual.ID != expected.ID {
			t.Errorf("found feature has incorrect ID: expected %d, got %d", expected.ID, actual.ID)
		}
	}

	// Verify that the actual prop feature name matches the original feature name
	if actual.Feature_Name != expected.Feature_Name {
		t.Errorf("found feature has incorrect Property_Name: expected %s, got %s", expected.Feature_Name, actual.Feature_Name)
	}
}
