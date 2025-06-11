
import { Music, List } from "lucide-react";
import { Track } from "@/types/music";

interface SidebarProps {
  tracks: Track[];
  onTrackSelect: (track: Track) => void;
}

export const Sidebar = ({ tracks, onTrackSelect }: SidebarProps) => {
  // Group tracks by album for better organization
  const albumGroups = tracks.reduce((acc, track) => {
    const albumId = track.album_id;
    if (!acc[albumId]) {
      acc[albumId] = [];
    }
    acc[albumId].push(track);
    return acc;
  }, {} as Record<string, Track[]>);

  const albums = Object.keys(albumGroups);

  return (
    <div className="w-80 bg-black/50 backdrop-blur-lg border-r border-white/10 flex flex-col">
      <div className="p-6 border-b border-white/10">
        <div className="flex items-center space-x-3">
          <Music className="text-purple-400" size={32} />
          <div>
            <h1 className="text-xl font-bold text-white">Music Player</h1>
            <p className="text-gray-400 text-sm">{tracks.length} tracks</p>
          </div>
        </div>
      </div>

      <div className="flex-1 overflow-y-auto p-4">
        <div className="space-y-4">
          <div className="flex items-center space-x-2 text-purple-400">
            <List size={16} />
            <span className="text-sm font-medium">Albums ({albums.length})</span>
          </div>
          
          {albums.map((albumId) => (
            <div key={albumId} className="space-y-1">
              <h3 className="text-white font-medium text-sm px-2 py-1 bg-white/5 rounded">
                {albumId}
              </h3>
              <div className="space-y-1 pl-2">
                {albumGroups[albumId].map((track) => (
                  <button
                    key={track.id}
                    onClick={() => onTrackSelect(track)}
                    className="w-full text-left p-2 rounded hover:bg-white/5 transition-colors group"
                  >
                    <div className="text-gray-300 text-sm truncate group-hover:text-white">
                      {track.title}
                    </div>
                    <div className="text-gray-500 text-xs truncate">
                      {track.artist_id}
                    </div>
                  </button>
                ))}
              </div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
};
