FROM golang:1.13

RUN apt-get update && apt-get install -y \
  make curl jq tmux vim

COPY . /ethermint

WORKDIR /ethermint

RUN make install
RUN cd /ethermint/docker/benchmarking && go get && go build

ENTRYPOINT ["/bin/bash"]
