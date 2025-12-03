package validation

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type ValidationError struct {
	Field   string
	Tag     string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

type ValidationErrors []ValidationError

func (e ValidationErrors) Error() string {
	var msgs []string
	for _, err := range e {
		msgs = append(msgs, err.Error())
	}
	return strings.Join(msgs, "; ")
}

func Validate(s interface{}) error {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return fmt.Errorf("validate: expected struct, got %s", v.Kind())
	}

	var errs ValidationErrors
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		tag := fieldType.Tag.Get("validate")

		if tag == "" || tag == "-" {
			continue
		}

		fieldName := fieldType.Name
		rules := strings.Split(tag, ",")

		for _, rule := range rules {
			rule = strings.TrimSpace(rule)
			if err := validateRule(field, fieldName, rule); err != nil {
				errs = append(errs, *err)
			}
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func validateRule(field reflect.Value, fieldName, rule string) *ValidationError {
	parts := strings.SplitN(rule, "=", 2)
	ruleName := parts[0]
	var ruleValue string
	if len(parts) > 1 {
		ruleValue = parts[1]
	}

	switch ruleName {
	case "required":
		return validateRequired(field, fieldName)
	case "min":
		return validateMin(field, fieldName, ruleValue)
	case "max":
		return validateMax(field, fieldName, ruleValue)
	case "len":
		return validateLen(field, fieldName, ruleValue)
	case "email":
		return validateEmail(field, fieldName)
	case "url":
		return validateURL(field, fieldName)
	case "oneof":
		return validateOneOf(field, fieldName, ruleValue)
	case "gt":
		return validateGt(field, fieldName, ruleValue)
	case "gte":
		return validateGte(field, fieldName, ruleValue)
	case "lt":
		return validateLt(field, fieldName, ruleValue)
	case "lte":
		return validateLte(field, fieldName, ruleValue)
	}

	return nil
}

func validateRequired(field reflect.Value, fieldName string) *ValidationError {
	if isZero(field) {
		return &ValidationError{
			Field:   fieldName,
			Tag:     "required",
			Message: "field is required",
		}
	}
	return nil
}

func validateMin(field reflect.Value, fieldName, min string) *ValidationError {
	switch field.Kind() {
	case reflect.String:
		minLen, _ := strconv.Atoi(min)
		if len(field.String()) < minLen {
			return &ValidationError{
				Field:   fieldName,
				Tag:     "min",
				Message: fmt.Sprintf("must be at least %s characters", min),
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		minVal, _ := strconv.ParseInt(min, 10, 64)
		if field.Int() < minVal {
			return &ValidationError{
				Field:   fieldName,
				Tag:     "min",
				Message: fmt.Sprintf("must be at least %s", min),
			}
		}
	case reflect.Float32, reflect.Float64:
		minVal, _ := strconv.ParseFloat(min, 64)
		if field.Float() < minVal {
			return &ValidationError{
				Field:   fieldName,
				Tag:     "min",
				Message: fmt.Sprintf("must be at least %s", min),
			}
		}
	}
	return nil
}

func validateMax(field reflect.Value, fieldName, max string) *ValidationError {
	switch field.Kind() {
	case reflect.String:
		maxLen, _ := strconv.Atoi(max)
		if len(field.String()) > maxLen {
			return &ValidationError{
				Field:   fieldName,
				Tag:     "max",
				Message: fmt.Sprintf("must be at most %s characters", max),
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		maxVal, _ := strconv.ParseInt(max, 10, 64)
		if field.Int() > maxVal {
			return &ValidationError{
				Field:   fieldName,
				Tag:     "max",
				Message: fmt.Sprintf("must be at most %s", max),
			}
		}
	case reflect.Float32, reflect.Float64:
		maxVal, _ := strconv.ParseFloat(max, 64)
		if field.Float() > maxVal {
			return &ValidationError{
				Field:   fieldName,
				Tag:     "max",
				Message: fmt.Sprintf("must be at most %s", max),
			}
		}
	}
	return nil
}

func validateLen(field reflect.Value, fieldName, length string) *ValidationError {
	if field.Kind() == reflect.String {
		expectedLen, _ := strconv.Atoi(length)
		if len(field.String()) != expectedLen {
			return &ValidationError{
				Field:   fieldName,
				Tag:     "len",
				Message: fmt.Sprintf("must be exactly %s characters", length),
			}
		}
	}
	return nil
}

func validateEmail(field reflect.Value, fieldName string) *ValidationError {
	if field.Kind() == reflect.String {
		email := field.String()
		emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
		if !emailRegex.MatchString(email) {
			return &ValidationError{
				Field:   fieldName,
				Tag:     "email",
				Message: "must be a valid email address",
			}
		}
	}
	return nil
}

func validateURL(field reflect.Value, fieldName string) *ValidationError {
	if field.Kind() == reflect.String {
		url := field.String()
		urlRegex := regexp.MustCompile(`^https?://[^\s/$.?#].[^\s]*$`)
		if !urlRegex.MatchString(url) {
			return &ValidationError{
				Field:   fieldName,
				Tag:     "url",
				Message: "must be a valid URL",
			}
		}
	}
	return nil
}

func validateOneOf(field reflect.Value, fieldName, values string) *ValidationError {
	if field.Kind() == reflect.String {
		val := field.String()
		options := strings.Split(values, " ")
		for _, opt := range options {
			if val == opt {
				return nil
			}
		}
		return &ValidationError{
			Field:   fieldName,
			Tag:     "oneof",
			Message: fmt.Sprintf("must be one of: %s", values),
		}
	}
	return nil
}

func validateGt(field reflect.Value, fieldName, value string) *ValidationError {
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val, _ := strconv.ParseInt(value, 10, 64)
		if field.Int() <= val {
			return &ValidationError{
				Field:   fieldName,
				Tag:     "gt",
				Message: fmt.Sprintf("must be greater than %s", value),
			}
		}
	case reflect.Float32, reflect.Float64:
		val, _ := strconv.ParseFloat(value, 64)
		if field.Float() <= val {
			return &ValidationError{
				Field:   fieldName,
				Tag:     "gt",
				Message: fmt.Sprintf("must be greater than %s", value),
			}
		}
	}
	return nil
}

func validateGte(field reflect.Value, fieldName, value string) *ValidationError {
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val, _ := strconv.ParseInt(value, 10, 64)
		if field.Int() < val {
			return &ValidationError{
				Field:   fieldName,
				Tag:     "gte",
				Message: fmt.Sprintf("must be greater than or equal to %s", value),
			}
		}
	case reflect.Float32, reflect.Float64:
		val, _ := strconv.ParseFloat(value, 64)
		if field.Float() < val {
			return &ValidationError{
				Field:   fieldName,
				Tag:     "gte",
				Message: fmt.Sprintf("must be greater than or equal to %s", value),
			}
		}
	}
	return nil
}

func validateLt(field reflect.Value, fieldName, value string) *ValidationError {
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val, _ := strconv.ParseInt(value, 10, 64)
		if field.Int() >= val {
			return &ValidationError{
				Field:   fieldName,
				Tag:     "lt",
				Message: fmt.Sprintf("must be less than %s", value),
			}
		}
	case reflect.Float32, reflect.Float64:
		val, _ := strconv.ParseFloat(value, 64)
		if field.Float() >= val {
			return &ValidationError{
				Field:   fieldName,
				Tag:     "lt",
				Message: fmt.Sprintf("must be less than %s", value),
			}
		}
	}
	return nil
}

func validateLte(field reflect.Value, fieldName, value string) *ValidationError {
	switch field.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val, _ := strconv.ParseInt(value, 10, 64)
		if field.Int() > val {
			return &ValidationError{
				Field:   fieldName,
				Tag:     "lte",
				Message: fmt.Sprintf("must be less than or equal to %s", value),
			}
		}
	case reflect.Float32, reflect.Float64:
		val, _ := strconv.ParseFloat(value, 64)
		if field.Float() > val {
			return &ValidationError{
				Field:   fieldName,
				Tag:     "lte",
				Message: fmt.Sprintf("must be less than or equal to %s", value),
			}
		}
	}
	return nil
}

func isZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.String() == ""
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan, reflect.Func:
		return v.IsNil()
	}
	return false
}

// Example usage:
/*
package main

import (
	"fmt"
	"your-package/validator"
)

type User struct {
	Name     string  `validate:"required,min=3,max=50"`
	Email    string  `validate:"required,email"`
	Age      int     `validate:"required,gte=18,lte=100"`
	Website  string  `validate:"url"`
	Role     string  `validate:"oneof=admin user guest"`
	Balance  float64 `validate:"gte=0"`
}

func main() {
	user := User{
		Name:    "Jo",
		Email:   "invalid-email",
		Age:     15,
		Website: "not-a-url",
		Role:    "superadmin",
		Balance: -10.5,
	}

	if err := validator.Validate(user); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			for _, e := range errs {
				fmt.Printf("Field: %s, Error: %s\n", e.Field, e.Message)
			}
		}
	}
}
*/
