package utility

import(
	"crypto/rand"
	"math/big"
	"strings"
)

func CreateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	charsetLen := big.NewInt(int64(len(charset)))
	var result strings.Builder

	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, charsetLen)
		if err != nil {
			return "", err
		}
		result.WriteByte(charset[randomIndex.Int64()])
	}

	return result.String(), nil

}