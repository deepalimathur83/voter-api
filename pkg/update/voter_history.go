package update

import (
	"time"
)

// This is part of the create Port!!!
//
// The Port is "create", VoterHistory is a part of that port. The Ports
// responsibility is to create entries in a repository

// Notice that this is modeling the data we need to create a new
// Voter History. It has everything history needs.
type VoterHistory struct {
	PollId   int       `json:"id"`
	VoteId   int       `json:"vote_id"`
	VoteDate time.Time `json:"vote_date"`
}
