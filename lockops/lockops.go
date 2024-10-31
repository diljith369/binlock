package lockops

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"golang.org/x/crypto/pbkdf2"
)

func BinLockWithoutSalt(inputFile, password string) ([]byte, error) {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return nil, getError("error reading file", err)
	}

	// Create cipher
	block, err := aes.NewCipher([]byte(padPassword(password)))
	if err != nil {
		return nil, getError("error creating cipher", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, getError("error creating GCM", err)
	}

	// Create nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, getError("error creating nonce", err)
	}

	// Encrypt the binary
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func BinLockWithSalt(inputFile, password string) ([]byte, error) {

	data, err := os.ReadFile(inputFile)
	if err != nil {
		return nil, getError("error reading file", err)
	}

	// Generate random salt
	salt := make([]byte, 32)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	// Derive key using PBKDF2
	key := deriveKey(password, salt)

	// Create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Generate random nonce
	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	// Create GCM cipher
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Encrypt data
	ciphertext := gcm.Seal(nil, nonce, data, nil)

	// Combine salt, nonce, and ciphertext
	result := make([]byte, 0, len(salt)+len(nonce)+len(ciphertext))
	result = append(result, salt...)
	result = append(result, nonce...)
	result = append(result, ciphertext...)

	return result, nil
}

func CreateBinLockerFile(outputFile string, ciphertext []byte) error {
	// Write the encrypted binary
	err := os.WriteFile(outputFile, ciphertext, 0644)
	if err != nil {
		return getError("error writing file", err)
	}
	return nil
}

func UnlockwithoutSalt(protectedPath, password string, debug bool) ([]byte, error) {
	debugPrint := func(format string, args ...interface{}) {
		if debug {
			fmt.Printf(format+"\n", args...)
		}
	}

	debugPrint("Reading protected file: %s", protectedPath)

	// Read the protected binary
	data, err := os.ReadFile(protectedPath)
	if err != nil {
		return nil, getError("failed to read protected file", err)
	}

	debugPrint("Protected file size: %d bytes", len(data))
	if debug {
		fmt.Printf("First 16 bytes of protected file: %s\n", hex.EncodeToString(data[:16]))
	}

	// Create cipher
	block, err := aes.NewCipher(padPassword(password))
	if err != nil {
		return nil, getError("failed to create cipher", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, getError("failed to create GCM", err)
	}

	// Get nonce size
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, getError("encrypted file is too short", nil)
	}

	// Split nonce and ciphertext
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	debugPrint("Nonce size: %d bytes", nonceSize)
	debugPrint("Ciphertext size: %d bytes", len(ciphertext))

	// Decrypt the binary
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, getError("failed to decrypt", err)
	}

	debugPrint("Decrypted binary size: %d bytes", len(plaintext))
	if debug {
		fmt.Printf("First 16 bytes of decrypted binary: %s\n", hex.EncodeToString(plaintext[:16]))
	}

	return plaintext, nil
}

func UnlockWithSalt(protectedPath, password string, debug bool) ([]byte, error) {

	encryptedData, err := os.ReadFile(protectedPath)
	if err != nil {
		return nil, getError("failed to read protected file", err)
	}
	// Extract salt, nonce, and ciphertext
	salt := encryptedData[:32]
	nonce := encryptedData[32:44]
	ciphertext := encryptedData[44:]

	// Derive key using PBKDF2 with the extracted salt
	key := deriveKey(password, salt)

	// Create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Create GCM cipher
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Decrypt data
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

func padPassword(password string) []byte {
	if len(password) > 32 {
		password = password[:32]
	}
	result := make([]byte, 32)
	copy(result, password)
	return result
}

func getError(errMessage string, err error) error {
	return fmt.Errorf("%s %v", errMessage, err)
}

func RunProtectedBinary(filecontent []byte, protectedPath string, binargs []string, debug bool) error {

	debugPrint := func(format string, args ...interface{}) {
		if debug {
			fmt.Printf(format+"\n", args...)
		}
	}
	// Create temporary file with appropriate extension
	ext := ""
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}

	tempDir := os.TempDir()
	tempPath := filepath.Join(tempDir, fmt.Sprintf("exec_%d%s", os.Getpid(), ext))

	debugPrint("Creating temporary file: %s", tempPath)

	// Write decrypted binary to temp file
	err := os.WriteFile(tempPath, filecontent, 0755)
	if err != nil {
		return getError("failed to write temp file", err)
	}

	// Make sure the file is executable
	if runtime.GOOS != "windows" {
		if err := os.Chmod(tempPath, 0755); err != nil {
			return getError("failed to set executable permissions", err)
		}
	}

	debugPrint("Executing binary: %s", tempPath)

	// Execute the decrypted binary
	cmd := exec.Command(tempPath, binargs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Dir = filepath.Dir(protectedPath) // Set working directory to the same as protected file

	// Clean up the temporary file
	defer func() {
		debugPrint("Cleaning up temporary file: %s", tempPath)
		os.Remove(tempPath)
	}()

	if err := cmd.Run(); err != nil {
		return getError("failed to run binary", err)
	}
	return nil
}

func deriveKey(password string, salt []byte) []byte {
	iterations := 100000 // High number of iterations for PBKDF2
	keyLen := 32         // For AES-256
	return pbkdf2.Key([]byte(password), salt, iterations, keyLen, sha256.New)
}
