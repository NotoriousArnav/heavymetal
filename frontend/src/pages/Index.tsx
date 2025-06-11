
import { useState, useEffect } from "react";
import { MusicPlayer } from "@/components/MusicPlayer";
import { TrackList } from "@/components/TrackList";
import { SearchBar } from "@/components/SearchBar";
import { Sidebar } from "@/components/PlayerSidebar";
import { Track } from "@/types/music";

const Index = () => {
  const [currentTrack, setCurrentTrack] = useState<Track | null>(null);
  const [tracks, setTracks] = useState<Track[]>([]);
  const [searchResults, setSearchResults] = useState<Track[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [searchQuery, setSearchQuery] = useState("");
  
  // Fetch all tracks on component mount
  useEffect(() => {
    const fetchTracks = async () => {
      try {
        const response = await fetch("http://localhost:8080/tracks/all");
        if (response.ok) {
          const data = await response.json();
          setTracks(data);
          setSearchResults(data);
        }
      } catch (error) {
        console.error("Failed to fetch tracks:", error);
      } finally {
        setIsLoading(false);
      }
    };

    fetchTracks();
  }, []);

  const handleTrackSelect = (track: Track) => {
    setCurrentTrack(track);
  };

  const handleSearch = async (query: string) => {
    setSearchQuery(query);
    if (!query.trim()) {
      setSearchResults(tracks);
      return;
    }

    try {
      const response = await fetch(`http://localhost:8080/search/${encodeURIComponent(query)}`);
      if (response.ok) {
        const data = await response.json();
        setSearchResults(data);
      }
    } catch (error) {
      console.error("Search failed:", error);
      setSearchResults([]);
    }
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900 flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-purple-400 mx-auto mb-4"></div>
          <p className="text-white text-lg">Loading your music...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900 text-white">
      <div className="flex h-screen">
        {/* Sidebar */}
        <Sidebar tracks={tracks} onTrackSelect={handleTrackSelect} />
        
        {/* Main Content */}
        <div className="flex-1 flex flex-col overflow-hidden">
          {/* Header */}
          <header className="p-6 border-b border-white/10">
            <SearchBar onSearch={handleSearch} />
          </header>

          {/* Track List */}
          <div className="flex-1 overflow-y-auto p-6">
            <TrackList 
              tracks={searchResults} 
              onTrackSelect={handleTrackSelect}
              currentTrack={currentTrack}
              searchQuery={searchQuery}
            />
          </div>
        </div>
      </div>

      {/* Music Player - Fixed at bottom */}
      {currentTrack && (
        <MusicPlayer 
          track={currentTrack} 
          tracks={tracks}
          onTrackChange={setCurrentTrack}
        />
      )}
    </div>
  );
};

export default Index;
