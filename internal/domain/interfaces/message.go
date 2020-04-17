package interfaces

// Message for HTTP proxying
type Message interface {
	Type() int
}
