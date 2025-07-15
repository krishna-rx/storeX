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
func InsertIntoMouse(db *sqlx.Tx, assetID string, userID string, mouseDetails models.MouseDetailsRequest) error {
	SQL := `INSERT INTO mouses (assets_id, name, IMEI_no_1, IMEI_no_2 ,price, buttons,created_by)
		VALUES ($1, $2, $3, $4, $5, $6,$7)`
	args := []interface{}{
		assetID,
		mouseDetails.MouseName,
		mouseDetails.IMEI1,
		mouseDetails.IMEI2,
		mouseDetails.Price,
		mouseDetails.Buttons,
		userID,
	}
	_, err := db.Exec(SQL, args...)
	return err
}
func InsertIntoMobile(db *sqlx.Tx, assetID string, userID string, mobile models.MobileDetailsRequest) error {
	SQL := `INSERT INTO laptops (assets_id, name, IMEI_no_1, IMEI_no_2 ,price, ram, os,created_by)
		VALUES ($1, $2, $3, $4, $5, $6,$7,$8)`
	args := []interface{}{
		assetID,
		mobile.MobileName,
		mobile.IMEI1,
		mobile.IMEI2,
		mobile.Price,
		mobile.RAM,
		mobile.OSVersion,
		userID,
	}
	_, err := db.Exec(SQL, args...)
	return err
}
func InsertIntoKeyboard(db *sqlx.Tx, assetID string, userID string, keyboard models.KeyboardDetailsRequest) error {
	SQL := `INSERT INTO laptops (assets_id, name, IMEI_no_1, IMEI_no_2 ,price, type,keys,is_wired,created_by)
		VALUES ($1, $2, $3, $4, $5, $6,$7,$8,$9)`
	args := []interface{}{
		assetID,
		keyboard.KeyboardName,
		keyboard.IMEI1,
		keyboard.IMEI2,
		keyboard.Price,
		keyboard.Type,
		keyboard.Keys,
		keyboard.IsWired,
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
func AssetAssignInsert(userID string, assetAssignReq models.AssignAssetRequest) error {
	SQL := `INSERT INTO asset_assignments (asset_id , assigned_to , assigned_by ) VALUES ($1, $2, $3)`
	args := []interface{}{
		assetAssignReq.AssetID,
		assetAssignReq.AssignTo,
		userID,
	}
	_, err := database.DB.Exec(SQL, args...)
	return err
}
func UpdateAssignStatus(assetID string) error {
	SQL := `UPDATE assets SET asset_status = 'assigned'
          WHERE id = $1`
	_, err := database.DB.Exec(SQL, assetID)
	return err
}
func LookUpAsset(brand string, model string, assetType string, limit int, page int) ([]models.GetAssetsRes, error) {
	SQL := `SELECT brand, model, asset_type, serial_no,purchased_at,asset_status,users.name,asset_assignments.assigned_at 
            FROM assets
            JOIN asset_assignments ON asset_assignments.asset_id = assets.id
            JOIN users ON asset_assignments.assigned_by = users.id
            WHERE 
            ($1 = '' OR brand ILIKE $1 || '%') AND
            ($2 = '' OR model ILIKE $2 || '%') AND
            ($3 = '' OR asset_type::text ILIKE $3 || '%') -- as ILIKE only work on text on integer so to make it work on enum you can type case it to text
            ORDER BY purchased_at ASC
            LIMIT $4 OFFSET $5
            `
	getAssets := make([]models.GetAssetsRes, 0)
	err := database.DB.Select(&getAssets, SQL, brand, model, assetType, limit, page)
	return getAssets, err
}
func UpdateEmployeeRecord(updateEmployeeRecord models.UpdateEmployeeInfoRequest) error {
	SQL := `UPDATE users SET
            name = $1,
            email = $2,
            phone_no = $3,
            updated_by = $4,
            user_type = $5::user_types
            WHERE id = $6`

	args := []interface{}{
		updateEmployeeRecord.Name,
		updateEmployeeRecord.Email,
		updateEmployeeRecord.PhoneNumber,
		updateEmployeeRecord.UpdatedBy,
		updateEmployeeRecord.UserTypes,
		updateEmployeeRecord.ID,
	}
	_, err := database.DB.Exec(SQL, args...)
	return err
}
func LookUpEmployeeInfo(email string, name string, phoneNumber string, userType string) ([]models.GetEmployeeInfoRes, error) {
	SQL := `SELECT name , email , phone_no , user_type ,user_roles.role,assets.asset_status
            FROM users
            JOIN user_roles ON users.id = user_roles.user_id
            JOIN asset_assignments ON users.id = asset_assignments.assigned_by 
            JOIN assets ON assets.id = asset_assignments.asset_id
            WHERE 
                ($1 = '' OR name ILIKE $1 || '%') AND
                ($2 = '' OR email ILIKE $2 || '%') AND
                ($3 = '' OR phone_no ILIKE $3 || '%') AND
                ($4 = '' OR user_type::text ILIKE $4 || '%')
                ORDER BY users.created_at ASC
                `
	employeeInfo := make([]models.GetEmployeeInfoRes, 0)
	err := database.DB.Select(&employeeInfo, SQL, name, email, phoneNumber, userType)
	return employeeInfo, err
}
func IsAssetAssignedToUser(retrieveAsset models.RetrieveAssetRequest) (bool, error) {
	SQL := `SELECT 1
	FROM assets
	WHERE id = $1 AND asset_type = $2
	LIMIT 1`
	var check bool
	args := []interface{}{
		retrieveAsset.AssetID,
		retrieveAsset.AssetType,
	}
	checkErr := database.DB.Get(&check, SQL, args...)
	return check, checkErr
}
func RetrieveAsset(retrieveAssetReq models.RetrieveAssetRequest) error {
	SQL := `UPDATE asset_assignments
	       SET retrieved_at = NOW()
	       WHERE 
		   asset_assignments.asset_id = $1 ;
            `
	args := []interface{}{
		retrieveAssetReq.AssetID,
	}
	_, err := database.DB.Exec(SQL, args...)
	return err
}
func UserAssetInfo(userID string) ([]models.UserAssetTimelineRes, error) {
	SQL := `SELECT 
             asset_assignments.assigned_at, 
             asset_assignments.retrieved_at,
             assets.brand,
             assets.model,
             assets.asset_type
             FROM asset_assignments 
             JOIN assets ON assets.id = asset_assignments.asset_id
             WHERE asset_assignments.assigned_by = $1
               AND assets.asset_status::text ILIKE 'assigned'
             ORDER BY asset_assignments.assigned_at DESC;`
	timeline := make([]models.UserAssetTimelineRes, 0)
	err := database.DB.Select(&timeline, SQL, userID)
	return timeline, err
}
func UpdateAsset(updateAsset models.UpdateUserAssetRequest, userID string) error {
	SQL := `
    UPDATE assets SET
            serial_no = COALESCE(NULLIF($1, ''), assets.serial_no),
            model = COALESCE(NULLIF($2, ''), assets.model),
            brand = COALESCE(NULLIF($3, ''), assets.brand),
            asset_status = COALESCE(NULLIF($4, ''), assets.asset_status::text)::asset_statuses
            updated_by = $5
    WHERE id = $5;`
	args := []interface{}{
		updateAsset.SerialNumber,
		updateAsset.ModelType,
		updateAsset.Brand,
		updateAsset.AssetStatus,
		updateAsset.ID,
		userID,
	}
	_, err := database.DB.Exec(SQL, args...)
	return err
}
func GetAssetIfExists(userID string) (int, error) {
	var count int
	err := database.DB.Get(&count, `
		SELECT COUNT(*) FROM asset_assignments
		WHERE assigned_to = $1 AND retrieved_at IS NULL
	`, userID)
	return count, err
}
func DeleteUserInfo(userID string) error {
	SQL := `UPDATE users SET archived_at = NOW() WHERE id = $1 AND archived_at IS NULL;`
	_, err := database.DB.Exec(SQL, userID)
	return err
}
func RetrieveAssetByUserID(userID string) error {
	SQL := `
	     	UPDATE asset_assignments
			SET retrieved_at = NOW()
			WHERE assigned_to = $1 AND retrieved_at IS NULL
		`
	_, err := database.DB.Exec(SQL, userID)
	return err
}
func GetCurrentStatus(AssetID string) (models.AssetStatus, error) {
	SQL := `SELECT asset_status FROM assets 
            WHERE id = $1;`
	var status models.AssetStatus
	err := database.DB.Get(&status, SQL, AssetID)
	return status, err
}
func CreateServiceRecord(serviceDetails models.ServiceDetailsRequest) error {
	SQL := `INSERT INTO services (
            asset_id,
            cost,
            service_type,
            description
            ) VALUES ($1, $2, $3, $4);
            `
	args := []interface{}{
		serviceDetails.AssetID,
		serviceDetails.Cost,
		serviceDetails.ServiceType,
		serviceDetails.Description,
	}
	_, err := database.DB.Exec(SQL, args...)
	return err
}
func UpdateAssetStatus(status models.AssetStatus, AssetID string) error {
	SQL := `UPDATE assets SET 
            asset_status = $1
            WHERE id = $2;`
	_, err := database.DB.Exec(SQL, status, AssetID)
	return err
}
func GetAssetDetails(AssetID string) ([]models.UserServiceDetailsRes, error) {
	SQL := `
            SELECT 
              asset_assignments.assigned_at,
              asset_assignments.retrieved_at,
              assets.brand,
              assets.model,
              assets.asset_type
            FROM assets 
            JOIN asset_assignments  ON assets.id = asset_assignments.asset_id
            WHERE assets.id = $1;
            ORDER BY asset_assignments.assigned_at DESC;`
	var assetDetails []models.UserServiceDetailsRes
	err := database.DB.Select(&assetDetails, SQL, AssetID)
	return assetDetails, err
}
