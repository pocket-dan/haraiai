.PHONY: build deploy

# functions
FUNC_BOT = func/bot


build: clean generate
	cd $(FUNC_BOT) && go mod vendor && go build

clean:
	rm -rf $(FUNC_BOT)/vendor

prepare:
	go install github.com/rakyll/gotest

generate:
	cd pkg && rm -rf mock && go generate ./...

test: generate
	cd pkg && PHASE=test gotest -v ./...

develop:
	PHASE=local go run tool/develop/main.go

deploy: build
	cd deploy && terraform apply
	cd front && npm run build && firebase deploy

ngrok:
	ngrok http 8080
