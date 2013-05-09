#!/bin/bash

# Simple MAC-address collector for WhoIsInTheHub2
#
# Author: rekoil (adrian@bjugard.se)
# Date: 2013-05-09
#
# Depends on:
# 	aircrack-ng
# Arguments:
# 	1. How long to run the script for

# -- Change this to your interface --
IFACE="wlp2s0"
# -----------------------------------

# Set variables
if [[ "$1" =~ ^[0-9]+$ ]]; then
	RUNTIME=$1
else
	# Set default runtime because NaN was given
	RUNTIME=15
fi
STRING="Station*"
MON="mon0"

# Make sure interface is available
if [ "$(ifconfig | grep -o $MON)" != "$MON" ]; then
	/usr/sbin/airmon-ng start $IFACE
fi

# Generate data
(/usr/bin/sleep $RUNTIME && /usr/bin/killall airodump-ng) &
/usr/sbin/airodump-ng --output-format csv mon0 -w out

# Clean output
/usr/bin/sed -i "1,/$STRING/d" out-01.csv
/usr/bin/cut -f1 -d"," out-01.csv > out.csv
/usr/bin/rm out-01.csv