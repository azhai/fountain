#!/bin/bash

GOARCH=amd64
GOOS=$(uname -s | tr [A-Z] [a-z])
if [ "$GOOS" == "darwin" ]; then
    GOBUILD="/usr/local/bin/go build"
    UPX=""
else
    GOBUILD="/usr/bin/go build"
    UPX="/usr/bin/upx"
    #UPX="/usr/bin/upx --ultra-brute"
fi

buildPlugin()
{
    NAME="$1"
    MOMENT=$(date +%FT%TZ)
    FLAGS="-s -w -pluginpath=$NAME.so.$MOMENT"
    rm -f "$NAME.so"
    $GOBUILD -buildmode=plugin -ldflags="$FLAGS" -o "$NAME.so" converter/"$NAME.go"
}

#buildPlugin markdown
rm -f fountain
$GOBUILD -ldflags="-s -w"

if [ -e "$UPX" ]; then
    $UPX fountain
fi
