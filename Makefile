WHISPER_INCLUDE_DIR := /Users/asparkoffire/.bin/whisper.cpp/include
WHISPER_LIBRARY_DIR := /Users/asparkoffire/.bin/whisper.cpp/lib

build:
	@mkdir -p bin
	@CGO_CFLAGS="-I$(WHISPER_INCLUDE_DIR)" CGO_LDFLAGS="-L$(WHISPER_LIBRARY_DIR)" go build -o bin/whisper main.go

run: build
	@DYLD_LIBRARY_PATH="$(WHISPER_LIBRARY_DIR)" ./bin/whisper
