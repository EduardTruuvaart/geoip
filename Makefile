.PHONY: build test run

build:
	go vet ./...
	GOARCH=arm64 GOOS=linux go build -o bootstrap ./cmd/lambda
	zip geoip.zip bootstrap
	rm -rf bootstrap
	mv geoip.zip ./build/

deploy:
	make build
	sam deploy --guided --no-fail-on-empty-changeset --no-confirm-changeset --region eu-west-1 --profile personal --stack-name geoip-test --template-file ./deployment/api.yml

test:
	go test ./... -v
