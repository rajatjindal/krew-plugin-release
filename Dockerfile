# Container image that runs your code
FROM golang:1.12 as builder

WORKDIR /go/src/github.com/rajatjindal/krew-plugin-release

COPY . .

RUN go build -o bin/krew-plugin-release main.go

FROM alpine:latest

COPY --from=builder /go/src/github.com/rajatjindal/krew-plugin-release/bin/krew-plugin-release /usr/local/bin/krew-plugin-release
ENTRYPOINT [ "krew-plugin-release" ]
