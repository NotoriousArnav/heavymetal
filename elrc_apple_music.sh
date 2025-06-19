#!/usr/bin/env bash

# Setup variables
SONG="$1"
ARTIST="$2"

# Check if both song and artist are provided
if [[ -z "$SONG" || -z "$ARTIST" ]]; then
  echo "Usage: $0 <song> <artist>"
  exit 1
fi

QUERY="$(printf "%s %s" "$SONG" "$ARTIST" | jq -s -R @uri)"  # URL-encode using jq

res=$(curl -s -G "https://paxsenix.alwaysdata.net/searchAppleMusic.php" --data-urlencode "q=$SONG $ARTIST" -H "Accept: application/json")

echo "🔎 Search Response (200 OK expected):"
# echo "$res" | jq '.'  # Pretty-print JSON response

# Extract ID of top result (using jq)
TOP_ID=$(echo $res | jq -r '.[0].id')
echo "✓ Top result ID: $TOP_ID"


# 2️⃣ Fetch synced lyrics for that song

curl -s "https://paxsenix.alwaysdata.net/getAppleMusicLyrics.php?id=$TOP_ID" \
     -H "Accept: text/plain" \
     -o "$SONG-$ARTIST.jsonc"

echo "📝 Raw Lyrics:"
head -n 10 "$SONG-$ARTIST.jsonc"  # Show first 10 lines of lyrics file


# 📌 Notes:
# - Use `QUERY=$(printf "%s %s" "$SONG" "$ARTIST" | jq -s -R @uri)` or
#   `--data-urlencode` to URL‑encode the combined song + artist.
# - The first request returns JSON array of AppleSearchResponse objects.
#   Helps replicate getSongInfo(...).
# - The second request returns raw synced-lyrics text (likely LRC).
#   Pass through PaxMusicHelper.formatWordByWordLyrics in your code.
# - Check HTTP status codes: status != 200–299 ⇒ treat as null.
# - Handle empty/invalid JSON or OOB indexes as nulls in your flow.

