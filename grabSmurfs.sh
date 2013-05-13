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
#	2. Address to POST to

# -- Change this to your interface --
IFACE="wlan0"
# -----------------------------------

# Exit if no site argument given
if [ -z "$2" ]; then
	echo "Usage: grabSmurfs.sh <time to scan> <address to POST to>"
	exit
fi

# Define utilities
AIRMON="/usr/sbin/airmon-ng"
AIRDUMP="/usr/sbin/airodump-ng"

SLEEP="/bin/sleep"
KILLALL="/usr/bin/killall"
SED="/bin/sed"
CUT="/usr/bin/cut"
RM="/bin/rm"
WGET="/usr/bin/wget"

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
	$AIRMON start $IFACE
fi

# Generate data
($SLEEP $RUNTIME && $KILLALL airodump-ng) &
$AIRDUMP --output-format csv mon0 -w out

# Clean output
$SED -i "1,/$STRING/d" out-01.csv
$CUT -f1 -d"," out-01.csv > out.txt
$SED -i -e '$d' out.txt

# Format JSON
JSON="$(echo -e "{\"mac\":[")"
while read -r line; do
	JSON+="$(echo -e "\"$line\",")"
done < out.txt

JSON="${JSON:0:${#JSON}-1}]}"
#JSON+="$(echo -e '\b]}')"

# Prints out JSON
# echo -e $JSON > out.json

# Sends JSON via wget
#echo $WGET "\"http://$2?\'site\'=\'addMacs\'&\'json\'=$JSON\""
addr="http://$2?\"site\"=\"addMacs\"&\"json\"=$JSON"

$WGET $addr

# Clean up
$RM out-01.csv
$RM out.txt
