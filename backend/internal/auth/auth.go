package auth

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/argon2"
)

const (
	argonMemory      = 64 * 1024
	argonIterations  = 3
	argonParallelism = 2
	argonSaltLength  = 16
	argonKeyLength   = 32
	tokenExpiry      = 60 * time.Minute
)

var (
	ErrInvalidHash         = errors.New("invalid hash format")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrTokenExpired        = errors.New("token expired")
	ErrTokenInvalid        = errors.New("invalid token")
	ErrUserInactive        = errors.New("user inactive")
)

type Claims struct {
	UserID    string   `json:"sub"`
	TenantID  string   `json:"tid"`
	Role      string   `json:"role"`
	Permissions []string `json:"permissions"`
	jwt.RegisteredClaims
}

type AuthConfig struct {
	Secret    string
	Issuer    string
	Audience  string
}

func NewAuthConfig(secret, issuer, audience string) *AuthConfig {
	return &AuthConfig{
		Secret:   secret,
		Issuer:   issuer,
		Audience: audience,
	}
}

func (c *AuthConfig) GenerateToken(userID, tenantID, role string, permissions []string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID:      userID,
		TenantID:    tenantID,
		Role:        role,
		Permissions: permissions,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			Issuer:    c.Issuer,
			Audience:  jwt.ClaimStrings{c.Audience},
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenExpiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(c.Secret))
}

func (c *AuthConfig) ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(c.Secret), nil
	}, jwt.WithIssuer(c.Issuer), jwt.WithAudience(c.Audience))

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

func HashPassword(password string) (string, error) {
	salt := make([]byte, argonSaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, argonIterations, argonMemory, argonParallelism, argonKeyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, argonMemory, argonIterations, argonParallelism, b64Salt, b64Hash), nil
}

func VerifyPassword(password, encodedHash string) error {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return ErrInvalidHash
	}

	var memory, iterations, parallelism uint32
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &memory, &iterations, &parallelism)
	if err != nil {
		return ErrInvalidHash
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return ErrInvalidHash
	}

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return ErrInvalidHash
	}

	computedHash := argon2.IDKey([]byte(password), salt, iterations, memory, uint8(parallelism), uint32(len(decodedHash)))

	if subtle.ConstantTimeCompare(computedHash, decodedHash) == 1 {
		return nil
	}

	return ErrInvalidCredentials
}

func DummyHash() {
	_ = argon2.IDKey([]byte("dummy"), make([]byte, argonSaltLength), argonIterations, argonMemory, argonParallelism, argonKeyLength)
}