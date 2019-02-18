############################
# STEP 1 build executable binary
############################
FROM golang:1.11-alpine3.7 AS builder
# Install git + SSL ca certificates.
# Git is required for fetching the dependencies.
# Ca-certificates is required to call HTTPS endpoints.
RUN apk update && apk add --no-cache git make ca-certificates tzdata && update-ca-certificates

# Create appuser
RUN adduser -D -g '' appuser

RUN go get -d -v github.com/coredns/coredns
WORKDIR /go/src/github.com/coredns/coredns
COPY . plugin/gslb/
RUN echo "gslb:gslb" >> plugin.cfg

#!Not merge to master, checkout pull request https://github.com/coredns/coredns/pull/2505
RUN git fetch origin pull/2505/head:build && git checkout build
# Using go mod.
# RUN go mod download
RUN go get -d -v

# Build the binary
RUN go generate
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/coredns

# Copy config file
COPY db.abcdefgh.fun /go/bin/db.abcdefgh.fun 
COPY Corefile /go/bin/Corefile

############################
# STEP 2 build a small image
############################
FROM scratch

# Import from builder.
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

# Use an unprivileged user.
USER appuser
WORKDIR /app

COPY --from=builder /go/bin/coredns .
COPY --from=builder /go/bin/db.abcdefgh.fun .
COPY --from=builder /go/bin/Corefile .

EXPOSE 53 53/udp
# Run the binary.
CMD [ "./coredns", "-conf", "./Corefile" ]
#ENTRYPOINT []
LABEL Name=coredns-gslb Version=0.0.1
