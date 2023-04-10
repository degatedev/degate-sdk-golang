#!/usr/bin/env bash

go build -buildmode=c-shared -o degate.so clibh.go main.go