// package models

// import (
// 	"gorm.io/datatypes"
// 	"gorm.io/gorm"
// )

// type User struct {
// 	gorm.Model
// 	FirstName           string         `json:"firstName"`
// 	LastName            string         `json:"lastName"`
// 	Email               string         `json:"email"`
// 	Password            string         `json:"password"`
// 	SocialLogin         bool           `json:"socialLogin"`
// 	SocialProvider      string         `json:"socialProvider"`
// 	Properties          []Property     `json:"properties"`
// 	SavedProperties     datatypes.JSON `json:"savedProperties"`
// 	PushTokens          datatypes.JSON `json:"pushTokens"`
// 	AllowsNotifications *bool          `json:"allowsNotifications"`
// 	isVerified *bool `json:"isVerified"`
// }

package models

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type MembershipTier string

const (
	FreeTier    MembershipTier = "Free"
	PremiumTier MembershipTier = "Premium"
	ProTier     MembershipTier = "Pro"
)
type User struct {
    gorm.Model
    FirstName           string         `json:"firstName"`
    LastName            string         `json:"lastName"`
    Email               string         `json:"email"`
    Password            string         `json:"password"`
    SocialLogin         bool           `json:"socialLogin"`
    SocialProvider      string         `json:"socialProvider"`
    Properties          []Property     `json:"properties"`
    SavedProperties     datatypes.JSON `json:"savedProperties"`
    PushTokens          datatypes.JSON `json:"pushTokens"`
    AllowsNotifications *bool          `json:"allowsNotifications"`
    IsVerified          *bool          `json:"isVerified"`
    MembershipTier      MembershipTier `json:"membershipTier" gorm:"type:membership_tier;default:'Free'"`
}