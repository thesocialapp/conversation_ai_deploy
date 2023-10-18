import io
import openai
import tempfile

class OpenAI:
    
    def __init__(self, apiKey: str):
        self.openai.api_key = apiKey
        
    def transcribe_audio(file: io.BytesIO) -> str:
        try:
            # Create a temp file
            with tempfile.NamedTemporaryFile(delete=False) as temp_file:
                temp_file.write(file.getvalue())
            response = openai.Audio.transcribe("whisper-1", temp_file.name)

            # Close the buffer
            file.close()
            return response["text"].strip()
        except:
            return ""