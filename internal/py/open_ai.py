import io
import openai
import tempfile


class OpenAI:

    def __init__(self, apiKey: str):
        print(f"Initializing OpenAI {apiKey}")
        openai.api_key = apiKey

    def transcribe_audio(self, file: io.BytesIO) -> str:
        try:
            # Create a temp file
            with tempfile.NamedTemporaryFile(delete=True, suffix=".ogg") as temp_file:
                temp_file.write(file.getvalue())

                with open(temp_file.name, "rb") as file_stream:
                    response = openai.Audio.transcribe(
                        "whisper-1", 
                        file=file_stream,
                        punctuate=True,
                    )

                    print(f'response {response}')
                    return response["text"].strip()
        except Exception as e:
            return str(e)
