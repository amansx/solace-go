#!/usr/bin/env bash

echo "=============================="
g++ --version
echo "=============================="
gvm version
echo "=============================="
go version
echo "=============================="

mkdir -p lib
mkdir -p test
mkdir -p bin && rm ./bin/*

make lib
make lib-tests
make binding
make binding-tests
make test