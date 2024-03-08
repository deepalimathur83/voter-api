package create

import (
	"errors"
	"strings"

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
	//Make sure you capitalize these to make them public or you
	//won't be able to use them!
	CreateVoter(Voter) error
	CreateVoterHistory(int, VoterHistory) error
}

/**
Step 2: Now that we have our service interface defined, things that
want to implement our service can use Adapter as the absrtraction.

For simplicity, we will also implement our own Adapter in the same
file. We know that our adapter needs to bridge rest and repository
Ports. our adapter will be using a repository. In order to keep this
code isolated from repository, we will define a reppository
interface. As long as the repository implements the method signature
we define in the interface, it will work without concerning itself
with how the data is created. Meaning you could theoretically use
Redis, an in- memory cache, Postgres, or any repository. For now,
focus on Redis...
**/

type Repository interface {
	//We will need a function that adds the voter.
	//remeber that on the storage Port, voter and
	//voter history are one item, so we only need
	//one addItem for both.
	AddItem(*storage.Voter) error

	//We will need a get function to check if the
	//voter exists for the history we are about to add.
	GetItem(int) (*storage.Voter, error)

	//We will need the update function to add new Voter
	//History. remember that there is only one Voter storage
	//object in redis and it has history as a component of that
	//object. so when we add new history we are actually fetching
	//an existing voter, editing the voterHistory map, and saving
	//voter which counts as an update.
	UpdateItem(*storage.Voter) error
}

// Now we create a struct to implement the Adapter interface
type adapter struct {
	r Repository
}

//Our struct will need a New function so that things
//can use it. Notice that we return Adapter (i.e., the interface)
//instead of adapter (i.e, the struct). This is because of the
//Liskov principle which states:
/**
The Liskov Substitution Principle (LSP) is one of the five SOLID
principles of object-oriented programming, named after Barbara
Liskov, who introduced it in the 1980s. The principle states:


"Subtypes must be substitutable for their base types."


In simpler terms, objects of a superclass should be able to be
replaced with objects of any of their subclasses without affecting
the correctness of the program. In other words, if a piece of code
works with a particular type, it should also work correctly when
the code is replaced with a subtype of that type.
**/
func New(r Repository) Adapter {
	return &adapter{r}
}

/**
Now that we have our new function you probably noticed that its
flagging an error:

cannot use &adapter{…} (value of type *adapter) as Adapter value in
return statement: *adapter does not implement Adapter
(missing method createVoter)compilerInvalidIfaceAssign

Thats to be expected. Remember, when we define an interface, the
compiler won't let us instantiate a concrete object unless all the
methods are defined. This same behavior occurs in other languages
like java. So lets do that right now, lets define the createVoter
and createVoterHistory functions.
**/

// Note: (a *adapter) is a way of saying that this function is a
// method of the adapter class. This is done so we can associate
// the createVoter function with the adapter struct. If we didn't
// do this the compiler wouldn't make the association and throw
// the same error we talked about above. Make sure you remeber this
// as it will pop up in all of your adapter! its very important!!!!
func (a *adapter) CreateVoter(voter Voter) error {
	//before we do anything, lets handle some validation
	//lets make sure that the id exists (i.e., != 0) and
	//that the Voter's name and email isn't blank.
	if voter.Id == 0 {
		return errors.New("invalid Voter Id")
	}

	//Note: there are two functions being used in tandem.
	//
	//len returns the length of a collection. Strings are
	//considered collections of characters.
	//
	//strings.TrimSpace is a function TrimSpace returns a
	//slice of the string (in this case voter.Name), with all
	//leading and trailing white space removed, as defined by
	//Unicode. If the string is empty or only has whitespace,
	//we will get a 0 as a result and end up throwing an error
	if len(strings.TrimSpace(voter.Name)) == 0 {
		return errors.New("voter name cannot be blank")
	}

	//Do the same for Email... Extra Credit. find a way to
	//validate that email is in the form of:
	//
	//<accountName>@<domain>
	//
	//hint: google "Regex"
	if len(strings.TrimSpace(voter.Email)) == 0 {
		return errors.New("voter email cannot be blank")
	}

	//Now that we have done some basic data validation and
	//we are confident that we have a valid Voter object, we
	//are in the clear to add it to a repository. So let's
	//convert it into a storage object

	storageObject := storage.Voter{
		Id:    voter.Id,
		Name:  voter.Name,
		Email: voter.Email,
	}

	//that now. Notice that we did this in the method signature
	//(a *adapter). Part of the reason why is so that the
	//compiler can associate this function with the adapter struct.
	//but it is also so this function, as a member of adapter can
	//access its private methods and variables... In this case
	//we want to use a to access r, the repsoitory. lets try it out.
	err := a.r.AddItem(&storageObject)
	if err != nil {
		return err
	}

	//Lets review what we did in this function. We accepted a voter.
	//we validated all the properties (i.e, id, name, email). and
	//we passed it to the repository function. whats next? If all
	//that works out we return nil as the error to let the rest
	//handler know that the operation was a succcess! Lets move on
	//to voter history.

	return nil
}

// Note: Please make sure you understand createVoter (above) before
// you read this. this function is going to be a tad lighter on
// explanations
func (a *adapter) CreateVoterHistory(voterId int, voterHistory VoterHistory) error {

	//first off, lets check if the Voter exists. If not, no need to
	//proceed with validation. we can use the repository to retrieve
	//the voter.
	targetVoter, err := a.r.GetItem(voterId)
	if err != nil {
		return err
	}

	//Now that we have our voter, we need to check if the pollId
	//already exists. If it does, since this is a create Port, we
	//throw an error. If it was an update port, we would just update
	//it.
	if _, exists := targetVoter.VoterHistory[voterHistory.PollId]; exists {
		return errors.New("the specified pollId allready exists inside the voter")
	}

	//Now that we know the pollId doesn't already exist within voter,
	//we just have to convert the history to the storage format, add
	//and add it to the voter.

	storageObject := storage.VoterHistory{
		PollId:   voterHistory.PollId,
		VoteId:   voterHistory.VoteId,
		VoteDate: voterHistory.VoteDate,
	}

	//We need to make sure the map is initialized. If not initialize it.
	if targetVoter.VoterHistory == nil {
		targetVoter.VoterHistory = make(storage.HistoryMap)
	}

	targetVoter.VoterHistory[voterHistory.PollId] = storageObject

	//now we just need to add it back into redis

	a.r.UpdateItem(targetVoter)

	//Now that we have defined our adapter on the create port,
	//we have to define an adapter on the redis port. lets jump
	//over to redis!

	return nil
}
