#!/usr/bin/env bash
check_dependency() {
    command -v curl >/dev/null 2>&1 || { echo >&2 "Error: curl is not installed. Please install it first."; exit 1; }
    command -v jq >/dev/null 2>&1 || { echo >&2 "Error: jq is not installed. Please install it first."; exit 1; }
    command -v mpv >/dev/null 2>&1 || { echo >&2 "Error: mpv is not installed. Please install it first."; exit 1; }
    command -v fzf >/dev/null 2>&1 || { echo >&2 "Error: fzf is not installed. Please install it first."; exit 1; }
}

check_dependency

# Run a While loop to continuously prompt for search queries and append that to $1 file as a playlist
# If user Ctrl+C or Ctrl+D, exit the loop and save the playlist to $1 file

echo "Press Ctrl+C or Ctrl+D to exit the loop and save the playlist."
read -p "Press any key to continue: "

while true; do
    if [ -z "$1" ]; then
        echo "No playlist file specified. Exiting."
        break
    fi

    # Search for tracks and select one using fzf
    track_id=$(curl "127.0.0.1:8080/tracks/all" |jq -r '.[] | "\(.id)::\(.title)::\(.file_path)"'| fzf | awk -F'::' '{print $1}')

    if [ -z "$track_id" ]; then
        echo "No track selected."
        break
    fi
    # Append the selected track to the playlist file
    echo $track_id >> "$1"
  done
