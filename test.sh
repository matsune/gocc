#!/bin/sh
ASM=asm
OUT=a.out
TESTFILE=testfile

RED='\033[0;31m'
GREEN='\033[0;32m'
CLEAR='\033[0m'

go build .
if [ $? -ne 0 ]; then
  exit 1
fi

if [ ! -d asm ]; then
  mkdir asm
fi

alert() {
  echo "${RED}${1}${CLEAR}"
}

PASSED=0
FAILED=0

# test [test number] [expect]
test() {
  if [ -z $1 ]; then
    alert "Test number is empty."
    exit 1
  fi
  echo "[test ${1}]"

  if [ -z $2 ]; then
    alert "Expect value is empty."
    exit 1
  fi
  FILE="test/${1}.c"
  if [ ! -e $FILE ]; then
    alert "${FILE} does not exist."
    exit 1
  fi
  ASM_FILE="${ASM}/${1}.s"
  ./gocc -o $ASM_FILE $FILE || return
  gcc -m32 $ASM_FILE -o $OUT
  ./$OUT
  res=$?
  cat $FILE
  if [ $res -eq $2 ]; then
    echo "=> ${res} ${GREEN}[OK]${CLEAR}"
    echo ""
    PASSED=$(( PASSED + 1 ))
  else
    alert "expected ${2}, but got ${res} [Failed]"
    FAILED=$(( FAILED + 1 ))
  fi
}

test 1 0
test 2 2
test 3 2
test 4 0
test 5 41

test 6 4
test 7 20
test 8 7
test 9 30
test 10 43

test 11 11
test 12 3
test 13 15
test 14 180

echo "Finished test."
echo "${GREEN}PASSED: ${PASSED}\t${RED}FAILED: ${FAILED}${CLEAR}"

rm $OUT
