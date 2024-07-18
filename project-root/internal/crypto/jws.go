package crypto

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Function to make an HTTP GET request to the specified URL and return the boolean status.
func GetStatusFromJWS(url string, publicKey *ecdsa.PublicKey) (bool, error) {
	// Make the HTTP GET request
	resp, err := http.Get(url)
	if err != nil {
		return false, fmt.Errorf("failed to make HTTP GET request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("failed to read response body: %v", err)
	}

	// Parse and validate the JWS
	token, err := jwt.Parse(string(body), func(token *jwt.Token) (interface{}, error) {
		// Check the signing method
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		return false, fmt.Errorf("failed to parse JWS: %v", err)
	}

	// Check token validity
	if !token.Valid {
		return false, errors.New("invalid JWS token")
	}

	// Extract and validate the claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return false, errors.New("failed to extract claims")
	}

	iat := int64(claims["iat"].(float64))
	exp := int64(claims["exp"].(float64))
	iss := claims["iss"].(string)

	if iat > time.Now().Unix() {
		return false, errors.New("token used before issued")
	}

	if exp < time.Now().Unix() {
		return false, errors.New("token expired")
	}

	expectedIss := fmt.Sprintf("http://example.com/api/status/%s", strings.Split(url, "/")[4])
	if iss != expectedIss {
		return false, fmt.Errorf("unexpected issuer: %s", iss)
	}

	// Extract the status from the payload
	statusPayload := claims["status"].(map[string]interface{})
	encodedList := statusPayload["encodedList"].(string)
	index := int(statusPayload["index"].(float64))

	decodedList, err := base64.StdEncoding.DecodeString(encodedList)
	if err != nil {
		return false, fmt.Errorf("failed to decode encoded list: %v", err)
	}

	byteIndex := index / 8
	bitIndex := index % 8

	if byteIndex >= len(decodedList) {
		return false, errors.New("index out of range")
	}

	status := (decodedList[byteIndex] & (1 << bitIndex)) != 0

	return status, nil
}

// Utility function to parse the public key from a PEM-encoded string
func ParseECDSAPublicKeyFromPEM(pemEncodedKey string) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemEncodedKey))
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, errors.New("failed to decode PEM block containing public key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse ECDSA public key: %v", err)
	}

	publicKey, ok := pub.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("not ECDSA public key")
	}

	return publicKey, nil
}
