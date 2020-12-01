FROM golang:1.15

RUN apt-get update && apt-get install -y \
  make curl jq tmux vim

COPY . /ethermint

WORKDIR /ethermint

RUN make install
RUN cd /ethermint/benchmarking && go get && go build

ENTRYPOINT ["/bin/bash"]
