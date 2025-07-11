package handlers

import (
	"main.go/database/dbHelper"
	"main.go/models"
	"main.go/utils"
	"net/http"
)

func AssignAsset(w http.ResponseWriter, r *http.Request) {
	var assignAsset models.AssignAssetRequest
	rawUserID := r.Context().Value("user_id")
	userID, ok := rawUserID.(string)
	if !ok {
		utils.ResponseError(w, http.StatusUnauthorized, "failed to get the user id from context")
		return
	}
	err := utils.ParseJSONBody(r, &assignAsset)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to parse request body")
		return
	}
	if assignAsset.AssignTo == userID {
		utils.ResponseError(w, http.StatusUnauthorized, "you can't assign asset to yourself")
	}
	err = dbHelper.AssestInsert(userID, assignAsset)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to insert into asset assignment table")
		return
	}
	err = utils.WriteJSONResponse(w, "assigned asset successfully")
}
