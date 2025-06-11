
import { useState } from "react";
import { Music } from "lucide-react";

interface AlbumArtProps {
  trackId: string;
  size?: "sm" | "md" | "lg";
}

export const AlbumArt = ({ trackId, size = "md" }: AlbumArtProps) => {
  const [imageError, setImageError] = useState(false);
  const [imageLoading, setImageLoading] = useState(true);

  const sizeClasses = {
    sm: "w-12 h-12",
    md: "w-16 h-16",
    lg: "w-32 h-32"
  };

  const iconSizes = {
    sm: 16,
    md: 20,
    lg: 32
  };

  const handleImageLoad = () => {
    setImageLoading(false);
  };

  const handleImageError = () => {
    setImageError(true);
    setImageLoading(false);
  };

  if (imageError || imageLoading) {
    return (
      <div className={`${sizeClasses[size]} bg-gradient-to-br from-purple-600 to-pink-600 rounded-lg flex items-center justify-center`}>
        <Music className="text-white/80" size={iconSizes[size]} />
      </div>
    );
  }

  return (
    <div className={`${sizeClasses[size]} rounded-lg overflow-hidden`}>
      <img
        src={`http://localhost:8080/cover/${trackId}`}
        alt="Album artwork"
        className="w-full h-full object-cover"
        onLoad={handleImageLoad}
        onError={handleImageError}
      />
    </div>
  );
};
