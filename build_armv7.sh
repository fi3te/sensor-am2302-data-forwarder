#!/bin/bash
env GOOS=linux GOARCH=arm GOARM=7 go build -o sensor-am2302-data-forwarder ./cmd
