# go-bah

Golang bluetooth audio hub

### should be able to use the debain docker image +bluez +someaudiolib to take care of any cma bluetooth issues

## using:

- https://github.com/tinygo-org/bluetooth
- https://github.com/DamonQin/portaudio

### Flow

- start bt server (bluetooth/hub(\_test).go)
  - scan for inputs
  - select an input
  - set up input stream (portaudio?) to get bytes
- start bt perif (bluetooth/node(\_test).go)
  - scan for perfis
  - on each select
  - set up output stream (portaudio) and register with the core (to get the bytes, fan channel?)

## Set up the pi

`export PI_ADDR=<your pi's address/domain>`

ssh
`ssh-copy-id pi@$PI_ADDR`

install go

```
# things to do
curl -L https://go.dev/dl/go1.17.5.linux-armv6l.tar.gz -o go.tar.gz
sudo tar -C /usr/local -xzf go.tar.gz
rm go.tar.gz
```

install bluetooth packages

```
apt install bluez bluez-tools pulseaudio-module-bluetooth --no-install-recommends -y
```

set up bluetooth

```
sudo nano /etc/systemd/system/bluetooth.target.wants/bluetooth.service
# change
# ExecStart=/usr/lib/bluetooth/bluetoothd
# to
# ExecStart=/usr/lib/bluetooth/bluetoothd --noplugin=sap
sudo systemctl daemon-reload
sudo service bluetooth restart
```

```
sudo nano /etc/bluetooth/audio.conf
```

```
Enable=Source,Sink,Media,Socket
```

```
sudo nano /etc/bluetooth/main.conf
Class = 0x41C

sudo usermod -a -G bluetooth pulse
echo "load-module module-bluetooth-policy" >> /etc/pulse/system.pa
echo "load-module module-bluetooth-discover" >> /etc/pulse/system.pa

```

# the prototype/POC that works

```
./proto clean
./proto disco (connect your audio source)
## ignore the error: /home/pi/go-bah/prototype/bash/gobah: line 188: warning: command substitution: ignored null byte in input
./proto scan (you have like 20s to get all your bluetooth speakers connected up, you can rerun this as needed)
### make note of the MAC addys of the audio outputs you want to add
./proto connect <insert the MAC addy here>
### ... as many as you like (up to the number of bluetooth adapters you have -1 for the source)
./proto start
### start the music
```

then add go to the ~/.bashrc
`echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc`

always good to `sudo reboot` after some installs

## using on arm6l (where there is no vscode)

`./pigo` will just ssh and enter the project directory of the pi (which defaults to `/home/pi/go-bah`) and execute `go $@`
`./pisync.sh` will rsync the local files over to the pi

## try it out dev-mode

```
./pigo run main.go scan
./pigo run main.go connect <your bt headphones> out
./pigo run main.go start
```

## to test audio

### if you can connect some bluetooth outputs this should play in at least 1 of them

on the pi...

```
sudo apt install alsa-utils
speaker-test
```
