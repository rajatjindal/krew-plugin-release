# Container image that runs your code
FROM golang:1.12

WORKDIR /go/src/github.com/rajatjindal/krew-plugin-release

COPY . .

RUN go build -o bin/krew-plugin-release main.go

ENTRYPOINT [ "/go/src/github.com/rajatjindal/krew-plugin-release/bin/krew-plugin-release" ]
