from flask import Flask
from ai import processing
import redis
from threading import Event
import logging
from documents import extract

def _message_handler(message: str, r: redis.StrictRedis):
    # Since the audio is coming as base64 string of bytes, we need to decode it
    # to a bytes object under base64 encoding
    try:
        # Convert message to bytes
        # audioData = message["data"].decode('utf-8')
        # audio = processing.synthesize_voice(audioData)
        # Convert audio to base64 string
        print(f"Type of message['data']: {type(message['data'])}")
        print(f"Length of message['data']: {len(message['data'])}")
        print(f"Message received: {message['data'].decode('utf-8')}")
        fileData = message["data"]

        print(f"The file data is {fileData}")
        
        # Extract text from the file
        text_content = extract.extract_text_content(fileData)
        print(f"Text content: {text_content}")
        r.publish('file-processed', text_content)
    except Exception as e:
        logging.error(f"Error receiving messages and transcribing them: {e}")

def event_stream(r: redis.StrictRedis, pi: Event):
    try:
        # Ensure we test that Redis can handle a PING PONG and 
        # continue once it does
        response = r.ping()
        if response:
            print("Redis server is up and responding.")
            pubsub = r.pubsub()
            pubsub.subscribe("file-document")

            # Setting this allows us to wait for redis to 
            # run and complete before flask can initiate on
            #  the main thread
            pi.set()

            # Listen for all data coming in
            for data in pubsub.listen():
                if data['type'] == 'message':
                    _message_handler(data, r)
        else:
            logging.error('Unable to connect to Redis successfully')
    except redis.ConnectionError as e:
        logging.error('Unable to connect to Redis successfully')