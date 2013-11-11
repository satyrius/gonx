# gonx

`gonx` is Nginx access log reader libriry for `Go`. In fact you can use it for any format.

## Usage

The library provides `Reader` type and two constructors for it.

Common constructor `NewReader` gets opened file (`io.Reader`) and log format (`string`) as argumets. format is in form os nginx `log_format` string.
	
	reader := gonx.NewReader(file, format)
	
`NewNginxReader` provides mo magic. It gets nginx config file (`io.Reader`) as second argument and `log_format` name (`string`) a third.

	reader := gonx.NewNginxReader(file, nginxConfig, format)

`Reader` implements `io.Reader`. Here is example usage

	for {
		rec, err := reader.Read()
		if err == io.EOF {
			break
		}
		// Process the record... e.g.
	}

See more examples in `example/*.go` sources.