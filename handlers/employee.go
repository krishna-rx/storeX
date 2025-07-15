package handlers

import (
	"main.go/database/dbHelper"
	"main.go/models"
	"main.go/utils"
	"net/http"
)

func UpdateEmployeeInfo(w http.ResponseWriter, r *http.Request) {
	var employeeInfo models.UpdateEmployeeInfoRequest

	err := utils.ParseJSONBody(r, &employeeInfo)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to parse request body")
		return
	}
	if !utils.IsAllowedDomain(employeeInfo.Email) {
		utils.ResponseError(w, http.StatusBadRequest, "email domain is not allowed")
		return
	}
	if !utils.IsValidUserType(string(employeeInfo.UserTypes)) {
		utils.ResponseError(w, http.StatusBadRequest, "userType is not valid provide a valid userType")
		return
	}
	rawUserID := r.Context().Value("user_id")
	userID, ok := rawUserID.(string)
	if !ok {
		utils.ResponseError(w, http.StatusBadRequest, "error getting the user id from context")
		return
	}
	employeeInfo.UpdatedBy = userID
	err = dbHelper.UpdateEmployeeRecord(employeeInfo)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "error updating employee info")
		return
	}
	err = utils.WriteJSONResponse(w, "data updated in the employee table successfully")
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to write response")
		return
	}
}
func GetEmployeeInfo(w http.ResponseWriter, r *http.Request) {
	phoneNumber := r.URL.Query().Get("phoneNumber")
	name := r.URL.Query().Get("name")
	email := r.URL.Query().Get("email")
	userType := r.URL.Query().Get("userType")
	employeeInfo, err := dbHelper.LookUpEmployeeInfo(email, name, phoneNumber, userType)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "error looking up employee info")
		return
	}
	if len(employeeInfo) == 0 {
		utils.ResponseError(w, http.StatusNotFound, "employee info not found")
		return
	}
	err = utils.WriteJSONResponse(w, employeeInfo)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to write response")
		return
	}

}
func UserAssetTimeline(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		utils.ResponseError(w, http.StatusBadRequest, "userID is required")
		return
	}
	userTimeline, err := dbHelper.UserAssetInfo(userID)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "error getting user asset info")
		return
	}
	if len(userTimeline) == 0 {
		utils.ResponseError(w, http.StatusNotFound, "user asset timeline not found")
		return
	}
	err = utils.WriteJSONResponse(w, userTimeline)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to write response")
		return
	}
}
func DeleteEmployeeInfo(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		utils.ResponseError(w, http.StatusBadRequest, "userID is required")
		return
	}
	count, err := dbHelper.GetAssetIfExists(userID)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "error getting user asset info")
		return
	}
	if count == 0 {
		err = dbHelper.DeleteUserInfo(userID)
		if err != nil {
			utils.ResponseError(w, http.StatusBadRequest, "error deleting user info")
			return
		}
	} else {
		err = dbHelper.RetrieveAssetByUserID(userID)
		if err != nil {
			utils.ResponseError(w, http.StatusBadRequest, "error getting user asset info by the id")
			return
		}
		err = dbHelper.DeleteUserInfo(userID)
		if err != nil {
			utils.ResponseError(w, http.StatusBadRequest, "error deleting user info by the id")
			return
		}
	}
	err = utils.WriteJSONResponse(w, "data deleted in the employee table successfully")
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to write response")
		return
	}
}
