package voteconf

import (
	"errors"
	"strconv"
	"strings"
)

func validateHashtag(hashtag string) (err error) {

	shorter := strings.TrimPrefix(hashtag, "#")

	if shorter == hashtag {
		err = errors.New("Hashtag has to Start with a #")
	}

	return
}

func validateVote(vote string) (voteValue int, err error) {

	voteValue, err = strconv.Atoi(vote)

	if voteValue < 0 || voteValue > 5 || err != nil {
		err = errors.New("Votes must be a number 1, 2, 3, 4 or 5")
	}
	return
}
