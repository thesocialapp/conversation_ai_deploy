import redis
import threading
from elevenlabs import set_api_key
from flask import Flask, jsonify, Response
from events import cai_events
import configs

app = Flask(__name__)

# Change environment based on configs
app.config['DEBUG'] = configs.env == 'development'

# This is used to ensure that the pubsub thread is running
pubsub_initialized = threading.Event()

# Initialize eleven labs api key
set_api_key(configs.elevenlabsKey)

# Initialize redis
parts = configs.redisAddress.split(':')
_host, _port = parts[0], int(parts[1])
r = redis.StrictRedis(host=_host, port=_port, db=0)

def start_rpubsub():
    """Start the redis pubsub thread"""
    cai_events.event_stream(r, pubsub_initialized)

@app.route('/healthy', methods=['GET'])
def health_check():
    return jsonify(status='OK')

@app.route('/stream')
def stream():
    return Response(start_rpubsub, mimetype='text/event-stream')


def start_app():
    app.run(host="0.0.0.0", port=configs.pyPort, debug=True, use_reloader=False)

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