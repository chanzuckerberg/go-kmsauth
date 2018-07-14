package util_test

import (
	"testing"

	"github.com/chanzuckerberg/go-kmsauth/kmsauth/util"
	"github.com/stretchr/testify/assert"
)

func TestStrPtrMapEqual(t *testing.T) {
	a := assert.New(t)

	foo := "foo"
	bar := "bar"
	foo2 := "foo"

	this := map[string]*string{
		"foo": &foo,
	}
	that := map[string]*string{
		"foo": &foo2,
	}
	// diff ptrs same values
	a.True(util.StrPtrMapEqual(this, that))
	a.True(util.StrPtrMapEqual(that, this))

	// same ptrs
	this = map[string]*string{
		"foo": &foo,
	}
	that = map[string]*string{
		"foo": &foo,
	}
	a.True(util.StrPtrMapEqual(this, that))
	a.True(util.StrPtrMapEqual(that, this))

	// diff values
	this = map[string]*string{
		"foo": &bar,
	}
	that = map[string]*string{
		"foo": &foo,
	}
	a.False(util.StrPtrMapEqual(this, that))
	a.False(util.StrPtrMapEqual(that, this))

	// one nil
	this = map[string]*string{
		"foo": nil,
	}
	that = map[string]*string{
		"foo": &foo,
	}
	a.False(util.StrPtrMapEqual(this, that))
	a.False(util.StrPtrMapEqual(that, this))

	// both nil
	this = map[string]*string{
		"foo": nil,
	}
	that = map[string]*string{
		"foo": nil,
	}
	a.True(util.StrPtrMapEqual(this, that))
	a.True(util.StrPtrMapEqual(that, this))

	// diff lengths
	this = map[string]*string{
		"foo": &foo,
		"bar": &bar,
	}
	that = map[string]*string{
		"foo": &foo,
	}
	a.False(util.StrPtrMapEqual(this, that))
	a.False(util.StrPtrMapEqual(that, this))

	// same lengths, diff keys
	this = map[string]*string{
		"foo": &foo,
	}
	that = map[string]*string{
		"bar": &foo,
	}
	a.True(len(this) == len(that))
	a.False(util.StrPtrMapEqual(this, that))
	a.False(util.StrPtrMapEqual(that, this))
}
