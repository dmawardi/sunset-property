package controller_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dmawardi/Go-Template/internal/db"
	"github.com/dmawardi/Go-Template/internal/models"
)

func TestPropertyAttachmentController_Upload(t *testing.T) {
	// Build test property
	propToCreate := &models.CreateProperty{
		Postcode:         14024,
		Property_Name:    "Renault No.2",
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

	// Create property for test
	createdProp, err := testConnection.properties.serv.Create(propToCreate)
	if err != nil {
		t.Fatalf("failed to create test property for find by id user service test: %v", err)
	}

	// Build upload data for request
	// Sample image file data
	var imageFileData = []byte{
		0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46, 0x00, 0x01, 0x01, 0x00, 0x00, 0x48, // ...
		// Add more byte values here
	}

	// Sample text file data
	var textFileData = []byte("Sample text file content")

	// Sample PDF file data
	var pdfFileData = []byte{
		0x25, 0x50, 0x44, 0x46, 0x2D, 0x31, 0x2E, 0x34, 0x0A, 0x25, 0xD0, 0xD4, 0xC5, 0xD8, 0x0A, 0x34, // ...
		// Add more byte values here
	}

	var uploadTests = []struct {
		testName               string
		data                   []byte
		tokenToUse             string
		fileName               string
		expectedResponseStatus int
	}{
		{
			testName:               "(user) Upload image file",
			data:                   imageFileData,
			tokenToUse:             testConnection.accounts.user.token,
			fileName:               "image.jpg",
			expectedResponseStatus: http.StatusForbidden,
		},
		{
			testName:               "(admin) Upload image file",
			data:                   imageFileData,
			tokenToUse:             testConnection.accounts.admin.token,
			fileName:               "image.jpg",
			expectedResponseStatus: http.StatusCreated,
		},
		{
			testName:               "Upload text file",
			data:                   textFileData,
			tokenToUse:             testConnection.accounts.admin.token,
			fileName:               "text.txt",
			expectedResponseStatus: http.StatusCreated,
		},
		{
			testName:               "Upload pdf file",
			data:                   pdfFileData,
			tokenToUse:             testConnection.accounts.admin.token,
			fileName:               "document.pdf",
			expectedResponseStatus: http.StatusCreated,
		},
	}

	// Iterate through tests
	for _, v := range uploadTests {
		// Create a buffer to hold the file data
		fileBuf := bytes.NewBuffer(v.data)

		// Create a new multipart writer
		fileBody := &bytes.Buffer{}
		writer := multipart.NewWriter(fileBody)

		// Create a new form file field with the file data
		fileField, err := writer.CreateFormFile("file", "image.jpg")
		if err != nil {
			t.Fatalf("Failed to create form file field: %v", err)
		}

		// Copy the file data to the form file field
		_, err = io.Copy(fileField, fileBuf)
		if err != nil {
			t.Fatalf("Failed to copy file data: %v", err)
		}

		// Close the multipart writer to finalize the form data
		writer.Close()

		// Create a request url
		requestUrl := "/api/property-attach/" + fmt.Sprint(createdProp.ID)

		// Make new request with work type creation in body
		req, err := http.NewRequest("POST", requestUrl, fileBody)
		if err != nil {
			t.Fatal(err)
		}
		// Set content type
		req.Header.Set("Content-Type", writer.FormDataContentType())

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

		// If expected response is successful, delete the created attachment
		if v.expectedResponseStatus == http.StatusCreated {
			// Delete created property attachment
			err = deleteCreatedPropertyAttachment(&testConnection)
			if err != nil {
				t.Fatalf("failed to delete created property attachment: %v", err)
			}
		}

	}

	// Delete property
	deleteResult := testConnection.dbClient.Delete(createdProp)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded property: %v", deleteResult.Error)
	}
}

func TestPropertyAttachmentController_Download(t *testing.T) {
	// Build test property
	propToCreate := &models.CreateProperty{
		Postcode:         14024,
		Property_Name:    "Gazeria No.2",
		Suburb:           "Kelapa Gading",
		City:             "Jakarta Utara",
		Street_Address_1: "Jl. Kintamani Raya no. 2",
		Street_Address_2: "Bukit Gading Villa",
		Bedrooms:         5,
		Bathrooms:        6,
		Land_Area:        400,
		Land_Metric:      "sqm",
		Description:      "A family home",
		Notes:            "The Bing slayer",
	}

	// Create property for test
	createdProp, err := testConnection.properties.serv.Create(propToCreate)
	if err != nil {
		t.Fatalf("failed to create test property for find by id user service test: %v", err)
	}

	attachmentToCreate1 := &db.PropertyAttachment{
		Label:     "Test Attachment 1",
		FileName:  "test1.jpg",
		ObjectKey: "properties/1/attachments/apricot.jpg",
		FileSize:  100,
		Property:  *createdProp,
		ETag:      "7d219e22bacfe3a56f5db68a58750361",
		FileType:  "jpg",
	}
	createdAttachments := []db.PropertyAttachment{*attachmentToCreate1}
	createResult := testConnection.dbClient.Create(createdAttachments)
	if createResult.Error != nil {
		t.Fatal("Failed to create attachments for test: ", createResult.Error)
	}

	var downloadTests = []struct {
		testName               string
		tokenToUse             string
		fileName               string
		expectedResponseStatus int
	}{
		{
			testName:               "(user) Download file",
			tokenToUse:             testConnection.accounts.user.token,
			expectedResponseStatus: http.StatusForbidden,
		},
		{
			testName:               "(admin) Download file",
			tokenToUse:             testConnection.accounts.admin.token,
			expectedResponseStatus: http.StatusOK,
		},
	}

	// Iterate through tests
	for _, v := range downloadTests {

		// Create a request url
		requestUrl := "/api/property-attach/" + fmt.Sprint(createdAttachments[0].ID)

		// Make new request with work type creation in body
		req, err := http.NewRequest("GET", requestUrl, nil)
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

	}

	// Delete created property attachment
	deleteResult := testConnection.dbClient.Delete(&db.PropertyAttachment{}, createdAttachments[0].ID)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded property attachment: %v", deleteResult.Error)
	}
	// Delete created property
	deleteResult = testConnection.dbClient.Delete(createdProp)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded property: %v", deleteResult.Error)
	}
}

func TestPropertyAttachmentController_FindAll(t *testing.T) {
	// Test setup
	// Create a property
	propToCreate := &db.Property{
		Postcode:         14024,
		Property_Name:    "Quetiau No.2",
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
	createdProps := []db.Property{*propToCreate}
	createResult := testConnection.dbClient.Create(createdProps)
	if createResult.Error != nil {
		t.Fatal("Failed to create properties for test: ", createResult.Error)
	}
	// Create a property attachment
	attachToCreate1 := &db.PropertyAttachment{
		Label:     "Test Attachment 1",
		FileName:  "apricot.jpg",
		FileSize:  100,
		FileType:  "jpg",
		ETag:      "7d219e22bacfe3a56f5db68a58750361",
		ObjectKey: "properties/1/attachments/apricot.jpg",
		Property:  db.Property{ID: createdProps[0].ID},
	}
	attachToCreate2 := &db.PropertyAttachment{
		Label:     "Test Attachment 2",
		FileName:  "apple.jpg",
		FileSize:  100,
		FileType:  "jpg",
		ETag:      "7d219e22bacfe3a56f5db68a58750361",
		ObjectKey: "properties/1/attachments/apple.jpg",
		Property:  db.Property{ID: createdProps[0].ID},
	}

	createdAttachments := []db.PropertyAttachment{*attachToCreate1, *attachToCreate2}
	createResult = testConnection.dbClient.Create(createdAttachments)
	if createResult.Error != nil {
		t.Fatal("Failed to create attachments for test: ", createResult.Error)
	}

	// Create a new request
	req, err := http.NewRequest("GET", "/api/property-attachments?limit=10&offset=0&order=", nil)
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
	var body []db.PropertyAttachment
	json.Unmarshal(rr.Body.Bytes(), &body)

	// Check length of array (should be two with seeded assets)
	if len(body) != len(createdAttachments) {
		t.Errorf("Array length check in findAll failed: expected %d, got %d", len(createdAttachments), len(body))
	}

	// Iterate through array received
	for _, actualAttachment := range body {
		// Iterate through prior created items to determine a match
		for _, createdAttachment := range createdAttachments {
			// If match found
			if actualAttachment.ID == createdAttachment.ID {
				// Check the details of the transaction match
				checkPropertyAttachmentDetails(&actualAttachment, &createdAttachment, t, false)
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
		request := fmt.Sprintf("/api/property-attachments?limit=%v&offset=%v&order=%v", v.limit, v.offset, v.order)
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
	deleteResult := testConnection.dbClient.Delete(createdAttachments)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded attachments: %v", deleteResult.Error)
	}
	deleteResult = testConnection.dbClient.Delete(createdProps)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded properties: %v", deleteResult.Error)
	}
}

func TestPropertyAttachmentController_Find(t *testing.T) {
	// Test setup
	// Create a property
	propToCreate := &db.Property{
		Postcode:         14024,
		Property_Name:    "Sugauoi No.2",
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
	createdProps := []db.Property{*propToCreate}
	createResult := testConnection.dbClient.Create(createdProps)
	if createResult.Error != nil {
		t.Fatal("Failed to create properties for test: ", createResult.Error)
	}
	// Create a property attachment
	attachToCreate1 := &db.PropertyAttachment{
		Label:     "Test Attachment 1",
		FileName:  "apricot.jpg",
		FileSize:  100,
		FileType:  "jpg",
		ETag:      "7d219e22bacfe3a56f5db68a58750361",
		ObjectKey: "properties/1/attachments/apricot.jpg",
		Property:  db.Property{ID: createdProps[0].ID},
	}

	createdAttachments := []db.PropertyAttachment{*attachToCreate1}
	createResult = testConnection.dbClient.Create(createdAttachments)
	if createResult.Error != nil {
		t.Fatal("Failed to create attachments for test: ", createResult.Error)
	}

	// Create a request with an "id" URL parameter
	requestUrl := fmt.Sprintf("/api/property-attachments/%v", createdAttachments[0].ID)
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		t.Fatal(err)
	}

	// Add auth token to header
	req.Header.Set("Authorization", fmt.Sprintf("bearer %v", testConnection.accounts.admin.token))
	// Create a response recorder
	rr := httptest.NewRecorder()

	// // Serve request using recorder and created request
	testConnection.router.ServeHTTP(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v. Error: %v",
			status, http.StatusOK, rr.Body.String())
	}

	// Extract the response body
	var body db.PropertyAttachment
	json.Unmarshal(rr.Body.Bytes(), &body)
	checkPropertyAttachmentDetails(&body, &createdAttachments[0], t, true)

	// Cleanup
	// Delete the created fixtures
	deleteResult := testConnection.dbClient.Delete(createdAttachments)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded attachments: %v", deleteResult.Error)
	}
	deleteResult = testConnection.dbClient.Delete(createdProps)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded properties: %v", deleteResult.Error)
	}
}

func TestPropertyAttachmentController_Delete(t *testing.T) {
	// Test setup
	// Create a property
	propToCreate := &db.Property{
		Postcode:         14024,
		Property_Name:    "Sunabu No.2",
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
	createdProps := []db.Property{*propToCreate}
	createResult := testConnection.dbClient.Create(createdProps)
	if createResult.Error != nil {
		t.Fatal("Failed to create properties for test: ", createResult.Error)
	}
	// Create a property attachment
	attachToCreate1 := &db.PropertyAttachment{
		Label:     "Test Attachment 1",
		FileName:  "apricot.jpg",
		FileSize:  100,
		FileType:  "jpg",
		ETag:      "7d219e22bacfe3a56f5db68a58750361",
		ObjectKey: "properties/1/attachments/apricot.jpg",
		Property:  db.Property{ID: createdProps[0].ID},
	}

	createdAttachments := []db.PropertyAttachment{*attachToCreate1}
	createResult = testConnection.dbClient.Create(createdAttachments)
	if createResult.Error != nil {
		t.Fatal("Failed to create attachments for test: ", createResult.Error)
	}

	// Create a request with an "id" URL parameter
	requestUrl := fmt.Sprintf("/api/property-attachments/%v", createdAttachments[0].ID)
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
			t.Errorf("Prop attachment deletion test (%v): got %v want %v.", test.testName,
				status, test.expectedResponseStatus)
		}
	}
	// Cleanup
	// Delete the created fixtures
	deleteResult := testConnection.dbClient.Delete(createdAttachments)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded attachments: %v", deleteResult.Error)
	}
	deleteResult = testConnection.dbClient.Delete(createdProps)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded properties: %v", deleteResult.Error)
	}
}

func TestPropertyAttachmentController_Update(t *testing.T) {
	// Test setup
	// Create a property
	propToCreate := &db.Property{
		Postcode:         14024,
		Property_Name:    "Furaibo No.2",
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
	createdProps := []db.Property{*propToCreate}
	createResult := testConnection.dbClient.Create(createdProps)
	if createResult.Error != nil {
		t.Fatal("Failed to create properties for test: ", createResult.Error)
	}
	// Create a property attachment
	attachToCreate1 := &db.PropertyAttachment{
		Label:     "Test Attachment 1",
		FileName:  "apricot.jpg",
		FileSize:  100,
		FileType:  "jpg",
		ETag:      "7d219e22bacfe3a56f5db68a58750361",
		ObjectKey: "properties/1/attachments/apricot.jpg",
		Property:  db.Property{ID: createdProps[0].ID},
	}

	createdAttachments := []db.PropertyAttachment{*attachToCreate1}
	createResult = testConnection.dbClient.Create(createdAttachments)
	if createResult.Error != nil {
		t.Fatal("Failed to create attachments for test: ", createResult.Error)
	}

	// Build test array
	var updateTests = []struct {
		data                   models.UpdatePropertyAttachment
		tokenToUse             string
		expectedResponseStatus int
		checkDetails           bool
		testName               string
	}{
		// Test of update failure: basic user
		{models.UpdatePropertyAttachment{
			Label: "Swahili",
		}, testConnection.accounts.user.token, http.StatusForbidden, false, "Basic user update test"},
		// Update should be allowed: admin
		{models.UpdatePropertyAttachment{
			Label:    "Swahili",
			ETag:     "7d219e22bacfe3a56f5db68a58750361",
			FileName: "apricot43.jpg",
		}, testConnection.accounts.admin.token, http.StatusOK, true, "Admin update test"},
		// Update should be disallowed due to being under valid length
		{models.UpdatePropertyAttachment{
			Label:     "ili",
			ObjectKey: "ape",
			ETag:      "df834",
			FileType:  "d",
		}, testConnection.accounts.admin.token, http.StatusBadRequest, false, "Invalid type update test"},
		// User should be forbidden before validating rather than Bad Request
		{models.UpdatePropertyAttachment{
			Label:     "ili",
			ObjectKey: "ape",
			ETag:      "df834",
			FileType:  "d",
		}, testConnection.accounts.user.token, http.StatusForbidden, false, "Invalid type update when basic user test"},
	}

	// Create a request url with an "id" URL parameter
	requestUrl := fmt.Sprintf("/api/property-attachments/%v", createdAttachments[0].ID)
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
		var body db.PropertyAttachment
		json.Unmarshal(rr.Body.Bytes(), &body)

		// Check response expected vs received
		if status := rr.Code; status != v.expectedResponseStatus {
			t.Errorf("Update test (%v): got %v want %v. \nBody: %v", v.testName, status, v.expectedResponseStatus, rr.Body.String())
		}

		// If need to check details
		if v.checkDetails == true {
			// Get property attachment details from database
			var expected db.PropertyAttachment
			findResult := testConnection.dbClient.Find(&expected, createdAttachments[0].ID)
			if findResult.Error != nil {
				t.Errorf("Error finding updated work type: %v", findResult.Error)
			}

			// Check task log details using updated object
			checkPropertyAttachmentDetails(&body, &expected, t, false)
		}
	}

	// Check for failure if incorrect ID parameter detected
	//
	var failUpdateTests = []struct {
		urlExtension           string
		expectedResponseStatus int
		testName               string
	}{
		// alpha character instead
		{urlExtension: "x", expectedResponseStatus: http.StatusForbidden, testName: "Alpha character test"},
		// Index out of bounds
		{urlExtension: "99", expectedResponseStatus: http.StatusBadRequest, testName: "Index out of bounds test"},
	}
	for _, v := range failUpdateTests {
		// Make new request with task log update in body
		req, err := http.NewRequest("PUT", fmt.Sprint("/api/property-attachments/"+v.urlExtension), buildReqBody(&db.PropertyAttachment{
			Label: "Swahimi",
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
			t.Errorf("Fail update test (%v): got %v want %v. Body:\n%v\n", v.testName,
				status, v.expectedResponseStatus, rr.Body.String())
		}
	}

	// Cleanup
	// Delete the created fixtures
	deleteResult := testConnection.dbClient.Delete(createdAttachments)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded attachments: %v", deleteResult.Error)
	}
	deleteResult = testConnection.dbClient.Delete(createdProps)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded properties: %v", deleteResult.Error)
	}
}

func TestPropertyAttachmentController_Create(t *testing.T) {
	// Test setup
	// Create a property
	propToCreate := &db.Property{
		Postcode:         14024,
		Property_Name:    "Fugayzi No.2",
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
	createdProps := []db.Property{*propToCreate}
	createResult := testConnection.dbClient.Create(createdProps)
	if createResult.Error != nil {
		t.Fatal("Failed to create properties for test: ", createResult.Error)
	}

	var createTests = []struct {
		data                   models.CreatePropertyAttachment
		expectedResponseStatus int
		tokenToUse             string
		testName               string
	}{
		// Should fail due to user role status of basic
		{models.CreatePropertyAttachment{
			Label:     "Test Attachment 1",
			FileName:  "apricot.jpg",
			FileSize:  100,
			FileType:  "jpg",
			ETag:      "7d219e22bacfe3a56f5db68a58750361",
			ObjectKey: "properties/1/attachments/apricot.jpg",
			Property:  db.Property{ID: createdProps[0].ID},
		}, http.StatusForbidden, testConnection.accounts.user.token, "basic user create"},
		// Should pass as user is admin
		{models.CreatePropertyAttachment{
			Label:     "Test Attachment 1",
			FileName:  "apricot.jpg",
			FileSize:  100,
			FileType:  "jpg",
			ETag:      "7d219e22bacfe3a56f5db68a58750361",
			ObjectKey: "properties/1/attachments/apricot.jpg",
			Property:  db.Property{ID: createdProps[0].ID},
		}, http.StatusCreated, testConnection.accounts.admin.token, "admin create"},
		// Create should be disallowed due to invalid details
		{models.CreatePropertyAttachment{
			Label:     "ili",
			ObjectKey: "ape",
			ETag:      "df834",
			FileType:  "d",
			FileSize:  100,
			Property:  db.Property{ID: createdProps[0].ID},
		}, http.StatusBadRequest, testConnection.accounts.admin.token, "invalid detail length create"},
	}

	// Create a request url with an "id" URL parameter
	requestUrl := "/api/property-attachments"

	for _, v := range createTests {

		// Make new request with work type creation in body
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
		var body db.PropertyAttachment
		var foundAttachment db.PropertyAttachment
		// Grab ID from response body
		json.Unmarshal(rr.Body.Bytes(), &body)

		// Find the created transaction (to obtain full data with ID)
		testConnection.dbClient.Find(foundAttachment, uint(body.ID))

		// Compare found details with those found in returned body
		checkPropertyAttachmentDetails(&body, &foundAttachment, t, true)

		// If the task log was created successfully, check that it's deleted after test
		if v.expectedResponseStatus == http.StatusCreated {
			// Cleanup
			//
			// Delete the created property attachment
			deleteMainResult := testConnection.dbClient.Delete(&db.PropertyAttachment{}, uint(body.ID))
			if deleteMainResult.Error != nil {
				t.Fatalf("Couldn't clean up created attachments: %v", deleteMainResult.Error)
			}
		}

	}

	// Cleanup
	deleteResult := testConnection.dbClient.Delete(createdProps)
	if deleteResult.Error != nil {
		t.Fatalf("Couldn't clean up seeded properties: %v", deleteResult.Error)
	}

}

// Check the property attachment details
func checkPropertyAttachmentDetails(actual *db.PropertyAttachment, expected *db.PropertyAttachment, t *testing.T, checkId bool) {
	// Only check ID if parameter checkId is true
	if checkId == true {
		// Verify that the actual id matches the created id
		if actual.ID != expected.ID {
			t.Errorf("found property attachment has incorrect ID: expected %d, got %d", expected.ID, actual.ID)
		}
	}

	// Check details
	if actual.ETag != expected.ETag {
		t.Errorf("found etag has incorrect value: expected %s, got %s", expected.ETag, actual.ETag)
	}
	if actual.FileName != expected.FileName {
		t.Errorf("found filename has incorrect value: expected %s, got %s", expected.FileName, actual.FileName)
	}
	if actual.FileType != expected.FileType {
		t.Errorf("found file type has incorrect value: expected %s, got %s", expected.FileType, actual.FileType)
	}
	if actual.FileSize != expected.FileSize {
		t.Errorf("found file size has incorrect value: expected %d, got %d", expected.FileSize, actual.FileSize)
	}
	if actual.Label != expected.Label {
		t.Errorf("found label has incorrect value: expected %s, got %s", expected.Label, actual.Label)
	}
	if actual.ObjectKey != expected.ObjectKey {
		t.Errorf("found object key has incorrect value: expected %s, got %s", expected.ObjectKey, actual.ObjectKey)
	}
}

// Deletes the first found property attachment
func deleteCreatedPropertyAttachment(testConnection *TestDbRepo) error {
	propertyAttachments := []db.PropertyAttachment{}
	result := testConnection.dbClient.Find(&propertyAttachments)
	if result.Error != nil {
		return result.Error
	}
	// Delete created property attachment
	result = testConnection.dbClient.Delete(&propertyAttachments, &propertyAttachments[0].ID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
