package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"math"
	"os"

	"golang.org/x/crypto/scrypt"
)

/* func HashPassword(password string, saltBytes []byte) (string, error) {
	hash, err := scrypt.Key()
	if err != nil {
		return "", err
	}
	return string(hash), nil
} */

func HashInput(input interface{}, saltBytes []byte) (string, error) {
	// Adapt the input if needed
	var adaptedInput []byte
	switch v := input.(type) {
	case string:
		adaptedInput = []byte(v)
	case []float32:
		adaptedInput = make([]byte, len(v)*4)
		for i, f := range v {
			binary.LittleEndian.PutUint32(adaptedInput[i*4:], math.Float32bits(f))
		}
	default:
		return "", fmt.Errorf("invalid input type")
	}

	// Hash the input
	hash := sha256.New()
	hash.Write(adaptedInput)
	return hex.EncodeToString(hash.Sum(nil)), nil
}

func GetSalt(saltPath string) (string, error) {
	salt, err := os.ReadFile(saltPath)
	if err != nil {
		return "", err
	}
	return string(salt), nil
}

func GenerateSalt(saltPath string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	if err := os.WriteFile(saltPath, salt, 0644); err != nil {
		return "", err
	}
	return string(salt), nil
}

// GenerateKey derives a key from a hash using scrypt
func GenerateKey(hash string, salt string) ([]byte, error) {
	saltBytes := []byte(salt) // should be securely stored/retrieved
	key, err := scrypt.Key([]byte(hash), saltBytes, 32768, 8, 1, 32)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// EncryptFile encrypts the input file and writes to the output file
func EncryptFile(inputFile, outputFile, hash, salt string) error {
	// Generate encryption key
	key, err := GenerateKey(hash, salt)
	if err != nil {
		return err
	}

	// Read the input file
	plaintext, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}

	// Create a new AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	// Wrap the cipher block in GCM
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	// Create a nonce
	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	// Encrypt the data
	ciphertext := aead.Seal(nonce, nonce, plaintext, nil)

	// Write the encrypted data to the output file
	if err := os.WriteFile(outputFile, ciphertext, 0644); err != nil { // Remember to change it back to 444
		return err
	}

	return nil
}

// DecryptFile decrypts the input file and writes to the output file
func DecryptFile(inputFile string, outputFile string, input interface{}, salt string) error {
	// Hash the input
	hash, err := HashInput(input, []byte(salt))
	if err != nil {
		return err
	}

	// Generate encryption key
	key, err := GenerateKey(hash, salt)
	if err != nil {
		return err
	}

	// Read the input file
	ciphertext, err := os.ReadFile(inputFile)
	if err != nil {
		return err
	}

	// Create a new AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	// Wrap the cipher block in GCM
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	// Extract the nonce
	nonceSize := aead.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// Decrypt the data
	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return err
	}

	// Write the decrypted data to the output file
	if err := os.WriteFile(outputFile, plaintext, 0644); err != nil {
		return err
	}

	return nil
}
