

# Initialize redis
import threading
import redis
import configs
from events import cai_events


parts = configs.redisAddress.split(':')
_host, _port = parts[0], int(parts[1])
r = redis.StrictRedis(host=_host, port=_port, db=0)

pubsub_initialized = threading.Event()

def run_redis_thread():
    pubsub_thread = threading.Thread(target=start_rpubsub)
    pubsub_thread.daemon = True # Break app if it breaks
    pubsub_thread.start()

    pubsub_initialized.wait()

def start_rpubsub():
    """Start the redis pubsub thread"""
    cai_events.event_stream(r, pubsub_initialized)