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
