#!/usr/bin/env bash
g++ --version
gvm version
go version

mkdir -p bin

make lib
make lib-tests
make bindings
make examples