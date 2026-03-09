package models

import (
	"time"
)

type User struct {
	UserID        string    `json:"user_id"                  dynamodbav:"user_id"`
	Email         string    `json:"email"                    dynamodbav:"email"`
	Phone         string    `json:"phone,omitempty"          dynamodbav:"phone,omitempty"`
	FirstName     string    `json:"first_name"               dynamodbav:"first_name"`
	MiddleInitial string    `json:"middle_initial,omitempty" dynamodbav:"middle_initial,omitempty"`
	LastName      string    `json:"last_name"                dynamodbav:"last_name"`
	AvatarURL     string    `json:"avatar_url,omitempty"     dynamodbav:"avatar_url,omitempty"`
	CreatedAt     time.Time `json:"created_at"               dynamodbav:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"               dynamodbav:"updated_at"`
}

type Address struct {
	AddressID string `json:"address_id"      dynamodbav:"address_id"`
	Alias     string `json:"alias,omitempty" dynamodbav:"alias,omitempty"`
	Line1     string `json:"line1"           dynamodbav:"line1"`
	Line2     string `json:"line2,omitempty" dynamodbav:"line2,omitempty"`
	City      string `json:"city"            dynamodbav:"city"`
	State     string `json:"state"           dynamodbav:"state"`
	ZipCode   string `json:"zip_code"        dynamodbav:"zip_code"`
	Country   string `json:"country"         dynamodbav:"country"`
	IsDefault bool   `json:"is_default"      dynamodbav:"is_default"`
}

// PaymentMethod stores a tokenized reference to a payment instrument.
// Token is write-only at the DB layer and never serialized in API responses.
type PaymentMethod struct {
	PaymentID   string `json:"payment_id"   dynamodbav:"payment_id"`
	Provider    string `json:"provider"     dynamodbav:"provider"`
	Token       string `json:"-"            dynamodbav:"token"`
	Last4       string `json:"last4"        dynamodbav:"last4"`
	Brand       string `json:"brand"        dynamodbav:"brand"`
	ExpiryMonth int    `json:"expiry_month" dynamodbav:"expiry_month"`
	ExpiryYear  int    `json:"expiry_year"  dynamodbav:"expiry_year"`
	IsDefault   bool   `json:"is_default"   dynamodbav:"is_default"`
}

type Preferences struct {
	Language      string            `json:"language"                  dynamodbav:"language"`
	Timezone      string            `json:"timezone"                  dynamodbav:"timezone"`
	Communication map[string]bool   `json:"communication"             dynamodbav:"communication"`
	Marketing     map[string]string `json:"marketing"                 dynamodbav:"marketing"`
}
