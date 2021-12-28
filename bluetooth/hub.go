package bluetooth

import (
	"fmt"

	"github.com/jothflee/bluetooth"
	"github.com/jothflee/go-bah/config"
)

type Hub struct {
	adapter *bluetooth.Adapter
}

func NewHub(adapterInterface string) Hub {
	return Hub{
		adapter: bluetooth.NewAdapter(adapterInterface),
	}
}

func (h Hub) Scan() {
	h.ScanToConnect("")
}

func (h Hub) Connect(d *config.Device) {
	d.Device, d.Name = h.ScanToConnect(d.MAC)

}
func (h Hub) ScanToConnect(mac string) (device *bluetooth.Device, name string) {
	fmt.Println("starting scan")
	adapter := h.adapter
	// Enable BLE interface.
	must("enable BLE stack", adapter.Enable())

	ch := make(chan bluetooth.ScanResult, 1)
	func() {
		// Start scanning.
		fmt.Println("scanning...")
		adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
			if mac == "" {
				fmt.Println("found device:", result.Address.String(), result.RSSI, result.LocalName())
			} else {
				if result.Address.String() == mac {
					fmt.Println("setting up", result.LocalName(), mac)

					adapter.StopScan()
					ch <- result
				}
			}
		})
	}()

	// var device *bluetooth.Device
	if mac != "" {

		result := <-ch
		fmt.Println("connecting to", result.LocalName(), result.Address)
		var err error
		for {
			device, err = adapter.Connect(result.Address, bluetooth.ConnectionParams{})
			if err != nil {
				fmt.Println("error", err.Error())
			} else {
				break
			}
		}
		name = result.LocalName()
		fmt.Println("connected to", name, result.Address.String())
	}

	return device, name

}
