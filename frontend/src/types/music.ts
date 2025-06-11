
export interface Track {
  id: string;
  title: string;
  artist_id: string;
  album_id: string;
  file_path: string;
}

export interface LyricLine {
  time: number;
  text: string;
}

export interface Lyrics {
  synced: LyricLine[];
  plain: string;
}
