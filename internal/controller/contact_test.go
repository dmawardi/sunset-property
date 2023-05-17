package controller_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/helpers"
)

func TestContactController_FindAll(t *testing.T) {
	// Test setup
	// Build features
	var createContact1 = db.Contact{
		FirstName:    "John",
		LastName:     "Doe",
		Email:        "sig@waygo.com",
		Phone:        "123456789",
		Mobile:       "123456789",
		ContactType:  "Buyer",
		ContactNotes: "This is a note",
	}
	var createContact2 = db.Contact{
		FirstName:    "John",
		LastName:     "Doe",
		Email:        "swole@waygo.com",
		Phone:        "124433",
		Mobile:       "2312345",
		ContactType:  "Buyer",
		ContactNotes: "This is a note",
	}
	var createdContacts = []db.Contact{createContact1, createContact2}
	// Create contacts in database
	seedErr := testConnection.dbClient.Create(createdContacts)
	if seedErr.Error != nil {
		t.Errorf("Error seeding database: %v", seedErr.Error)
	}

	// Create a new request
	req, err := http.NewRequest("GET", "/api/contacts?limit=10&offset=0&order=", nil)
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
	var body []db.Contact
	json.Unmarshal(rr.Body.Bytes(), &body)

	// Check length of contacts array (should be two with seeded assets)
	if len(body) != len(createdContacts) {
		t.Errorf("Contacts array in findAll failed: expected %d, got %d", len(createdContacts), len(body))
	}

	// Iterate through contacts array received
	for _, actual := range body {
		// Iterate through created contacts to determine a match
		for _, created := range createdContacts {
			// If match found
			if actual.ID == created.ID {
				// Check the details of the feature
				checkContactDetails(&actual, &created, t, false)
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
		request := fmt.Sprintf("/api/contacts?limit=%v&offset=%v&order=%v", v.limit, v.offset, v.order)
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
	deleteResult := testConnection.dbClient.Delete(createdContacts)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded features for find all: %v", deleteResult.Error)
	}
}

func TestContactController_Find(t *testing.T) {
	// Test setup
	var createContact = db.Contact{
		FirstName:    "Calvin",
		LastName:     "Kajole",
		Email:        "what@saeli.com",
		Phone:        "123456789",
		Mobile:       "123456789",
		ContactType:  "Buyer",
		ContactNotes: "This is a note",
	}

	var createdContacts = []db.Contact{createContact}
	// Create contacts in database
	seedErr := testConnection.dbClient.Create(createdContacts)
	if seedErr.Error != nil {
		t.Errorf("Error seeding database: %v", seedErr.Error)
	}

	// Create a request with an "id" URL parameter
	requestUrl := fmt.Sprintf("/api/contacts/%v", createdContacts[0].ID)
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
	// Extract the response body
	var body db.Contact
	json.Unmarshal(rr.Body.Bytes(), &body)
	checkContactDetails(&body, &createdContacts[0], t, true)

	// Cleanup
	// Delete the created fixtures
	deleteResult := testConnection.dbClient.Delete(createdContacts)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded features for find all: %v", deleteResult.Error)
	}
}

func TestContactController_Delete(t *testing.T) {
	// Test setup
	var createContact = db.Contact{
		FirstName:    "Calvin",
		LastName:     "Kajole",
		Email:        "hereoin@sangalaki.com",
		Phone:        "123456789",
		Mobile:       "123456789",
		ContactType:  "Buyer",
		ContactNotes: "This is a note",
	}

	var createdContacts = []db.Contact{createContact}
	// Create contacts in database
	seedErr := testConnection.dbClient.Create(createdContacts)
	if seedErr.Error != nil {
		t.Errorf("Error seeding database: %v", seedErr.Error)
	}

	// Create a request with an "id" URL parameter
	requestUrl := fmt.Sprintf("/api/contacts/%v", createdContacts[0].ID)
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
		{testName: "Contacts basic user delete test", tokenToUse: testConnection.accounts.user.token, expectedResponseStatus: http.StatusForbidden},
		// Must be last
		// Tests of deletion success using admin token
		{testName: "Contacts admin delete test", tokenToUse: testConnection.accounts.admin.token, expectedResponseStatus: http.StatusOK},
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

func TestContactController_Update(t *testing.T) {
	// Test setup
	var createContact = db.Contact{
		FirstName:    "Brokie",
		LastName:     "Kajole",
		Email:        "creat@gmail.com",
		Phone:        "123456789",
		Mobile:       "123456789",
		ContactType:  "Buyer",
		ContactNotes: "This is a note",
	}

	var createdContacts = []db.Contact{createContact}
	// Create contacts in database
	seedErr := testConnection.dbClient.Create(createdContacts)
	if seedErr.Error != nil {
		t.Errorf("Error seeding database: %v", seedErr.Error)
	}

	// Build test array
	var updateTests = []struct {
		data                   db.Contact
		tokenToUse             string
		expectedResponseStatus int
		checkDetails           bool
		testName               string
	}{
		{db.Contact{FirstName: "Kanjing", LastName: "Blister", Email: "swag@gmail.com", Phone: "87987239487"}, testConnection.accounts.user.token, http.StatusForbidden, false, "Contacts basic user update test"},
		{db.Contact{FirstName: "Kanjing", LastName: "Blister", Email: "gila@gmail.com", Phone: "87987239487"}, testConnection.accounts.admin.token, http.StatusOK, true, "Contacts admin update test"},
		// Update should be disallowed due to being too short
		{db.Contact{FirstName: "a", LastName: "b", Email: "thedog@gmail.com"}, testConnection.accounts.admin.token, http.StatusBadRequest, false, "Contacts admin too short fail test"},
		// Update should be disallowed due to not being proper email
		{db.Contact{Email: "abdul"}, testConnection.accounts.admin.token, http.StatusBadRequest, false, "Contacts admin bad email fail test"},
		// Update should be disallowed due to not being proper phone number
		{db.Contact{Phone: "a978081234b"}, testConnection.accounts.admin.token, http.StatusBadRequest, false, "Contacts admin bad phone fail test"},
		// User should be forbidden before validating rather than Bad Request
		{db.Contact{FirstName: "b"}, testConnection.accounts.user.token, http.StatusForbidden, false, "Contacts basic user too short fail test"},
	}

	// Create a request url with an "id" URL parameter
	requestUrl := fmt.Sprintf("/api/contacts/%v", createdContacts[0].ID)

	// Iterate through update tests
	for _, v := range updateTests {
		// Make new request with contact update in body
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
		var body db.Contact
		json.Unmarshal(rr.Body.Bytes(), &body)
		// Check response expected vs received
		if status := rr.Code; status != v.expectedResponseStatus {
			t.Errorf("Contact update (%v): got %v want %v.", v.testName,
				status, v.expectedResponseStatus)
		}

		// If need to check details
		if v.checkDetails == true {
			// Get contact details from database
			var expected db.Contact
			findResult := testConnection.dbClient.Find(&expected, createdContacts[0].ID)
			if findResult.Error != nil {
				t.Errorf("Error finding updated contact: %v", findResult.Error)
			}

			// Check contact details using updated object
			checkContactDetails(&body, &expected, t, false)
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
		req, err := http.NewRequest("PUT", fmt.Sprint("/api/contacts/"+v.urlExtension), buildReqBody(&db.Contact{
			FirstName: "Scrappy Kid",
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

	// Delete the created fixtures
	deleteResult := testConnection.dbClient.Delete(createdContacts)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded fixtures: %v", deleteResult.Error)
	}
}

func TestContactController_Create(t *testing.T) {
	// Setup
	//
	var createTests = []struct {
		data                   db.Contact
		expectedResponseStatus int
		tokenToUse             string
	}{
		// Should fail due to user role status of basic
		{db.Contact{FirstName: "Kanjing", LastName: "Blister", Email: "swag@gmail.com", Phone: "87987239487", ContactType: "owner"}, http.StatusForbidden, testConnection.accounts.user.token},
		// Should pass as user is admin
		{db.Contact{FirstName: "Kanjing", LastName: "Blister", Email: "swag@gmail.com", Phone: "87987239487", ContactType: "owner"}, http.StatusCreated, testConnection.accounts.admin.token},
		// Create should be disallowed due to being too short
		{db.Contact{FirstName: "a", LastName: "b", Email: "swag@gmail.com", Phone: "87987239487"}, http.StatusBadRequest, testConnection.accounts.admin.token},
		// Should be a bad request due to invalid email
		{db.Contact{FirstName: "Magat", LastName: "Swagger", Email: "swag", Phone: "87987239487"}, http.StatusBadRequest, testConnection.accounts.admin.token},
		// Should be a bad request due to invalid phone
		{db.Contact{FirstName: "Membra", LastName: "Sercra", Email: "heylow@swag.com", Phone: "83adf7239487"}, http.StatusBadRequest, testConnection.accounts.admin.token},
	}

	// Create a request url with an "id" URL parameter
	requestUrl := "/api/contacts"

	for _, v := range createTests {
		// Make new request with contact creation in body
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
			t.Errorf("Contact create test (%v): got %v want %v.", v.data.FirstName,
				status, v.expectedResponseStatus)
		}

		// Init body for response extraction
		var body db.Contact
		var foundContact db.Contact
		// Grab ID from response body
		json.Unmarshal(rr.Body.Bytes(), &body)

		// Find the created contact (to obtain full data with ID)
		testConnection.dbClient.Find(foundContact, uint(body.ID))

		// Compare found contact details with those found in returned body
		checkContactDetails(&body, &foundContact, t, true)

		// Delete the created fixtures
		delResult := testConnection.dbClient.Delete(&db.Contact{}, uint(body.ID))
		if delResult.Error != nil {
			t.Fatalf("Issue encountered deleting seeded assets for contact create test (%v): %v", v.data.FirstName, delResult.Error)
		}
	}
}

// Updates the parameter contact struct with the updated values
func updateContactChangesOnly(createdContact *db.Contact, updatedFields map[string]string) error {
	// Iterate through map and change struct values
	for k, v := range updatedFields {
		// Update each struct field using map
		err := helpers.UpdateStructField(createdContact, k, v)
		if err != nil {
			return err
		}
	}
	return nil
}

// Check if contact details match expected
func checkContactDetails(actual *db.Contact, expected *db.Contact, t *testing.T, checkId bool) {

	// Only check ID if parameter checkId is true
	if checkId == true {
		// Verify that the actual contact id matches the created contact
		if actual.ID != expected.ID {
			t.Errorf("found feature has incorrect ID: expected %d, got %d", expected.ID, actual.ID)
		}
	}

	// Verify that the actual contact details match the expected contact details
	if actual.FirstName != expected.FirstName {
		t.Errorf("Contact has incorrect First Name: expected %s, got %s", expected.FirstName, actual.FirstName)
	}
	if actual.LastName != expected.LastName {
		t.Errorf("Contact has incorrect Last Name: expected %s, got %s", expected.LastName, actual.LastName)
	}
	if actual.ContactType != expected.ContactType {
		t.Errorf("Contact has incorrect Contact Type: expected %s, got %s", expected.ContactType, actual.ContactType)
	}
	if actual.Email != expected.Email {
		t.Errorf("Contact has incorrect Email: expected %s, got %s", expected.Email, actual.Email)
	}
	if actual.Phone != expected.Phone {
		t.Errorf("Contact has incorrect Phone: expected %s, got %s", expected.Phone, actual.Phone)
	}
	if actual.Mobile != expected.Mobile {
		t.Errorf("Contact has incorrect Mobile: expected %s, got %s", expected.Mobile, actual.Mobile)
	}
	if actual.ContactNotes != expected.ContactNotes {
		t.Errorf("Contact has incorrect Contact Notes: expected %s, got %s", expected.ContactNotes, actual.ContactNotes)
	}
}
