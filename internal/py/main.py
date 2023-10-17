import redis
import threading
from flask import Flask, jsonify
from decouple import config

app = Flask(__name__)

@app.route('/healthy')
def health_check():
    return jsonify(status='OK')

def message_handler(message):
    print("Message received: " + str(message))



if __name__ == "__main__":
    redisPort = config('REDIS_PORT', default=6379, cast=int)
    serverPort = config('PY_PORT', default=4401, cast=int)

    r = redis.Redis(host=config('REDIS_HOST', default='redis'), port=redisPort, db=0)
    #  Test connection first
    try:
        response = r.ping()
        if response:
            print("Redis server is up and responding.")
        else:
            print("Redis server is not responding.")

        pubsub = r.pubsub()
        pubsub.psubscribe('*')
        
        # Run Flask on a separate thread
        flask_thread = threading.Thread(target=lambda: app.run(host='0.0.0.0', port=serverPort, debug=True, use_reloader=False))
        flask_thread.start()
        
        for data in pubsub.listen():
            if data['type'] == 'audio':
                message_handler(data)
        
    except redis.ConnectionError as e:
        print(f"Failed to connect to Redis {e}")

   
    # pubsub.subscribe(**{'audio': message_handler})