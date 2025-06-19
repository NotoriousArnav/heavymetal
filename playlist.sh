#!/usr/bin/env bash
check_dependency() {
    command -v curl >/dev/null 2>&1 || { echo >&2 "Error: curl is not installed. Please install it first."; exit 1; }
    command -v jq >/dev/null 2>&1 || { echo >&2 "Error: jq is not installed. Please install it first."; exit 1; }
    command -v mpv >/dev/null 2>&1 || { echo >&2 "Error: mpv is not installed. Please install it first."; exit 1; }
    command -v fzf >/dev/null 2>&1 || { echo >&2 "Error: fzf is not installed. Please install it first."; exit 1; }
}

check_dependency

playlist_file=$1
if [[ -z "$playlist_file" ]]; then
    echo "Usage: $0 <playlist_file>"
    exit 1
fi

EXTRA_ARGS=${2:-}

for track in $(cat $playlist_file); do
    mpv \
      $EXTRA_ARGS \
      http://127.0.0.1:8080/stream/$track
done
