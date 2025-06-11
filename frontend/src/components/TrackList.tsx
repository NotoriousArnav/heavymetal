
import { Play } from "lucide-react";
import { Track } from "@/types/music";
import { AlbumArt } from "./AlbumArt";

interface TrackListProps {
  tracks: Track[];
  onTrackSelect: (track: Track) => void;
  currentTrack: Track | null;
  searchQuery: string;
}

export const TrackList = ({ tracks, onTrackSelect, currentTrack, searchQuery }: TrackListProps) => {
  const highlightText = (text: string, query: string) => {
    if (!query) return text;
    
    const parts = text.split(new RegExp(`(${query})`, 'gi'));
    return parts.map((part, index) =>
      part.toLowerCase() === query.toLowerCase() ? (
        <span key={index} className="bg-purple-500/30 text-purple-200">
          {part}
        </span>
      ) : (
        part
      )
    );
  };

  return (
    <div className="space-y-1">
      <h2 className="text-2xl font-bold mb-6 text-white">
        {searchQuery ? `Search Results for "${searchQuery}"` : "Your Music"}
      </h2>
      
      {tracks.length === 0 ? (
        <div className="text-center py-12">
          <p className="text-gray-400 text-lg">
            {searchQuery ? "No tracks found" : "No music available"}
          </p>
        </div>
      ) : (
        <div className="grid gap-2">
          {tracks.map((track) => (
            <div
              key={track.id}
              onClick={() => onTrackSelect(track)}
              className={`group flex items-center space-x-4 p-3 rounded-lg cursor-pointer transition-all duration-200 hover:bg-white/5 ${
                currentTrack?.id === track.id 
                  ? "bg-purple-600/20 ring-1 ring-purple-400/30" 
                  : ""
              }`}
            >
              <div className="relative">
                <AlbumArt trackId={track.id} size="md" />
                <div className="absolute inset-0 bg-black/40 opacity-0 group-hover:opacity-100 transition-opacity duration-200 rounded-lg flex items-center justify-center">
                  <Play className="text-white" size={24} />
                </div>
              </div>
              
              <div className="flex-1 min-w-0">
                <h3 className="text-white font-medium truncate">
                  {highlightText(track.title, searchQuery)}
                </h3>
                <p className="text-gray-400 text-sm truncate">
                  {highlightText(track.artist_id, searchQuery)}
                </p>
                <p className="text-gray-500 text-xs truncate">
                  {track.album_id}
                </p>
              </div>
              
              {currentTrack?.id === track.id && (
                <div className="text-purple-400">
                  <div className="flex space-x-1">
                    <div className="w-1 h-4 bg-purple-400 rounded-full animate-pulse"></div>
                    <div className="w-1 h-4 bg-purple-400 rounded-full animate-pulse" style={{ animationDelay: "0.1s" }}></div>
                    <div className="w-1 h-4 bg-purple-400 rounded-full animate-pulse" style={{ animationDelay: "0.2s" }}></div>
                  </div>
                </div>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
};
