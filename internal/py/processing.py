from elevenlabs import generate
from decouple import config
from langchain.llms import OpenAI

def synthesize_voice(text: str):
    """Convert transcription to audio"""
    try:
        audio = generate(
            text=text,
            voice="Bella",
            model="eleven_multilingual_v2"
        )
        return audio
    except Exception as e:
        print(e)
        raise e
    
def synthesize_voice_openai(text: str):
    """Convert transcription to audio"""
    try:
        llm = OpenAI(openai_api_key=config('OPENAI_API_KEY'))
        
        return audio
    except Exception as e:
        print(e)
        raise e