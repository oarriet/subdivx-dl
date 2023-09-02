package utils

import (
	"fmt"
	"strings"
)

// FormatIntWithCommasAndPoints formats an integer with points as thousands separators and commas as decimal points.
func FormatIntWithCommasAndPoints(num int) string {
	// Convert the integer to a string
	numStr := fmt.Sprintf("%d", num)

	// Split the integer part and the decimal part
	parts := strings.Split(numStr, ".")

	// Format the integer part with points as thousands separators
	integerPart := parts[0]
	formattedInteger := ""
	for i, c := range integerPart {
		if i > 0 && (len(integerPart)-i)%3 == 0 {
			formattedInteger += "."
		}
		formattedInteger += string(c)
	}

	// Combine the formatted integer part and the decimal part
	formattedNum := formattedInteger
	if len(parts) > 1 {
		formattedNum += "," + parts[1]
	}

	return formattedNum
}
