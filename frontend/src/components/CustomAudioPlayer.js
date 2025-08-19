import React, { useState, useRef, useEffect } from 'react';
import './CustomAudioPlayer.css';

const CustomAudioPlayer = ({ songId, onEnded }) => {
  const [audioBlob, setAudioBlob] = useState(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState(null);
  const [isPlaying, setIsPlaying] = useState(false);
  const [currentTime, setCurrentTime] = useState(0);
  const [duration, setDuration] = useState(0);
  const audioRef = useRef(null);

  useEffect(() => {
    fetchAudio();
  }, [songId]);

  const fetchAudio = async () => {
    try {
      setIsLoading(true);
      setError(null);

      const token = localStorage.getItem('authToken');
      if (!token) {
        throw new Error('No authentication token found');
      }

      const response = await fetch(`http://localhost:8080/api/v1/songs/${songId}/stream`, {
        headers: {
          'Authorization': `Bearer ${token}`,
        },
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      const blob = await response.blob();
      setAudioBlob(blob);
      
      // Create blob URL
      const blobUrl = URL.createObjectURL(blob);
      if (audioRef.current) {
        audioRef.current.src = blobUrl;
      }
    } catch (err) {
      console.error('Failed to fetch audio:', err);
      setError(err.message);
    } finally {
      setIsLoading(false);
    }
  };

  const handlePlayPause = () => {
    if (audioRef.current) {
      if (isPlaying) {
        audioRef.current.pause();
      } else {
        audioRef.current.play();
      }
      setIsPlaying(!isPlaying);
    }
  };

  const handleTimeUpdate = () => {
    if (audioRef.current) {
      setCurrentTime(audioRef.current.currentTime);
      setDuration(audioRef.current.duration);
    }
  };

  const handleSeek = (e) => {
    const seekTime = (e.target.value / 100) * duration;
    if (audioRef.current) {
      audioRef.current.currentTime = seekTime;
      setCurrentTime(seekTime);
    }
  };

  const handleEnded = () => {
    setIsPlaying(false);
    setCurrentTime(0);
    if (onEnded) {
      onEnded();
    }
  };

  const formatTime = (time) => {
    if (isNaN(time)) return '0:00';
    const minutes = Math.floor(time / 60);
    const seconds = Math.floor(time % 60);
    return `${minutes}:${seconds.toString().padStart(2, '0')}`;
  };

  if (isLoading) {
    return (
      <div className="custom-audio-player loading">
        <div className="loading-spinner">🎵</div>
        <span>Loading audio...</span>
      </div>
    );
  }

  if (error) {
    return (
      <div className="custom-audio-player error">
        <span>❌ Error: {error}</span>
        <button onClick={fetchAudio} className="retry-btn">Retry</button>
      </div>
    );
  }

  return (
    <div className="custom-audio-player">
      <audio
        ref={audioRef}
        onTimeUpdate={handleTimeUpdate}
        onEnded={handleEnded}
        onLoadedMetadata={handleTimeUpdate}
        preload="metadata"
      />
      
      <div className="audio-controls">
        <button 
          onClick={handlePlayPause} 
          className={`play-pause-btn ${isPlaying ? 'playing' : ''}`}
        >
          {isPlaying ? '⏸️' : '▶️'}
        </button>
        
        <div className="time-display">
          <span>{formatTime(currentTime)}</span>
          <span>/</span>
          <span>{formatTime(duration)}</span>
        </div>
        
        <div className="progress-container">
          <input
            type="range"
            min="0"
            max="100"
            value={duration > 0 ? (currentTime / duration) * 100 : 0}
            onChange={handleSeek}
            className="progress-bar"
          />
        </div>
      </div>
    </div>
  );
};

export default CustomAudioPlayer; 