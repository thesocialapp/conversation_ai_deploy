# test_socketio_client.py

import socketio

sio = socketio.Client()

@sio.event
def connect():
    print('Connection established')
    sio.emit('transcript_received', 'This is a test transcript.')

@sio.event
def response_generated(data):
    print('Response received:', data)

sio.connect('http://localhost:8000')
sio.wait()
