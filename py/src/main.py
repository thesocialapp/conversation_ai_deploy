from elevenlabs import set_api_key
import configs
from ai.vdb import VectorDb
from app import start_app
from redis_init import run_redis_thread


# Initialize eleven labs api key
set_api_key(configs.elevenlabsKey)


if __name__ == "__main__":
    # Start redis thread
    run_redis_thread()
    # Start flask app
    start_app()

    # Initialize the vector database
    VectorDb.get_instance()