
from flask import Flask
from flask_socketio import SocketIO
from socket_manager import SocketManager
from routes import documents_route, health_route, query_route
import configs

app = Flask(__name__)

# Change environment based on configs
app.config['DEBUG'] = configs.env == 'development'

socketio = SocketIO()

# Register blueprints
app.register_blueprint(documents_route.files_bp, url_prefix='/api')
app.register_blueprint(health_route.health_bp, url_prefix='/api')
app.register_blueprint(query_route.query_bp, url_prefix='/api')

SocketIO.init_app(socketio, app=app, path="/io", async_mode="eventlet", cors_allowed_origins="*")
manager = SocketManager(socketio)


def start_app():
    socketio.run(app, host="0.0.0.0", port=configs.pyPort, debug=True, use_reloader=False, log_output=True)