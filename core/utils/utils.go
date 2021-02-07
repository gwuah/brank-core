package utils

import (
	crypto_rand "crypto/rand"
	"encoding/binary"
	"fmt"
	"log"
	"net/http"

	"brank/core"
	"math/rand"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/rs/xid"
)

type CustomClaims struct {
	AppID    int `json:"app_id"`
	ClientID int `json:"client_id"`
	jwt.StandardClaims
}

func GenerateUUID() string {
	return xid.New().String()
}

func ConvertToUint64(num string) uint64 {
	id64, _ := strconv.ParseUint(num, 10, 64)
	return id64
}

func ConvertToInt(num string) int {
	id64, _ := strconv.ParseInt(num, 10, 64)
	return int(id64)
}
func Bool(b bool) *bool {
	temp := b
	return &temp
}

func String(v string) *string {
	return &v
}

func StringValue(v *string) string {
	if v != nil {
		return *v
	}
	return ""
}

func GenerateTopic(topic string) string {
	return fmt.Sprintf("%s%s", core.Get("CLOUDKARAFKA_TOPIC_PREFIX", ""), topic)
}

func Seed() int64 {
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		panic("cannot Seed math/rand package with cryptographically secure random number generator")
	}
	return int64(binary.LittleEndian.Uint64(b[:]))
}

func StringWithCharset(prefix string, length int) string {
	const charset = "0123456789"
	b := make([]byte, length)
	var seededRand = rand.New(rand.NewSource(Seed()))
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return fmt.Sprintf("%s%s", prefix, string(b))
}

func Error(err error, m *string, code int) core.BrankResponse {
	if err != nil {
		log.Println(err)
	}

	var message string
	if m == nil {
		message = "failure"
	} else {
		message = StringValue(m)
	}

	return core.BrankResponse{
		Error: true,
		Code:  code,
		Meta: core.BrankMeta{
			Data:    nil,
			Message: message,
		},
	}
}

func Success(data *map[string]interface{}, m *string) core.BrankResponse {

	var message string
	if m == nil {
		message = "success"
	} else {
		message = StringValue(m)
	}

	return core.BrankResponse{
		Error: false,
		Code:  http.StatusOK,
		Meta: core.BrankMeta{
			Data:    data,
			Message: message,
		},
	}
}
