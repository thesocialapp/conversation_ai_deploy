# main.py

import json
import redis
import httpx
import logging.handlers
from fastapi import FastAPI, HTTPException, Request, UploadFile, File
from pydantic import BaseModel, UUID4
from processing import process_text, transcribe_audio, save_context
from config import ELEVENLABS_API_KEY, VOICE_ID
from history import get_history

# Initialize Redis
r = redis.Redis()

# Parameters
CACHE_EXPIRATION_TIME = 3600  # 1 hour

# Custom Exception for Cache Hits
class CacheHit(Exception):
    pass

# Custom Exception for Invalid Transcript
class InvalidTranscript(Exception):
    pass

# Custom Logging Formatter
class AppLoggingFormatter(logging.Formatter):
    def format(self, record):
        log_object = {
            "level": record.levelname,
            "timestamp": self.formatTime(record, self.datefmt),
            "message": record.getMessage(),
            "trace_id": record.trace_id if hasattr(record, 'trace_id') else None
        }
        return json.dumps(log_object)

# Initialize Logger
logger = logging.getLogger(__name__)
handler = logging.handlers.HTTPHandler('localhost:9200', '/app/logs', method='POST')
formatter = AppLoggingFormatter()
handler.setFormatter(formatter)
logger.addHandler(handler)

# Request and Response Models
class TextRequest(BaseModel):
    user_id: UUID4
    text: str

class TextResponse(BaseModel):
    response: str
    transcript: str
    audio: bytes

app = FastAPI()

@app.middleware("http")
async def add_trace_id(request: Request, call_next):
    trace_id = request.headers.get("X-Trace-ID", "TraceIDNotFound")
    response = await call_next(request)
    response.headers["X-Trace-ID"] = trace_id
    return response

# ElevenLabs Text to Speech Function
async def text_to_speech(text: str, voice_id: str, api_key: str) -> bytes:
    url = f"https://api.elevenlabs.io/v1/text-to-speech/{voice_id}"
    headers = {
        'xi-api-key': api_key,
        'Content-Type': 'application/json'
    }
    data = {'text': text}
    async with httpx.AsyncClient() as client:
        response = await client.post(url, headers=headers, json=data)
    if response.status_code != 200:
        error_message = response.json().get('message', 'Unknown error')
        raise Exception(f"Error from ElevenLabs API: {error_message}")
    return response.content

@app.post("/processText/", response_model=TextResponse)
async def process_text_endpoint(request: Request, audio: UploadFile = File(...)):
    trace_id = request.headers.get("X-Trace-ID", "TraceIDNotFound")
    try:
        audio_bytes = await audio.read()
        text = transcribe_audio(audio_bytes)
        if not text:
            raise InvalidTranscript("Invalid transcription result")
        request_data = TextRequest(text=text, user_id="some_user_id")
        logger.info(f"Transcript from Deepgram: {text}", extra={'trace_id': trace_id})
    except Exception as e:
        logger.error(f"Error: {e}", extra={'trace_id': trace_id})
        raise HTTPException(status_code=400, detail=str(e))

    user_id_str = str(request_data.user_id)
    key = f"user_data:{user_id_str}"

    try:
        history_entities = get_history(user_id_str)
        if history_entities is None:
            history_entities = {}
        cached_data = r.get(key)
        if cached_data:
            logger.info(f"Cache hit for user: {user_id_str}, key: {key}", extra={'trace_id': trace_id})
            raise CacheHit()
        else:
            processed_text, context = process_text(text, user_id_str, history_entities)
            cache_data = json.dumps({
                "processed_text": processed_text,
                "audio_bytes": audio_bytes.hex(),
                "transcript": text,
                "request_data": request_data.dict()
            })
            r.setex(key, CACHE_EXPIRATION_TIME, cache_data)
    except CacheHit:
        cache_data = json.loads(cached_data)
        processed_text = cache_data['processed_text']
        transcript = cache_data['transcript']
        logger.info(f"Transcript on Cache Hit: {transcript}", extra={'trace_id': trace_id})
        return TextResponse(response=processed_text, transcript=transcript)
    except redis.RedisError as e:
        logger.error(f"Redis Error: {e}", extra={'trace_id': trace_id})
        raise HTTPException(status_code=500, detail=str(e))
    except Exception as e:
        logger.error(f"Error processing text for user {user_id_str}: {e}", extra={'trace_id': trace_id})
        raise HTTPException(status_code=500, detail=str(e))

    audio_data = await text_to_speech(text, VOICE_ID, ELEVENLABS_API_KEY)
    return TextResponse(response=processed_text, transcript=text, audio=audio_data)

@app.post("/converse")
async def converse(user_id: UUID4, text: str):
    history_entities = get_history(str(user_id))
    if history_entities is None:
        history_entities = {}
    processed_text, context = process_text(text, str(user_id), history_entities)
    save_context(str(user_id), context)
    return {"response": processed_text}


# Unit tests can be added as per the testing framework being used, example using pytest:
# 
# def test_transcribe_audio():
#     # Your testing logic here...
#     pass
# 
# def test_process_text():
#     # Your testing logic here...
#     pass
# 
# Ensure to run these tests as part of your CI/CD pipeline to validate the application behavior before deployment.
