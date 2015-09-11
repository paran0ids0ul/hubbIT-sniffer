#!/usr/bin/env bash

iface=wlan0
if [[ -n "$1" ]]; then
    iface=$1
fi

if [ $(id -u) -ne 0 ]; then 
    >&2 echo "You need to be root. Otherwise you will only get leftover scan data"
    exit 1
fi

# Print the mac addresses of the nearby access points
iwlist $iface scan | grep Address | awk '{print $5}'
