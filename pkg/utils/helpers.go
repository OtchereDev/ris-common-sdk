package users

import (
	"fmt"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const numberBytes = "0123456789"

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func CompareUserPassword(password string, userPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(password))
}

func GenerateRandomPassword(passwordLen int) string {
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, passwordLen)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}

	password := string(b)

	return password
}

func GenerateOTP(otpLen int) string {
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, otpLen)
	for i := range b {
		b[i] = numberBytes[rand.Intn(len(numberBytes))]
	}

	otp := string(b)

	return otp
}

func CreateAuditLog(userId string, activity string, timestamp time.Time) (log map[string]interface{}) {
	log = make(map[string]interface{})
	log["userId"] = userId
	log["activity"] = fmt.Sprintf(activity)
	log["timestamp"] = timestamp.Format(time.RFC3339)

	return
}
