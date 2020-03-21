#!/usr/bin/env bash
set -x
go run main.go -e 'bash -c "sleep 1 && exit 0"' -c 2
go run main.go -e 'bash -c "sleep 5 && exit 0"' -c 2
go run main.go -e 'bash -c "sleep 1 && exit 1"' -c 2
go run main.go -e 'sh -c "sleep 10 && exit 1"' -c 2
go run main.go -e 'bash -c "if [ -f lock ]; then exit 1; fi; sleep 10 && touch lock && exit 1"' -c 1
go run main.go -e 'bash -c "if [ -f lock ]; then exit 1; fi; sleep 10 && touch lock && exit 1"' -c 1