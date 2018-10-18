#!/bin/sh

DIR=`dirname "$0"`
FILE=bundled.go
BIN=`go env GOPATH`/bin

cd $DIR
rm $FILE

$BIN/bundler -package bugs -name codeDark code-dark.svg > $FILE
$BIN/bundler -package bugs -name codeLight -append code-light.svg >> $FILE
$BIN/bundler -package bugs -name bugDark -append bug-dark.svg >> $FILE
$BIN/bundler -package bugs -name bugLight -append bug-light.svg >> $FILE
$BIN/bundler -package bugs -name flagDark -append flag-dark.svg >> $FILE
$BIN/bundler -package bugs -name flagLight -append flag-light.svg >> $FILE

