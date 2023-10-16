package main

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
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

		var messages []string

		if requiredTag := tag.Get(RequiredTag); requiredTag != "" {
			if fieldValue.IsZero() {
				messages = append(messages, "Field is required")
			}
		}

		if maxTag := tag.Get(MaxTag); maxTag != "" {
			maxValue, err := strconv.Atoi(maxTag)
			if err == nil && fieldValue.Int() > int64(maxValue) {
				messages = append(messages, "Value exceeds the maximum allowed")
			}
		}

		if minTag := tag.Get(MinTag); minTag != "" {
			minValue, err := strconv.Atoi(minTag)
			if err == nil && fieldValue.Int() < int64(minValue) {
				messages = append(messages, "Value is below the minimum allowed")
			}
		}

		if maxLenTag := tag.Get(MaxLenTag); maxLenTag != "" {
			maxLength, err := strconv.Atoi(maxLenTag)
			if err == nil && len(fieldValue.String()) > maxLength {
				messages = append(messages, "Exceeds maximum length allowed")
			}
		}

		if minLenTag := tag.Get(MinLenTag); minLenTag != "" {
			minLength, err := strconv.Atoi(minLenTag)
			if err == nil && len(fieldValue.String()) < minLength {
				messages = append(messages, "Below the minimum length allowed")
			}
		}

		if emailTag := tag.Get(EmailTag); emailTag != "" {
			if !isValidEmail(fieldValue.String()) {
				messages = append(messages, "Invalid email format")
			}
		}

		// Field is valid if messages are empty
		isValid := len(messages) == 0

		results[i] = ValidationResult{
			Field:   fieldName,
			IsValid: isValid,
			Message: fmt.Sprintf("%s: %s", fieldName, strings.Join(messages, ", ")),
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
