package utils

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"main.go/models"
	"net/http"
	"os"
	"strconv"
	"strings"
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

var validUserTypes = map[string]bool{
	"freelancer": true,
	"intern":     true,
	"full_time":  true,
}

func IsValidUserType(value string) bool {
	_, ok := validUserTypes[strings.ToLower(value)]
	return ok
}

// validate domain
func IsAllowedDomain(email string) bool {
	email = strings.ToLower(email)
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	domain := parts[1]
	if strings.EqualFold(domain, "remotestate.com") {
		return true
	}
	return false
}
func Pagination(r *http.Request) (int, int) {
	page := 1
	limit := 10
	queryParams := r.URL.Query()
	if pageValue := queryParams.Get("page"); pageValue != "" {
		if p, err := strconv.Atoi(queryParams.Get("page")); err == nil && p > 0 {
			page = p
		}
	}
	if limitValue := queryParams.Get("limit"); limitValue != "" {
		if l, err := strconv.Atoi(queryParams.Get("limit")); err == nil && l > 0 {
			limit = l
		}
	}
	offset := (page - 1) * limit

	return limit, offset
}

var validStatusTypes = map[models.AssetStatus][]models.AssetStatus{
	models.StatusAvailable:        {models.StatusAssigned, models.StatusService},
	models.StatusAssigned:         {models.StatusAvailable, models.StatusDamaged},
	models.StatusDamaged:          {models.StatusService},
	models.StatusService:          {models.StatusWaitingForRepair},
	models.StatusWaitingForRepair: {models.StatusAvailable, models.StatusDamaged},
}

func CheckStatus(currentStatus, requestStatus models.AssetStatus) bool {
	validNextStatuses, ok := validStatusTypes[currentStatus]
	if !ok {
		return false
	}
	for _, status := range validNextStatuses {
		if status == requestStatus {
			return true
		}
	}
	return false
}
