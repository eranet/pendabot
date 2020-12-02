#!/bin/bash

./server/gnatsd  &
python3 hw/send_command.py &
go run hw/encoders_from_arduino.go
