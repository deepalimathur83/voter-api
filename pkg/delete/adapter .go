package delete

import (
	"errors"

	"drexel.edu/voter-api/pkg/storage"
)

type Adapter interface {
	DeleteVoter(int) error
	DeleteVoterHistory(int, int) error
	//DeleteAllVoterHistory(int) error
	DeleteAllVoters() error
}

type Repository interface {
	GetItem(int) (*storage.Voter, error)
	UpdateItem(item *storage.Voter) error
	DeleteItem(int) error
	DeleteVoterHistory(int, int) error
	//DeleteAllVoterHistory(int) error
	DeleteAllVoters() error
}

// Now we create a struct to implement the Adapter interface

type adapter struct {
	r Repository
}

func New(r Repository) Adapter {
	return &adapter{r}
}

func (a *adapter) DeleteVoter(voterId int) error {
	if voterId < 1 {
		return errors.New("invalid Voter Id")
	}
	return a.r.DeleteItem(voterId)
}

// Delete voter history

func (a *adapter) DeleteVoterHistory(voterId int, pollId int) error {

	if voterId < 1 || pollId < 1 {

		return errors.New("invalid Voter Id or Poll Id")
	}

	targetVoter, err := a.r.GetItem(voterId)
	if err != nil {
		return err
	}

	if _, exists := targetVoter.VoterHistory[pollId]; !exists {
		return errors.New("the specified pollId does not exist in the voter's history")
	}

	delete(targetVoter.VoterHistory, pollId)

	err = a.r.UpdateItem(targetVoter)

	if err != nil {
		return err
	}

	return nil
}

// Delete All

func (a *adapter) DeleteAllVoters() error {

	return a.r.DeleteAllVoters()
}
