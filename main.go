package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
)

// Struct tag constants
const (
	RequiredTag = "required"
	MaxTag      = "max"
	MinTag      = "min"
	MaxLenTag   = "maxLen"
	MinLenTag   = "minLen"
	EmailTag    = "email"
)

// ValidationResult stores the result of validation.
type ValidationResult struct {
	Field   string
	IsValid bool
	Message string
}

// ValidateStruct melakukan validasi terhadap objek struct berdasarkan tag.
func ValidateStruct(data interface{}) []ValidationResult {
	valueOf := reflect.ValueOf(data)
	typeOf := valueOf.Type()
	numFields := valueOf.NumField()
	results := make([]ValidationResult, numFields)

	for i := 0; i < numFields; i++ {
		fieldValue := valueOf.Field(i)
		fieldType := typeOf.Field(i)
		fieldName := fieldType.Name
		tag := fieldType.Tag

		var message string

		if requiredTag := tag.Get(RequiredTag); requiredTag != "" && fieldValue.IsZero() {
			message = "is required"
		}

		if maxTag := tag.Get(MaxTag); maxTag != "" {
			if intValue, err := strconv.Atoi(maxTag); err == nil && fieldValue.Int() > int64(intValue) {
				message = "exceeds the maximum allowed"
			}
		}

		if minTag := tag.Get(MinTag); minTag != "" {
			if intValue, err := strconv.Atoi(minTag); err == nil && fieldValue.Int() < int64(intValue) {
				message = "is below the minimum allowed"
			}
		}

		if maxLenTag := tag.Get(MaxLenTag); maxLenTag != "" {
			if maxLength, err := strconv.Atoi(maxLenTag); err == nil && len(fieldValue.String()) > maxLength {
				message = "exceeds maximum length allowed"
			}
		}

		if emailTag := tag.Get(EmailTag); emailTag != "" && !isValidEmail(fieldValue.String()) {
			message = "has an invalid email format"
		}

		// Field is valid if message is empty
		isValid := message == ""

		results[i] = ValidationResult{
			Field:   fieldName,
			IsValid: isValid,
			Message: fmt.Sprintf("%s %s", fieldName, message),
		}
	}

	return results
}

// isValidEmail memeriksa apakah sebuah string adalah alamat email yang valid menggunakan ekspresi reguler.
func isValidEmail(email string) bool {
	// Gunakan ekspresi reguler yang sederhana untuk memeriksa format email
	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,4}$`
	return regexp.MustCompile(emailPattern).MatchString(email)
}

// Contoh penggunaan:
type AvengersMember struct {
	Name  string `required:"true" maxLen:"50"`
	Age   int    `min:"18" max:"100"`
	Email string `required:"true" email:"true"`
}

func main() {
	member := AvengersMember{
		Name:  "Iron Man",
		Age:   45,
		Email: "ironman@example.com",
	}

	validationResults := ValidateStruct(member)
	for _, result := range validationResults {
		if !result.IsValid {
			fmt.Printf("Validation Error: %s\n", result.Message)
		}
	}
}
