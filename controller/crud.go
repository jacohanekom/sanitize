package controller

import (
	"errors"
	"log"
	"runtime/debug"
	"sanitize/data"
	"strings"
)

type crudOperation int

const (
	SELECT crudOperation = iota
	INSERT
	UPDATE
	DELETE
)

func doCrudOperation(operation crudOperation, request SanitizeWord, database *data.SanitizeDB) (result SanitizeWord, crudError error) {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			crudError = errors.New("an error occurred, while performing crud operation")
		}
	}()

	records, err := database.ListRecords()
	if err != nil {
		crudError = err
	}

	switch operation {
	case SELECT:
		for _, v := range records {
			result.Words = append(result.Words, v)
		}

		return result, nil
	case INSERT:
		for _, rec := range request.Words {
			_, err = database.AddEntry(rec)
			if err != nil {
				log.Printf("Unable to add %s reason %v", rec, err)
			} else {
				result.Words = append(result.Words, rec)
			}
		}

		if len(result.Words) == 0 {
			crudError = errors.New("unable to add entries")
			return
		}
	case UPDATE:
		var removed bool

		if len(request.Words)%2 > 0 {
			crudError = errors.New("update parameters not correctly specified")
			return
		}

		for _, rec := range request.Words {
			if !removed {
				result, err = doCrudOperation(DELETE, SanitizeWord{Words: []string{rec}}, database)
				if err != nil && result.Words == nil || len(result.Words) == 0 {
					crudError = errors.New("unable to remove entries")
					return
				}
				removed = true
			} else {
				result, err = doCrudOperation(INSERT, SanitizeWord{Words: []string{rec}}, database)
				if err != nil && result.Words == nil || len(result.Words) == 0 {
					crudError = errors.New("unable to add entries")
					return
				}
				removed = false
			}
		}
	case DELETE:
		for _, rec := range request.Words {
			for id, currentWord := range records {
				if currentWord == strings.ToUpper(rec) {
					deleteErr := database.RemoveEntry(id)
					if deleteErr != nil {
						log.Printf("Unable to remove entry %v: %v", id, deleteErr)
					} else {
						result.Words = append(result.Words, rec)
					}
				}
			}
		}

		if len(result.Words) == 0 {
			crudError = errors.New("unable remove entries")
			return
		}
	default:
		log.Println("Unsupported operation")
		return SanitizeWord{}, errors.New("unhandled default case")
	}

	return result, nil
}
