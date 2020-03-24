@ECHO OFF

del fountain-windows-x64.exe
go.exe build -ldflags="-s -w" -o fountain-windows-x64.exe .
