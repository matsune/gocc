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
    echo "=> ${res} ${GREEN}[OK]${CLEAR}"
  else
    echo ${RED}
    echo "[Failed]'${1}' expected ${2}, but got ${res}"
    echo ${CLEAR}
  fi
}

expect "int main() { 1; }" 0
expect "int main() { return 1; }" 1
expect "int main() { return 2; }" 2

expect "int main() { return 1+2+3+4; }" 10
expect "int main() { return 5-8+10-2; }" 5

expect "int main() { return 3*4; }" 12
expect "int main() { return 5/3; }" 1
expect "int main() { return 4%3; }" 1
expect "int main() { return (3 * 4) % 2; }" 0

expect "int main() { return (4 * 5 / 2 + 4) * 3 - 1; }" 41

expect "int main() { int a = 4; return a; }" 4
expect "int main() { int a = 2 * (5 + 10 / 2); return a; }" 20
expect "int main() { int a = 3; int b = a + 4; return b; }" 7
expect "int main() { int a = 3; int b = 4 + a; return b; }" 7
expect "int main() { int a = 5; int b = a + 8; int c = a + b + 12; return c; }" 30
expect "int main() { int a = 1 + 2; int b = 3 + a + 5; int c = 10 + a * b; return c; }" 43

expect "int main() { return 1 + 2; }" 3

expect "int a() { return 3; } int main() { return a(); }" 3
expect "int a() { return 3 + 4 * 2; } int main() { return a(); }" 11

expect "int sum(int a, int b) { return a + b; } int main() { return sum(1, 2); }" 3
expect "int add10(int a) { return a + 10; } int main() { int a = 5; return add10(a); }" 15
expect "int sum(int a, int b, int c, int d, int e, int f, int g, int h) { return a+b+c+d+e+f+g+h; } int main() { int s = sum(1,2,3,4,5,6,7,8); return s * 5; }" 180

rm $TESTFILE
rm $OUT
