from elevenlabs import generate
import llm
from langchain.schema import HumanMessage
from decouple import config


def synthesize_voice(text: str):
    """Convert transcription to audio"""
    try:
        prediction = _synthesize_response(llm=llm, text=text)
        audio = generate(
            text=prediction,
            voice="Bella",
            model="eleven_multilingual_v2"
        )
        return audio
    except Exception as e:
        print(e)
        raise e
    

def _synthesize_response(llm, text: str) -> str:
    try:
        messages = [HumanMessage(text=text)]
        answer = llm.predict_messages(messages=messages)
        print("Synthesizing response...")
        return answer
    except Exception as e:
        print(e)
        raise e
    

