package controller

import (
	"errors"
	"maps"
	"runtime/debug"
	"sanitize/data"
	"sanitize/logic"
	"slices"
)

func doSanitize(request Sanitize, database *data.SanitizeDB) (result Sanitize, returnError error) {
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			returnError = errors.New("an error occurred, while sanitizing text")
		}
	}()

	currentWords, err := database.ListRecords()
	if err != nil {
		returnError = err
		return
	}

	result.Sentences, err = logic.SanitizeText(request.Sentences, slices.Collect(maps.Values(currentWords)))
	if err != nil {
		returnError = err
		return
	}

	return
}
