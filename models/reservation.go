// models/reservation.go
package models

import (
	"time"

	"gorm.io/gorm"
)

type Reservation struct {
    gorm.Model
    UserID     uint      `json:"userId"`
    User       User      `json:"user" gorm:"foreignKey:UserID"`
    PropertyID uint      `json:"propertyId"`
    Property   Property  `json:"property" gorm:"foreignKey:PropertyID"`
    StartDate  time.Time `json:"startDate"`
    EndDate    time.Time `json:"endDate"`
    Status     string    `json:"status"` // "pending", "confirmed", "cancelled"
    TotalPrice float64   `json:"totalPrice"`
}