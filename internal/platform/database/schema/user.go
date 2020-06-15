package schema

import (
	"encoding/json"
	"regexp"
	"strings"
	"time"

	"github.com/CRoasSanhez/yofio-test/internal/utils"
	"github.com/jinzhu/gorm"
)

// User ...
type User struct {
	gorm.Model
	ID            int       `gorm:"primary_key"`
	FullName      string    `gorm:"type:varchar(150)" json:"name"`
	Phone         string    `gorm:"type:varchar(13);unique_index" json:"phone"`
	Email         string    `gorm:"type:varchar(100);unique_index" json:"email"`
	Pwd           string    `gorm:"type:varchar(100)" json:"password"`
	LoginFailures int       `gorm:"type:int(1)" json:"-" default:"0"`
	InsertedAt    time.Time ` json:"-"`
	UpdatedAt     time.Time ` json:"-"`
	IsBlocked     bool      `json:"is_blocked" default:"false"`
	IsDeleted     bool      ` json:"is_deleted" default:"false"`
}

// MarshalJSON return bot struct as JSON.
func (u User) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
		Email string `json:"email"`
	}{
		Name:  u.FullName,
		Phone: u.Phone,
		Email: u.Email,
	})
}

// AddLoginAttempts ...
func (u User) AddLoginAttempts(db *gorm.DB) {
	u.UpdatedAt = time.Now()
	db.Update("LoginFailures", "UpdatedAt")
}

// IsValidFullName ...
func (u User) IsValidFullName() bool {
	var isValid = true

	// Validate more than 1 word
	splitted := strings.Split(u.FullName, " ")
	if len(splitted) <= 1 {
		return !isValid
	}

	// Looking for not camelcase words
	for _, w := range splitted {
		if !utils.IsCamelCase(w) {
			isValid = false
			break
		}
	}

	return isValid
}

// IsValidPassword ...
func (u User) IsValidPassword() bool {
	reLow := regexp.MustCompile(`[[:lower:]]{1,}`)
	reDigit := regexp.MustCompile(`[[:digit:]]{1,}`)
	reUpper := regexp.MustCompile(`[[:upper:]]{1,}`)

	if !reLow.MatchString(u.Pwd) || !reDigit.MatchString(u.Pwd) || !reUpper.MatchString(u.Pwd) || len(u.Pwd) <= 8 {
		return false
	}
	return true
}
