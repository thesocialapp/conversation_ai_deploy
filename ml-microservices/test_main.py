import pytest
from main import process_text_endpoint, transcribe_audio, process_text, InvalidTranscript
from fastapi import HTTPException
from fastapi.testclient import TestClient
from main import app, TextResponse
from fastapi.encoders import jsonable_encoder

client = TestClient(app)

# Mocking the audio file for testing
class MockAudio:
    async def read(self):
        return b'audio_data'

# Mocking Request for testing
class MockRequest:
    headers = {}

# Test transcribe_audio function
def test_transcribe_audio():
    audio = MockAudio()
    result = transcribe_audio(audio)
    assert isinstance(result, str), "Result must be a string."

# Test process_text function
def test_process_text():
    text = "This is a test."
    result = process_text(text)
    assert isinstance(result, str), "Result must be a string."

# Test process_text_endpoint with valid data
@pytest.mark.asyncio
async def test_process_text_endpoint_valid():
    request = MockRequest()
    audio = MockAudio()
    response: TextResponse = await process_text_endpoint(request, audio)
    response_dict = jsonable_encoder(response)
    assert "response" in response_dict
    assert "transcript" in response_dict

# Test process_text_endpoint with invalid data
@pytest.mark.asyncio
async def test_process_text_endpoint_invalid():
    request = MockRequest()
    audio = None
    with pytest.raises(HTTPException) as e:
        await process_text_endpoint(request, audio)
    assert str(e.value.status_code) in str(e.value.detail), "Invalid request should return status code 400."

# Test transcribe_audio_invalid function
@pytest.mark.asyncio
async def test_transcribe_audio_invalid():
    with pytest.raises(InvalidTranscript):
        await process_text_endpoint(audio=MockAudio())

# Test the FastAPI endpoint
def test_process_text_endpoint_success():
    response = client.post(
        "/processText/",
        files={"audio": ("test_audio.wav", b"Some binary data", "audio/wav")},
    )
    assert response.status_code == 200
    assert "response" in response.json()
    assert "transcript" in response.json()
