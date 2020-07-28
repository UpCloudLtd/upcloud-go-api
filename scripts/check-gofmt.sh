#!/usr/bin/env bash

DIR="."
if [[ ! -z "$1" ]]; then
  DIR="$1"
fi

go_files=$(gofmt -l -s $(find . -type f -name \*.go))
if [[ ! -z "$go_files" ]]; then
  echo "The following files need formatting:"
  echo $go_files
  exit 1
fi

