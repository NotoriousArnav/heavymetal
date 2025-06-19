#!/usr/bin/env bash
check_dependency() {
    command -v curl >/dev/null 2>&1 || { echo >&2 "Error: curl is not installed. Please install it first."; exit 1; }
    command -v jq >/dev/null 2>&1 || { echo >&2 "Error: jq is not installed. Please install it first."; exit 1; }
    command -v mpv >/dev/null 2>&1 || { echo >&2 "Error: mpv is not installed. Please install it first."; exit 1; }
    command -v fzf >/dev/null 2>&1 || { echo >&2 "Error: fzf is not installed. Please install it first."; exit 1; }
}

check_dependency

query=$(echo "$1" | sed 's/ /%20/g')
if [ -z "$query" ]; then
    echo "Usage: $0 <search_query>"
    exit 1
fi

track_id=$(curl "127.0.0.1:8080/search/track/$query" |jq -r '.[] | "\(.id)::\(.title)::\(.file_path)"'| fzf | awk -F'::' '{print $1}')

if [ -z "$track_id" ]; then
    echo "No track selected."
    exit 1
fi

EXTRA_ARGS="$2"

mpv \
  $EXTRA_ARGS \
  http://127.0.0.1:8080/stream/$track_id


