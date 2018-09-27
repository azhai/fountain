#!/bin/bash

export GOARCH=amd64
export GOOS=$(uname -s | tr [A-Z] [a-z])
if [ "$GOOS" == "darwin" ]; then
    export GOBUILD="/usr/local/bin/go build"
else
    export GOBUILD="/usr/bin/go build"
fi

buildPlugin()
{
    NAME="$1"
    MOMENT=$(date +%FT%TZ)
    FLAGS="-pluginpath=$NAME.so.$MOMENT"
    $GOBUILD -buildmode=plugin --ldflags="$FLAGS" -o "$NAME.so" converter/"$NAME.go"
}

rm -f fountain *.so
buildPlugin markdown
$GOBUILD -ldflags "-w -s"
