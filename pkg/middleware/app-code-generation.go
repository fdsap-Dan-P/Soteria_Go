package middleware

import (
	"fmt"
	"soteria_go/pkg/models/response"
	"soteria_go/pkg/utils/go-utils/database"
	"strings"
	"time"
)

func AppCodeGeneration(appName string) (string, error) {
	rowCount := response.Total{}
	words := strings.Fields(appName)
	acronym := ""

	// Create acronym from the first letter of each word
	for _, word := range words {
		acronym += strings.ToUpper(string(word[0]))
	}

	// Get the total number of rows in the table
	if fetchErr := database.DBConn.Raw("SELECT COUNT(*) FROM public.applications").Scan(&rowCount).Error; fetchErr != nil {
		return "", fetchErr
	}

	rowID := rowCount.Count + 1

	// Ensure RowID is 4 digits, pad with leading zeros if necessary
	formattedID := fmt.Sprintf("%04d", rowID)

	// Get the current Unix timestamp
	timestamp := fmt.Sprintf("%d", time.Now().Unix())

	// Combine acronym, formattedID, and timestamp
	generatedAppCode := fmt.Sprintf("%s%s-%s", acronym, formattedID, timestamp)

	return generatedAppCode, nil
}
