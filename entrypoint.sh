#!/bin/sh

if [ "$1" = "-V" ] || [ "$1" = "--version" ]; then
         ./drive-scanner -V
elif [ "$1" = "-h" ] || [ "$1" = "--help" ]; then
        ./drive-scanner -h
else
        ./drive-scanner | jq
fi