run: main.exe
	.\build\main.exe

main.exe: $(wildcard .\cmd\server\*.go)
	go build -o .\build\main.exe .\cmd\server\main.go