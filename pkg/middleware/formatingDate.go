package middleware

import (
	"soteria_go/pkg/models/response"
	"time"
)

func FormatingDate(dateString, username, instiCode, appCode, moduleName, funcName, methodUsed, endpoint string) response.ReturnModel {

	// List of supported date layouts
	dateLayouts := []string{
		"2006-01-02",      // YYYY-MM-DD
		"01-02-2006",      // MM-DD-YYYY
		"01/02/2006",      // MM/DD/YYYY
		"02-Jan-2006",     // DD-Mon-YYYY (e.g., 02-Feb-2006)
		"January 2 2006",  // Month day YYYY (e.g., February 2 2006)
		"January 2, 2006", // Month day YYYY (e.g., February 2, 2006)
		"Jan 2 2006",      // Mon day YYYY (e.g., Feb 2 2006)
		"Jan 2, 2006",     // Mon day YYYY (e.g., Feb 2, 2006)
		"Jan. 2, 2006",    // Mon day YYYY (e.g., Feb. 2, 2006)
		"02 Jan 2006",     // DD Mon YYYY (e.g., 02 Feb 2006)
		"2 Jan 2006",      // D Mon YYYY (e.g., 2 Feb 2006)
	}

	var parsedDate time.Time
	var formatedDatErr error
	for _, layout := range dateLayouts {
		parsedDate, formatedDatErr = time.Parse(layout, dateString)
		if formatedDatErr == nil {
			break
		}
	}

	if formatedDatErr != nil {
		returnMessage := ResponseData(username, instiCode, appCode, moduleName, funcName, "107", methodUsed, endpoint, []byte(dateString), []byte(""), "", nil, nil)
		if !returnMessage.Data.IsSuccess {
			return (returnMessage)
		}
	}

	// Formatting the parsed date into desired format
	formattedDate := parsedDate.Format("2006-01-02")

	ActivityLogger(username, instiCode, appCode, moduleName, funcName, "200", methodUsed, endpoint, []byte(dateString), []byte(formattedDate), "Successful", "", nil)
	return response.ReturnModel{
		RetCode: "200",
		Message: "Successful",
		Data: response.DataModel{
			Message:   formattedDate,
			IsSuccess: true,
			Error:     nil,
		},
	}
}
