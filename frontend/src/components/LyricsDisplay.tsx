
import { useEffect, useRef } from "react";
import { Lyrics } from "@/types/music";

interface LyricsDisplayProps {
  lyrics: Lyrics;
  currentTime: number;
}

export const LyricsDisplay = ({ lyrics, currentTime }: LyricsDisplayProps) => {
  const lyricsRef = useRef<HTMLDivElement>(null);

  const getCurrentLyricIndex = () => {
    if (!lyrics.synced.length) return -1;
    
    for (let i = lyrics.synced.length - 1; i >= 0; i--) {
      if (currentTime >= lyrics.synced[i].time) {
        return i;
      }
    }
    return -1;
  };

  const currentIndex = getCurrentLyricIndex();

  useEffect(() => {
    if (lyricsRef.current && currentIndex >= 0) {
      const currentLine = lyricsRef.current.children[currentIndex] as HTMLElement;
      if (currentLine) {
        currentLine.scrollIntoView({
          behavior: 'smooth',
          block: 'center'
        });
      }
    }
  }, [currentIndex]);

  if (!lyrics.synced.length && !lyrics.plain) {
    return (
      <div className="mt-4 p-4 bg-white/5 rounded-lg">
        <p className="text-gray-400 text-center">No lyrics available</p>
      </div>
    );
  }

  return (
    <div className="mt-4 p-4 bg-white/5 rounded-lg max-h-48 overflow-y-auto">
      <div ref={lyricsRef} className="space-y-2">
        {lyrics.synced.length > 0 ? (
          lyrics.synced.map((line, index) => (
            <p
              key={index}
              className={`text-center transition-all duration-300 ${
                index === currentIndex
                  ? "text-purple-300 font-medium text-lg scale-105"
                  : index < currentIndex
                  ? "text-gray-500"
                  : "text-gray-400"
              }`}
            >
              {line.text}
            </p>
          ))
        ) : (
          <div className="text-gray-400 text-center whitespace-pre-line">
            {lyrics.plain}
          </div>
        )}
      </div>
    </div>
  );
};
