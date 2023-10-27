import redis
import threading
from elevenlabs import set_api_key
from flask import Flask, jsonify, Response
from events import cai_events
from decouple import config

app = Flask(__name__)
pubsub_initialized = threading.Event()

redisPort = config('REDIS_PORT', default=6379, cast=int)
serverPort = config('PY_PORT', default=4401, cast=int)
host = config('REDIS_HOST', default='redis')
elevenlabsKey = config('ELEVEN_LABS_APIKEY')
set_api_key(elevenlabsKey)

r = redis.StrictRedis(host=host, port=redisPort, db=0)

def start_rpubsub():
    cai_events.event_stream(app, r, pubsub_initialized)

@app.route('/healthy', methods=['GET'])
def health_check():
    return jsonify(status='OK')

@app.route('/stream')
def stream():
    return Response(start_rpubsub, mimetype='text/event-stream')


def start_app():
    app.run(host='0.0.0.0', port=serverPort, debug=True, use_reloader=False)

def run_redis_thread():
    pubsub_thread = threading.Thread(target=start_rpubsub)
    pubsub_thread.daemon = True # Break app if it breaks
    pubsub_thread.start()

    pubsub_initialized.wait()

if __name__ == "__main__":
    # Start redis thread
    run_redis_thread()
    # Start flask app
    start_app()