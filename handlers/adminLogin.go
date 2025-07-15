package handlers

import (
	"fmt"
	"main.go/database/dbHelper"
	"main.go/models"
	"main.go/utils"
	"net/http"
	"strings"
)

func AdminLogin(w http.ResponseWriter, r *http.Request) {
	loginWithRole(w, r, "admin")
}
func AssetManagerLogin(w http.ResponseWriter, r *http.Request) {
	loginWithRole(w, r, "asset_manager")
}
func EmployeeManagerLogin(w http.ResponseWriter, r *http.Request) {
	loginWithRole(w, r, "employee_manager")
}
func loginWithRole(w http.ResponseWriter, r *http.Request, role string) {
	var loginReq models.UserLoginRequest
	err := utils.ParseJSONBody(r, &loginReq)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to parsed request query")
		return
	}
	loginReq.Email = strings.TrimSpace(loginReq.Email)
	exist, existErr := dbHelper.IsUserExist(loginReq.Email)
	if existErr != nil {
		utils.ResponseError(w, http.StatusBadRequest, "error checking user exist")
	}
	if !exist {
		utils.ResponseError(w, http.StatusBadRequest, "user does not exist")
		return
	}
	err = dbHelper.LookUpUserRole(&loginReq)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to look up user role")
		fmt.Println(err)
		return
	}
	if loginReq.Role != role {
		utils.ResponseError(w, http.StatusBadRequest, "user doesnt have the correct role")
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
		return
	}
}
func UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	var updateReq models.UpdateUserDetails
	err := utils.ParseJSONBody(r, &updateReq)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to parsed request query")
		return
	}
	rawRole := r.Context().Value("role")
	role, rawOk := rawRole.(string)
	if !rawOk {
		utils.ResponseError(w, http.StatusBadRequest, "role is not a string")
		return
	}
	if role != "admin" {
		utils.ResponseError(w, http.StatusBadRequest, "current user dont have permission to update role")
		return
	}
	err = dbHelper.UpdateRoleByAdmin(updateReq)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to update the role")
		return
	}
	err = utils.WriteJSONResponse(w, "user updated successfully by admin")
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to encode the json response")
		return
	}
}
