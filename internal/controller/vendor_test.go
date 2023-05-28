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

func TestVendorController_FindAll(t *testing.T) {
	// Test setup
	vendorToCreate1 := &db.Vendor{
		CompanyName:      "PT Wijaya Karya",
		NPWP:             "09098.2314.234.234",
		Email:            "gustav@wijayakarya.com",
		Phone:            "081234567890",
		NIB:              "1234567890",
		Street_Address_1: "Jl. Raya Jatinegara Barat No. 179",
		City:             "Jakarta Timur",
		Province:         "DKI Jakarta",
		Postal_Code:      "13310",
		Suburb:           "Jatinegara",
	}
	vendorToCreate2 := &db.Vendor{
		CompanyName:      "PT Wijaya Karyo",
		NPWP:             "09098.2314.234.234",
		Email:            "gustav@wijayakarya.com",
		Phone:            "081234567890",
		NIB:              "1234567890",
		Street_Address_1: "Jl. Raya Jatinegara Barat No. 179",
		City:             "Jakarta Timur",
		Province:         "DKI Jakarta",
		Postal_Code:      "13310",
		Suburb:           "Jatinegara"}
	createdVendors := []db.Vendor{*vendorToCreate1, *vendorToCreate2}
	createResult := testConnection.dbClient.Create(createdVendors)
	if createResult.Error != nil {
		t.Fatal("Failed to create vendors for test: ", createResult.Error)
	}

	// Create a new request
	req, err := http.NewRequest("GET", "/api/vendors?limit=10&offset=0&order=", nil)
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
	var body []db.Vendor
	json.Unmarshal(rr.Body.Bytes(), &body)

	// Check length of array (should be two with seeded assets)
	if len(body) != len(createdVendors) {
		t.Errorf("Array length check in findAll failed: expected %d, got %d", len(createdVendors), len(body))
	}

	// Iterate through array received
	for _, actual := range body {
		// Iterate through prior created items to determine a match
		for _, expected := range createdVendors {
			// If match found
			if actual.ID == expected.ID {
				// Check the details of the vendor match
				checkVendorDetails(&actual, &expected, t, false)
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
		request := fmt.Sprintf("/api/vendors?limit=%v&offset=%v&order=%v", v.limit, v.offset, v.order)
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
	deleteResult := testConnection.dbClient.Delete(createdVendors)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded work types: %v", deleteResult.Error)
	}
}

func TestVendorController_Find(t *testing.T) {
	// Test setup
	vendorToCreate1 := &db.Vendor{
		CompanyName:      "PT Wijaya Suroyaarya",
		NPWP:             "09098.2314.234.234",
		Email:            "gustav@wijayakarya.com",
		Phone:            "081234567890",
		NIB:              "1234567890",
		Street_Address_1: "Jl. Raya Jatinegara Barat No. 179",
		City:             "Jakarta Timur",
		Province:         "DKI Jakarta",
		Postal_Code:      "13310",
		Suburb:           "Jatinegara",
	}
	createdVendors := []db.Vendor{*vendorToCreate1}
	createResult := testConnection.dbClient.Create(createdVendors)
	if createResult.Error != nil {
		t.Fatal("Failed to create vendors for test: ", createResult.Error)
	}

	// Create a request with an "id" URL parameter
	requestUrl := fmt.Sprintf("/api/vendors/%v", createdVendors[0].ID)
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
		t.Errorf("handler returned wrong status code: got %v want %v. Error: %v",
			status, http.StatusOK, rr.Body.String())
	}

	// Extract the response body
	var body db.Vendor
	json.Unmarshal(rr.Body.Bytes(), &body)
	checkVendorDetails(&body, &createdVendors[0], t, true)

	// Cleanup
	// Delete the created fixtures
	deleteResult := testConnection.dbClient.Delete(createdVendors)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded vendors: %v", deleteResult.Error)
	}
}

func TestVendorController_Delete(t *testing.T) {
	// Test setup
	vendorToCreate1 := &db.Vendor{
		CompanyName:      "PT Wijaya Kartaa",
		NPWP:             "09098.2314.234.234",
		Email:            "gustav@wijayakarya.com",
		Phone:            "081234567890",
		NIB:              "1234567890",
		Street_Address_1: "Jl. Raya Jatinegara Barat No. 179",
		City:             "Jakarta Timur",
		Province:         "DKI Jakarta",
		Postal_Code:      "13310",
		Suburb:           "Jatinegara",
	}
	createdVendors := []db.Vendor{*vendorToCreate1}
	createResult := testConnection.dbClient.Create(createdVendors)
	if createResult.Error != nil {
		t.Fatal("Failed to create vendors for test: ", createResult.Error)
	}

	// Create a request with an "id" URL parameter
	requestUrl := fmt.Sprintf("/api/vendors/%v", createdVendors[0].ID)
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
		{testName: "Basic user delete test", tokenToUse: testConnection.accounts.user.token, expectedResponseStatus: http.StatusForbidden},
		// Must be last
		// Tests of deletion success using admin privileges
		{testName: "Admin delete test", tokenToUse: testConnection.accounts.admin.token, expectedResponseStatus: http.StatusOK},
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
			t.Errorf("Task log deletion test (%v): got %v want %v.", test.testName,
				status, test.expectedResponseStatus)
		}
	}
	// Cleanup
	// Delete the created fixtures
	deleteResult := testConnection.dbClient.Delete(createdVendors)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded vendors: %v", deleteResult.Error)
	}
}

func TestVendorController_Update(t *testing.T) {
	// Test setup
	vendorToCreate1 := &db.Vendor{
		CompanyName:      "PT Wijaya proauaa",
		NPWP:             "09098.2314.234.234",
		Email:            "gustav@wijayakarya.com",
		Phone:            "081234567890",
		NIB:              "1234567890",
		Street_Address_1: "Jl. Raya Jatinegara Barat No. 179",
		City:             "Jakarta Timur",
		Province:         "DKI Jakarta",
		Postal_Code:      "13310",
		Suburb:           "Jatinegara",
	}
	createdVendors := []db.Vendor{*vendorToCreate1}
	createResult := testConnection.dbClient.Create(createdVendors)
	if createResult.Error != nil {
		t.Fatal("Failed to create vendors for test: ", createResult.Error)
	}

	// Build test array
	var updateTests = []struct {
		data                   models.UpdateVendor
		tokenToUse             string
		expectedResponseStatus int
		checkDetails           bool
		testName               string
	}{
		// Test of update failure: basic user
		{models.UpdateVendor{
			City: "Jakarta Timur",
		}, testConnection.accounts.user.token, http.StatusForbidden, false, "Basic user update test"},
		// Update should be allowed: admin
		{models.UpdateVendor{
			City: "Jakarta Timur",
		}, testConnection.accounts.admin.token, http.StatusOK, true, "Admin update test"},
		// Update should be disallowed due to being invalid value for province
		{models.UpdateVendor{
			Province: "Pantat",
		}, testConnection.accounts.admin.token, http.StatusBadRequest, false, "Invalid province update test"},
		// User should be forbidden before validating rather than Bad Request
		{models.UpdateVendor{
			Province: "Pantat",
		}, testConnection.accounts.user.token, http.StatusForbidden, false, "Invalid type update when basic user test"},
	}

	// Create a request url with an "id" URL parameter
	requestUrl := fmt.Sprintf("/api/vendors/%v", createdVendors[0].ID)
	fmt.Printf("Request URL: %v", requestUrl)

	// Iterate through update tests
	for _, v := range updateTests {
		// Make new request with update in body
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
		var body db.Vendor
		json.Unmarshal(rr.Body.Bytes(), &body)

		// Check response expected vs received
		if status := rr.Code; status != v.expectedResponseStatus {
			t.Errorf("Update test: got %v want %v. \nBody: %v", status, v.expectedResponseStatus, rr.Body.String())
		}

		// If need to check details
		if v.checkDetails == true {
			// Get vendor details from database
			var expected db.Vendor
			findResult := testConnection.dbClient.Find(&expected, createdVendors[0].ID)
			if findResult.Error != nil {
				t.Errorf("Error finding updated vendor: %v", findResult.Error)
			}

			// Check details using updated object
			checkVendorDetails(&body, &expected, t, false)
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
		// Make new request with task log update in body
		req, err := http.NewRequest("PUT", fmt.Sprint("/api/vendors/"+v.urlExtension), buildReqBody(&db.Vendor{
			Email: "Jembut@walau.com",
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
	// Delete the created fixtures
	deleteResult := testConnection.dbClient.Delete(createdVendors)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded vendors: %v", deleteResult.Error)
	}
}

func TestVendorController_Create(t *testing.T) {
	var createTests = []struct {
		data                   models.CreateVendor
		expectedResponseStatus int
		tokenToUse             string
		testName               string
	}{
		// Should fail due to user role status of basic
		{models.CreateVendor{
			CompanyName: "PT widodo Jokowow",
			NPWP:        "09098.2314.234.234",
		}, http.StatusForbidden, testConnection.accounts.user.token, "basic user create"},
		// Should pass as user is admin
		{models.CreateVendor{
			CompanyName: "PT widodo Grw",
			NPWP:        "09098.2314.234.234",
		}, http.StatusCreated, testConnection.accounts.admin.token, "admin create"},
		// Should fail due to incorrect province value
		{models.CreateVendor{
			CompanyName: "PT widodo Grlkjfaww",
			NPWP:        "09098.2314.234.234",
			Province:    "Pantat",
		}, http.StatusBadRequest, testConnection.accounts.admin.token, "admin create"},
		// Create should be disallowed due to note being too long
		{models.CreateVendor{
			CompanyName: "PT widodo Slow",
			NPWP:        "09098.2314.234.234",
			Notes:       "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec euismod, nisl eget ultricies ultricies, nisl nisl luctus nisl, vitae aliquam nislLorem ipsum dolor sit amet, consectetur adipiscing elit. Donec euismod, nisl eget ultricies ultricies, nisl nisl luctus nisl, vitae aliquam nislLorem ipsum dolor sit amet, consectetur adipiscing elit. Donec euismod, nisl eget ultricies ultricies, nisl nisl luctus nisl, vitae aliquam nislLorem ipsum dolor sit amet, consectetur adipiscing elit. Donec euismod, nisl eget ultricies ultricies, nisl nisl luctus nisl, vitae aliquam nislLorem ipsum dolor sit amet, consectetur adipiscing elit. Donec euismod, nisl eget ultricies ultricies, nisl nisl luctus nisl, vitae aliquam nislLorem ipsum dolor sit amet, consectetur adipiscing elit. Donec euismod, nisl eget ultricies ultricies, nisl nisl luctus nisl, vitae aliquam nislLorem ipsum dolor sit amet, consectetur adipiscing elit. Donec euismod, nisl eget ultricies ultricies, nisl nisl luctus nisl, vitae aliquam nislLorem ipsum dolor sit amet, consectetur adipiscing elit. Donec euismod, nisl eget ultricies ultricies, nisl nisl luctus nisl, vitae aliquam nisl",
		}, http.StatusBadRequest, testConnection.accounts.admin.token, "invalid notes length create"},
	}

	// Create a request url with an "id" URL parameter
	requestUrl := "/api/vendors"

	for _, v := range createTests {

		// Make new request with vendor creation in body
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
			t.Errorf("response test (%v): got %v want %v. \nBody: %v\n", v.testName,
				status, v.expectedResponseStatus, rr.Body.String())
		}

		// Init body for response extraction
		var body db.Vendor
		var foundVendor db.Vendor
		// Grab ID from response body
		json.Unmarshal(rr.Body.Bytes(), &body)

		// Find the created transaction (to obtain full data with ID)
		testConnection.dbClient.Find(foundVendor, uint(body.ID))

		// Compare found details with those found in returned body
		checkVendorDetails(&body, &foundVendor, t, true)

		// If the task log was created successfully, check that it's deleted after test
		if v.expectedResponseStatus == http.StatusCreated {
			// Cleanup
			//
			// Delete the created vendor
			deleteMainResult := testConnection.dbClient.Delete(&db.Vendor{}, uint(body.ID))
			if deleteMainResult.Error != nil {
				t.Fatalf("Couldn't clean up created vendors: %v", deleteMainResult.Error)
			}
		}

	}
}

// Check the vendor details
func checkVendorDetails(actual *db.Vendor, expected *db.Vendor, t *testing.T, checkId bool) {
	// Only check ID if parameter checkId is true
	if checkId == true {
		// Verify that the actual id matches the created id
		if actual.ID != expected.ID {
			t.Errorf("found vendor has incorrect ID: expected %d, got %d", expected.ID, actual.ID)
		}
	}

	// Check details
	if actual.CompanyName != expected.CompanyName {
		t.Errorf("found vendor has incorrect company name: expected %s, got %s", expected.CompanyName, actual.CompanyName)
	}
	if actual.NPWP != expected.NPWP {
		t.Errorf("found vendor has incorrect NPWP: expected %s, got %s", expected.NPWP, actual.NPWP)
	}
	if actual.Email != expected.Email {
		t.Errorf("found vendor has incorrect email: expected %s, got %s", expected.Email, actual.Email)
	}
	if actual.Phone != expected.Phone {
		t.Errorf("found vendor has incorrect phone: expected %s, got %s", expected.Phone, actual.Phone)
	}
	if actual.NIB != expected.NIB {
		t.Errorf("found vendor has incorrect NIB: expected %s, got %s", expected.NIB, actual.NIB)
	}
	if actual.Street_Address_1 != expected.Street_Address_1 {
		t.Errorf("found vendor has incorrect street address 1: expected %s, got %s", expected.Street_Address_1, actual.Street_Address_1)
	}
	if actual.Street_Address_2 != expected.Street_Address_2 {
		t.Errorf("found vendor has incorrect street address 2: expected %s, got %s", expected.Street_Address_2, actual.Street_Address_2)
	}
	if actual.City != expected.City {
		t.Errorf("found vendor has incorrect city: expected %s, got %s", expected.City, actual.City)
	}
	if actual.Province != expected.Province {
		t.Errorf("found vendor has incorrect province: expected %s, got %s", expected.Province, actual.Province)
	}
	if actual.Postal_Code != expected.Postal_Code {
		t.Errorf("found vendor has incorrect postal code: expected %s, got %s", expected.Postal_Code, actual.Postal_Code)
	}
	if actual.Suburb != expected.Suburb {
		t.Errorf("found vendor has incorrect suburb: expected %s, got %s", expected.Suburb, actual.Suburb)
	}
}
