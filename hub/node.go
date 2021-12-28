package hub

// NID - node id is a random and unique value
type NID string

// Node is an interface between a Bluetooth Device and the Hub
type Node interface {
	// ID returns the NID
	ID() NID
	// Connected is this Node connected
	Connected() bool
}
