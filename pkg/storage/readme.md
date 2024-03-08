For simplicity, all Ports under storage should inherit from
voter.go and voter_history.go and use them as value objects.
This means that the storage package is technically not conforming
to a true hexagonal design pattern, but the alternative is
having to create a ton of adapters for each Port which is way
less maintainable. In this setup incoming Ports have the 
responsibility of data validation and converting the data to
the storage format. It should be possible to use composition
to deal with situations where voter.go and voter_history.go 
require polymorphic behavior. In that circumstance something
like a Visitor Pattern may be effective.