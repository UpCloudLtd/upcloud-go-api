#!/usr/bin/env bash

DIR="."
if [[ ! -z "$1" ]]; then
  DIR="$1"
fi

go_files=$(goimports -l $(find . -type f -name \*.go))
if [[ ! -z "$go_files" ]]; then
  echo "The following files need imports fixed:"
  for f in $go_files; do
    echo $f
  done
  exit 1
fi

