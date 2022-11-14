#!/bin/bash
#cd $(dirname $(readlink -m $BASH_SOURCE))

echo Building and compiling aad-auth-proxy

# Install golang.
# https://golang.org/doc/install#install
wget -nv -c https://go.dev/dl/go1.19.2.linux-amd64.tar.gz
[ -d /usr/local/go ] && rm -rf /usr/local/go
tar -C /usr/local -xzf go1.19.2.linux-amd64.tar.gz
# Installing tests dependencies.
export PATH=$PATH:/usr/local/go/bin
export GOBIN=/usr/local/go/bin

echo Restoring packages

go get -v -t -d
if [ $? != 0 ]; then
  printf "Error : [%d] when executing command: go get" $?
  exit $?
fi

# Note: Building in ubuntu machine, but using alpine image, so we need to statically link binary to libraries
# https://stackoverflow.com/questions/58205781/dockerfile-error-standard-init-linux-go207-exec-user-process-caused-no-such
env CGO_ENABLED=0 go build -o main -ldflags "-w -s -v"

if [ $? != 0 ]; then
  printf "Error building main"
  exit 1
fi

tar cvf aad-auth-proxy.tar main

if [ $? != 0 ]; then
  printf "Error packing: aad-auth-proxy.tar"
  exit 1
fi
