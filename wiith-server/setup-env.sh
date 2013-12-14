#!/bin/bash
# To use: source setup-env.sh
# Written by Emil 'Eda' Edholm (digIT 2013)

# The path to the script, i.e. the project folder
export GOPATH=$(cd `dirname "${BASH_SOURCE[0]}"` && pwd)
export PATH="$PATH:$GOPATH/bin"
