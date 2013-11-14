test: deps
	go test -v .

deps:
	go get github.com/stretchr/testify

dev-deps:
	go get github.com/nsf/gocode
	go get code.google.com/p/rog-go/exp/cmd/godef
