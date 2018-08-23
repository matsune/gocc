#!/bin/sh
ASM=asm
OUT=a.out
TESTFILE=testfile
APP=app

RED='\033[0;31m'
GREEN='\033[0;32m'
CLEAR='\033[0m'

go build -o $APP .
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
COUNT=0

# test [test number] [expect]
test() {
  COUNT=$(( COUNT + 1 ))
  if [ -z $1 ]; then
    alert "Test number is empty."
    exit 1
  fi
  echo "[test ${1}]"

  if [ -z $2 ]; then
    alert "Expect value is empty."
    exit 1
  fi
  FILE="c/${1}.c"
  if [ ! -e $FILE ]; then
    alert "${FILE} does not exist."
    exit 1
  fi
  ASM_FILE="${ASM}/${1}.s"
  ./$APP -S -o $ASM_FILE $FILE || return
  gcc $ASM_FILE -o $OUT
  ./$OUT
  res=$?
  cat $FILE
  if [ $res -eq $2 ]; then
    echo "=> ${res} ${GREEN}[OK]${CLEAR}"
    echo ""
    PASSED=$(( PASSED + 1 ))
  else
    alert "expected ${2}, but got ${res} [Failed]"
  fi
}

test return_int 2
test add_sub 2
test mul_div 0
test add_sub_mul_div 41

test var_def 4
test var_def_expr 20
test var_def_2 7
test var_def_3 43

test call_no_arg 11
test call_2_args 3
test call_10_args 110

test char 98
test char_arg 195
test char_int_args 118

test pointer 3
test pointer2 20

test array 16

test if_stmt 1
test if_else 10
test if_else_if 1
test if_binaries 5

test inc_dec 11
test inc_pointer 1

test for_stmt 10

echo "Finished test."
FAILED=$(( COUNT - PASSED ))
echo "${GREEN}PASSED: ${PASSED}\t${RED}FAILED: ${FAILED}${CLEAR}"

rm $OUT
