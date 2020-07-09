#!/usr/bin/env bash

echo "=============================="
g++ --version
echo "=============================="
gvm version
echo "=============================="
go version
echo "=============================="

mkdir -p lib
mkdir -p bin

find bin/* -print0  | xargs -0  rm -rf
find lib/* -print0  | xargs -0  rm -rf

make lib