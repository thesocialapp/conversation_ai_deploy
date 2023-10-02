// actions.js
export const SET_TRANSCRIPT = 'SET_TRANSCRIPT';
export const SET_LOADING = 'SET_LOADING';
export const SET_CONVERSATION = 'SET_CONVERSATION';
export const ADD_MESSAGE = 'ADD_MESSAGE';
export const SET_ERROR = 'SET_ERROR';

export const setTranscript = transcript => ({
  type: SET_TRANSCRIPT,
  transcript,
});

export const setLoading = loading => ({
  type: SET_LOADING,
  loading,
});

export const setConversation = conversation => ({
  type: SET_CONVERSATION,
  conversation,
});

export const addMessage = message => ({
  type: ADD_MESSAGE,
  message,
});

export const setError = error => ({
  type: SET_ERROR,
  error,
});
