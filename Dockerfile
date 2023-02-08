FROM golang:1.19 as buildbase

WORKDIR /go/src/gitlab.com/q-dev/q-id/issuer
COPY . .

ARG BUILD_TOKEN
ARG USERNAME=oauth2
RUN git config --global url."https://${USERNAME}:${BUILD_TOKEN}@gitlab.com/".insteadOf https://gitlab.com/
RUN echo "machine gitlab.com login gitlab-ci-token password ${BUILD_TOKEN}" > ~/.netrc
RUN go env -w GOPRIVATE=gitlab.com/*
RUN go env -w GONOSUMDB=gitlab.com/*
RUN go env -w GONOPROXY=gitlab.com/*


RUN CGO_ENABLED=1 GOOS=linux go build -o /usr/local/bin/issuer /go/src/gitlab.com/q-dev/q-id/issuer

FROM ubuntu:20.04

WORKDIR /

RUN apt-get update && apt-get -y install libgomp1 && apt-get -y install ca-certificates
COPY --from=buildbase /go/pkg/mod/github.com/wasmerio/wasmer-go@v1.0.4/wasmer/packaged/lib/linux-amd64/libwasmer.so /go/pkg/mod/github.com/wasmerio/wasmer-go@v1.0.4/wasmer/packaged/lib/linux-amd64/libwasmer.so

COPY --from=buildbase /usr/local/bin/issuer /usr/local/bin/issuer
COPY --from=buildbase /go/src/gitlab.com/q-dev/q-id/issuer/circuits /circuits

ENTRYPOINT ["issuer"]
