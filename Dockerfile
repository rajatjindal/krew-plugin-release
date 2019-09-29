# Container image that runs your code
FROM golang:1.12

WORKDIR /go/src/github.com/rajatjindal/krew-plugin-release-go

COPY . .

RUN go build -o bin/krew-plugin-release-go main.go

ENTRYPOINT [ "/go/src/github.com/rajatjindal/krew-plugin-release-go/bin/krew-plugin-release-go" ]
