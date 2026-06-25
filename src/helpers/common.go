package helpers

import (
	"crypto/rand"
	"encoding/hex"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yourorg/go-mvc-app/src/config"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword returns a bcrypt hash of the plain-text password.
func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(b), err
}

// CheckPassword compares a bcrypt hash against a plain-text password.
func CheckPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// GenerateToken returns a cryptographically random hex token of byteLen bytes.
func GenerateToken(byteLen int) (string, error) {
	b := make([]byte, byteLen)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// IsValidEmail checks email format via a basic regex.
func IsValidEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// ParsePagination extracts page and page_size query params with defaults.
func ParsePagination(c *gin.Context) (page, pageSize int) {
	page, _ = strconv.Atoi(c.DefaultQuery("page", strconv.Itoa(config.DefaultPage)))
	pageSize, _ = strconv.Atoi(c.DefaultQuery("page_size", strconv.Itoa(config.DefaultPageSize)))
	if page < 1 {
		page = config.DefaultPage
	}
	if pageSize < 1 || pageSize > config.MaxPageSize {
		pageSize = config.DefaultPageSize
	}
	return
}

// Offset calculates the DB offset from page and pageSize.
func Offset(page, pageSize int) int {
	return (page - 1) * pageSize
}

// SanitizeString trims whitespace and lowercases a string.
func SanitizeString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

// ContainsString checks whether slice contains target.
func ContainsString(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}
