FROM golang:stretch as build-env

# Install minimum necessary dependencies
ENV PACKAGES curl make git libc-dev bash gcc
RUN apt-get update && apt-get upgrade -y && \
    apt-get install -y $PACKAGES

# Set working directory for the build
WORKDIR /go/src/github.com/ChainSafe/ethermint

# Add source files
COPY . .

# build Ethermint
RUN make build-ethermint-linux

# Final image
FROM alpine:edge

# Install ca-certificates
RUN apk add --update ca-certificates
WORKDIR /root

# Copy over binaries from the build-env
COPY --from=build-env /go/src/github.com/ChainSafe/ethermint/build/emintd /usr/bin/emintd
COPY --from=build-env /go/src/github.com/ChainSafe/ethermint/build/emintcli /usr/bin/emintcli

EXPOSE 26656 26657 1317

# Run emintd by default, omit entrypoint to ease using container with emintcli
CMD ["emintd"]
