package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"strings"
	"time"

	"github.com/OtchereDev/ris-common-sdk/pkg/proto/audit"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const numberBytes = "0123456789"

func HashPassword(password string) ([]byte, error) {
	if len(strings.Trim(password, " ")) == 0 {
		return nil, fmt.Errorf("password cannot be empty ")
	}
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func CompareUserPassword(password string, userPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(password))
}

func GenerateRandomPassword(passwordLen int) (string, error) {
	b := make([]byte, passwordLen)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letterBytes))))
		if err != nil {
			return "", err
		}
		b[i] = letterBytes[num.Int64()]
	}
	return string(b), nil
}

func GenerateOTP(otpLen int) (string, error) {
	b := make([]byte, otpLen)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(numberBytes))))
		if err != nil {
			return "", err
		}
		b[i] = numberBytes[num.Int64()]
	}
	return string(b), nil
}

func CreateAuditLog(userId string, activity string, timestamp time.Time, action string) (r proto.Message) {
	r = &audit.Audit{
		User:      userId,
		Activity:  activity,
		Action:    action,
		Timestamp: timestamppb.New(timestamp),
	}

	return
}
