#!/bin/sh
ASM=asm
OUT=out
TESTFILE=testfile

go build .

if [ ! -d asm ]; then
  mkdir asm
fi

TEST_NUM=0

expect() {
  TEST_NUM=`expr $TEST_NUM + 1`
  echo $1 > $TESTFILE
  ./gocc -o "${ASM}/${TEST_NUM}.s" $TESTFILE
  cc "${ASM}/${TEST_NUM}.s" -o $OUT
  ./$OUT
  res=$?
  if [ $res -eq $2 ]; then
    echo "OK ${1}"
  else
    echo "expected ${1}, but got ${res}"
  fi
}

expect "1" 1
expect "2" 2

expect "1+2" 3
expect "1+2+3+4" 10
expect "5-8+10-2" 5
expect "2-1" 1

expect "3*4" 12

rm $TESTFILE
rm $OUT
