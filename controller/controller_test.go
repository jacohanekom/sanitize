package controller

import (
	"encoding/json"
	"log"
	"os"
	"sanitize/data"
	td "sanitize/testdata"
	"testing"
)

const databaseName = "test.db"
const sampleDatabase = "sqlite;test.db"

func TestInsert(t *testing.T) {
	if _, err := os.Stat(databaseName); err == nil {
		os.Remove(databaseName)
	}

	db, err := data.Initialize(sampleDatabase)
	if err != nil {
		t.Error(err)
	}

	result, err := doCrudOperation(INSERT, SanitizeWord{Words: []string{"TestInsert", "TestInsert2"}}, &db)
	if err != nil {
		t.Error(err)
	}

	if len(result.Words) != 2 {
		t.Error("Expected 2 words, got ", len(result.Words))
	}
}

func TestInsertPrimaryKey(t *testing.T) {
	if _, err := os.Stat(databaseName); err == nil {
		os.Remove(databaseName)
	}

	db, err := data.Initialize(sampleDatabase)
	if err != nil {
		t.Error(err)
	}

	result, err := doCrudOperation(INSERT, SanitizeWord{Words: []string{"TestInsert", "TestInsert"}}, &db)
	if err != nil {
		t.Error(err)
	}

	if len(result.Words) != 1 {
		t.Error("Expected 1 word, got ", len(result.Words))
	}
}

func TestDelete(t *testing.T) {
	jsonString, _ := json.Marshal(td.Data)
	err := os.WriteFile("sql_sensitive_list.json", jsonString, 0644)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(databaseName); err == nil {
		os.Remove(databaseName)
	}

	db, err := data.Initialize(sampleDatabase)
	if err != nil {
		t.Error(err)
	}

	result, err := doCrudOperation(DELETE, SanitizeWord{Words: []string{"ACTION"}}, &db)
	if err != nil {
		t.Error(err)
	}

	if len(result.Words) != 1 {
		t.Error("Expected 1 words, got ", len(result.Words))
	}

	result, err = doCrudOperation(SELECT, SanitizeWord{}, &db)
	if len(result.Words) != 227 {
		t.Error("Expected 227 words, got ", len(result.Words))
	}
}

func TestUpdate(t *testing.T) {
	if _, err := os.Stat(databaseName); err == nil {
		os.Remove(databaseName)
	}

	db, err := data.Initialize(sampleDatabase)
	if err != nil {
		t.Error(err)
	}

	TestInsertPrimaryKey(t)
	result, err := doCrudOperation(UPDATE, SanitizeWord{Words: []string{"TestInsert", "TestInsertUPDATED"}}, &db)
	if err != nil {
		t.Error(err)
	}

	result, err = doCrudOperation(SELECT, SanitizeWord{}, &db)
	present := false
	for _, word := range result.Words {
		if word == "TestInsertUPDATED" {
			present = true
			break
		}
	}

	if !present {
		t.Error("Expected TestInsertUPDATED not present")
	}
}

func TestUpdateInvalidParameter(t *testing.T) {
	if _, err := os.Stat(databaseName); err == nil {
		os.Remove(databaseName)
	}

	db, err := data.Initialize(sampleDatabase)
	if err != nil {
		t.Error(err)
	}

	_, err = doCrudOperation(UPDATE, SanitizeWord{Words: []string{"TestInsert"}}, &db)
	if err == nil {
		t.Error("validation fails on incorrect parameters")
	}
}
