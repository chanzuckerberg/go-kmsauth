package kmsauth

import (
	"encoding/json"
	"time"

	"github.com/chanzuckerberg/tmp/pkg/aws"
	"github.com/pkg/errors"
)

// TokenVersion is a token version
type TokenVersion int

const (
	// TokenVersion1 is a token version
	TokenVersion1 = 1
	// TokenVersion2 is a token version
	TokenVersion2 = 2
)

// TokenGenerator generates a token
type TokenGenerator struct {
	// AuthKey the key_arn or alias to use for authentication
	AuthKey        string
	AuthContext    AuthContext
	Region         string
	TokenVersion   TokenVersion
	TokenLifetime  time.Duration
	TokenCacheFile *string
	AwsClient      *aws.Client
}

// NewTokenGenerator returns a new kmsauth token generator
func NewTokenGenerator() *TokenGenerator {
	return &TokenGenerator{}
}

// getCachedToken tries to fetch a token from the cache
func (tg *TokenGenerator) getCachedToken() (*Token, error) {
	return nil, nil
}

// cacheToken caches a token
func (tg *TokenGenerator) cacheToken(token *Token) error {
	return nil
}

// GetToken gets a token
func (tg *TokenGenerator) GetToken() (*Token, error) {
	now := time.Now()
	// Start the notBefore x time in the past to avoid clock skew
	notBefore := now.Add(-1 * timeSkew)
	// Set the notAfter x time in the future but account for timeSkew
	notAfter := now.Add(tg.TokenLifetime - timeSkew)

	token, err := tg.getCachedToken()
	if err != nil {
		return nil, err
	}

	token = &Token{
		NotBefore: TokenTime(notBefore),
		NotAfter:  TokenTime(notAfter),
	}
	return token, nil
}

// GetEncryptedToken returns the encrypted kmsauth token
func (tg *TokenGenerator) GetEncryptedToken() (*EncryptedToken, error) {
	token, err := tg.GetToken()
	if err != nil {
		return nil, err
	}

	tokenBytes, err := json.Marshal(token)
	if err != nil {
		return nil, errors.Wrap(err, "Could not marshal token")
	}

	encrypted, err := tg.AwsClient.KMS.EncryptBytes(tg.AuthKey, tokenBytes, tg.AuthContext.GetKMSContext())
}
