package handlers

import (
	"github.com/jmoiron/sqlx"
	"main.go/database"
	"main.go/database/dbHelper"
	"main.go/models"
	"main.go/utils"
	"net/http"
	"strings"
)

func UserLogin(w http.ResponseWriter, r *http.Request) {
	var loginReq models.UserLoginRequest
	var registerReq models.UserRegisterRequest
	err := utils.ParseJSONBody(r, &loginReq)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to parsed request query")
		return
	}
	loginReq.Email = strings.TrimSpace(loginReq.Email)
	checkMail := strings.Split(loginReq.Email, "@")
	if len(checkMail) != 2 && strings.EqualFold(checkMail[1], "remotestate.com") {
		utils.ResponseError(w, http.StatusBadRequest, "email is not valid")
		return
	}
	exist, existErr := dbHelper.IsUserExist(loginReq.Email)
	if existErr != nil {
		utils.ResponseError(w, http.StatusBadRequest, "getting error in lookup user")
		return
	}
	if exist {
		err = dbHelper.LookUpUserRole(&loginReq)
		if err != nil {
			utils.ResponseError(w, http.StatusInternalServerError, "failed to look up user role")
			return
		}
		tokenString, err := utils.GenerateJWT(loginReq)
		if err != nil {
			utils.ResponseError(w, http.StatusBadRequest, "failed to generate JWT token")
			return
		}
		err = utils.WriteJSONResponse(w, tokenString)
		if err != nil {
			utils.ResponseError(w, http.StatusBadRequest, "failed to encode the json response")
		}
	} else {
		registerReq.Name = strings.ReplaceAll(checkMail[0], ".", " ")
		registerReq.Email = loginReq.Email
		registerReq.Role = "employee"
		registerReq.UserType = "full_time"

		txErr := database.Tx(func(tx *sqlx.Tx) error {
			// create user
			userID, err := dbHelper.CreateUser(tx, registerReq)
			if err != nil {
				return err
			}

			// create  user role
			err = dbHelper.CreateUserRole(tx, userID, registerReq.Role)
			if err != nil {
				return err
			}

			return nil
		})
		if txErr != nil {
			utils.ResponseError(w, http.StatusInternalServerError, "failed to create user")
			return
		}
		err = utils.WriteJSONResponse(w, "user created successfully")
		if err != nil {
			utils.ResponseError(w, http.StatusBadRequest, "failed to encode the json response")
			return
		}
	}
}
func UserRegister(w http.ResponseWriter, r *http.Request) {
	var registerReq models.UserRegisterRequest
	err := utils.ParseJSONBody(r, &registerReq)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to parsed request query")
		return
	}
	registerReq.Role = "employee"
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		// create user
		userID, err := dbHelper.CreateUser(tx, registerReq)
		if err != nil {
			return err
		}
		// create  user role
		err = dbHelper.CreateUserRole(tx, userID, registerReq.Role)
		if err != nil {
			return err
		}
		return nil
	})
	if txErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to create user")
		return
	}
	err = utils.WriteJSONResponse(w, "user created successfully")
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to encode the json response")
	}
}
