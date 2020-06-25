FROM golang:1.13

RUN apt-get update && apt-get install -y \
  make curl jq

COPY . /ethermint

WORKDIR /ethermint

RUN make install

ENTRYPOINT ["/bin/bash"]
