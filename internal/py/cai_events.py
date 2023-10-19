

import base64
import io
from flask import Flask
import redis
from threading import Event
from open_ai import OpenAI


def message_handler(app: Flask, message: str, apiKey: str):
    # print("Message received: " + str(message['data']))
    print(f"api key {apiKey}\n")
    audio_str = message['data']
    # Since the audio is coming as base64 string of bytes, we need to decode it
    # to a bytes object under base64 encoding
    try:
        decoded_audio = base64.b64decode(audio_str)
        bytes_audio = io.BytesIO(decoded_audio)
        open_ai = OpenAI(apiKey=apiKey)
        # 
        # Transcribe the audio we got
        tr = open_ai.transcribe_audio(bytes_audio)
        print(f"Transcribed audio {tr}")
    except Exception as e:
        print(f"Error {e}")
    
    app.logger.info(f"We were able to transcribe the data {tr}")

def event_stream(app: Flask, r: redis.StrictRedis, pi: Event, openApiKey: str):
    try:
        # Ensure we test that Redis can handle a PING PONG and 
        #  continue once it does
        response = r.ping()
        if response:
            print("Redis server is up and responding.")
            pubsub = r.pubsub()
            pubsub.subscribe('audio')

            # Setting this allows us to wait for redis to 
            # run and complete before flask can initiate on
            #  the main thread
            pi.set()

            # Listen for all data coming in
            for data in pubsub.listen():
                if data['type'] == 'message':
                    message_handler(app, data, openApiKey)
        else:
            app.logger.error('Unable to connect to Redis successfully')
    except redis.ConnectionError as e:
        app.logger.error('Unable to connect to Redis successfully')