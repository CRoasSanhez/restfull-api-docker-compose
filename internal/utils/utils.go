package utils

import (
	"errors"
	"log"
	"reflect"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

// CopyStruct ...
func CopyStruct(in interface{}, out interface{}) error {
	v := reflect.ValueOf(in).Elem()

	if !isStructType(in) || !isStructType(out) {
		return errors.New("In and out must be structs")
	}

	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := v.Type().Field(i)

		outFieldValue := reflect.ValueOf(out).Elem()

		f := outFieldValue.FieldByName(fieldType.Name)
		if f.IsValid() && f.CanSet() {
			f.Set(fieldValue)
		}
	}

	return nil
}

func isStructType(s interface{}) bool {
	t := reflect.TypeOf(s)
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	return t.Kind() == reflect.Struct
}

// IsValidPhone ...
func IsValidPhone(phone string) bool {
	re := regexp.MustCompile(`^([+]?\d{1,2}[.-\s]?)?(\d{3}[.-]?){2}\d{4}`)
	return re.MatchString(phone)
}

// IsValidCard ...
func IsValidCard(cardNumber string) bool {
	re := regexp.MustCompile(`^(?:4[0-9]{12}(?:[0-9]{3})?|[25][1-7][0-9]{14}|6(?:011|5[0-9][0-9])[0-9]{12}|3[47][0-9]{13}|3(?:0[0-5]|[68][0-9])[0-9]{11}|(?:2131|1800|35\d{3})\d{11})$`)
	return re.MatchString(cardNumber)
}

// IsValidEmail ...
func IsValidEmail(email string) bool {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return re.MatchString(email)
}

// IsValidAmunt ...
func IsValidAmunt(amount int) bool {
	return amount <= 100000
}

// IsValidExpirationDate ...
func IsValidExpirationDate(mmyy string) bool {
	re := regexp.MustCompile(`^(0[1-9]|1[0-2])\/?([0-9]{4}|[0-9]{2})$/`)
	return re.MatchString(mmyy)
}

// IsCamelCase ...
func IsCamelCase(word string) bool {
	re := regexp.MustCompile("^[a-z]+(?:[A-Z][a-z]+)*$")
	return re.MatchString(word)
}

// IsEven ...
func IsEven(number int) bool {
	return number%2 == 0
}

// IsOdd ...
func IsOdd(number int) bool {
	return !IsEven(number)
}

// HashPassword ...
func HashPassword(password string) (hashedPassword string) {

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	return string(hash)

}

// CheckPassword ...
func CheckPassword(hashedPassword string, password string) (isPasswordValid bool) {

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return false
	}

	return true
}
