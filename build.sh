#!/bin/bash

GOARCH=amd64
GOOS=$(uname -s | tr [A-Z] [a-z])
if [ "$GOOS" == "darwin" ]; then
    GOBUILD="/usr/local/bin/go build --mod=vendor"
    UPX=""
else
    GOBUILD="/usr/bin/go build --mod=vendor"
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

#buildPlugin rst

rm -f fountain-linux-amd64
$GOBUILD -ldflags="-s -w" -o fountain-linux-amd64 *.go

if [ -e "$UPX" ]; then
    sudo -H $UPX fountain-linux-amd64
fi
