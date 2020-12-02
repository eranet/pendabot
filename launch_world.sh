#!/bin/bash

./server/gnatsd  &
python3 hw/send_command.py
