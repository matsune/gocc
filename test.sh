#!/bin/sh
go build -o gocc main.go 

./gocc 1
cc main.s
./a.out
if [ $? -eq 1 ]; then
  echo OK
else
  echo failed
  exit 1
fi

./gocc 2
cc main.s
./a.out
if [ $? -eq 2 ]; then
  echo OK
else
  echo failed
  exit 1
fi
