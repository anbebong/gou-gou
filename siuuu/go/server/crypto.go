package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base32"
	"encoding/base64"

	// "encoding/hex"
	"fmt"
	"io"
	"math/big"
	"time"

	"github.com/pquerna/otp/totp"
)

// Khóa mã hóa (phải là 16, 24, hoặc 32 byte)
var aesKey = []byte("1234567897654321") // 16 byte cho AES-128

// encrypt mã hóa văn bản bằng AES
func encrypt(text string) (string, error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}

	plaintext := []byte(text)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// decrypt giải mã văn bản bằng AES
func decrypt(cryptoText string) (string, error) {
	ciphertext, err := base64.URLEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}

	if len(ciphertext) < aes.BlockSize {
		return "", fmt.Errorf("ciphertext quá ngắn")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return string(ciphertext), nil
}

// generateOTP tạo một mã OTP ngẫu nhiên gồm 6 chữ số.
// HÀM NÀY SẼ ĐƯỢC THAY THẾ BẰNG LOGIC TOTP
func generateOTP() (string, error) {
	const otpLength = 6
	max := new(big.Int)
	max.Exp(big.NewInt(10), big.NewInt(otpLength), nil) // 10^6

	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	// Định dạng số thành chuỗi 6 chữ số, thêm số 0 ở đầu nếu cần
	return fmt.Sprintf("%0*d", otpLength, n), nil
}

// generateOTPSecretFromClientID thực hiện thuật toán bạn yêu cầu để tạo OTPSecret từ ClientID.
func generateOTPSecretFromClientID(clientID string) string {
	// 1. Tính SHA1 của ClientID, làm việc với raw bytes
	hasher := sha1.New()
	hasher.Write([]byte(clientID))
	sha1Bytes := hasher.Sum(nil) // Mảng 20 bytes

	// 2. Chia mảng 20 bytes thành 5 phần, mỗi phần 4 bytes
	p1 := sha1Bytes[0:4]
	p2 := sha1Bytes[4:8]
	p3 := sha1Bytes[8:12]
	p4 := sha1Bytes[12:16]
	p5 := sha1Bytes[16:20]

	// 3. Ghép các phần byte lại theo thứ tự 5 + 1 + 4 + 2 + 3
	var secretBytes []byte
	secretBytes = append(secretBytes, p5...)
	secretBytes = append(secretBytes, p1...)
	secretBytes = append(secretBytes, p4...)
	secretBytes = append(secretBytes, p2...)
	secretBytes = append(secretBytes, p3...)

	// 4. Mã hóa kết quả sang Base32 để tương thích với thư viện TOTP
	secret := base32.StdEncoding.EncodeToString(secretBytes)
	return secret
}

// generateTOTP tạo mã TOTP hợp lệ dựa trên một secret.
func generateTOTP(secret string) (string, error) {
	return totp.GenerateCode(secret, time.Now())
}
