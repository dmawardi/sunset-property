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

func TestTransactionController_FindAll(t *testing.T) {
	// Test setup
	// Create property
	propertyToCreate := &db.Property{
		Property_Name:    "Test Property3",
		Postcode:         80361,
		Suburb:           "Test Suburb",
		City:             "Test City",
		Street_Address_1: "Test Street Address 1",
		Bedrooms:         3,
		Bathrooms:        2,
		Description:      "Test Description",
		Managed:          true,
	}
	createdProperties := []db.Property{*propertyToCreate}
	createResult := testConnection.dbClient.Create(createdProperties)
	if createResult.Error != nil {
		t.Fatal("Failed to create properties for transaction find all test: ", createResult.Error)
	}
	// Create task
	taskToCreate1 := &db.Task{
		TaskName: "Test Task",
		Type:     "Transaction",
		Notes:    "Yohoo",
	}
	taskToCreate2 := &db.Task{
		TaskName: "Test Task 2",
		Type:     "Transaction",
		Notes:    "Yohoo",
	}
	createdTasks := []db.Task{*taskToCreate1, *taskToCreate2}
	createResult = testConnection.dbClient.Create(createdTasks)
	if createResult.Error != nil {
		t.Fatal("Failed to create tasks for transaction find all test: ", createResult.Error)
	}
	// Create transactions
	transactionToCreate1 := &db.Transaction{
		TaskID:           createdTasks[0].ID,
		Type:             "Transaction",
		Agency:           "Own",
		AgencyName:       "Test Agency Name",
		IsLease:          true,
		Fee:              3.5,
		TransactionNotes: "This is a note",
		TenancyType:      "Monthly",
		Property:         db.Property{ID: createdProperties[0].ID},
	}
	transactionToCreate2 := &db.Transaction{
		TaskID:           createdTasks[1].ID,
		Type:             "Transaction",
		Agency:           "Own",
		AgencyName:       "Test Agency Name",
		IsLease:          true,
		Fee:              3.5,
		TransactionNotes: "This is a note",
		TenancyType:      "Monthly",
		Property:         *propertyToCreate,
	}
	createdTransactions := []db.Transaction{*transactionToCreate1, *transactionToCreate2}
	// Create task logs in db
	createResult = testConnection.dbClient.Create(createdTransactions)
	if createResult.Error != nil {
		t.Fatal("Failed to create transactions for transactions find all test: ", createResult.Error)
	}

	// Create a new request
	req, err := http.NewRequest("GET", "/api/transactions?limit=10&offset=0&order=", nil)
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
	var body []db.Transaction
	json.Unmarshal(rr.Body.Bytes(), &body)

	// Check length of array (should be two with seeded assets)
	if len(body) != len(createdTransactions) {
		t.Errorf("Array length check in findAll failed: expected %d, got %d", len(createdTransactions), len(body))
	}

	// Iterate through array received
	for _, actualTrans := range body {
		// Iterate through created transactions to determine a match
		for _, createdTrans := range createdTransactions {
			// If match found
			if actualTrans.ID == createdTrans.ID {
				// Check the details of the transaction match
				checkTransactionDetails(&actualTrans, &createdTrans, t, true)
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
		request := fmt.Sprintf("/api/transactions?limit=%v&offset=%v&order=%v", v.limit, v.offset, v.order)
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
	deleteResult := testConnection.dbClient.Delete(createdTransactions)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded transactions: %v", deleteResult.Error)
	}
	deleteResult = testConnection.dbClient.Delete(createdTasks)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded tasks: %v", deleteResult.Error)
	}
	deleteResult = testConnection.dbClient.Delete(createdProperties)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded props: %v", deleteResult.Error)
	}
}

func TestTransactionController_Find(t *testing.T) {
	// Test setup
	// Create property
	propertyToCreate := &db.Property{
		Property_Name:    "Test Property",
		Postcode:         80361,
		Suburb:           "Test Suburb",
		City:             "Test City",
		Street_Address_1: "Test Street Address 1",
		Bedrooms:         3,
		Bathrooms:        2,
		Description:      "Test Description",
		Managed:          true,
	}
	createdProperties := []db.Property{*propertyToCreate}
	createResult := testConnection.dbClient.Create(createdProperties)
	if createResult.Error != nil {
		t.Fatal("Failed to create properties for transaction find all test: ", createResult.Error)
	}
	// Create task
	taskToCreate1 := &db.Task{
		TaskName: "Test Task",
		Type:     "Transaction",
		Notes:    "Yohoo",
	}
	createdTasks := []db.Task{*taskToCreate1}
	createResult = testConnection.dbClient.Create(createdTasks)
	if createResult.Error != nil {
		t.Fatal("Failed to create tasks for transaction find all test: ", createResult.Error)
	}
	// Create transactions
	transactionToCreate1 := &db.Transaction{
		TaskID:           createdTasks[0].ID,
		Type:             "Transaction",
		Agency:           "Own",
		AgencyName:       "Test Agency Name",
		IsLease:          true,
		Fee:              3.5,
		TransactionNotes: "This is a note",
		TenancyType:      "Monthly",
		Property:         db.Property{ID: createdProperties[0].ID},
	}
	createdTransactions := []db.Transaction{*transactionToCreate1}
	// Create task logs in db
	createResult = testConnection.dbClient.Create(createdTransactions)
	if createResult.Error != nil {
		t.Fatal("Failed to create transactions for transactions find all test: ", createResult.Error)
	}

	// Create a request with an "id" URL parameter
	requestUrl := fmt.Sprintf("/api/transactions/%v", createdTransactions[0].ID)
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
	var body db.Transaction
	json.Unmarshal(rr.Body.Bytes(), &body)
	checkTransactionDetails(&body, &createdTransactions[0], t, true)

	// Cleanup
	// Delete the created fixtures
	deleteResult := testConnection.dbClient.Delete(createdTransactions)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded transactions: %v", deleteResult.Error)
	}
	deleteResult = testConnection.dbClient.Delete(createdTasks)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded tasks: %v", deleteResult.Error)
	}
	deleteResult = testConnection.dbClient.Delete(createdProperties)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded props: %v", deleteResult.Error)
	}
}

func TestTransactionController_Delete(t *testing.T) {
	// Test setup
	// Create property
	propertyToCreate := &db.Property{
		Property_Name:    "Test Property2",
		Postcode:         80361,
		Suburb:           "Test Suburb",
		City:             "Test City",
		Street_Address_1: "Test Street Address 1",
		Bedrooms:         3,
		Bathrooms:        2,
		Description:      "Test Description",
		Managed:          true,
	}
	createdProperties := []db.Property{*propertyToCreate}
	createResult := testConnection.dbClient.Create(createdProperties)
	if createResult.Error != nil {
		t.Fatal("Failed to create properties for transaction find all test: ", createResult.Error)
	}
	// Create task
	taskToCreate1 := &db.Task{
		TaskName: "Test Task",
		Type:     "Transaction",
		Notes:    "Yohoo",
	}
	createdTasks := []db.Task{*taskToCreate1}
	createResult = testConnection.dbClient.Create(createdTasks)
	if createResult.Error != nil {
		t.Fatal("Failed to create tasks for transaction find all test: ", createResult.Error)
	}
	// Create transactions
	transactionToCreate1 := &db.Transaction{
		TaskID:           createdTasks[0].ID,
		Type:             "Transaction",
		Agency:           "Own",
		AgencyName:       "Test Agency Name",
		IsLease:          true,
		Fee:              3.5,
		TransactionNotes: "This is a note",
		TenancyType:      "Monthly",
		Property:         db.Property{ID: createdProperties[0].ID},
	}
	createdTransactions := []db.Transaction{*transactionToCreate1}
	// Create task logs in db
	createResult = testConnection.dbClient.Create(createdTransactions)
	if createResult.Error != nil {
		t.Fatal("Failed to create transactions for transactions find all test: ", createResult.Error)
	}

	// Create a request with an "id" URL parameter
	requestUrl := fmt.Sprintf("/api/transactions/%v", testConnection.propertyLogs.created[0].ID)
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
		{testName: "Transaction basic user delete test", tokenToUse: testConnection.accounts.user.token, expectedResponseStatus: http.StatusForbidden},
		// Must be last
		// Tests of deletion success using admin privileges
		{testName: "Transaction admin delete test", tokenToUse: testConnection.accounts.admin.token, expectedResponseStatus: http.StatusOK},
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
	deleteResult := testConnection.dbClient.Delete(createdTransactions)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded transactions: %v", deleteResult.Error)
	}
	deleteResult = testConnection.dbClient.Delete(createdTasks)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded tasks: %v", deleteResult.Error)
	}
	deleteResult = testConnection.dbClient.Delete(createdProperties)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded props: %v", deleteResult.Error)
	}
}

func TestTransactionController_Update(t *testing.T) {
	// Test setup
	// Create property
	propertyToCreate := &db.Property{
		Property_Name:    "Test Property4",
		Postcode:         80361,
		Suburb:           "Test Suburb",
		City:             "Test City",
		Street_Address_1: "Test Street Address 1",
		Bedrooms:         3,
		Bathrooms:        2,
		Description:      "Test Description",
		Managed:          true,
	}
	createdProperties := []db.Property{*propertyToCreate}
	createResult := testConnection.dbClient.Create(createdProperties)
	if createResult.Error != nil {
		t.Fatal("Failed to create properties for transaction find all test: ", createResult.Error)
	}
	// Create task
	taskToCreate1 := &db.Task{
		TaskName: "Test Task",
		Type:     "Transaction",
		Notes:    "Yohoo",
	}
	createdTasks := []db.Task{*taskToCreate1}
	createResult = testConnection.dbClient.Create(createdTasks)
	if createResult.Error != nil {
		t.Fatal("Failed to create tasks for transaction find all test: ", createResult.Error)
	}
	// Create transactions
	transactionToCreate1 := &db.Transaction{
		TaskID:           createdTasks[0].ID,
		Type:             "Transaction",
		Agency:           "Own",
		AgencyName:       "Test Agency Name",
		IsLease:          true,
		Fee:              3.5,
		TransactionNotes: "This is a note",
		TenancyType:      "Monthly",
		Property:         db.Property{ID: createdProperties[0].ID},
	}
	createdTransactions := []db.Transaction{*transactionToCreate1}
	createResult = testConnection.dbClient.Create(createdTransactions)
	if createResult.Error != nil {
		t.Fatal("Failed to create transactions for transactions find all test: ", createResult.Error)
	}

	// Build test array
	var updateTests = []struct {
		data                   models.UpdateTransaction
		tokenToUse             string
		expectedResponseStatus int
		checkDetails           bool
	}{
		// Test of update failure: basic user
		{models.UpdateTransaction{
			Fee:              4.5,
			TransactionValue: 35000000,
		}, testConnection.accounts.user.token, http.StatusForbidden, false},
		// Update should be allowed: admin
		{models.UpdateTransaction{
			Fee:              4.5,
			TransactionValue: 35000000,
		}, testConnection.accounts.admin.token, http.StatusOK, true},
		// Update should be disallowed due to being invalid value for agency
		{models.UpdateTransaction{
			Agency:           "Insane",
			Fee:              4.5,
			TransactionValue: 35000000,
		}, testConnection.accounts.admin.token, http.StatusBadRequest, false},
		// Update should be disallowed due to being invalid value for type
		{models.UpdateTransaction{
			Type:             "Insane",
			Fee:              4.5,
			TransactionValue: 35000000,
		}, testConnection.accounts.admin.token, http.StatusBadRequest, false},
		// User should be forbidden before validating rather than Bad Request
		{models.UpdateTransaction{
			Type: "Insane",
		}, testConnection.accounts.user.token, http.StatusForbidden, false},
	}

	// Create a request url with an "id" URL parameter
	requestUrl := fmt.Sprintf("/api/transactions/%v", createdTransactions[0].ID)

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
		var body db.Transaction
		json.Unmarshal(rr.Body.Bytes(), &body)
		// Check response expected vs received
		if status := rr.Code; status != v.expectedResponseStatus {
			t.Errorf("Update test: got %v want %v. \nBody: %v", status, v.expectedResponseStatus, body)

		}

		// If need to check details
		if v.checkDetails == true {
			// Get task details from database
			var expected db.Transaction
			findResult := testConnection.dbClient.Find(&expected, createdTransactions[0].ID)
			if findResult.Error != nil {
				t.Errorf("Error finding updated transaction: %v", findResult.Error)
			}

			// Check task log details using updated object
			checkTransactionDetails(&body, &expected, t, false)
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
		req, err := http.NewRequest("PUT", fmt.Sprint("/api/transactions/"+v.urlExtension), buildReqBody(&db.Transaction{
			Agency: "Own",
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
	deleteResult := testConnection.dbClient.Delete(createdTransactions)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded transactions: %v", deleteResult.Error)
	}
	deleteResult = testConnection.dbClient.Delete(createdTasks)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded tasks: %v", deleteResult.Error)
	}
	deleteResult = testConnection.dbClient.Delete(createdProperties)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded props: %v", deleteResult.Error)
	}
}

func TestTransactionController_Create(t *testing.T) {
	// Setup
	//
	// Create property
	propertyToCreate := &db.Property{
		Property_Name:    "Test Property5",
		Postcode:         80361,
		Suburb:           "Test Suburb",
		City:             "Test City",
		Street_Address_1: "Test Street Address 1",
		Bedrooms:         3,
		Bathrooms:        2,
		Description:      "Test Description",
		Managed:          true,
	}
	createdProperties := []db.Property{*propertyToCreate}
	createResult := testConnection.dbClient.Create(createdProperties)
	if createResult.Error != nil {
		t.Fatal("Failed to create properties for transaction find all test: ", createResult.Error)
	}

	var createTests = []struct {
		data                   models.CreateTransaction
		expectedResponseStatus int
		tokenToUse             string
	}{
		// Should fail due to user role status of basic
		{models.CreateTransaction{
			Type:             "Sale",
			Agency:           "Own",
			AgencyName:       "Test Agency Name",
			IsLease:          false,
			Fee:              3.5,
			TransactionNotes: "This is a note",
			TenancyType:      "Monthly",
			Property:         db.Property{ID: createdProperties[0].ID},
		}, http.StatusForbidden, testConnection.accounts.user.token},
		// Should pass as user is admin
		{models.CreateTransaction{
			Type:             "Sale",
			Agency:           "Own",
			AgencyName:       "Test Agency Name",
			IsLease:          false,
			Fee:              3.5,
			TransactionNotes: "This is a note",
			TenancyType:      "Monthly",
			Property:         db.Property{ID: createdProperties[0].ID},
		}, http.StatusCreated, testConnection.accounts.admin.token},
		// Create should be disallowed due to invalid agency value
		{models.CreateTransaction{
			Type:             "Lease",
			Agency:           "Sakra",
			AgencyName:       "Test Agency Name",
			IsLease:          false,
			Fee:              3.5,
			TransactionNotes: "This is a note",
			TenancyType:      "Monthly",
			Property:         db.Property{ID: createdProperties[0].ID},
		}, http.StatusBadRequest, testConnection.accounts.admin.token},
		// Create should be disallowed due to invalid type value
		{models.CreateTransaction{
			Type:             "Crazy",
			Agency:           "Own",
			AgencyName:       "Test Agency Name",
			IsLease:          false,
			Fee:              3.5,
			TransactionNotes: "This is a note",
			TenancyType:      "Monthly",
			Property:         db.Property{ID: createdProperties[0].ID},
		}, http.StatusBadRequest, testConnection.accounts.admin.token},
		// Create should be disallowed due to invalid Tenancy type value
		{models.CreateTransaction{
			Type:             "Lease",
			Agency:           "Own",
			AgencyName:       "Test Agency Name",
			IsLease:          false,
			Fee:              3.5,
			TransactionNotes: "This is a note",
			TenancyType:      "Iglesias",
			Property:         db.Property{ID: createdProperties[0].ID},
		}, http.StatusBadRequest, testConnection.accounts.admin.token},
		// Create should be disallowed due to notes being too long
		{models.CreateTransaction{
			Type:             "Lease",
			Agency:           "Own",
			AgencyName:       "Test Agency Name",
			IsLease:          false,
			Fee:              3.5,
			TransactionNotes: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Suspendisse vulputate, nunc sit amet efficitur bibendum, sapien odio auctor nisi, a interdum magna nisl ac purus. Fusce condimentum malesuada mi at eleifend. Sed laoreet varius risus, id mattis libero tristique nec. Sed eget malesuada magna. Morbi feugiat sapien euismod neque commodo suscipit. Vivamus vehicula euismod dui, id imperdiet elit lacinia non. Integer hendrerit, enim ac gravida malesuada, dolor leo dictum purus, nec bibendum velit est vel nulla. Nulla sagittis nulla non elit imperdiet convallis. Sed bibendum sollicitudin nunc, vel facilisis nulla convallis a. Nunc id ex feugiat, finibus magna sit amet, ultricies lacus.",
			TenancyType:      "Iglesias",
			Property:         db.Property{ID: createdProperties[0].ID},
		}, http.StatusBadRequest, testConnection.accounts.admin.token},
	}

	// Create a request url with an "id" URL parameter
	requestUrl := "/api/transactions"

	for _, v := range createTests {
		// Each create test setup
		// Create task (for each test)
		taskToCreate1 := &db.Task{
			TaskName: "Test Task",
			Type:     "Transaction",
			Notes:    "Yohoo",
		}
		createdTasks := []db.Task{*taskToCreate1}
		createResult = testConnection.dbClient.Create(createdTasks)
		if createResult.Error != nil {
			t.Fatal("Failed to create tasks for transaction create test: ", createResult.Error)
		}
		// Update v.data to include newly created task ID
		v.data.Task.ID = createdTasks[0].ID

		// Make new request with transaction creation in body
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
			t.Errorf("response test (%v): got %v want %v. \nBody: %v\n", v.data.Type,
				status, v.expectedResponseStatus, rr.Body.String())
		}

		// Init body for response extraction
		var body db.Transaction
		var foundTransaction db.Transaction
		// Grab ID from response body
		json.Unmarshal(rr.Body.Bytes(), &body)

		// Find the created transaction (to obtain full data with ID)
		testConnection.dbClient.Find(foundTransaction, uint(body.ID))

		// Compare found details with those found in returned body
		checkTransactionDetails(&body, &foundTransaction, t, true)

		// If the task log was created successfully, check that it's deleted after test
		if v.expectedResponseStatus == http.StatusCreated {
			// Cleanup
			//
			// Delete the created task logs
			deleteResult := testConnection.dbClient.Delete(&db.TaskLog{}, uint(body.ID))
			if deleteResult.Error != nil {
				t.Fatalf("Couldn't clean up created transactions: %v", deleteResult.Error)
			}
			// Delete the created fixtures
			deleteResult = testConnection.dbClient.Delete(createdTasks)
			if deleteResult.Error != nil {
				t.Fatalf("Couldn't clean up seeded tasks: %v", deleteResult.Error)
			}
		}
	}

	// cleanup

	deleteResult := testConnection.dbClient.Delete(createdProperties)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded props: %v", deleteResult.Error)
	}
}

// Check the transaction details
func checkTransactionDetails(actual *db.Transaction, expected *db.Transaction, t *testing.T, checkId bool) {

	// Only check ID if parameter checkId is true
	if checkId == true {
		// Verify that the actual transaction id matches the created transaction id
		if actual.ID != expected.ID {
			t.Errorf("found transaction has incorrect ID: expected %d, got %d", expected.ID, actual.ID)
		}
	}

	// Verify that the actual details matches the expected details
	if actual.Type != expected.Type {
		t.Errorf("found transaction has incorrect type: expected %s, got %s", expected.Type, actual.Type)
	}
	if actual.IsLease != expected.IsLease {
		t.Errorf("found transaction has incorrect is lease: expected %t, got %t", expected.IsLease, actual.IsLease)
	}
	if actual.TenancyType != expected.TenancyType {
		t.Errorf("found transaction has incorrect tenancy type: expected %s, got %s", expected.TenancyType, actual.TenancyType)
	}
	if actual.Fee != expected.Fee {
		t.Errorf("found transaction has incorrect fee: expected %f, got %f", expected.Fee, actual.Fee)
	}
	if actual.TransactionNotes != expected.TransactionNotes {
		t.Errorf("found transaction has incorrect transaction notes: expected %s, got %s", expected.TransactionNotes, actual.TransactionNotes)
	}
	// Agency
	if actual.Agency != expected.Agency {
		t.Errorf("found transaction has incorrect agency: expected %s, got %s", expected.Agency, actual.Agency)
	}
	if actual.AgencyName != expected.AgencyName {
		t.Errorf("found transaction has incorrect agency name: expected %s, got %s", expected.AgencyName, actual.AgencyName)
	}
	// Transaction completion
	if actual.TransactionCompletion != expected.TransactionCompletion {
		t.Errorf("found transaction has incorrect transaction completion: expected %s, got %s", expected.TransactionCompletion, actual.TransactionCompletion)
	}
	if actual.TransactionValue != expected.TransactionValue {
		t.Errorf("found transaction has incorrect transaction value: expected %f, got %f", expected.TransactionValue, actual.TransactionValue)
	}
}

// func buildTransactionFixtures(propertiesToCreate []db.Property, tasksToCreate []db.Task, transactionsToCreate []db.Transaction, t *testing.T) ([]db.Property, []db.Task, []db.Transaction) {

// 	createResult := testConnection.dbClient.Create(propertiesToCreate)
// 	if createResult.Error != nil {
// 		t.Fatal("Failed to create properties for transaction find all test: ", createResult.Error)
// 	}
// 	// Create tasks
// 	createResult = testConnection.dbClient.Create(tasksToCreate)
// 	if createResult.Error != nil {
// 		t.Fatal("Failed to create tasks for transaction find all test: ", createResult.Error)
// 	}

// 	// Create transactions
// 	createResult = testConnection.dbClient.Create(transactionsToCreate)
// 	if createResult.Error != nil {
// 		t.Fatal("Failed to create transactions for transactions find all test: ", createResult.Error)
// 	}

// 	return transactionsToCreate
// }
