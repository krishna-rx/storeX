package models

type RegisterRes struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
type GetAssetsRes struct {
	SerialNumber string `json:"serialNumber" db:"serial_no"`
	Model        string `json:"model" db:"model"`
	Brand        string `json:"brand" db:"brand"`
	AssetType    string `json:"assetType" db:"asset_type"`
}
