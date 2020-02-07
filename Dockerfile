FROM golang:1.14rc1

WORKDIR /go/src/github.com/mfojtik/bugtrend
COPY . .
RUN go get -d -v ./... && go install -v ./...
RUN chmod -R 0777 /go

CMD ["bugtrend"]
