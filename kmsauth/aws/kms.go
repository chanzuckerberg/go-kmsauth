package aws

import (
	"encoding/base64"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
	"github.com/aws/aws-sdk-go/service/kms/kmsiface"
	"github.com/pkg/errors"
)

// KMS is a kms client
type KMS struct {
	Svc kmsiface.KMSAPI
}

// NewKMS returns a KMS client
func NewKMS(s *session.Session) KMS {
	return KMS{kms.New(s)}
}

// EncryptBytes encrypts the plaintext using the keyID key and the given context, result is base64 encoded
func (k *KMS) EncryptBytes(keyID string, plaintext []byte, context map[string]*string) (string, error) {
	input := &kms.EncryptInput{}
	input.SetKeyId(keyID).SetPlaintext(plaintext).SetEncryptionContext(context)
	response, err := k.Svc.Encrypt(input)
	if err != nil {
		return "", errors.Wrap(err, "KMS encryption failed")
	}
	return base64.StdEncoding.EncodeToString(response.CiphertextBlob), nil
}

// Decrypt decrypts a b64 string
func (k *KMS) Decrypt(ciphertext []byte, context map[string]*string) ([]byte, string, error) {

	input := &kms.DecryptInput{}
	input.SetCiphertextBlob(ciphertext).SetEncryptionContext(context)
	response, err := k.Svc.Decrypt(input)
	if err != nil {
		return nil, "", errors.Wrap(err, "KMS decryption failed")
	}
	if response.KeyId == nil {
		return nil, "", errors.New("Nil KMS keyID returned from AWS")
	}
	return response.Plaintext, *response.KeyId, nil
}
