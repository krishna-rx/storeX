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

var assetHandlers = map[string]AssetHandler{
	"laptop": LaptopHandler{},
}

func CreateAsset(w http.ResponseWriter, r *http.Request) {
	var assetReq models.CreateAssetRequest
	rawUserID := r.Context().Value("user_id")
	userID, ok := rawUserID.(string)
	if !ok {
		utils.ResponseError(w, http.StatusUnauthorized, "failed to get the user id from context")
		return
	}
	err := utils.ParseJSONBody(r, &assetReq)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to parse the json body ")
		return
	}
	txErr := database.Tx(func(tx *sqlx.Tx) error {
		assetID, err := dbHelper.CreateAsset(tx, assetReq)
		assetReq.ID = assetID
		if err != nil {
			return err
		}
		handler, ok := assetHandlers[strings.ToLower(assetReq.AssetType)]
		if !ok {
			utils.ResponseError(w, http.StatusBadRequest, "unknown asset type")
			return err
		}
		config, err := handler.UnmarshalConfig(assetReq.Config)
		if err != nil {
			utils.ResponseError(w, http.StatusBadRequest, "failed to parse the asset types ")
			return err
		}
		if err := handler.Validate(config); err != nil {
			utils.ResponseError(w, http.StatusBadRequest, "invalid json for the asset type")
			return err
		}
		if err := handler.Insert(tx, assetReq.ID, userID, config); err != nil {
			utils.ResponseError(w, http.StatusInternalServerError, "failed to insert the asset type")
			return err
		}
		return nil
	})
	if txErr != nil {
		utils.ResponseError(w, http.StatusInternalServerError, txErr.Error())
		return
	}
	err = utils.WriteJSONResponse(w, "asset created successfully")
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to write response")
	}
}
func GetAsset(w http.ResponseWriter, r *http.Request) {
	brand := r.URL.Query().Get("brand")
	model := r.URL.Query().Get("model")
	assetType := r.URL.Query().Get("type")
	limit, offset := utils.Pagination(r)
	assets, err := dbHelper.LookUpAsset(brand, model, assetType, limit, offset)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to get the assets by the logged user")
		return
	}
	err = utils.WriteJSONResponse(w, assets)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to write response")
		return
	}
}
func RetrieveEmployeeAsset(w http.ResponseWriter, r *http.Request) {
	var retrieveAsset models.RetrieveAssetRequest
	err := utils.ParseJSONBody(r, &retrieveAsset)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to parse request body")
		return
	}
	if !utils.IsAllowedDomain(retrieveAsset.Email) {
		utils.ResponseError(w, http.StatusBadRequest, "email domain is not allowed")
		return
	}
	if retrieveAsset.AssetID == "" {
		utils.ResponseError(w, http.StatusBadRequest, "assetId is required to retrieve asset")
		return
	}

	exist, err := dbHelper.IsAssetAssignedToUser(retrieveAsset)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "error checking asset assigned to user exist")
		return
	}
	if !exist {
		utils.ResponseError(w, http.StatusNotFound, "asset does not exist to the user")
		return
	}
	err = dbHelper.RetrieveAsset(retrieveAsset)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "error retrieving asset from the user")
		return
	}
	err = utils.WriteJSONResponse(w, "retrieved the asset from the user")
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to write response")
		return
	}
}
func UpdateUserAsset(w http.ResponseWriter, r *http.Request) {
	var updateAsset models.UpdateUserAssetRequest
	err := utils.ParseJSONBody(r, &updateAsset)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to parse the request body")
		return
	}
	if *updateAsset.ID == "" {
		utils.ResponseError(w, http.StatusBadRequest, "id is required to update the asset")
		return
	}
	currentStatus, err := dbHelper.GetCurrentStatus(*updateAsset.ID)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to get current status")
		return
	}
	if !utils.CheckStatus(models.AssetStatus(currentStatus), models.AssetStatus(updateAsset.AssetStatus)) {
		utils.ResponseError(w, http.StatusBadRequest, "asset status can't be changed due to restrictions")
		return
	}
	if updateAsset.AssetStatus == models.Service {
		var retrieveAsset models.RetrieveAssetRequest
		retrieveAsset.AssetID = *updateAsset.ID
		retrieveAsset.AssetType = string(models.Service)
		err = dbHelper.RetrieveAsset(retrieveAsset)
		if err != nil {
			utils.ResponseError(w, http.StatusBadRequest, "error retrieving asset from the user")
			return
		}
	}
	rawUserID := r.Context().Value("user_id")
	userID, ok := rawUserID.(string)
	if !ok {
		utils.ResponseError(w, http.StatusUnauthorized, "failed to get the user id from context")
		return
	}
	err = dbHelper.UpdateAsset(updateAsset, userID)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to update the asset")
		return
	}
	err = utils.WriteJSONResponse(w, "asset updated successfully")
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to write response")
		return
	}
}
func GetAssetTimeline(w http.ResponseWriter, r *http.Request) {
	AssetID := r.URL.Query().Get("assetID")
	if AssetID == "" {
		utils.ResponseError(w, http.StatusBadRequest, "assetID is required to get assetTimeline")
		return
	}
	AssetDetails, err := dbHelper.GetAssetDetails(AssetID)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to get the asset details by the logged user")
		return
	}

	err = utils.WriteJSONResponse(w, AssetDetails)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to write response")
		return
	}
}
func SendToService(w http.ResponseWriter, r *http.Request) {
	var serviceDetails models.ServiceDetailsRequest
	err := utils.ParseJSONBody(r, &serviceDetails)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to parse the request body")
		return
	}
	currentStatus, err := dbHelper.GetCurrentStatus(serviceDetails.AssetID)
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "error checking asset status from the user")
		return
	}
	if currentStatus != models.StatusAvailable {
		utils.ResponseError(w, http.StatusBadRequest, "asset status is not valid to send it to service")
		return
	}
	err = dbHelper.CreateServiceRecord(serviceDetails)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to create the service record")
		return
	}
	err = dbHelper.UpdateAssetStatus(models.StatusService, serviceDetails.AssetID)
	if err != nil {
		utils.ResponseError(w, http.StatusInternalServerError, "failed to update the service record")
		return
	}
	err = utils.WriteJSONResponse(w, "the service record created successfully")
	if err != nil {
		utils.ResponseError(w, http.StatusBadRequest, "failed to write response")
		return
	}
}
