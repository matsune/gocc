#!/bin/sh
ASM=asm
S=$ASM/gocc.s
OUT=out
TESTFILE=testfile

go build .

if [ ! -d asm ]; then
  mkdir asm
fi

expect() {
  echo $1 > $TESTFILE
  ./gocc -o $S $TESTFILE
  cc $S -o $OUT
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

rm $TESTFILE
rm -rf $ASM
rm $OUT
