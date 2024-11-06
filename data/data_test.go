package data

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	td "sanitize/testdata"
)

const databaseName = "test.db"
const sampleDatabase = "sqlite;test.db"

var entry uint

func TestSetup(t *testing.T) {
	//Make sure we are running in a clean environment, delete the sql lite database
	jsonString, _ := json.Marshal(td.Data)
	err := os.WriteFile("sql_sensitive_list.json", jsonString, 0644)
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(databaseName); err == nil {
		os.Remove(databaseName)
	}

	db, err := Initialize(sampleDatabase)
	defer func() {
		os.Remove("sql_sensitive_list.json")
		err := db.removeDatabaseFile(sampleDatabase)
		if err != nil {
			log.Fatal("Unable to remove test database")
		}
	}()

	if err != nil {
		t.Fatalf(err.Error())
	}

	loaded, err := db.ListRecords()
	if err != nil {
		t.Fatalf(err.Error())
	}

	if len(loaded) != 228 {
		t.Fatalf("Expected %d records, got %d", 228, len(loaded))
	}
}

func TestAddEntry(t *testing.T) {
	db, err := Initialize(sampleDatabase)
	defer func() {
		err := db.removeDatabaseFile(sampleDatabase)
		if err != nil {
			log.Fatal("Unable to remove test database")
		}
	}()

	if err != nil {
		t.Fatalf(err.Error())
	}

	entry, err = db.AddEntry("TestAddEntry")
	if err != nil {
		t.Fatalf(err.Error())
	}

	loaded, err := db.ListRecords()
	if err != nil {
		t.Fatalf(err.Error())
	}

	if _, ok := loaded[entry]; !ok {
		t.Fatalf("Expected entry to exist")
	}
}

func TestAddEntryUnique(t *testing.T) {
	db, err := Initialize(sampleDatabase)
	defer func() {
		err := db.removeDatabaseFile(sampleDatabase)
		if err != nil {
			log.Fatal("Unable to remove test database")
		}
	}()

	if err != nil {
		t.Fatalf(err.Error())
	}

	entry, err = db.AddEntry("TestAddEntry")
	if err != nil {
		t.Fatalf(err.Error())
	}

	entry, err = db.AddEntry("TestAddEntry")
	if err == nil {
		t.Fatalf("duplicate test not valid")
	}
}

func TestAddRemoveEntry(t *testing.T) {
	db, err := Initialize(sampleDatabase)
	defer func() {
		err := db.removeDatabaseFile(sampleDatabase)
		if err != nil {
			log.Fatal("Unable to remove test database")
		}
	}()

	if err != nil {
		t.Fatalf(err.Error())
	}

	entry, err = db.AddEntry("TestAddEntry")
	if err != nil {
		t.Fatalf(err.Error())
	}

	err = db.RemoveEntry(entry)
	if err != nil {
		t.Fatalf(err.Error())
	}

	loaded, err := db.ListRecords()
	if _, ok := loaded[entry]; ok {
		t.Fatalf("Expected entry to exist")
	}
}
