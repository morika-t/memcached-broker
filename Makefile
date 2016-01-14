build: generate
	go build

generate:
	./scripts/generate-app
	counterfeiter storage Storage

test: generate
	ginkgo -r
