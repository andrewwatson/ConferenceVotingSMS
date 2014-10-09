package voteconf

import (
	"github.com/stretchr/testify/assert"
	// "os"

	"testing"
)

func TestValidateGoodTag(t *testing.T) {

	err := validateHashtag("#hashtag")
	assert.Nil(t, err, "err was not nil")
}

func TestValidateBadTag(t *testing.T) {
	err := validateHashtag("hashtag")

	assert.NotNil(t, err, "hashtag should have failed")
}

func TestValidateGoodVote(t *testing.T) {

	voteOutput, err := validateVote("1")
	assert.Nil(t, err, "err happen")
	assert.Equal(t, 1, voteOutput)
}

func TestValidateBadVote(t *testing.T) {

	_, err := validateVote("7")
	assert.NotNil(t, err, "should have been caught")
}

func TestValidateNonNumberVote(t *testing.T) {
	_, err := validateVote("A")
	assert.NotNil(t, err, "should have been caught")
}
