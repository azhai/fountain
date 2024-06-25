@ECHO OFF

del fountain.exe
go.exe build -ldflags="-s -w" -o fountain.exe main.go
