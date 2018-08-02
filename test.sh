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
  if [ $? -eq $2 ]; then
    echo "OK ${1}"
  else
    echo "failed ${1}"
  fi
}

expect "1" 1
expect "2" 2

expect "1+2" 3

rm $TESTFILE
rm -rf $ASM
rm $OUT
