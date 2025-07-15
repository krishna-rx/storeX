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
func (h MouseHandler) UnmarshalConfig(raw jsoniter.RawMessage) (any, error) {
	var config models.MouseDetailsRequest
	err := jsoniter.Unmarshal(raw, &config)
	return config, err
}
func (h MouseHandler) Validate(config any) error {
	mouse := config.(models.MouseDetailsRequest)
	if mouse.Buttons > 1 || mouse.Price > 0 || mouse.MouseName == "" || mouse.IMEI1 == "" || mouse.IMEI2 == "" {
		return fmt.Errorf("invalid mouse config")
	}
	return nil
}
func (h MouseHandler) Insert(tx *sqlx.Tx, assetID string, userID string, config any) error {
	mouse := config.(models.MouseDetailsRequest)
	err := dbHelper.InsertIntoMouse(tx, assetID, userID, mouse)
	if err != nil {
		return err
	}
	return nil
}
func (h MobileHandler) UnmarshalConfig(raw jsoniter.RawMessage) (any, error) {
	var config models.MobileDetailsRequest
	err := jsoniter.Unmarshal(raw, &config)
	return config, err
}
func (h MobileHandler) Validate(config any) error {
	mobile := config.(models.MobileDetailsRequest)
	if mobile.RAM == "" || mobile.OSVersion == "" || mobile.MobileName == "" || mobile.IMEI1 == "" || mobile.IMEI2 == "" {
		return fmt.Errorf("invalid mobile config")
	}
	return nil
}
func (h MobileHandler) Insert(tx *sqlx.Tx, assetID string, userID string, config any) error {
	mobile := config.(models.MobileDetailsRequest)
	err := dbHelper.InsertIntoMobile(tx, assetID, userID, mobile)
	if err != nil {
		return err
	}
	return nil
}
func (h KeyboardHandler) UnmarshalConfig(raw jsoniter.RawMessage) (any, error) {
	var config models.KeyboardDetailsRequest
	err := jsoniter.Unmarshal(raw, &config)
	return config, err
}
func (h KeyboardHandler) Validate(config any) error {
	keyboard := config.(models.KeyboardDetailsRequest)
	if keyboard.Price > 0 || keyboard.KeyboardName == "" || keyboard.IMEI1 == "" || keyboard.IMEI2 == "" || keyboard.Type == "" {
		return fmt.Errorf("invalid keyboard config")
	}
	return nil
}
func (h KeyboardHandler) Insert(tx *sqlx.Tx, assetID string, userID string, config any) error {
	keyboard := config.(models.KeyboardDetailsRequest)
	err := dbHelper.InsertIntoKeyboard(tx, assetID, userID, keyboard)
	if err != nil {
		return err
	}
	return nil
}
