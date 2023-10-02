/**
 * Importing necessary libraries and components.
 */
import React, { useEffect, useRef } from 'react';
import axios from 'axios';
import SpeechRecognition from 'react-speech-recognition';
import Astrofox from 'astrofox';
import { Socket } from 'socket.io-client';
import Chat from './Chat';
import useAudio from './useAudio';
import Loader from 'react-loader-spinner';
import { useSelector, useDispatch } from 'react-redux';
import './styles.css';
import {
  setTranscript,
  setLoading,
  addMessage,
  setError,
  setListening,
} from './actions';
import {
  getTranscript,
  getLoading,
  getConversation,
  getError,
  getListening,
} from './selectors';

/**
 * Main App component
 */
function App() {
  // References to manage Astrofox and audio components.
  const foxRef = useRef();
  const { audioRef, playAudio } = useAudio();

  // Getting states from the Redux store.
  const transcript = useSelector(getTranscript);
  const loading = useSelector(getLoading);
  const conversation = useSelector(getConversation);
  const error = useSelector(getError);
  const listening = useSelector(getListening); // new listening state

  // useDispatch hook to dispatch actions.
  const dispatch = useDispatch();

  /**
   * Handle transcript.
   * Dispatches addMessage action and sends the transcript to server.
   * @param {string} transcript - The transcript to handle.
   */
  const handleTranscript = (transcript) => {
    if (transcript && transcript.trim() !== '') {
      dispatch(addMessage({ me: transcript }));
      axios.post('/api/speech', { transcript });
    }
  };

  /**
   * Handle microphone click.
   * Sends audio data to server and updates the transcript in the state.
   */
  const handleMicClick = async () => {
    try {
      dispatch(setLoading(true));
      const audioData = getAudioData();
      const response = await axios.post(
        'http://localhost:actual-server-port/api/speech',
        audioData
      );
      if (response.status !== 200) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
      dispatch(setTranscript(response.data)); // Dispatching transcript update on mic click
    } catch (error) {
      dispatch(setError(error.toString()));
    } finally {
      dispatch(setLoading(false));
    }
  };

  /**
   * Start recording audio.
   */
  const startRecording = () => {
    dispatch(setListening(true)); // Dispatching listening state
    SpeechRecognition.startListening({ onTranscript: handleTranscript });
  };

  /**
   * useEffect hook for handling socket connection and disconnection.
   */
  useEffect(() => {
    const socketConnection = Socket();
    return () => {
      socketConnection.disconnect(); // Disconnecting socket on App unmount
    };
  }, []);

  /**
   * useEffect hook for getting response and handling socket events.
   */
  useEffect(() => {
    const getResponse = async () => {
      try {
        dispatch(setLoading(true)); // setting loading before API call
        const { data } = await axios.get('/api/response');
        dispatch(addMessage({ speaker: 'bot', audio: data }));
      } catch (error) {
        dispatch(setError(error.toString()));
      } finally {
        dispatch(setLoading(false)); // reseting loading after API call
      }
    };

    getResponse();

    // Establishing socket connection.
    const socket = Socket();
    socket.on('audio_response', (audioBlob) => {
      dispatch(addMessage({ bot: audioBlob }));
      playAudio(audioBlob);
      foxRef.current.listen();
    });

    // Handling socket error.
    socket.on('error', (errorMsg) => {
      dispatch(setError(errorMsg));
    });

    return () => {
      socket.off('audio_response');
      socket.off('error');
    };
  }, [conversation]);

  /**
   * Rendering the main App component.
   */
  return (
    <div>
      <Astrofox ref={foxRef} />
      <audio ref={audioRef} controls />
      {error && <p className="error">Error: {error}</p>}
      <button
        className="button"
        onClick={startRecording}
        disabled={loading || listening}
      >
        Start
      </button>
      <button className="button" onClick={handleMicClick} disabled={listening}>
        Mic
      </button>
      {loading && (
        <Loader type="Puff" color="#00BFFF" height={100} width={100} />
      )}
      <div>
        <p>Transcript: {transcript}</p>
        <Chat conversation={conversation} />
      </div>
    </div>
  );
}

export default App;

/**
 * Function to get audio data.
 * This function simulates getting some audio data.
 * @returns {Blob} - A new Blob object.
 */
function getAudioData() {
  return new Blob();
}
