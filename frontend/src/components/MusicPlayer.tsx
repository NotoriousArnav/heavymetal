
import { useState, useEffect, useRef } from "react";
import { Play, Pause, SkipBack, SkipForward, Volume2, VolumeOff } from "lucide-react";
import { Track, Lyrics } from "@/types/music";
import { Howl } from "howler";
import { LyricsDisplay } from "./LyricsDisplay";
import { AlbumArt } from "./AlbumArt";
import { ProgressBar } from "./ProgressBar";

interface MusicPlayerProps {
  track: Track;
  tracks: Track[];
  onTrackChange: (track: Track) => void;
}

export const MusicPlayer = ({ track, tracks, onTrackChange }: MusicPlayerProps) => {
  const [isPlaying, setIsPlaying] = useState(false);
  const [currentTime, setCurrentTime] = useState(0);
  const [duration, setDuration] = useState(0);
  const [volume, setVolume] = useState(0.7);
  const [isMuted, setIsMuted] = useState(false);
  const [lyrics, setLyrics] = useState<Lyrics | null>(null);
  const [showLyrics, setShowLyrics] = useState(false);
  const howlRef = useRef<Howl | null>(null);
  const intervalRef = useRef<NodeJS.Timeout | null>(null);

  // Initialize audio when track changes
  useEffect(() => {
    if (howlRef.current) {
      howlRef.current.unload();
    }

    const sound = new Howl({
      src: [`http://localhost:8080/stream/${track.id}`],
      html5: true,
      volume: isMuted ? 0 : volume,
      onload: () => {
        setDuration(sound.duration());
      },
      onplay: () => {
        setIsPlaying(true);
        startProgressUpdate();
      },
      onpause: () => {
        setIsPlaying(false);
        stopProgressUpdate();
      },
      onend: () => {
        handleNext();
      },
    });

    howlRef.current = sound;
    setCurrentTime(0);
    
    // Fetch lyrics
    fetchLyrics();

    return () => {
      if (howlRef.current) {
        howlRef.current.unload();
      }
      stopProgressUpdate();
    };
  }, [track]);

  const fetchLyrics = async () => {
    try {
      // Using LRCLib.net API for lyrics
      const response = await fetch(
        `https://lrclib.net/api/search?artist_name=${encodeURIComponent(track.artist_id)}&track_name=${encodeURIComponent(track.title)}`
      );
      
      if (response.ok) {
        const data = await response.json();
        if (data.length > 0) {
          const lyricData = data[0];
          if (lyricData.syncedLyrics) {
            const synced = parseLRC(lyricData.syncedLyrics);
            setLyrics({
              synced,
              plain: lyricData.plainLyrics || ""
            });
          } else if (lyricData.plainLyrics) {
            setLyrics({
              synced: [],
              plain: lyricData.plainLyrics
            });
          }
        }
      }
    } catch (error) {
      console.error("Failed to fetch lyrics:", error);
      setLyrics(null);
    }
  };

  const parseLRC = (lrcString: string) => {
    const lines = lrcString.split('\n');
    const lyrics = [];
    
    for (const line of lines) {
      const match = line.match(/\[(\d{2}):(\d{2})\.(\d{2})\](.*)/);
      if (match) {
        const minutes = parseInt(match[1]);
        const seconds = parseInt(match[2]);
        const centiseconds = parseInt(match[3]);
        const time = minutes * 60 + seconds + centiseconds / 100;
        const text = match[4].trim();
        
        if (text) {
          lyrics.push({ time, text });
        }
      }
    }
    
    return lyrics.sort((a, b) => a.time - b.time);
  };

  const startProgressUpdate = () => {
    intervalRef.current = setInterval(() => {
      if (howlRef.current && howlRef.current.playing()) {
        setCurrentTime(howlRef.current.seek());
      }
    }, 1000);
  };

  const stopProgressUpdate = () => {
    if (intervalRef.current) {
      clearInterval(intervalRef.current);
      intervalRef.current = null;
    }
  };

  const togglePlay = () => {
    if (!howlRef.current) return;
    
    if (isPlaying) {
      howlRef.current.pause();
    } else {
      howlRef.current.play();
    }
  };

  const handleNext = () => {
    const currentIndex = tracks.findIndex(t => t.id === track.id);
    const nextIndex = (currentIndex + 1) % tracks.length;
    onTrackChange(tracks[nextIndex]);
  };

  const handlePrevious = () => {
    const currentIndex = tracks.findIndex(t => t.id === track.id);
    const prevIndex = currentIndex === 0 ? tracks.length - 1 : currentIndex - 1;
    onTrackChange(tracks[prevIndex]);
  };

  const handleSeek = (time: number) => {
    if (howlRef.current) {
      howlRef.current.seek(time);
      setCurrentTime(time);
    }
  };

  const handleVolumeChange = (newVolume: number) => {
    setVolume(newVolume);
    if (howlRef.current) {
      howlRef.current.volume(isMuted ? 0 : newVolume);
    }
  };

  const toggleMute = () => {
    setIsMuted(!isMuted);
    if (howlRef.current) {
      howlRef.current.volume(isMuted ? volume : 0);
    }
  };

  const formatTime = (time: number) => {
    const minutes = Math.floor(time / 60);
    const seconds = Math.floor(time % 60);
    return `${minutes}:${seconds.toString().padStart(2, '0')}`;
  };

  return (
    <div className="fixed bottom-0 left-0 right-0 bg-black/95 backdrop-blur-lg border-t border-white/10 p-4">
      <div className="max-w-screen-2xl mx-auto">
        <div className="flex items-center justify-between mb-2">
          {/* Track Info */}
          <div className="flex items-center space-x-4 flex-1 min-w-0">
            <AlbumArt trackId={track.id} size="sm" />
            <div className="min-w-0">
              <h3 className="text-white font-medium truncate">{track.title}</h3>
              <p className="text-gray-400 text-sm truncate">{track.artist_id}</p>
            </div>
            {lyrics && (
              <button
                onClick={() => setShowLyrics(!showLyrics)}
                className="text-purple-400 hover:text-purple-300 text-sm px-3 py-1 rounded-full border border-purple-400/30 hover:border-purple-300/50 transition-colors"
              >
                {showLyrics ? "Hide Lyrics" : "Show Lyrics"}
              </button>
            )}
          </div>

          {/* Controls */}
          <div className="flex flex-col items-center space-y-2 flex-1 max-w-md">
            <div className="flex items-center space-x-4">
              <button
                onClick={handlePrevious}
                className="text-gray-400 hover:text-white transition-colors"
              >
                <SkipBack size={20} />
              </button>
              
              <button
                onClick={togglePlay}
                className="bg-purple-600 hover:bg-purple-700 text-white p-3 rounded-full transition-colors"
              >
                {isPlaying ? <Pause size={20} /> : <Play size={20} />}
              </button>
              
              <button
                onClick={handleNext}
                className="text-gray-400 hover:text-white transition-colors"
              >
                <SkipForward size={20} />
              </button>
            </div>
            
            <ProgressBar
              currentTime={currentTime}
              duration={duration}
              onSeek={handleSeek}
            />
            
            <div className="flex items-center space-x-2 text-xs text-gray-400">
              <span>{formatTime(currentTime)}</span>
              <span>/</span>
              <span>{formatTime(duration)}</span>
            </div>
          </div>

          {/* Volume */}
          <div className="flex items-center space-x-3 flex-1 justify-end">
            <button onClick={toggleMute} className="text-gray-400 hover:text-white">
              {isMuted ? <VolumeOff size={20} /> : <Volume2 size={20} />}
            </button>
            <input
              type="range"
              min="0"
              max="1"
              step="0.1"
              value={isMuted ? 0 : volume}
              onChange={(e) => handleVolumeChange(parseFloat(e.target.value))}
              className="w-24 accent-purple-600"
            />
          </div>
        </div>

        {/* Lyrics Display */}
        {showLyrics && lyrics && (
          <LyricsDisplay lyrics={lyrics} currentTime={currentTime} />
        )}
      </div>
    </div>
  );
};
