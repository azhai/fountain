@ECHO OFF

del fountain-win7-x64.exe
go build -mod=vendor -ldflags="-s -w" -o fountain-win7-x64.exe main.go

PAUSE
