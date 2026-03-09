package security

import "golang.org/x/crypto/bcrypt"

// HashPassword takes a plain text password and returns its bcript hash
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), err
}

// CheckPassword compares a bcrypt hashed password with its plaintext equivalent
func CheckPassword(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
