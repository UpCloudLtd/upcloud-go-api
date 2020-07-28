#!/usr/bin/env bash

go_files=$(go mod tidy -v 2>&1)
if [[ ! -z "$go_files" ]]; then
  echo "The following modules need tidied."
  echo $go_files
  echo "If run locally they will have been tidied for you."
  exit 1
fi

