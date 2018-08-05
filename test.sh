#!/bin/sh

ASM=asm
OUT=out
TESTFILE=testfile

RED='\033[0;31m'
SET='\033[0m'

go build .
if [ $? -ne 0 ]; then
  exit 1
fi


if [ ! -d asm ]; then
  mkdir asm
fi

TEST_NUM=0

expect() {
  TEST_NUM=`expr $TEST_NUM + 1`
  echo "${1}" > $TESTFILE
  ./gocc -o "${ASM}/${TEST_NUM}.s" $TESTFILE || return
  cc "${ASM}/${TEST_NUM}.s" -o $OUT
  ./$OUT
  res=$?
  if [ $res -eq $2 ]; then
    echo "[OK test${TEST_NUM}] '${1}' => ${res}"
  else
    echo ${RED}
    echo "[Failed test${TEST_NUM}] '${1}' expected ${2}, but got ${res}"
    echo ${SET}
  fi
}

expect "1" 1
expect "2" 2

expect "1+2" 3
expect "1+2+3+4" 10
expect "5-8+10-2" 5
expect "2-1" 1

expect "3*4" 12
expect "5/3" 1
expect "4%3" 1
expect "(3 * 4) % 2" 0

expect "5 * 3 - 4" 11
expect "2 + 3 * 4 - 6 / 3" 12
expect "(4 * 5 / 2 + 4) * 3 - 1" 41

expect "int a = 4;" 4
expect "int a = 2 * (5 + 10 / 2);" 20
expect "int a = 3; int b = a + 4;" 7
expect "int a = 3; int b = 4 + a;" 7
expect "int a = 5; int b = a + 8; int c = a + b + 12;" 30


rm $TESTFILE
rm $OUT
