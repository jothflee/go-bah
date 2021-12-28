package bluetooth

import "testing"

func TestNodeConnect(t *testing.T) {
	n := NewNode("")

	n.Connect()
}
