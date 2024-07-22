package crypto

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Status represents the status structure in the payload
type Status struct {
	EncodedList string `json:"encodedList"`
	Index       int    `json:"index"`
}

// JWSResponse represents the structure of the JWS payload
type JWSResponse struct {
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
	Issuer    string `json:"iss"`
	Status    Status `json:"status"`
}

// ParseJWSResponse parses the JWS response body and validates the signature
func ParseJWSResponse(body []byte, publicKey *ecdsa.PublicKey) (Status, error) {
	// Convert the body to a string
	bodyStr := string(body)

	// Parse the JWS token
	token, err := jwt.Parse(bodyStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return publicKey, nil
	})

	if err != nil {
		return Status{}, err
	}

	// Validate the token
	if !token.Valid {
		return Status{}, errors.New("invalid token")
	}

	// Extract the claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return Status{}, errors.New("invalid claims")
	}

	// Convert claims to JSON
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return Status{}, err
	}

	// Unmarshal the JSON into JWSResponse struct
	var response JWSResponse
	err = json.Unmarshal(claimsJSON, &response)
	if err != nil {
		return Status{}, err
	}

	return response.Status, nil
}

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
