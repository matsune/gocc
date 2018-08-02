#!/bin/sh
go build -o gocc main.go

expect() {
  echo $1 > testfile
  ./gocc testfile
  cc main.s
  ./a.out
  if [ $? -eq $1 ]; then
    echo "OK ${1}"
  else
    echo "failed ${1}"
  fi
}

expect "1" 1
expect "2" 2
