run: main.exe
	.\build\main.exe

main.exe: .\cmd\server\main.go
	go build -o .\build\main.exe .\cmd\server\main.go