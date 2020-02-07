all: build
.PHONY: all

build:
	go build -mod=vendor -trimpath .

push:
	docker build -t mfojtik/bugtrend:v0.0.2 && \
	docker push mfojtik/bugtrend:v0.0.2