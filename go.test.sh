#!/usr/bin/env bash

set -e
go test -coverprofile=coverage.txt -covermode=atomic github.com/wtetsu/gaze/pkg/...

