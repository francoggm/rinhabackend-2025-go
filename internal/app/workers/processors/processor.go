package processors

type Processor interface {
	ProcessEvent(event any) error
}
