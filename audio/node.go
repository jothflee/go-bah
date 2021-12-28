package audio

import (
	"fmt"

	"github.com/gordonklaus/portaudio"
)

const sampleRate = 44100
const seconds = 1

func ListAudioDevices() []*portaudio.DeviceInfo {
	portaudio.Initialize()
	defer portaudio.Terminate()
	devices, err := portaudio.Devices()
	if err != nil {
		panic(err)
	}

	for _, d := range devices {
		fmt.Println(d.Name)
	}

	return devices
}
func Output(intf string) []byte {
	return []byte{}
}

func Input(intf string, data []byte) {

}
