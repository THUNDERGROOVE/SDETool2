all:
	go install github.com/THUNDERGROOVE/SDETool2/args
	go install github.com/THUNDERGROOVE/SDETool2/sde
	go install github.com/THUNDERGROOVE/SDETool2/web
	go install github.com/THUNDERGROOVE/SDETool2/types
	go install github.com/THUNDERGROOVE/SDETool2/market
	go install github.com/THUNDERGROOVE/SDETool2/log
	go build -v
deps:
	go get github.com/mattn/go-sqlite3
	go get github.com/gorilla/mux
	go get github.com/gorilla/handlers
	go get github.com/atotto/clipboard
	go get github.com/lucasb-eyer/go-colorful