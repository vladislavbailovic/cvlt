#!/bin/bash

cat /tmp/cvlt.fifo &
go run .

wait
