test:
	go test -v .

bench:
	go test -bench .

deps:
	go get github.com/stretchr/testify

dev-deps:
	go get github.com/nsf/gocode
	go get code.google.com/p/rog-go/exp/cmd/godef
	go install code.google.com/p/rog-go/exp/cmd/godef

examples:
	go run ./example/common/common.go
	go run ./example/nginx/nginx.go
