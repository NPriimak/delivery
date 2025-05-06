package ddd

type AggregateRoot interface {
	GetDomainEvents() []DomainEvent
	ClearDomainEvents()
	RaiseDomainEvent(DomainEvent)
}

type BaseAggregate struct {
	domainEvents []DomainEvent
}

func NewBaseAggregate() *BaseAggregate {
	return &BaseAggregate{
		domainEvents: make([]DomainEvent, 0),
	}
}

func (a *BaseAggregate) ClearDomainEvents() {
	a.domainEvents = nil
}

func (a *BaseAggregate) GetDomainEvents() []DomainEvent {
	return a.domainEvents
}

func (a *BaseAggregate) RaiseDomainEvent(event DomainEvent) {
	a.domainEvents = append(a.domainEvents, event)
}
