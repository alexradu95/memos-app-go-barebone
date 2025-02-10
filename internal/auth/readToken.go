package auth

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type Claims map[string]interface{}

func ExtractClaims(token string) (Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("Invalid token format: token must have three parts.")
	}

	// Get the claims part (second segment)
	claimsPart := parts[1]

	// Add padding if needed
	if l := len(claimsPart) % 4; l > 0 {
		claimsPart += strings.Repeat("=", 4-l)
	}

	// Decode base64
	decoded, err := base64.URLEncoding.DecodeString(claimsPart)
	if err != nil {
		return nil, errors.New("Failed to decode claims: " + err.Error())
	}

	// Parse JSON
	var claims Claims
	if err := json.Unmarshal(decoded, &claims); err != nil {
		return nil, errors.New("Failed to parse claims JSON: " + err.Error())
	}

	return claims, nil
}

func ValidateClaims(claims Claims) error {
	now := time.Now().Unix()

	// Check expiration time
	if exp, ok := claims["exp"].(float64); ok {
		if int64(exp) < now {
			return errors.New("Token has expired.")
		}
	}

	// Check issued at time
	if iat, ok := claims["iat"].(float64); ok {
		if int64(iat) > now {
			return errors.New("Token issued in the future")
		}
	}

	// Check if claims include user identifier
	if _, ok := claims["accountId"]; !ok {
		return errors.New("Token does not include user identifier.")
	}

	return nil
}
