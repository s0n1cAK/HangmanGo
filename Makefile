build:
	go build -o bin/Hangman main.go

run:
	go run main.go

compile:
	echo "Compiling for every OS and Platform"
	GOOS=freebsd GOARCH=386 go build -o bin/Hangman-freebsd main.go
	GOOS=linux GOARCH=386 go build -o bin/Hangman-linux main.go
#	GOOS=windows GOARCH=386 go build -o bin/Hangman-windows main.go
