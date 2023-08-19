.PHONY: build deploy

# functions
FUNC_BOT = func/bot


build: clean prepare generate
	cd $(FUNC_BOT) && go mod tidy && go mod vendor && go build

clean:
	rm -rf $(FUNC_BOT)/vendor
	rm -rf $(FUNC_BOT)/messaging
	rm -rf pkg/mock

prepare:
	cd pkg \
		&& go install go.uber.org/mock/mockgen@latest \
		&& go install github.com/gotesttools/gotestfmt/v2/cmd/gotestfmt@latest \
		&& go install github.com/google/wire/cmd/wire@latest

generate:
	cd pkg \
		&& GOFLAGS=-mod=mod go generate ./...

test: generate
	cd pkg && PHASE=test gotest -v ./... $(ARGS)

ci-test: prepare build
	cd pkg && PHASE=test go test -json ./... | tee /tmp/gotest.log | gotestfmt

develop:
	cd tool/develop && PHASE=local go run main.go

develop-fe:
	cd front && npm run dev

before-deploy:
	# copy static resources
	rm -rf $(FUNC_BOT)/images $(FUNC_BOT)/flexmessage || true
	mkdir -p $(FUNC_BOT)/messaging/static && cp -r pkg/messaging/static $(FUNC_BOT)/messaging/

deploy: build before-deploy
	cd deploy && terraform apply

ngrok:
	ngrok http 8080
