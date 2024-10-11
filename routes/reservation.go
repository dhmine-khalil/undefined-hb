// routes/reservation.go
package routes

import (
	"fmt"
	"habitat-server/models"
	"habitat-server/storage"
	"net/http"
	"time"

	"github.com/kataras/iris/v12"
	"gorm.io/gorm"
)

func CreateReservation(ctx iris.Context) {
    var reservation models.Reservation
    if err := ctx.ReadJSON(&reservation); err != nil {
        ctx.StatusCode(http.StatusBadRequest)
        ctx.JSON(iris.Map{"error": "Invalid request payload"})
        return
    }

    userID := ctx.Values().Get("userID").(uint)
    reservation.UserID = userID

    if err := validateReservation(&reservation); err != nil {
        ctx.StatusCode(http.StatusBadRequest)
        ctx.JSON(iris.Map{"error": err.Error()})
        return
    }

    if err := checkAvailability(&reservation); err != nil {
        ctx.StatusCode(http.StatusConflict)
        ctx.JSON(iris.Map{"error": err.Error()})
        return
    }

    reservation.Status = "pending"
    reservation.PaymentStatus = "pending"
    reservation.TotalPrice = calculateTotalPrice(&reservation)

    if err := storage.DB.Create(&reservation).Error; err != nil {
        ctx.StatusCode(http.StatusInternalServerError)
        ctx.JSON(iris.Map{"error": "Failed to create reservation"})
        return
    }

    ctx.StatusCode(http.StatusCreated)
    ctx.JSON(reservation)
}

func GetReservation(ctx iris.Context) {
    id := ctx.Params().GetUintDefault("id", 0)
    if id == 0 {
        ctx.StatusCode(http.StatusBadRequest)
        ctx.JSON(iris.Map{"error": "Invalid reservation ID"})
        return
    }

    var reservation models.Reservation
    if err := storage.DB.Preload("User").Preload("Property").First(&reservation, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            ctx.StatusCode(http.StatusNotFound)
            ctx.JSON(iris.Map{"error": "Reservation not found"})
        } else {
            ctx.StatusCode(http.StatusInternalServerError)
            ctx.JSON(iris.Map{"error": "Failed to retrieve reservation"})
        }
        return
    }

    ctx.JSON(reservation)
}

func GetReservationsByUserID(ctx iris.Context) {
    userID := ctx.Values().Get("userID").(uint)

    var reservations []models.Reservation
    if err := storage.DB.Where("user_id = ?", userID).Preload("Property").Find(&reservations).Error; err != nil {
        ctx.StatusCode(http.StatusInternalServerError)
        ctx.JSON(iris.Map{"error": "Failed to retrieve reservations"})
        return
    }

    ctx.JSON(reservations)
}

func UpdateReservation(ctx iris.Context) {
    id := ctx.Params().GetUintDefault("id", 0)
    if id == 0 {
        ctx.StatusCode(http.StatusBadRequest)
        ctx.JSON(iris.Map{"error": "Invalid reservation ID"})
        return
    }

    var reservation models.Reservation
    if err := storage.DB.First(&reservation, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            ctx.StatusCode(http.StatusNotFound)
            ctx.JSON(iris.Map{"error": "Reservation not found"})
        } else {
            ctx.StatusCode(http.StatusInternalServerError)
            ctx.JSON(iris.Map{"error": "Failed to retrieve reservation"})
        }
        return
    }

    userID := ctx.Values().Get("userID").(uint)
    if reservation.UserID != userID {
        ctx.StatusCode(http.StatusForbidden)
        ctx.JSON(iris.Map{"error": "You don't have permission to update this reservation"})
        return
    }

    var updateData struct {
        StartDate       time.Time `json:"startDate"`
        EndDate         time.Time `json:"endDate"`
        Status          string    `json:"status"`
        PaymentStatus   string    `json:"paymentStatus"`
        GuestCount      int       `json:"guestCount"`
        SpecialRequests string    `json:"specialRequests"`
    }

    if err := ctx.ReadJSON(&updateData); err != nil {
        ctx.StatusCode(http.StatusBadRequest)
        ctx.JSON(iris.Map{"error": "Invalid request payload"})
        return
    }

    if !updateData.StartDate.IsZero() {
        reservation.StartDate = updateData.StartDate
    }
    if !updateData.EndDate.IsZero() {
        reservation.EndDate = updateData.EndDate
    }
    if updateData.Status != "" {
        reservation.Status = updateData.Status
    }
    if updateData.PaymentStatus != "" {
        reservation.PaymentStatus = updateData.PaymentStatus
    }
    if updateData.GuestCount != 0 {
        reservation.GuestCount = updateData.GuestCount
    }
    if updateData.SpecialRequests != "" {
        reservation.SpecialRequests = updateData.SpecialRequests
    }

    if err := validateReservation(&reservation); err != nil {
        ctx.StatusCode(http.StatusBadRequest)
        ctx.JSON(iris.Map{"error": err.Error()})
        return
    }

    if err := checkAvailability(&reservation); err != nil {
        ctx.StatusCode(http.StatusConflict)
        ctx.JSON(iris.Map{"error": err.Error()})
        return
    }

    reservation.TotalPrice = calculateTotalPrice(&reservation)

    if err := storage.DB.Save(&reservation).Error; err != nil {
        ctx.StatusCode(http.StatusInternalServerError)
        ctx.JSON(iris.Map{"error": "Failed to update reservation"})
        return
    }

    ctx.JSON(reservation)
}

func DeleteReservation(ctx iris.Context) {
    id := ctx.Params().GetUintDefault("id", 0)
    if id == 0 {
        ctx.StatusCode(http.StatusBadRequest)
        ctx.JSON(iris.Map{"error": "Invalid reservation ID"})
        return
    }

    var reservation models.Reservation
    if err := storage.DB.First(&reservation, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            ctx.StatusCode(http.StatusNotFound)
            ctx.JSON(iris.Map{"error": "Reservation not found"})
        } else {
            ctx.StatusCode(http.StatusInternalServerError)
            ctx.JSON(iris.Map{"error": "Failed to retrieve reservation"})
        }
        return
    }

    userID := ctx.Values().Get("userID").(uint)
    if reservation.UserID != userID {
        ctx.StatusCode(http.StatusForbidden)
        ctx.JSON(iris.Map{"error": "You don't have permission to delete this reservation"})
        return
    }

    if err := storage.DB.Delete(&reservation).Error; err != nil {
        ctx.StatusCode(http.StatusInternalServerError)
        ctx.JSON(iris.Map{"error": "Failed to delete reservation"})
        return
    }

    ctx.StatusCode(http.StatusNoContent)
}

func validateReservation(reservation *models.Reservation) error {
    if reservation.StartDate.After(reservation.EndDate) {
        return fmt.Errorf("start date must be before end date")
    }
    if reservation.StartDate.Before(time.Now()) {
        return fmt.Errorf("start date must be in the future")
    }
    if reservation.GuestCount <= 0 {
        return fmt.Errorf("guest count must be greater than 0")
    }
    return nil
}

func checkAvailability(reservation *models.Reservation) error {
    var overlappingReservations int64
    storage.DB.Model(&models.Reservation{}).
        Where("id != ? AND property_id = ? AND status != 'cancelled' AND ((start_date <= ? AND end_date >= ?) OR (start_date <= ? AND end_date >= ?) OR (start_date >= ? AND end_date <= ?))",
            reservation.ID, reservation.PropertyID, reservation.StartDate, reservation.StartDate, reservation.EndDate, reservation.EndDate, reservation.StartDate, reservation.EndDate).
        Count(&overlappingReservations)

    if overlappingReservations > 0 {
        return fmt.Errorf("the property is not available for the selected dates")
    }
    return nil
}

func calculateTotalPrice(reservation *models.Reservation) float64 {
    var property models.Property
    storage.DB.First(&property, reservation.PropertyID)
    days := reservation.EndDate.Sub(reservation.StartDate).Hours() / 24
    return float64(property.RentHigh) * float64(days)
}