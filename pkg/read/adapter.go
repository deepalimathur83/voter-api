package read

import (
	"errors"

	"drexel.edu/voter-api/pkg/storage"
)

/**
The Design Pattern. Also try googling it if there are any questions

This code uses a design pattern called Hexagaonal design. The goal
of the design is to make it so that every piece of the code is as
independant as possible. How do we avhieve this? We create Ports
and adapters to achieve a seperation of concerns and isolate each
system.

What is a Port?

A port is an abstract definition of how the application communicates
with the outside world. It provides an interface that allows the
core application to interact with external systems, such as
databases, user interfaces, external APIs, messaging systems, or
any other external dependencies.

In simple terms Port is an abstraction that does something
(i.e., a verb). In this case, our ports correspond to database
operations create, read, update, and delete. Each of the ports
is responsible for its respective operation.

What is an adapter?

Adapters act as the bridge between the core application and the
outside world. They translate the core application's requests and
responses into a format that the external systems or components
understand, and vice versa. Adapters are specific implementations
that adhere to the interfaces (or ports) defined by the core
application.


There are two types of adapters within the hexagonal architecture:


Inbound Adapters / Driven Adapters:

These adapters handle the communication between the external
systems/components and the core application. They are responsible
for receiving requests from external systems, adapting or
transforming those requests, and passing them to the core
application. Inbound adapters translate external system-specific
details into a format that the core application can understand.

Outbound Adapters / Driving Adapters:

These adapters are responsible for connecting the core application
to external systems/components. They receive the core application's
requests, adapt them into a format suitable for the external system,
and communicate with the external system to perform the requested
operations. Outbound adapters handle the translation of the core
application's responses into a format that the external system can
understand.

What is the simple explanation?

Adapters bridge communications between Ports. In this case Create is
going to Bridge the rest Port, the repository port, as well as itself.
It does this by defining an interface. Any object can use create as a
bridge and achieve polymorphic behavior as long as they implement
the interface.

**/

/**
Step 1: Define the adapter interface. What do you want your adapter
to do? In our case we have two objects, Voter and Voter History and
we want to create them. So lets create an interface for that.
*/

type Adapter interface {
	ReadVoter(int) (Voter, error)
	ReadVoterHistory(int, int) (VoterHistory, error)
	ReadAllVoter() ([]*Voter, error)
	ReadAllVoterHistory(int) ([]*VoterHistory, error)
}

type Repository interface {
	GetItem(int) (*storage.Voter, error)

	UpdateItem(*storage.Voter) error

	GetAllItems() ([]storage.Voter, error)
}

// Now we create a struct to implement the Adapter interface

type adapter struct {
	r Repository
}

func New(r Repository) Adapter {
	return &adapter{r}
}

// Get Voter

func (a *adapter) ReadVoter(voterId int) (Voter, error) {

	if voterId < 1 {

		return Voter{}, errors.New("invalid Voter Id")
	}

	voter, err := a.r.GetItem(voterId)

	if err != nil {

		return Voter{}, err
	}

	voterHistory := make(HistoryMap)

	for _, item := range voter.VoterHistory {

		voterHistory[item.PollId] = VoterHistory{
			PollId:   item.PollId,
			VoteId:   item.VoteId,
			VoteDate: item.VoteDate,
		}
	}

	returnVoter := Voter{
		Id:           voter.Id,
		Name:         voter.Name,
		Email:        voter.Email,
		VoterHistory: voterHistory,
	}

	return returnVoter, nil
}

// Get Voter History

func (a *adapter) ReadVoterHistory(voterId int, pollId int) (VoterHistory, error) {

	targetVoter, err := a.r.GetItem(voterId)
	if err != nil {
		return VoterHistory{}, err
	}

	targetHistory, exists := targetVoter.VoterHistory[pollId]

	if !exists {

		return VoterHistory{}, errors.New("the specified pollId does not exists inside the voter")

	}

	returnObject := VoterHistory{
		PollId:   targetHistory.PollId,
		VoteId:   targetHistory.VoteId,
		VoteDate: targetHistory.VoteDate,
	}

	return returnObject, nil

}

// Get all Voters

func (a *adapter) ReadAllVoter() ([]*Voter, error) {

	allVoters, err := a.r.GetAllItems()

	if err != nil {
		return nil, err
	}

	var voters []*Voter
	for _, voter := range allVoters {
		voterHistory := make(HistoryMap)
		for _, item := range voter.VoterHistory {
			voterHistory[item.PollId] = VoterHistory{
				PollId:   item.PollId,
				VoteId:   item.VoteId,
				VoteDate: item.VoteDate,
			}
		}

		voterObj := &Voter{
			Id:           voter.Id,
			Name:         voter.Name,
			Email:        voter.Email,
			VoterHistory: voterHistory,
		}
		voters = append(voters, voterObj)
	}

	return voters, nil
}

// Get all voter history

func (a *adapter) ReadAllVoterHistory(voterId int) ([]*VoterHistory, error) {

	// Assuming GetItem is used to fetch a single voter by ID

	voter, err := a.r.GetItem(voterId)
	if err != nil {
		return nil, err
	}

	var voterHistories []*VoterHistory
	for _, history := range voter.VoterHistory {
		voterHistory := VoterHistory{
			PollId:   history.PollId,
			VoteId:   history.VoteId,
			VoteDate: history.VoteDate,
		}
		voterHistories = append(voterHistories, &voterHistory)
	}

	return voterHistories, nil
}
