FROM golang:1.13-alpine3.11 AS builder

ARG NETRC_ENABLED="false"
ARG NETRC_MACHINE="github.com"
ARG NETRC_LOGIN
ARG NETRC_PASSWORD

RUN test "$NETRC_ENABLED" && printf "machine ${NETRC_MACHINE}\nlogin ${NETRC_LOGIN}\npassword ${NETRC_PASSWORD}\n" >> /root/.netrc \
	&& chmod 600 /root/.netrc

RUN apk --update add \
		ca-certificates \
		gcc \
		git \
		musl-dev \
		tzdata

RUN echo 'nobody:x:65534:65534:nobody:/:' > /tmp/passwd \
	&& echo 'nobody:x:65534:' > /tmp/group

COPY go.mod go.sum /go/src/github.com/juli3nk/doorbell/
WORKDIR /go/src/github.com/juli3nk/doorbell

ENV GO111MODULE on
RUN go mod download

COPY . .

RUN go build -ldflags "-linkmode external -extldflags -static -s -w" -o /tmp/doorbell


FROM scratch

COPY --from=builder /tmp/group /tmp/passwd /etc/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /tmp/doorbell /doorbell

USER nobody:nobody

ENTRYPOINT ["/doorbell"]
