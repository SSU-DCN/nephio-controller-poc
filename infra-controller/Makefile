.PHONY: all
all:
	@echo "**********************************************************"
	@echo "**                    Provider API Service Builder*     **"
	@echo "**********************************************************"


.PHONY: run
run:
	go run .
.PHONY: image
image:
	docker build . -t ntnguyencse/infra-controller:lastest
.PHONY: push-image
image:
	docker push ntnguyencse/infra-controller:lastest