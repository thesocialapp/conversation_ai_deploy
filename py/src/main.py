from elevenlabs import set_api_key
from flask import Flask
import configs
from routes import documents_route, health_route

from redis_init import run_redis_thread

app = Flask(__name__)

# Register blueprints
app.register_blueprint(documents_route.files_bp, url_prefix='/api')
app.register_blueprint(health_route.health_bp, url_prefix='/api')

# Change environment based on configs
app.config['DEBUG'] = configs.env == 'development'

def start_app():
    app.run(host="0.0.0.0", port=configs.pyPort, debug=True, use_reloader=False)

# Initialize eleven labs api key
set_api_key(configs.elevenlabsKey)


if __name__ == "__main__":
    # Start redis thread
    run_redis_thread()
    # Start flask app
    start_app()