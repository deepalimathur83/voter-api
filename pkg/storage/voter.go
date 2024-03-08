package storage

type HistoryMap map[int]VoterHistory

// This is part of the redis Port
//
// This models the storage format for Voter. Notice that
// unlike create, we have a VoterHistoryMap. This is because we
// Don't maintain two seperate tables for Voter and Voter History
// We combine them into one Voter Object and save them as a Voter
type Voter struct {
	Id           int        `json:"id"`
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	VoterHistory HistoryMap `json:"history"`
}
