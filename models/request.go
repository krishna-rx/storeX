package models

import (
	jsoniter "github.com/json-iterator/go"
	"time"
)

type UserLoginRequest struct {
	ID    string `json:"id"   db:"id"`
	Email string `json:"email" db:"email"`
	Role  string `json:"role" db:"role"`
}
type UserRegisterRequest struct {
	ID       string `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	Email    string `json:"email" db:"email"`
	Role     string `json:"role" db:"role"`
	UserType string `json:"userType" db:"user_type"`
}
type UpdateUserDetails struct {
	ID          string `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Email       string `json:"email" db:"email"`
	Role        string `json:"role" db:"role"`
	UserType    string `json:"userType" db:"user_type"`
	PhoneNumber string `json:"phoneNo" db:"phone_number"`
	UpdatedBy   string `json:"updatedBy" db:"updated_by"`
}
type CreateAssetRequest struct {
	ID           string              `json:"id" db:"id"`
	SerialNumber string              `json:"serialNumber" db:"serial_no"`
	ModelType    string              `json:"modelType" db:"model"`
	Brand        string              `json:"brand" db:"brand"`
	AssetType    string              `json:"assetType" db:"asset_type"`
	PurchaseDate time.Time           `json:"purchaseDate" db:"purchase_at"`
	Config       jsoniter.RawMessage `json:"config"`
}
type LaptopDetailsRequest struct {
	AssetID    string  `json:"assetId" db:"asset_id"`
	LaptopName string  `json:"laptopName" db:"laptop_name"`
	IMEI1      string  `json:"IMEI1" db:"imei_no_2"`
	IMEI2      string  `json:"IMEI2" db:"imei_no_2"`
	Price      float64 `json:"price" db:"price"`
	RAM        string  `json:"ram" db:"ram"`
	OSVersion  string  `json:"osVersion" db:"os"`
}
type AssignAssetRequest struct {
	AssetID     string    `json:"assetId" db:"asset_id"`
	AssignTo    string    `json:"assignTo" db:"assigned_to"`
	AssignBy    string    `json:"assignBy" db:"assigned_by"`
	RetrievedAt time.Time `json:"retrievedAt" db:"retrieved_at"`
	Description string    `json:"description" db:"description"`
}
