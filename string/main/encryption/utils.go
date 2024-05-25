package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
	"os"

	"golang.org/x/crypto/scrypt"
)

func GetSalt(saltFile string) (string, error) {
	salt, err := os.ReadFile(saltFile)
	if err != nil {
		return "", err
	}
	return string(salt), nil
}

func GenerateSalt(saltFile string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	if err := os.WriteFile(saltFile, salt, 0644); err != nil {
		return "", err
	}
	return string(salt), nil
}

// GenerateKey derives a key from a passphrase using scrypt
func GenerateKey(passphrase, saltStr string) ([]byte, error) {
	salt := []byte(saltStr) // should be securely stored/retrieved
	key, err := scrypt.Key([]byte(passphrase), salt, 32768, 8, 1, 32)
	if err != nil {
		return nil, err
	}
	return key, nil
}

// EncryptFile encrypts the input file and writes to the output file
func EncryptFile(inputFile, outputFile, passphrase, saltFile string) error {
	// Generate salt
	salt, err := GenerateSalt(saltFile)
	if err != nil {
		return err
	}

	// Generate encryption key
	key, err := GenerateKey(passphrase, salt)
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
	if err := os.WriteFile(outputFile, ciphertext, 0444); err != nil {
		return err
	}

	return nil
}

// DecryptFile decrypts the input file and writes to the output file
func DecryptFile(inputFile, outputFile, passphrase, saltFile string) error {
	// Read the salt
	salt, err := GetSalt(saltFile)
	if err != nil {
		return err
	}

	// Generate encryption key
	key, err := GenerateKey(passphrase, salt)
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
