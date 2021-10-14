# Makefile is esaier for me - i am sure others can do this better
#

pkg_dir = $(GOPATH)/pkg/darwin_amd64
bin_dir = $(GOPATH)/bin
sunname = sunrisesunset

LOGMSG = $(pkg_dir)/github.com/seldonsmule/logmsg.a

SUNBIN = $(bin_dir)/$(sunname)

all: $(SUNBIN)

clean:
	go clean
	echo $(SUNBIN)
	rm -f $(SUNBIN)

$(SUNBIN): sunriseset.go
	go build -o $(sunname)
	go install

$(LOGMSG): 
	make deps

deps:
	go get github.com/seldonsmule/logmsg
	go get github.com/seldonsmule/restapi
	go get github.com/mattn/go-sqlite3
	go get github.com/denisbrodbeck/machineid
	go get golang.org/x/crypto/ssh/terminal

rmdeps:
	@rm -f $(LOGMSG)
	@rm -rf $(GOPATH)/src/github.com/seldonsmule/logmsg
	@rm -f $(pkg_dir)/github.com/seldonsmule/restapi.a
	@rm -rf $(GOPATH)/src/github.com/seldonsmule/restapi
	@rm -f $(pkg_dir)/github.com/mattn/go-sqlite3.a
	@rm -rf $(GOPATH)/src/github.com/mattn/go-sqlite3
	@rm -f $(pkg_dir)/github.com/denisbrodbeck/machineid.a
	@rm -rf $(GOPATH)/src/github.com/denisbrodbeck/machineid

usage:
	@echo make all - builds application
	@echo make clean - cleans out all builds


default: usage
