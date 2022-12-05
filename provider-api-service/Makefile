.PHONY: all
all:
	@echo "**********************************************************"
	@echo "**                    Provider API Service Builder*     **"
	@echo "**********************************************************"


.PHONY: run
run:
	go run .

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o provider main.go
.PHONY: clean
clean: 
	rm *.txt provider
.PHONY: image
image:
	docker build . -t ntnguyencse/provider-api-service:lastest
.PHONY: push-image
image:
	docker push ntnguyencse/provider-api-service:lastest