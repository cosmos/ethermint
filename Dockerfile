FROM golang:1.14
RUN apt-get update && apt-get install -y \
  make curl jq tmux vim 
RUN apt-get update && apt-get install -f -y \
  npm protobuf-compiler
COPY . $HOME/ethermint
WORKDIR $HOME/ethermint
RUN make install
ENTRYPOINT ["/bin/bash"]