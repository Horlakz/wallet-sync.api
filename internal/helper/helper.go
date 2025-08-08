package helper

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/snowflake"
)

var (
	perm []int
	mu   sync.Mutex
)

func GenerateSnowflakeID() (int64, error) {
	mu.Lock()
	defer mu.Unlock()

	if len(perm) == 0 {
		perm = rand.Perm(1023)

		for i := range perm {
			perm[i]++
		}
	}

	nodeNum := perm[len(perm)-1]
	perm = perm[:len(perm)-1]

	node, err := snowflake.NewNode(int64(nodeNum))
	if err != nil {
		return 0, err
	}

	// Sleep for 1 milliseond to prevent generating the same ID.
	time.Sleep(1 * time.Millisecond)

	// Generate a snowflake ID.
	id := node.Generate()

	return id.Int64(), nil
}

func GenerateRandomDigits(length int) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	code := make([]byte, length)
	for i := range code {
		code[i] = byte(rand.Intn(10) + 48)
	}

	return string(code)
}

func GenerateRandomString(length int) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ123456790")

	code := make([]rune, length)
	for i := range code {
		code[i] = chars[rand.Intn(len(chars))]
	}

	return string(code)
}

func GenerateRandomAlphaNumeric(length int) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz1234567890")

	code := make([]rune, length)
	for i := range code {
		code[i] = chars[rand.Intn(len(chars))]
	}

	return string(code)
}

func GenerateRandomHexStr(length int) string {
	rand.New(rand.NewSource(time.Now().UnixNano()))

	code := make([]byte, length)
	rand.Read(code)

	return hex.EncodeToString(code)
}

func GenerateSlug(title string) string {
	// Convert title to lowercase
	title = strings.ToLower(title)

	// Replace spaces with hyphens
	title = strings.ReplaceAll(title, " ", "-")

	// Remove special characters
	title = removeSpecialChars(title)

	return title
}

func removeSpecialChars(title string) string {
	// Define special characters to remove
	specialChars := []string{"!", "@", "#", "$", "%", "^", "&", "*", "(", ")", "_", "+", "=", "{", "}", "[", "]", "|", "\\", ":", ";", "\"", "'", "<", ">", ",", ".", "?", "/"}

	// Remove special characters from the title
	for _, char := range specialChars {
		title = strings.ReplaceAll(title, char, "")
	}

	return title
}

func Int64ToString(i int64) string {
	return strconv.FormatInt(i, 10)
}

func StringToInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func ShortCode(name string, numLetters int) (string, error) {
	if len(name) == 0 {
		return "", errors.New("empty name provided")
	}

	if numLetters <= 0 || numLetters > len(name) {
		return "", errors.New("invalid number of letters")
	}

	return name[:numLetters], nil
}

func GenerateTimestamp() string {
	return time.Now().Format("20060102150405")
}

func StringToPointer(s string) *string {
	return &s
}

func GenerateAccountNumber() string {
	prefix := "50"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	num := r.Intn(100000000) // 8 digits
	return prefix + fmt.Sprintf("%08d", num)
}
