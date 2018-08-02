#!/bin/sh
test -d asm && rm -rf asm
test -f *.out && rm *.out
test -f ./gocc && rm gocc
test -f testfile && rm testfile
