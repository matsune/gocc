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

TEST_NUM=0

expect() {
  TEST_NUM=`expr $TEST_NUM + 1`
  echo "[test${TEST_NUM}] ${1}"
  echo "${1}" > $TESTFILE
  ./gocc -o "${ASM}/${TEST_NUM}.s" $TESTFILE || return
  cc -m32 "${ASM}/${TEST_NUM}.s" -o $OUT
  ./$OUT
  res=$?
  if [ $res -eq $2 ]; then
    echo "${GREEN}=> ${res} [OK]${CLEAR}"
  else
    echo ${RED}
    echo "[Failed]'${1}' expected ${2}, but got ${res}"
    echo ${CLEAR}
  fi
}

expect "int main() { 1; }" 1
expect "int main() { 2; }" 2

expect "int main() { 1+2+3+4; }" 10
expect "int main() { 5-8+10-2; }" 5

expect "int main() { 3*4; }" 12
expect "int main() { 5/3; }" 1
expect "int main() { 4%3; }" 1
expect "int main() { (3 * 4) % 2; }" 0

expect "int main() { (4 * 5 / 2 + 4) * 3 - 1; }" 41

expect "int main() { int a = 4; }" 4
expect "int main() { int a = 2 * (5 + 10 / 2); }" 20
expect "int main() { int a = 3; int b = a + 4; }" 7
expect "int main() { int a = 3; int b = 4 + a; }" 7
expect "int main() { int a = 5; int b = a + 8; int c = a + b + 12; }" 30
expect "int main() { int a = 1 + 2; int b = 3 + a + 5; int c = 10 + a * b; }" 43

expect "int main() { return 1 + 2; }" 3

expect "int a() { return 3; } int main() { return a(); }" 3
expect "int a() { return 3 + 4 * 2; } int main() { return a(); }" 11

rm $TESTFILE
rm $OUT
