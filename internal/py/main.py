import redis
import threading
from flask import Flask, jsonify, Response
from decouple import config

app = Flask(__name__)
pubsub_initialized = threading.Event()

redisPort = config('REDIS_PORT', default=6379, cast=int)
serverPort = config('PY_PORT', default=4401, cast=int)

r = redis.StrictRedis(host=config('REDIS_HOST', default='redis'), port=redisPort, db=0)

def message_handler(message):
    app.logger.info(f"Received message: {message}")
    print("Message received: " + str(message))
    # yield message

def event_stream():
    try:
        response = r.ping()
        if response:
            print("Redis server is up and responding.")
        else:
            print("Redis server is not responding.")

        pubsub = r.pubsub()
        pubsub.subscribe('audio')
        pubsub_initialized.set()
        for data in pubsub.listen():
            print(f'We are listening {data}')
            if data['type'] == 'message':
                message_handler(data)
    except redis.ConnectionError as e:
        print(f"Failed to connect to Redis {e}")
    


@app.route('/healthy', methods=['GET'])
def health_check():
    return jsonify(status='OK')

@app.route('/stream')
def stream():
    return Response(event_stream, mimetype='text/event-stream')


def start_app():
    app.run(host='0.0.0.0', port=serverPort, debug=True, use_reloader=False)
    # flask_initialized.set()

if __name__ == "__main__":
    pubsub_thread = threading.Thread(target=event_stream)
    pubsub_thread.daemon = True # Break app if it breaks
    pubsub_thread.start()

    pubsub_initialized.wait()

    # Start flask app
    start_app()