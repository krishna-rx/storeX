package models

import "time"

type RegisterRes struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
type GetAssetsRes struct {
	SerialNumber string `json:"serialNumber"`
	Model        string `json:"model"`
	Brand        string `json:"brand"`
	AssetType    string `json:"assetType"`
	AssignedTo   string `json:"assignedTo"`
	AssignedAt   string `json:"assignedAt"`
	AssetStatus  string `json:"assetStatus"`
	PurchasedAt  string `json:"purchasedAt"`
}

type GetEmployeeInfoRes struct {
	Name        string   `json:"name"`
	Email       string   `json:"email"`
	PhoneNumber string   `json:"phoneNumber"`
	UserType    UserType `json:"userType"`
	AssetStatus string   `json:"assetStatus"`
	Role        string   `json:"role"`
}

type AssetType string

const (
	Laptop   AssetType = "laptop"
	Mouse    AssetType = "mouse"
	Keyboard AssetType = "keyboard"
	Phone    AssetType = "phone"
)

type UserAssetTimelineRes struct {
	Brand       string     `json:"brand" db:"brand"`
	Model       string     `json:"model" db:"model"`
	AssetTypes  AssetType  `db:"asset_type"`
	AssignedAt  time.Time  `json:"assignedAt" db:"assigned_at"`
	RetrievedAt *time.Time `json:"retrievedAt" db:"retrieved_at"`
}

type UserServiceDetailsRes struct {
	StartDate   time.Time `json:"startDate" db:"start_date"`
	EndDate     time.Time `json:"endDate" db:"end_date"`
	Cost        float64   `json:"cost" db:"cost"`
	Description string    `json:"description" db:"description"`
	ServiceType string    `json:"serviceType" db:"service_type"`
}
