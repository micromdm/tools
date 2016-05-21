GO    		:= go
glide    	:= glide
release    	:= ./release.sh
.PHONY: all clean appmanifest poke

all: clean buildall

clean: 
	rm -rf ./build/*

buildall: appmanifest poke

appmanifest: 
	rm -rf ./build/appmanifest
	@echo ">> building appmanifest"
	cd ./appmanifest && $(release)

poke: 
	rm -rf ./build/poke
	@echo ">> building poke"
	cd ./poke && $(release)
