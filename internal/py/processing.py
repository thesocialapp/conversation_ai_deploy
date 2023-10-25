from elevenlabs import generate

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