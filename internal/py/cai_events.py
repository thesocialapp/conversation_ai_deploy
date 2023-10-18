

from flask import Flask
import redis
import audio
from threading import Event
from open_ai import OpenAI


def message_handler(app: Flask, message: str, apiKey: str):
    app.logger.info(f"Received message: {message}")
    print("Message received: " + str(message))
    inmem = audio.ogg_to_mp4(message)
    
    open_ai = OpenAI(apiKey=apiKey)
    # Transcribe the audio we got
    tr = open_ai.transcribe_audio(inmem)
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