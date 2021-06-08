#!/bin/bash

if [ "$1" = "start" ]; then
	tsc ui/js/indigo.ts --target ES6
	go run main.go &
	node index.js &
	echo "Project Indigo started!"
elif [ "$1" = "stop" ]; then
	ps -ef | grep "node" | grep -v grep | awk '{print $2}' | xargs kill
	ps -ef | grep "go" | grep -v grep | awk '{print $2}' | xargs kill
	echo "Project Indigo stopped!"
else
	echo "Valid options: start and stop"
fi