package security

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
	// Default bcrypt cost is fine for most development setups.
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func VerifyPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

