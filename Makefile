run:
	weaver generate ./...
	go build -o books .
	weaver multi deploy weaver.toml