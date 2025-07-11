package handlers

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	jsoniter "github.com/json-iterator/go"
	"main.go/database/dbHelper"
	"main.go/models"
)

type AssetHandler interface {
	UnmarshalConfig(raw jsoniter.RawMessage) (any, error)
	Validate(config any) error
	Insert(tx *sqlx.Tx, assetID string, userID string, config any) error
}

type LaptopHandler struct{}
type MouseHandler struct{}
type MobileHandler struct{}
type KeyboardHandler struct{}

func (h LaptopHandler) UnmarshalConfig(raw jsoniter.RawMessage) (any, error) {
	var config models.LaptopDetailsRequest
	err := jsoniter.Unmarshal(raw, &config)
	return config, err
}
func (h LaptopHandler) Validate(config any) error {
	laptop := config.(models.LaptopDetailsRequest)
	if laptop.RAM == "" || laptop.OSVersion == "" || laptop.LaptopName == "" || laptop.IMEI1 == "" || laptop.IMEI2 == "" {
		return fmt.Errorf("invalid laptop config")
	}
	return nil
}
func (h LaptopHandler) Insert(tx *sqlx.Tx, assetID string, userID string, config any) error {
	laptop := config.(models.LaptopDetailsRequest)
	err := dbHelper.InsertIntoLaptop(tx, assetID, userID, laptop)
	if err != nil {
		return err
	}
	return nil
}
