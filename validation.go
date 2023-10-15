package main

import (
	"fmt"
	"strings"
)

type Validator struct {
	Errors map[string][]string
}

type ValidUsername string
type ValidPassword string

// NewValidator returns a *Validator struct with an initialized map of Errors
func NewValidator() *Validator{
	return &Validator{
		Errors: make(map[string][]string),
	}
}

// IsValid checks if the Validator contains any entries into the map of Errors
func (v *Validator) IsValid() bool {
	return len(v.Errors) < 1
}

// AddError adds an error to the map of Errors, with the given key and message
func (v *Validator) AddError(errorMessage, errorKey string){
	_, ok := v.Errors[errorKey]
	if ok{
		v.Errors[errorKey] = append(v.Errors[errorKey], errorMessage)
	} else{
		v.Errors[errorKey] = []string{errorMessage}
	}
	
}

// CheckMinLength checks if a string length is less than a given minimum, and adds to the map of Errors if it is
func (v *Validator) CheckMinLength(minLength int, dataToValidate, errorKey string){
	if len(dataToValidate) < minLength{
		errorMessage := fmt.Sprintf("Should be at least %d characters long", minLength)
		v.AddError(errorMessage, errorKey)
	}
}

// CheckMinLength checks if a string length is more than a given maximum, and adds to the map of Errors if it is
func (v *Validator) CheckMaxLength(maxLength int, dataToValidate, errorKey string){
	if len(dataToValidate) > maxLength{
		errorMessage := fmt.Sprintf("Should be less than %d characters long", maxLength)
		v.AddError(errorMessage, errorKey)
	}
}

// CheckUnallowedCharacters checks if a string contains any of the characters in the given chars string, and adds to the map of Errors if it is
func (v *Validator) CheckUnallowedCharacters(dataToValidate, chars, errorKey string){
	if strings.ContainsAny(dataToValidate, chars) {
		charListString := strings.Join(strings.Split(chars, ""), ", ")
		errorMessage := fmt.Sprintf("Cannot contain certain characters (%s)", charListString)
		v.AddError(errorMessage, errorKey)
	}
}

// CheckRequiredCharacterGroup checks if a string contains at least one of the characters in the given chars string, and adds to the map of Errors if it does not
func (v *Validator) CheckRequiredCharacterGroup(dataToValidate, chars, errorMessage, errorKey string){
	if !strings.ContainsAny(dataToValidate, chars){
		v.AddError(errorMessage, errorKey)
	}
}

// CheckOnlyAllowedCharacters checks if a string contains only the characters in the chars string, and adds to the map of Errors if it does not
func (v *Validator) CheckOnlyAllowedCharacters(dataToValidate, chars, errorMessage, errorKey string){
	for _, char := range dataToValidate{
		if !strings.Contains(chars, string(char)){
			v.AddError(errorMessage, errorKey)
		}
	}
}

// ValidateUsername validates a username string given a standard set of rules
func (v *Validator) ValidateUsername(username string) {
	errorKey := "username"

	v.CheckMinLength(4, username, errorKey)
	v.CheckMaxLength(24, username, errorKey)
	v.CheckOnlyAllowedCharacters(username, "abcdefghijklmnopqrstuvwxyz1234567890._", "Can only contain alphanumeric characters, \".\" or \"_\"", errorKey)
}

// ValidatePassword validates a password given a standard set of rules
func (v *Validator) ValidatePassword(password string) {
	errorKey := "password"

	v.CheckMinLength(12, password, errorKey)
	v.CheckMaxLength(128, password, errorKey)

	//Check for 1 special character
	v.CheckRequiredCharacterGroup(password, "~`!@#$%^&*()_-+={[}]|\\:;\"'<,>.?/", "Must contain at least one special character", errorKey)
	//Check for 1 capital letter
	v.CheckRequiredCharacterGroup(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ", "Must contain at least one capital letter", errorKey)
	//Check for 1 lowercase
	v.CheckRequiredCharacterGroup(password, "abcdefghijklmnopqrstuvwxyz", "Must contain at least one lowercase letter", errorKey)
	//Check for 1 number
	v.CheckRequiredCharacterGroup(password, "1234567890", "Must contain at least 1 number", errorKey)
}