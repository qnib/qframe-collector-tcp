#!/bin/bash

sleep 2
echo "Test-$(date +%s)" | nc -w1  $1 11001
