package data

import (
	"encoding/json"
	"errors"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"io"
	"os"
	"strings"
)

/*
/*This package uses GORM, which is a ORM Library for GoLang, to assist CRUD operations and hide any
* database specific logic behind a obfuscation level. This also assist in testing
*/

// this is a private object definition used in the database by Gorm to build the table */
type sensitiveWord struct {
	ID        uint   `gorm:"primaryKey; autoIncrement:true;"`
	Sensitive string `gorm:"index:idx_sensitive,unique"`
}

type SanitizeDB struct {
	db *gorm.DB
}

const dataFileName = "sql_sensitive_list.json"

// Initialize initializes the database connection and provides the SanitizeDB object to perform CRUD operations
// on the database for management sensitive words. Both SQL Lite and MS SQL is supported through GORM, it is recommended that
// SQL Lite is only used for testing purposes. For database initialization and testing the function supports reading a file in
// a json format, representing an string list. The expect filename is sql_sensitive_list.json. This file will be pushed
// to the database and should be removed after the first successful startup.
func Initialize(connectionString string) (SanitizeDB, error) {
	//Reading a data intialize file, this will be used for testing and the initial launch to setup the database with correct data
	var words []string
	var result SanitizeDB

	if _, err := os.Stat(dataFileName); err == nil {
		jsonFile, err := os.Open(dataFileName)
		// if we os.Open returns an error then handle it
		if err != nil {
			return SanitizeDB{}, err
		}

		// defer the closing of our jsonFile so that we can parse it later on
		defer jsonFile.Close()

		byteValue, err := io.ReadAll(jsonFile)
		if err != nil {
			return SanitizeDB{}, err
		}

		//Converting the object into a object, so that it can be parsed
		err = json.Unmarshal(byteValue, &words)
		if err != nil {
			return SanitizeDB{}, err
		}
	}

	var err error
	if strings.Split(connectionString, ";")[0] == "sqlite" {
		result.db, err = gorm.Open(sqlite.Open(strings.Split(connectionString, ";")[1]), &gorm.Config{})
	} else if strings.Split(connectionString, ";")[0] == "sqlserver" {
		result.db, err = gorm.Open(sqlserver.Open(strings.Split(connectionString, ";")[1]), &gorm.Config{})
	}

	if err != nil {
		return SanitizeDB{}, err
	}

	// Migrate the schema
	err = result.db.AutoMigrate(&sensitiveWord{})
	if err != nil {
		return SanitizeDB{}, err
	}

	for _, word := range words {
		result.db.Create(&sensitiveWord{Sensitive: word})
	}

	return result, nil
}

// ListRecords returns a map with the current loaded sanitized keywords and their unique id.
// In the case where it cannot connect or an error occurred the method will return an error
func (sanitize *SanitizeDB) ListRecords() (map[uint]string, error) {
	var words []sensitiveWord

	_, err := sanitize.db.Find(&words).Rows()
	if err != nil {
		return nil, err
	}

	records := make(map[uint]string)
	for _, word := range words {
		records[word.ID] = strings.ToUpper(word.Sensitive)
	}

	return records, nil
}

// RemoveEntry provides the ability to remove an entry from the database
// The unique ID is required to complete the operation. Should no delete occur an "entry not found" error will
// be returned
func (sanitize *SanitizeDB) RemoveEntry(id uint) error {
	result := *sanitize.db.Delete(&sensitiveWord{}, id)
	if result.RowsAffected > 0 {
		return nil
	} else {
		return errors.New("entry not found")
	}
}

// AddEntry provides the ability to add an entry to the database
// There is an unique key index on the name, and should it already be present, or an error occurred
// the function will return "entry not added" error
func (sanitize *SanitizeDB) AddEntry(entry string) (uint, error) {
	insert := sensitiveWord{Sensitive: strings.ToUpper(entry)}

	sanitize.db.Create(&insert)
	if insert.ID == 0 {
		return 0, errors.New("entry not added")
	}
	return insert.ID, nil
}

// removeDatabase is an internal method, and should not be called. But in the case of testing
// this function is used to clean up once the testing has been completed
func (sanitize *SanitizeDB) removeDatabaseFile(connectionString string) error {
	if strings.Split(connectionString, ";")[0] == "sqlite" {
		err := os.Remove(strings.Split(connectionString, ";")[1])
		if err != nil {
			return err
		}
	}

	return nil
}
