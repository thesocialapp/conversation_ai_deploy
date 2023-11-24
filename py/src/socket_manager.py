from ai.documents import Document
import asyncio

# Socker Manager encapsulates all socket.io events
class SocketManager:
    def __init__(self, socketio) -> None:
        self.socketio = socketio
        self.doc = Document()

        # Add event handlers
        self.socketio.on_event("connect", self.on_connect)
        self.socketio.on_event("query", self.on_query)
        self.socketio.on_event("disconnect", self.on_disconnect)

    def on_connect(self):
        print("Connected to socket.io")
        self.socketio.emit("connected", {"data": "Connected to socket.io"})
        self.socketio.emit("test", {"data": "Test data"})

    # Label the event name and data

    def on_query(self, data):
        # Get question from data
        question = data.get("question")
        print(f"Question: {question}")
        # Get the answer
        answer = self.doc.query(question)
        print(f"Answer: {answer}")
        # Emit the answer
        self.socketio.emit("answer", {"data": answer})

        
    def on_disconnect(self):
        print("Disconnected from socket.io")
        self.socketio.emit("disconnected", {"data": "Disconnected from socket.io"})