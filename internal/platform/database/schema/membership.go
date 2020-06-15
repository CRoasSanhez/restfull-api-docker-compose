package schema

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/jinzhu/gorm"
)

// Membership ...
type Membership struct {
	gorm.Model
	ID         int    `gorm:"primary_key"`
	UserID     int    `gorm:"type:int(64)" json:"user_id"`
	CardNumber string `gorm:"type:varchar(16)" json:"card_number"`
	ExpDate    string `json:"exp_date"`
	Owner      string `gorm:"type:varchar(100)" json:"owner"`
	CVV        string
	Tier       string    `gorm:"type:varchar(10)"`
	Pricing    int       `gorm:"type:int(8)" default:"100000"`
	Attempts   int       `gorm:"type:int(8)" json:"attempts" default:"0"`
	Status     string    `gorm:"type:varchar(100)" default:"payment_pending"`
	IsBlocked  bool      `gorm:"type:TINYINT(1)"`
	BlockedAt  time.Time `json:"-"`
	InsertedAt time.Time ` json:"-"`
	UpdatedAt  time.Time ` json:"-"`
	IsDeleted  bool      `gorm:"type:TINYINT(1)" json:"is_deleted" default:"false"`
}

// MembershipTier ...
// TODO: Create membership tier table
type MembershipTier struct {
}

// NewMembership ...
func NewMembership(userID, pricing int, cardNumber, owner string) *Membership {
	return &Membership{
		UserID:     userID,
		Pricing:    pricing,
		Owner:      owner,
		CardNumber: cardNumber,
		IsBlocked:  false,
		IsDeleted:  false,
		Status:     "active",
		Attempts:   0,
		Tier:       "normal",
	}
}

// MarshalJSON return bot struct as JSON.
func (m Membership) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		CardNumber string    `json:"card_number"`
		ExpDate    string    `json:"exp_date"`
		Owner      string    `json:"owner"`
		Tier       string    `json:"tier"`
		Pricing    int       `json:"pricing"`
		IsBlocked  bool      `json:"is_blocked"`
		BlockedAt  time.Time `json:"blocked_at"`
	}{
		CardNumber: m.CardNumber,
		ExpDate:    m.ExpDate,
		Owner:      m.Owner,
		Tier:       m.Tier,
		IsBlocked:  m.IsBlocked,
		BlockedAt:  m.BlockedAt,
	})
}

// SaveAttempt ...
func (m Membership) SaveAttempt(db *gorm.DB, qty int) int {

	var response int

	if m.Attempts >= 3 {
		db.Model(&m).Update("Attempts")

		u := &User{
			ID:        m.UserID,
			IsBlocked: true,
		}
		db.Model(u).Update("ID", "IsBlocked")
		response = http.StatusTooManyRequests
	} else {
		response = http.StatusBadRequest
	}

	p := &Payment{
		MembershipID: m.ID,
		UserID:       m.UserID,
		Status:       "unsuccess",
		Amount:       qty,
		InsertedAt:   time.Now(),
	}

	go p.SavePayment(db)

	return response

}
