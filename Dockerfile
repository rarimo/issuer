FROM golang:1.18-alpine as buildbase

RUN apk add git build-base

WORKDIR /go/src/gitlab.com/q-dev/q-id/qid-issuer
COPY vendor .
COPY . .

RUN GOOS=linux go build  -o /usr/local/bin/qid-issuer /go/src/gitlab.com/q-dev/q-id/qid-issuer


FROM alpine:3.9

COPY --from=buildbase /usr/local/bin/qid-issuer /usr/local/bin/qid-issuer
RUN apk add --no-cache ca-certificates

ENTRYPOINT ["qid-issuer"]
