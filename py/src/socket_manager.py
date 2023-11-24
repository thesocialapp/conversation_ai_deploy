

# Socker Manager encapsulates all socket.io events
class SocketManager:
    def __init__(self, socketio) -> None:
        self.socketio = socketio

        # Add event handlers
        self.socketio.on_event("connect", self.on_connect)
        self.socketio.on_event("disconnect", self.on_disconnect)

    def on_connect(self):
        print("Connected to socket.io")
        self.socketio.emit("connected", {"data": "Connected to socket.io"})
        self.socketio.emit("test", {"data": "Test data"})

    def on_disconnect(self):
        print("Disconnected from socket.io")
        self.socketio.emit("disconnected", {"data": "Disconnected from socket.io"})