#!/usr/bin/env sh

BINARY=/emintd/${BINARY:-emintd}
ID=${ID:-0}
LOG=${LOG:-emintd.log}

if ! [ -f "${BINARY}" ]; then
	echo "The binary $(basename "${BINARY}") cannot be found. Please add the binary to the shared folder. Please use the BINARY environment variable if the name of the binary is not 'emintd'"
	exit 1
fi

BINARY_CHECK="$(file "$BINARY" | grep 'ELF 64-bit LSB executable, x86-64')"

if [ -z "${BINARY_CHECK}" ]; then
	echo "Binary needs to be OS linux, ARCH amd64"
	exit 1
fi

export EMINTDHOME="/emintd/node${ID}/emintd"

if [ -d "$(dirname "${EMINTDHOME}"/"${LOG}")" ]; then
  "${BINARY}" --home "${EMINTDHOME}" "$@" | tee "${EMINTDHOME}/${LOG}"
else
  "${BINARY}" --home "${EMINTDHOME}" "$@"
fi
