default:

requirements:
	go get 'github.com/stretchr/testify'

test: requirements
	go test -v .
