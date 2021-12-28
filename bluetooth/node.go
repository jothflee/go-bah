package bluetooth

import (
	"fmt"
	"time"

	"github.com/jothflee/bluetooth"
	"github.com/jothflee/go-bah/hub"
)

type Node struct {
	hub.Node
	nid     hub.NID
	adapter *bluetooth.Adapter
}

// var adapter = bluetooth.DefaultAdapter

func NewNode(adapterInterface string) *Node {
	return &Node{
		nid:     "a",
		adapter: bluetooth.NewAdapter(adapterInterface),
	}
}

func (n Node) NID() hub.NID {
	return n.nid
}

func (n Node) Connected() bool {
	return false
}

func (n *Node) Connect() {
	adapter := n.adapter
	// fmt.Println("setting up")
	// // Enable BLE interface.
	// must("enable BLE stack", adapter.Enable())

	// // Start scanning.
	// println("scanning...")
	// err := adapter.Scan(func(adapter *bluetooth.Adapter, device bluetooth.ScanResult) {
	// 	println("found device:", device.Address.String(), device.RSSI, device.LocalName())

	// })
	// must("start scan", err)
	disconnected := true
	adapter.SetConnectHandler(func(d bluetooth.Addresser, c bool) {
		connected := c
		fmt.Println("connection change", c)

		if !connected && !disconnected {

			disconnected = true
		}

		if connected {
			disconnected = false
			device, err := adapter.Connect(d, bluetooth.ConnectionParams{})
			if err != nil {
				println(err.Error())
				return
			}
			svc, _ := device.DiscoverServices(nil)
			fmt.Println(svc)

		}
	})

	must("enable BLE stack", adapter.Enable())
	adv := adapter.DefaultAdvertisement()
	must("config adv", adv.Configure(bluetooth.AdvertisementOptions{
		LocalName: "Hub Output",
	}))
	must("start adv", adv.Start())

	var audioCharasteristic bluetooth.Characteristic
	var (
		serviceUUID = [16]byte{0xa0, 0xb4, 0x00, 0x01, 0x92, 0x6d, 0x4d, 0x61, 0x98, 0xdf, 0x8c, 0x5c, 0x62, 0xee, 0x53, 0xb3}
		charUUID    = [16]byte{0xa0, 0xb4, 0x00, 0x02, 0x92, 0x6d, 0x4d, 0x61, 0x98, 0xdf, 0x8c, 0x5c, 0x62, 0xee, 0x53, 0xb3}
	)
	buffer := make([]byte, 1024) // start out with red
	must("add service", adapter.AddService(&bluetooth.Service{
		UUID: bluetooth.NewUUID(serviceUUID),
		Characteristics: []bluetooth.CharacteristicConfig{
			{
				Handle: &audioCharasteristic,
				UUID:   bluetooth.NewUUID(charUUID),
				Value:  buffer[:],
				Flags:  bluetooth.CharacteristicReadPermission,

				WriteEvent: func(client bluetooth.Connection, offset int, value []byte) {
					fmt.Printf("write event\nvalue: %d bytes\noffset: %d\n", len(value), offset)

					if offset != 0 || len(value) != 3 {
						return
					}
				},
			},
		},
	}))

	for {
		// set latency timing
		// fmt.Println(string(buffer[0:64]))
		time.Sleep(100 * time.Millisecond)
		// play music
	}

}

func must(action string, err error) {
	if err != nil {
		panic("failed to " + action + ": " + err.Error())
	}
}
