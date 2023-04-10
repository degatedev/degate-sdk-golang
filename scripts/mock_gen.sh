#!/usr/bin/env bash

go generate ./...

source_path=degate
dest_path=mocks

mock_file_in_dir() {
  for filename in ${1}/*
  do
    if [ -d "$filename" ]; then
      echo "gen dir $filename"
      mock_file_in_dir $filename
    elif [ -f "$filename" ]; then
      if [ "${filename##*.}" == "go" ]; then
        echo "gen file $filename"
        mockgen -source=${filename} -destination=${filename/${source_path}/${dest_path}}
      fi
    fi
  done
}

# clear existing mocks
rm -r ${dest_path}
# generate new mocks
mock_file_in_dir ${source_path}
# generate mock repo
#mockgen --build_flags=--mod=mod -destination=internal/pkg/domain/mocks/domain.go git.bybit.com/gtd/micro/gridbot/internal/pkg/domain Repository