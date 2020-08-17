FROM golang:1.14

RUN apt-get update && apt-get install -y \
  make curl jq tmux

RUN apt-get update && apt-get install -f -y \
  npm protobuf-compiler

# Install dependencies
RUN apk add --update $PACKAGES
RUN apk add linux-headers

RUN npm install -g solc
RUN mv /usr/local/bin/solcjs /usr/local/bin/solc

COPY . /ethermint

WORKDIR /ethermint

RUN make install

ENTRYPOINT ["/bin/bash"]
