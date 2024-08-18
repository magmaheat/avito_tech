package auth

import (
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

func GenerateRandomEmail() string {
	const charset = "abcdefghijklmnopqrstuvwxyz"
	const emailDomain = "yandex.ru"
	const emailLength = 10

	rand.Seed(time.Now().UnixNano())

	var randStr strings.Builder
	for i := 0; i < emailLength; i++ {
		randStr.WriteByte(charset[rand.Intn(len(charset))])
	}

	return fmt.Sprintf("%s@%s", randStr.String(), emailDomain)
}

func IsValidEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}
