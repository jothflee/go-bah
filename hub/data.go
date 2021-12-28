package hub

import "time"

// Data is the interface defining the communication payload of a Node
type Data interface {
	// the size of the data len(Size())
	Size() int64
	// the data
	Payload() []byte
	// the time the data was recieved from the source
	Time() time.Time
	// Source is the originating node of the data
	Source() NID
}
