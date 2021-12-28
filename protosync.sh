#! /bin/bash 

if [ -z "$PI_ADDR" ]; then echo "set the PI_ADDR env var with your pi's ip address"; exit 1; fi
rsync -a --delete ./prototype/ pi@$PI_ADDR:/home/pi/go-bah/prototype

