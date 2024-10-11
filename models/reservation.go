// models/reservation.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type Reservation struct {
    gorm.Model
    UserID      uint      `json:"userId"`
    User        User      `json:"user"`
    PropertyID  uint      `json:"propertyId"`
    Property    Property  `json:"property"`
    StartDate   time.Time `json:"startDate"`
    EndDate     time.Time `json:"endDate"`
    TotalPrice  float64   `json:"totalPrice"`
    Status      string    `json:"status"` // "pending", "confirmed", "cancelled"
    PaymentStatus string  `json:"paymentStatus"` // "pending", "paid", "refunded"
    GuestCount  int       `json:"guestCount"`
    SpecialRequests string `json:"specialRequests"`
}