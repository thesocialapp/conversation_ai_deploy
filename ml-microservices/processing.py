import logging
import json
import uuid
from time import perf_counter
from flask import Flask, request, jsonify
from transformers import AutoModelForCausalLM, AutoTokenizer
from config import TTS_API_KEY, VOICE_ID  # Importing from config.py
from elevenlabs import TTS  # Import ElevenLabs SDK
from main import dg_client  # Adjust the import accordingly.
from history import get_history  # Import get_history from history.py

# Initialize ElevenLabs TTS client
tts_client = TTS(TTS_API_KEY)

# Initialize Flask app
app = Flask(__name__)

# Initialize Falcon model and tokenizer.
model_name = "tiiuae/falcon-7b"
model = AutoModelForCausalLM.from_pretrained(model_name)
tokenizer = AutoTokenizer.from_pretrained(model_name)

# Middleware for request and response logging.
class RequestLogger:
    def __init__(self, app):
        self.app = app

    def __call__(self, environ, start_response):
        request_id = str(uuid.uuid4())
        request_data = {"request_id": request_id}
        logger.info(json.dumps({"event": "Request Received", "data": request_data}))
        start = perf_counter()

        def logging_start_response(status, headers, *args):
            response_time = perf_counter() - start
            request_data.update({"status": status, "response_time": response_time})
            logger.info(json.dumps({"event": "Sending Response", "data": request_data}))
            return start_response(status, headers, *args)

        return self.app(environ, logging_start_response)

# Register middleware.
app.wsgi_app = RequestLogger(app.wsgi_app)

# Configure JSON logging.
logging.basicConfig(level=logging.DEBUG, format='%(message)s')
logger = logging.getLogger(__name__)

# Conversation memory.
CONVERSATION_MEMORY = {}

# Function to transcribe audio using Deepgram.
def transcribe_audio(audio):
    try:
        transcript = dg_client.transcription.prerecorded(audio)
        if not transcript:
            raise ValueError("Invalid transcription result")
        logger.info(f"Transcript: {transcript}")
        return transcript
    except Exception as e:
        logger.error(f"Deepgram Error: {str(e)}")
        raise e

# Helper function to extract entities
def parse_entities(conversations):
    entities = []
    for conv in conversations:
        entities.extend(extract_entities(conv))  
    return entities

# Function to generate response text using Falcon model.
def generate_response(input_text, entities=None):
    try:
        context = f"{input_text} {' '.join(entities)}" if entities else input_text
        input_ids = tokenizer.encode(context, return_tensors='pt')
        output_ids = model.generate(input_ids)
        response = tokenizer.decode(output_ids[0], skip_special_tokens=True)
        return response
    except Exception as e:
        logger.error(json.dumps({
            "event": "Error generating response",
            "error": str(e)
        }))
        raise e

# Synthesize voice function.
def synthesize_voice(text):
    try:
        audio = tts_client.synthesize(text, voice=VOICE_ID)
        return audio
    except Exception as e:
        logger.error(f"ElevenLabs Error: {str(e)}")
        raise e

# Route to process the audio file and generate response text.
@app.route('/process_text', methods=['POST'])
def process_text_endpoint():
    request_id = str(uuid.uuid4())
    user_id = request.form.get('user_id')  # Assume each request contains user_id
    try:
        audio = request.files['audio']
        audio_bytes = audio.read()
        transcript = transcribe_audio(audio_bytes)
        logger.info(json.dumps({
            "event": "Processing text",
            "request_id": request_id
        }))

        # Get user's conversation history and parse entities
        history = get_history(user_id)
        entities = []
        if history:
            entities = parse_entities(history)

        # Generate response using the transcript and entities as context
        response = generate_response(transcript, entities)
        audio_response = synthesize_voice(response)
        CONVERSATION_MEMORY[request_id] = [transcript, response]
    except Exception as e:
        logger.error(json.dumps({
            "event": "Error processing text",
            "error": str(e),
            "request_id": request_id
        }))
        return jsonify({"error": str(e)}), 500

    logger.info(json.dumps({
        "event": "Output text",
        "output": response,
        "request_id": request_id
    }))

    return jsonify({"response": response, "audio": audio_response, "transcript": transcript, "request_id": request_id, "audio_bytes": audio_bytes.hex()}), 200

# Run the Flask app.
if __name__ == "__main__":
    app.run(debug=True, host="localhost", port=5000)
