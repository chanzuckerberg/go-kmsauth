package kmsauth_test

import (
	"testing"

	kmsauth "github.com/chanzuckerberg/go-kmsauth"
	"github.com/chanzuckerberg/go-kmsauth/kmsauth/util"
	"github.com/stretchr/testify/assert"
)

func TestAuthContextValidate(t *testing.T) {
	a := assert.New(t)
	ac := kmsauth.AuthContextV1{From: "foo"}
	a.NotNil(ac.Validate())
	ac = kmsauth.AuthContextV1{From: "foo", To: "bar"}
	a.Nil(ac.Validate())
	ac2 := kmsauth.AuthContextV2{From: "foo"}
	a.NotNil(ac2.Validate())
	ac2 = kmsauth.AuthContextV2{From: "foo", To: "bar"}
	a.NotNil(ac2.Validate())
	ac2 = kmsauth.AuthContextV2{From: "foo", To: "bar", UserType: "foobar"}
	a.Nil(ac2.Validate())
}

func TestAuthContextGetUsername(t *testing.T) {
	a := assert.New(t)
	ac := kmsauth.AuthContextV1{To: "foo"}
	a.Equal(ac.GetUsername(), "")
	ac = kmsauth.AuthContextV1{From: "foo", To: "bar"}
	a.Equal(ac.GetUsername(), "foo")
	ac2 := kmsauth.AuthContextV2{To: "foo"}
	a.Equal(ac2.GetUsername(), "2//")
	ac2 = kmsauth.AuthContextV2{From: "foo", To: "bar"}
	a.Equal(ac2.GetUsername(), "2//foo")
	ac2 = kmsauth.AuthContextV2{From: "foo", To: "bar", UserType: "gas"}
	a.Equal(ac2.GetUsername(), "2/gas/foo")
}

func TestAuthContextGetKSMContext(t *testing.T) {
	a := assert.New(t)

	foo := "foo"
	bar := "bar"
	ac := kmsauth.AuthContextV1{From: foo, To: bar}

	expected := map[string]*string{
		"foo": &foo,
		"bar": &bar,
	}

	a.True(util.StrPtrMapEqual(ac.GetKMSContext(), expected))

	ac2 := kmsauth.AuthContextV2{From: "foo", To: "bar", UserType: "gas"}
	a.Equal(ac2.GetUsername(), "2/gas/foo")
}
