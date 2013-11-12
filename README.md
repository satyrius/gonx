# gonx [![Build Status](https://travis-ci.org/satyrius/gonx.png)](https://travis-ci.org/satyrius/gonx)

`gonx` is Nginx access log reader library for `Go`. In fact you can use it for any format.

## Usage

The library provides `Reader` type and two constructors for it.

Common constructor `NewReader` gets opened file (`io.Reader`) and log format (`string`) as argumets. format is in form os nginx `log_format` string.

```go
reader := gonx.NewReader(file, format)
```

`NewNginxReader` provides mo magic. It gets nginx config file (`io.Reader`) as second argument and `log_format` name (`string`) a third.

```go
reader := gonx.NewNginxReader(file, nginxConfig, format_name)
```

`Reader` implements `io.Reader`. Here is example usage

```go
for {
	rec, err := reader.Read()
	if err == io.EOF {
		break
	}
	// Process the record... e.g.
}
```

See more examples in `example/*.go` sources.

## Format

As I said above this library is primary for nginx access log parsing, but it can be configured to parse any other format. `NewReader` accepts `format` argument, it will be transformed to regular expression and used for log line by line parsing. Format is nginx-like, here is example

	`$remote_addr [$time_local] "$request"`

It should contain variables inn form `$name`. The regular expression will be created using this string format representation

	`^(?P<remote_addr>[^ ]+) \[(?P<time_local>[^]]+)\] "(?P<request>[^"]+)"$`

If log line does not match this format, the `Reader.Read` returns an `error`. Otherwise you will get the record of type `Entry` (which is customized `map[string][string]`) with `remote_addr`, `time_local` and `request` keys filled with parsed values.

## Stability

This library API and internal representation can be changed at any moment, but I guarantee that backward capability will be supported for the following public interfaces.

* `func NewReader(logFile io.Reader, format string) *Reader`
* `func NewNginxReader(logFile io.Reader, nginxConf io.Reader, formatName string) (reader *Reader, err error)`
* `func (r *Reader) Read() (record Entry, err error)`

## Changelog

All major changes will be noticed in [release notes](https://github.com/satyrius/gonx/releases).

## Contributing

Fork the repo, create a feature branch then send me pull request. Feel free to create new issues or contact me using email.
