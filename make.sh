#!/bin/bash -e

errEcho() { echo "$@" 1>&2; }

# https://stackoverflow.com/questions/5947742/how-to-change-the-output-color-of-echo-in-linux
red=`tput setaf 1`
green=`tput setaf 2`
yellow=`tput setaf 3`
pink=`tput setaf 5`
blue=`tput setaf 6`
reset=`tput sgr0`

FILENAME="CheckMirror"
VERSION=$(git describe --always)
# -s -w reduce the size of the binary. Add --long to "git describe" to always get the commit sha
LDFLAGS="-X writeameer/checkmirror/core.VERSION=$VERSION -s -w"

build_linux() {
  echo -e "*** Building Linux binary $VERSION in: ${green}$FILENAME${reset}"
  GOOS=linux GOARCH=amd64 go build -tags "netgo" -ldflags "$LDFLAGS" -o $FILENAME *.go
}

build_mac() {
  echo -e "*** Building darwin binary $VERSION in: ${green}$FILENAME${reset}"
  GOOS=darwin GOARCH=amd64 go build -tags "netgo" -ldflags "$LDFLAGS" -o $FILENAME *.go
}

build_windows() {
  echo -e "*** Building Windows binary $VERSION in: ${green}$FILENAME.exe${reset}"
  GOOS=windows GOARCH=amd64 go build -tags "netgo" -ldflags "$LDFLAGS" -o "$FILENAME.exe" *.go
}


errEcho
errEcho "*** ${green}$0${reset} executed with params: ${blue}$1 $2${reset}"

SUBCMD=$1
test "$SUBCMD" = "build-windows" && build_windows
test "$SUBCMD" = "build-linux" && build_linux
test "$SUBCMD" = "build-mac" && build_mac
exit 0
