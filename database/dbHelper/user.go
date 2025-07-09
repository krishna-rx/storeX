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
func LookUpUserRole(RegisterReq *models.UserRegisterRequest) error {
	query := `
           SELECT ur.role
           FROM users u
           JOIN user_roles ur ON ur.user_id = u.id
           WHERE u.email = $1
          `
	err := database.DB.Get(&RegisterReq.Role, query, RegisterReq.Email)
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
