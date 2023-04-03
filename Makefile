APPNAME="simple-ecommerce"

all: test build run

build:
	@echo "Building application..."
	@go build -v -o ${APPNAME} cmd/*.go

run:
	@go run cmd/*.go

run-script-report:
	@go run cmd/*.go run-script-report

run-cron:
	@go run cmd/*.go run-cron

test:
	@go test ./... -cover -race