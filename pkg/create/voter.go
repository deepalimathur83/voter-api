package create

//This is part of the create Port!!!
//
//The Port is "create", Voter is a part of that port. The Ports
//responsibility is to create entries in a repository

// Notice that this is modeling the data we need to create a new
// Voter. It has everything a voter needs.

type HistoryMap map[int]VoterHistory

type Voter struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	VoterHistory HistoryMap
}
