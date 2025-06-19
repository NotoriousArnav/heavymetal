#!/usr/bin/env bash

# Set variables
SONG="$1"
ARTIST="$2"

if [[ -z "$SONG" || -z "$ARTIST" ]]; then
  echo "Usage: $0 <song> <artist>"
  exit 1
fi

QUERY="$(printf "%s %s" "$SONG" "$ARTIST" | jq -s -R @uri)"

# cURL request to search endpoint
ID=$(curl -s -G "https://lrclib.net/api/search" \
  --data-urlencode "track_name=$SONG" \
  --data-urlencode "artist_name=$ARTIST" \
  -H "Accept: application/json" \
  -o - | jq '.[0].id' -)

echo "🔎 Search Response (200 OK expected):"
echo "✓ Top result ID: $ID"

RESPONSE=$(curl -s -G https://lrclib.net/api/get/$ID | jq .syncedLyrics -r)

# Check if the response is empty
if [[ -z "$RESPONSE" ]]; then
  echo "❌ No synced lyrics found for $SONG by $ARTIST."
  exit 1
fi

# Save the response to a file
OUTPUT_FILE="$SONG-$ARTIST.lrc"
echo "$RESPONSE" > "$OUTPUT_FILE"
echo "📝 Synced lyrics saved to $OUTPUT_FILE"




