package schema

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm"
)

// Payment ...
type Payment struct {
	gorm.Model
	ID           int       `gorm:"primary_key"`
	MembershipID int       `gorm:"type:int(64)" json:"membership_id"`
	UserID       int       `gorm:"type:int(64)" json:"user_id"`
	Status       string    `gorm:"type:varchar(100)" json:"-"`
	Amount       int       `gorm:"type:int(8)" default:"100000"`
	InsertedAt   time.Time ` json:"-"`
	UpdatedAt    time.Time ` json:"-"`
	IsDeleted    bool      `gorm:"type:TINYINT(1)" json:"is_deleted" default:"false"`
}

// MarshalJSON return bot struct as JSON.
func (p Payment) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		MembershipID int       `json:"membership_id"`
		Status       string    `json:"status"`
		Amount       int       `json:"amount"`
		InsertedAt   time.Time `json:"inserted_at"`
	}{
		MembershipID: p.MembershipID,
		Status:       p.Status,
		Amount:       p.Amount,
		InsertedAt:   p.InsertedAt,
	})
}

// SavePayment ...
func (p Payment) SavePayment(db *gorm.DB) {
	db.Create(p)
}
