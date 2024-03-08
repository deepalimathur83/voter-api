package storage

import (
	"time"
)

// This is part of the redis Port
//
// Notice that it models the storage format for voter history. Take
// Note however, that there is no actual Voter History repository.
// Rather Voter History will be used within a map inside of Voter.go \
// so we only need to maintain one object in Redis. See redis/voter.go
type VoterHistory struct {
	PollId   int       `json:"poll_id"`
	VoteId   int       `json:"vote_id"`
	VoteDate time.Time `json:"vote_date"`
}
