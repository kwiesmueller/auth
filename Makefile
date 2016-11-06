all: test install run
install:
	GOBIN=$(GOPATH)/bin GO15VENDOREXPERIMENT=1 go install bin/auth_server/*.go
test:
	GO15VENDOREXPERIMENT=1 go test -cover `glide novendor`
vet:
	go tool vet .
	go tool vet --shadow .
lint:
	golint -min_confidence 1 ./...
errcheck:
	errcheck -ignore '(Close|Write)' ./...
check: lint vet errcheck
runledis:
	ledis-server \
	-addr=localhost:5555 \
	-databases=1
run:
	auth_server \
	-logtostderr \
	-v=2 \
	-port=6666 \
	-ledisdb-address=localhost:5555 \
	-auth-application-password=test123 \
	-prefix "/auth"
open:
	open http://localhost:6666/auth/
format:
	find . -name "*.go" -exec gofmt -w "{}" \;
	goimports -w=true .
prepare:
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/Masterminds/glide
	go get -u github.com/siddontang/ledisdb/cmd/ledis-server
	go get -u github.com/golang/lint/golint
	go get -u github.com/kisielk/errcheck
	glide install
update:
	glide up
clean:
	rm -rf var vendor
