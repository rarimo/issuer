FROM golang:1.19 as buildbase

WORKDIR /go/src/gitlab.com/rarimo/identity/issuer
COPY . .

ARG BUILD_TOKEN
ARG USERNAME=oauth2
RUN git config --global url."https://${USERNAME}:${BUILD_TOKEN}@gitlab.com/".insteadOf https://gitlab.com/
RUN echo "machine gitlab.com login gitlab-ci-token password ${BUILD_TOKEN}" > ~/.netrc
RUN go env -w GOPRIVATE=gitlab.com/*
RUN go env -w GONOSUMDB=gitlab.com/*
RUN go env -w GONOPROXY=gitlab.com/*

RUN CGO_ENABLED=1 GOOS=linux go build -o /usr/local/bin/issuer /go/src/gitlab.com/rarimo/identity/issuer

FROM alpine:3.16.0

RUN apk --update add --no-cache libstdc++ gcompat libgomp ca-certificates

WORKDIR /

COPY --from=buildbase "/go/pkg/mod/github.com/wasmerio/wasmer-go@v1.0.4/wasmer/packaged/lib/linux-amd64/libwasmer.so" \
"/go/pkg/mod/github.com/wasmerio/wasmer-go@v1.0.4/wasmer/packaged/lib/linux-amd64/libwasmer.so"

COPY --from=buildbase /usr/local/bin/issuer /usr/local/bin/issuer
COPY --from=buildbase /go/src/gitlab.com/rarimo/identity/issuer/circuits /circuits

ENTRYPOINT ["issuer"]
