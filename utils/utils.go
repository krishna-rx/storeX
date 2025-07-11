package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"main.go/models"
	"net/http"
	"os"
	"time"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type clientError struct {
	StatusCode    int    `json:"statusCode"`
	MessageToUser string `json:"messageToUser"`
}

func ResponseError(w http.ResponseWriter, statusCode int, messageToUser string) {
	logrus.Errorf("status : %d, message : %s", statusCode, messageToUser)
	clientErr := &clientError{
		StatusCode:    statusCode,
		MessageToUser: messageToUser,
	}
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(clientErr); err != nil {
		logrus.Errorf("failed to send the error %+v", err)
	}
}
func GenerateJWT(req models.UserLoginRequest) (string, error) {
	claims := jwt.MapClaims{
		"user_id": req.ID,
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
		"role":    req.Role,
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := []byte(os.Getenv("JWT_SECRET"))
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
func ParseJSONBody(r *http.Request, dst interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return fmt.Errorf("failed to decode JSON body: %w", err)
	}
	return nil
}
func WriteJSONResponse(w http.ResponseWriter, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")

	encoder := json.NewEncoder(w)
	err := encoder.Encode(data)
	if err != nil {
		return fmt.Errorf("failed to encode JSON response: %w", err)
	}
	return nil
}
func validField(text string) bool {
	if text == "" {
		return false
	}
	return true
}
