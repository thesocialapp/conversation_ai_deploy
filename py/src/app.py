from flask import Flask
from flask_socketio import SocketIO
from socket_manager import SocketManager
from routes import documents_route, health_route, query_route

app = Flask()

# Initialize socket.io under port 4042
# socketio = SocketIO(app, cors_allowed_origins="*", path='/py/io')

# Register blueprints
app.register_blueprint(documents_route.files_bp, url_prefix='/api')
app.register_blueprint(health_route.health_bp, url_prefix='/api')
app.register_blueprint(query_route.query_bp, url_prefix='/api')

# socket_manager = SocketManager(socketio)