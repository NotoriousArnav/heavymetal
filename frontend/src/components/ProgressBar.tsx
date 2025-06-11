
import { useState, useRef } from "react";

interface ProgressBarProps {
  currentTime: number;
  duration: number;
  onSeek: (time: number) => void;
}

export const ProgressBar = ({ currentTime, duration, onSeek }: ProgressBarProps) => {
  const [isDragging, setIsDragging] = useState(false);
  const progressRef = useRef<HTMLDivElement>(null);

  const progress = duration > 0 ? (currentTime / duration) * 100 : 0;

  const handleMouseDown = (e: React.MouseEvent) => {
    setIsDragging(true);
    handleSeek(e);
  };

  const handleMouseMove = (e: React.MouseEvent) => {
    if (isDragging) {
      handleSeek(e);
    }
  };

  const handleMouseUp = () => {
    setIsDragging(false);
  };

  const handleSeek = (e: React.MouseEvent) => {
    if (!progressRef.current) return;
    
    const rect = progressRef.current.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const percentage = Math.max(0, Math.min(100, (x / rect.width) * 100));
    const seekTime = (percentage / 100) * duration;
    
    onSeek(seekTime);
  };

  return (
    <div
      ref={progressRef}
      className="w-full h-1 bg-gray-600 rounded-full cursor-pointer group"
      onMouseDown={handleMouseDown}
      onMouseMove={handleMouseMove}
      onMouseUp={handleMouseUp}
      onMouseLeave={handleMouseUp}
    >
      <div
        className="h-full bg-purple-500 rounded-full relative transition-all group-hover:bg-purple-400"
        style={{ width: `${progress}%` }}
      >
        <div className="absolute right-0 top-1/2 transform translate-x-1/2 -translate-y-1/2 w-3 h-3 bg-purple-500 rounded-full opacity-0 group-hover:opacity-100 transition-opacity"></div>
      </div>
    </div>
  );
};
