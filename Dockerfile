# Container image that runs your code
FROM golang:1.12.10 as builder

WORKDIR /go/src/github.com/rajatjindal/krew-plugin-release

COPY . .

RUN GOOS=linux CGO_ENABLED=0 go build  -ldflags "-w -s" -o bin/krew-plugin-release main.go

FROM scratch

COPY --from=builder /go/src/github.com/rajatjindal/krew-plugin-release/bin/krew-plugin-release /usr/local/bin/krew-plugin-release
ENTRYPOINT [ "krew-plugin-release" ]
