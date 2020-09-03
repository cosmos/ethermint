FROM golang:1.14
RUN apt-get update && apt-get install -y \
  make curl jq tmux vim
COPY . /ethermint
WORKDIR /ethermint
RUN make install
ENTRYPOINT ["/bin/bash"]