package kmsauth

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/chanzuckerberg/go-kmsauth/kmsauth/aws"
	"github.com/chanzuckerberg/go-kmsauth/kmsauth/util"
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
	Region         string
	TokenVersion   TokenVersion
	TokenLifetime  time.Duration
	TokenCacheFile *string
	AuthContext    AuthContext

	AwsClient *aws.Client

	mutex sync.RWMutex
}

// Validate validates the TokenGenerator
func (tg *TokenGenerator) Validate() error {
	if tg == nil {
		return errors.New("Nil token generator")
	}
	return tg.AuthContext.Validate()
}

// getCachedToken tries to fetch a token from the cache
func (tg *TokenGenerator) getCachedToken() (*Token, error) {
	if tg.TokenCacheFile == nil {
		log.Debug("No TokenCacheFile specified")
		return nil, nil
	}
	// lock for reading
	tg.mutex.RLock()
	defer tg.mutex.RUnlock()

	if _, err := os.Stat(*tg.TokenCacheFile); os.IsNotExist(err) {
		// token cache file does not exist
		return nil, nil
	}
	cacheBytes, err := ioutil.ReadFile(*tg.TokenCacheFile)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not open token cache file %s", *tg.TokenCacheFile))
	}

	tokenCache := &TokenCache{}
	err = json.Unmarshal(cacheBytes, tokenCache)
	if err != nil {
		return nil, errors.Wrap(err, "Could not unmarshal token cache")
	}

	// Compare token cache with current params
	ok := util.StrPtrMapEqual(tokenCache.AuthContext, tg.AuthContext.GetKMSContext())
	if !ok {
		log.Info("Cached token invalid")
		return nil, nil
	}
	now := time.Now()
	// subtract timeSkew to account for clock skew
	notAfter := tokenCache.Token.NotAfter.Add(-1 * timeSkew)
	if now.Before(notAfter) {
		// token is still valid, use it
		return &tokenCache.Token, nil
	}

	// otherwise expired, need new token
	return nil, nil
}

// cacheToken caches a token
func (tg *TokenGenerator) cacheToken(tokenCache *TokenCache) error {
	if tg.TokenCacheFile == nil {
		log.Debug("No TokenCacheFile specified")
		return nil
	}
	// lock for writing
	tg.mutex.Lock()
	defer tg.mutex.Unlock()

	dir := path.Dir(*tg.TokenCacheFile)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return errors.Wrapf(err, "Could not create cache directories %s", dir)
	}

	data, err := json.Marshal(tokenCache)
	if err != nil {
		return errors.Wrap(err, "Could not marshal token cache")
	}

	return errors.Wrap(
		ioutil.WriteFile(*tg.TokenCacheFile, data, 0644),
		"Could not write token to cache")
}

// GetToken gets a token
func (tg *TokenGenerator) GetToken() (*Token, error) {
	newToken := NewToken(tg.TokenLifetime)
	token, err := tg.getCachedToken()
	if err != nil {
		return nil, err
	}
	// If we could not find a token then return a new one
	if token == nil {
		token = newToken
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

	encryptedStr, err := tg.AwsClient.KMS.EncryptBytes(
		tg.AuthKey,
		tokenBytes,
		tg.AuthContext.GetKMSContext())

	encryptedToken := EncryptedToken(encryptedStr)
	return &encryptedToken, err
}
