export GO111MODULE=on

GOCMD=go
MAIN_DIR=.
PACKAGE_NAME="."

app-macos:
	GOOS=darwin $(GOCMD) build -i -o ./build/app_macos $(PACKAGE_NAME)

app-linux:
	GOOS=linux $(GOCMD) build -i -o ./build/app_linux $(PACKAGE_NAME)

app-win:
	GOOS=windows $(GOCMD) build -i -o ./build/app_linux $(PACKAGE_NAME)

run-macos: app-macos
	./build/app_macos