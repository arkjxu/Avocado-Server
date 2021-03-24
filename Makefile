CXX = go
CXXFLAGS = 
CXXENV = 

DISTENV = GOOS=linux GOARCH=amd64
TESTFLAGS = 

SRC = ./src/
BUILD = build
DIST = dist

EXEC = avocado_server

all: run

run:
	$(CXXENV) $(CXX) run $(CXXFLAGS) $(SRC)

build:
	@mkdir -p $(BUILD)
	$(CXXENV) $(CXX) build $(CXXFLAGS) -o $(BUILD)/$(EXEC) $(SRC)

dist:
	@mkdir -p $(DIST)
	$(DISTENV) $(CXX) build $(CXXFLAGS) -o $(DIST)/$(EXEC) $(SRC)

test:
	$(CXXENV) $(CXX) test -v -count=1 $(CXXFLAGS) $(SRC)

clean:
	rm -rf $(DIST) $(BUILD)
	go clean -testcache