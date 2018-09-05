package aws

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// Client is an AWS client
type Client struct {
	KMS KMS
}

// NewClient returns a new aws client
func NewClient(sess *session.Session) (*Client, error) {
	var err error
	if sess == nil {
		log.Debug("nil aws.Session provided, attempting to create one")
		sess, err = session.NewSession()
		if err != nil {
			return nil, errors.Wrap(err, "Could not create aws session")
		}
	}
	return &Client{KMS: NewKMS(sess)}, nil
}
