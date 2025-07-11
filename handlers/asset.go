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
	assets, err := dbHelper.GetAssetByRole(brand, model, assetType)
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
