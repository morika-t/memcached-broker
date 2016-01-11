build: ./app
	go build main.go

./app:
	./scripts/generate-app
