FROM golang:alpine AS build-env

# Install minimum necessary dependencies
ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python3
RUN apk add --no-cache $PACKAGES

# Set working directory for the build
WORKDIR /go/src/github.com/ChainSafe/ethermint

# Add source files
COPY . .

# build Cosmos SDK, remove packages
RUN make build-ethermint && \
  cp ./build/emint* /go/bin


# Final image
FROM alpine:edge

# Install ca-certificates
RUN apk add --update ca-certificates
WORKDIR /root

# Copy over binaries from the build-env
COPY --from=build-env /go/src/github.com/Chainsafe/ethermint/build/emintd /usr/bin/emintd
COPY --from=build-env /go/src/github.com/Chainsafe/ethermint/build/emintcli /usr/bin/emintcli

EXPOSE 26656 26657 1317

# Run emintd by default, omit entrypoint to ease using container with emintcli
CMD ["emintd"]
