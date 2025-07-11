package dbHelper

import (
	"github.com/jmoiron/sqlx"
	"main.go/database"
	"main.go/models"
)

func IsUserExist(email string) (bool, error) {
	SQL := `SELECT COUNT(id) > 0 as is_exist FROM users WHERE email = TRIM($1) AND archived_at IS NULL;`
	var check bool
	chkErr := database.DB.Get(&check, SQL, email)
	return check, chkErr
}
func LookUpUserRole(loginReq *models.UserLoginRequest) error {
	query := `
           SELECT ur.role, u.id
           FROM users u
           JOIN user_roles ur ON ur.user_id = u.id
           WHERE u.email = $1
          `
	err := database.DB.Get(loginReq, query, loginReq.Email)
	return err
}
func CreateUser(db *sqlx.Tx, registerReq models.UserRegisterRequest) (string, error) {
	// language=sql
	SQL := `INSERT INTO users (name, email , user_type) VALUES ($1, $2, $3) returning id`
	args := []interface{}{
		registerReq.Name,
		registerReq.Email,
		registerReq.UserType,
	}
	var id string
	err := db.Get(&id, SQL, args...)
	return id, err
}
func CreateUserRole(db *sqlx.Tx, userID, role string) error {
	// language=sql
	SQL := `INSERT INTO user_roles (user_id, role) VALUES ($1, $2)`
	args := []interface{}{
		userID,
		role,
	}
	_, err := db.Exec(SQL, args...)
	return err
}
func UpdateRoleByAdmin(updateReq models.UpdateUserDetails) error {
	SQL := `UPDATE user_roles
	         SET role = $1
	         FROM users
	         WHERE user_roles.user_id = users.id AND users.email = $2`
	_, err := database.DB.Exec(SQL, updateReq.Role, updateReq.Email)
	return err
}
func InsertIntoLaptop(db *sqlx.Tx, assetID string, userID string, laptop models.LaptopDetailsRequest) error {
	SQL := `INSERT INTO laptops (assets_id, name, IMEI_no_1, IMEI_no_2 ,price, ram, os,created_by)
		VALUES ($1, $2, $3, $4, $5, $6,$7,$8)`
	args := []interface{}{
		assetID,
		laptop.LaptopName,
		laptop.IMEI1,
		laptop.IMEI2,
		laptop.Price,
		laptop.RAM,
		laptop.OSVersion,
		userID,
	}
	_, err := db.Exec(SQL, args...)
	return err
}
func CreateAsset(db *sqlx.Tx, assetReq models.CreateAssetRequest) (string, error) {
	SQL := `INSERT INTO assets (serial_no,model,brand,asset_type,purchased_at) VALUES ($1, $2, $3, $4, $5) returning id`
	args := []interface{}{
		assetReq.SerialNumber,
		assetReq.ModelType,
		assetReq.Brand,
		assetReq.AssetType,
		assetReq.PurchaseDate,
	}
	var id string
	err := db.Get(&id, SQL, args...)
	return id, err
}
func AssestInsert(userID string, assetAssignReq models.AssignAssetRequest) error {
	SQL := `INSERT INTO asset_assignments (asset_id , assigned_to , assigned_by ) VALUES ($1, $2, $3)`
	args := []interface{}{
		assetAssignReq.AssetID,
		assetAssignReq.AssignTo,
		userID,
	}
	_, err := database.DB.Exec(SQL, args...)
	return err
}
func GetAssetByRole(brand string, model string, assetType string) ([]models.GetAssetsRes, error) {
	SQL := `SELECT 
			assets.brand,
			assets.model,
			assets.asset_type,
			assets.serial_no
		FROM assets
		JOIN asset_assignments ON asset_assignments.asset_id = assets.id
		WHERE !utils.validText($1) OR ILIKE $1 OR !utils.validText($2)  `
	var getAssets []models.GetAssetsRes
	err := database.DB.Select(&getAssets, SQL, brand, model, assetType)
	return getAssets, err
}
