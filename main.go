package main

import (
	"fmt"
	"os"
	"os/signal"
	"path"
	"strings"
	"syscall"

	"github.com/jothflee/go-bah/bluetooth"
	"github.com/jothflee/go-bah/config"
	"gopkg.in/yaml.v3"
)

var GlobalConfig *config.Config
var GlobalConfigPath = "~/.bah.config"

func init() {
	loadConfig()
}

func main() {
	args := os.Args[1:]
	fmt.Println(args)

	if len(args) > 0 {
		switch args[0] {
		case "scan":
			// scan is always 1 adapter
			hub := bluetooth.NewHub("")
			hub.Scan()
		case "connect":
			if len(args) == 3 {
				addNewDevice(args[1], config.DeviceType(args[2]))
			} else {
				panic("connect requires 2 args <device address> and <in|out>")
			}
		case "start":
			start()
		}
	} else {
		fmt.Printf(`
bluetooth audio hub (bah)
manages the onboard bluetooth drivers and audio drivers to create a
many in many out bluetooth audio mixing expierence

for now, the cli tool creates a tmp config file (located in ~/.bah.yml) to connect and mix the audio
later it will get more complicate.

Usage:
bah scan - scans for availble bluetooth
bah connect AA:BB:CC:DD:EE:FF in|out - connects to a bluetooth device as an input (in) or output (out)
bah start - starts up the audio hub, connecting to availble inputs and outputs and mixes them together
`)

	}
}

func addNewDevice(MAC string, t config.DeviceType) {
	d := config.NewDevice(MAC, t)
	GlobalConfig.AddDevice(&d)
	writeConfig()
}

func start() {
	fmt.Println("Starting BAH!")
	// TODO: generate a list of adapaters (hciconfig)
	adapters := []string{"hci1", "hci0"}
	// cannot have more adapters than devices in theory
	// shim for now
	hubs := map[string]bluetooth.Hub{}

	for _, adapter := range adapters {
		hubs[adapter] = bluetooth.NewHub(adapter)
	}

	adapterI := 0

	for addr, d := range GlobalConfig.Devices {
		if adapterI < len(adapters) {
			fmt.Println("Connecting to", addr)
			hubs[adapters[adapterI]].Connect(d)
		}
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	done := make(chan struct{}, 1)

	go func() {
		<-sigs
		done <- struct{}{}
	}()

	<-done
	for _, d := range GlobalConfig.Devices {
		d.Disconnect()
	}
	writeConfig()
}

func loadConfig() {
	fp := getGlobalConfigPath()

	b, err := os.ReadFile(fp)
	if err == nil {
		err = yaml.Unmarshal(b, &GlobalConfig)
		if err != nil {
			panic(err)
		}
	} else {

		GlobalConfig = config.NewConfig()
	}

	// this will update any changes to the config object
	// in theory we could migrate here, but this is a shim anyway so...
	// do this for now
	writeConfig()
}
func writeConfig() {
	fmt.Println("writing config out...")
	fp := getGlobalConfigPath()

	os.MkdirAll(path.Dir(fp), 0755)
	b, err := yaml.Marshal(GlobalConfig)
	if err != nil {
		panic(err)
	}
	os.WriteFile(fp, b, 0655)
}
func getGlobalConfigPath() string {
	hdir, _ := os.UserHomeDir()
	return strings.Replace(GlobalConfigPath, "~", hdir, 1)

}
