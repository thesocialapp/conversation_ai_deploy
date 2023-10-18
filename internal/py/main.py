import redis
import threading
from flask import Flask, jsonify, Response
from decouple import config
import cai_events

app = Flask(__name__)
pubsub_initialized = threading.Event()

redisPort = config('REDIS_PORT', default=6379, cast=int)
serverPort = config('PY_PORT', default=4401, cast=int)
openAIKey = config('OPENAI_API_KEY')
host = config('REDIS_HOST', default='redis')

r = redis.StrictRedis(host=host, port=redisPort, db=0)

def start_rpubsub():
    cai_events.event_stream(app, r, pubsub_initialized, openAIKey)

@app.route('/healthy', methods=['GET'])
def health_check():
    return jsonify(status='OK')

@app.route('/stream')
def stream():
    return Response(start_rpubsub, mimetype='text/event-stream')


def start_app():
    app.run(host='0.0.0.0', port=serverPort, debug=True, use_reloader=False)

if __name__ == "__main__":
    pubsub_thread = threading.Thread(target=start_rpubsub)
    pubsub_thread.daemon = True # Break app if it breaks
    pubsub_thread.start()

    pubsub_initialized.wait()

    # Start flask app
    start_app()