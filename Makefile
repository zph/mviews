run:
	bash -c "go generate src/mviews/main.go"
	bash -c "go run src/mviews/*.go"

generate:
	bash -c "go generate src/mviews/main.go"

build: $(clean)
	bash -c "go generate src/mviews/main.go"
	go build -o bin/mviews src/mviews/*.go

clean:
	bash -c "rm -rf bin/mviews src/mviews/generated_*.go"
