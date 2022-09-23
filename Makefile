.PHONY: build deploy

# functions
FUNC_BOT = func/bot
FUNC_API = func/api


build: clean generate
	cd $(FUNC_BOT) && go mod tidy && go mod vendor && go build
	cd $(FUNC_API) && go mod tidy && go mod vendor && go build

clean:
	rm -rf $(FUNC_BOT)/vendor

prepare:
	cd pkg && \
		go install github.com/golang/mock/mockgen@v1.6.0 && \
		go install github.com/haveyoudebuggedit/gotestfmt/v2/cmd/gotestfmt@latest

generate:
	cd pkg && rm -rf mock && go generate ./...

test: generate
	cd pkg && PHASE=test gotest -v ./...

ci-test: prepare build generate
	cd pkg && PHASE=test go test -json ./... | tee /tmp/gotest.log | gotestfmt

develop:
	cd tool/develop && PHASE=local go run main.go

develop-fe:
	cd front && npm run dev

deploy: build
	# copy static resources
	rm -rf func/bot/images || true
	cp -r pkg/images func/bot/
	# deploy
	cd deploy && terraform apply

ngrok:
	ngrok http 8080
