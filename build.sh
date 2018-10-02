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
    FLAGS="-s -w -pluginpath=$NAME.so.$MOMENT"
    rm -f "$NAME.so"
    $GOBUILD -buildmode=plugin -ldflags="$FLAGS" -o "$NAME.so" converter/"$NAME.go"
}

#buildPlugin markdown
rm -f fountain
$GOBUILD -ldflags="-s -w"
#upx --ultra-brute fountain
