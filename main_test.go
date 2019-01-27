package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var a Apps

func TestMain(m *testing.M) {
	a = Apps{}
	a.Initialize(os.Getenv("USERNSME"), os.Getenv("PASSWORD"), os.Getenv("DBSTRING"))
	ensureTableExists()

	code := m.Run()

	clearTable()

	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM contact")
	a.DB.Exec("ALTER TABLE contact AUTO_INCREMENT = 1")
}

const tableCreationQuery = `
CREATE TABLE IF NOT EXISTS contact
(
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL
)`

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/contact", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestGetNonExistentUser(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/contact/45", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Contact not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Contact not found'. Got '%s'", m["error"])
	}
}

func TestCreateUser(t *testing.T) {
	clearTable()

	payload := []byte(`{"name":"test user","email":"test@gmail.com"}`)

	req, _ := http.NewRequest("POST", "/contact", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "test user" {
		t.Errorf("Expected user name to be 'test user'. Got '%v'", m["name"])
	}

	if m["email"] != "test@gmail.com" {
		t.Errorf("Expected email to be 'test@gmail.com'. Got '%v'", m["email"])
	}

	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["id"] != 1.0 {
		t.Errorf("Expected contact ID to be '1'. Got '%v'", m["id"])
	}
}

func addUsers() {

	statement := fmt.Sprintf("INSERT INTO contact(name, email) VALUES('%s', '%s')", "User", "test@gmail.com")
	a.DB.Exec(statement)

}

func TestGetUser(t *testing.T) {
	clearTable()
	addUsers()

	req, _ := http.NewRequest("GET", "/contact/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateUser(t *testing.T) {
	clearTable()
	addUsers()

	req, _ := http.NewRequest("GET", "/contact/1", nil)
	response := executeRequest(req)
	var originalUser map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalUser)

	payload := []byte(`{"name":"test user - updated name","email":"test1@gmail.com"}`)

	req, _ = http.NewRequest("PUT", "/contact/1", bytes.NewBuffer(payload))
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["id"] != originalUser["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalUser["id"], m["id"])
	}

	if m["name"] == originalUser["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalUser["name"], m["name"], m["name"])
	}

	if m["email"] == originalUser["email"] {
		t.Errorf("Expected the email to change from '%v' to '%v'. Got '%v'", originalUser["email"], m["email"], m["email"])
	}
}

func TestDeleteUser(t *testing.T) {
	clearTable()
	addUsers()

	req, _ := http.NewRequest("GET", "/contact/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/contact/1", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/contact/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}
