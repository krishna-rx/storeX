package models

import (
	"github.com/aarondl/null/v9"
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
	ID          string `json:"id"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Role        string `json:"role"`
	UserType    string `json:"userType"`
	PhoneNumber string `json:"phoneNo"`
	UpdatedBy   string `json:"updatedBy"`
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
type MouseDetailsRequest struct {
	AssetID   string  `json:"assetId" db:"asset_id"`
	MouseName string  `json:"mouseName" db:"name"`
	IMEI1     string  `json:"imei1" db:"imei_1"`
	IMEI2     string  `json:"imei2" db:"imei_2"`
	Buttons   int     `json:"buttons" db:"buttons"`
	Price     float64 `json:"price" db:"price"`
}
type KeyboardDetailsRequest struct {
	AssetID      string  `json:"assetId" db:"asset_id"`
	KeyboardName string  `json:"keyboardName" db:"name"`
	IMEI1        string  `json:"imei1" db:"imei_no_1"`
	IMEI2        string  `json:"imei2" db:"imei_no_2"`
	Type         string  `json:"type" db:"type"`
	Price        float64 `json:"price" db:"price"`
	IsWired      bool    `json:"is_wired" db:"is_wired"`
	Keys         string  `json:"keys" db:"keys"`
}
type MobileDetailsRequest struct {
	AssetID    string  `json:"assetId" db:"asset_id"`
	MobileName string  `json:"mobileName" db:"name"`
	IMEI1      string  `json:"imei1" db:"imei_no_1"`
	IMEI2      string  `json:"imei2" db:"imei_no_2"`
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
type UserType string

const (
	Freelancer UserType = "freelancer"
	Intern     UserType = "intern"
	FullTime   UserType = "full_time"
)

type UpdateEmployeeInfoRequest struct {
	ID          *string     `json:"id"`
	Name        *string     `json:"name"`
	Email       string      `json:"email"`
	PhoneNumber null.String `json:"phoneNumber"`
	UpdatedBy   string      `json:"updatedBy"`
	UserTypes   UserType    `json:"userType"`
}

type RetrieveAssetRequest struct {
	Email     string `json:"email" db:"email"`
	AssetID   string `json:"assetId" db:"asset_id"`
	AssetType string `json:"assetType" db:"asset_type"`
}

type AssetStatuses string

const (
	Available        AssetStatuses = "available"
	WaitingForRepair AssetStatuses = "waiting_for_repair"
	Service          AssetStatuses = "service"
	damaged          AssetStatuses = "damaged"
)

type UpdateUserAssetRequest struct {
	ID           *string       `json:"id"`
	SerialNumber string        `json:"serialNumber" `
	ModelType    string        `json:"modelType" `
	Brand        string        `json:"brand"`
	AssetStatus  AssetStatuses `json:"asset_status"`
}

type Role string

const (
	Admin           Role = "admin"
	AssetManager    Role = "asset_manager"
	EmployeeManager Role = "employee_manager"
	Employee        Role = "employee"
)

type AssetStatus string

const (
	StatusAvailable        AssetStatus = "available"
	StatusWaitingForRepair AssetStatus = "waiting_for_repair"
	StatusAssigned         AssetStatus = "assigned"
	StatusDamaged          AssetStatus = "damaged"
	StatusService          AssetStatus = "service"
)

type ServiceDetailsRequest struct {
	AssetID     string    `json:"assetID" db:"asset_id"`
	StartDate   time.Time `json:"startDate" db:"start_date"`
	EndDate     time.Time `json:"endDate" db:"end_date"`
	Cost        int       `json:"cost" db:"cost"`
	Description *string   `json:"description" db:"description"`
	ServiceType string    `json:"serviceType" db:"service_type"`
}
