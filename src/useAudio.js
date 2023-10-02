// useAudio.js
import { useRef } from 'react';

const useAudio = () => {
  const audioRef = useRef(null);

  const playAudio = (audioBlob) => {
    const audioUrl = URL.createObjectURL(audioBlob);
    audioRef.current.src = audioUrl;
    audioRef.current.play();
  };

  return { audioRef, playAudio };
};

export default useAudio;
