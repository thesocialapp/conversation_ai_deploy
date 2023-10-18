
import base64
import io
from pydub import AudioSegment


# Converts the ogg to mp4. 
#  return inmem file
def ogg_to_mp4(message: str) -> io.BytesIO:
    ogg_audio_bytes = base64.b64decode(message)

    ogg_audio = AudioSegment.from_file(io.BytesIO(ogg_audio_bytes), format="ogg")
    mp4_audio = ogg_audio.export(ogg_audio, format="mp4")
    
    inmem_file = io.BytesIO(mp4_audio.read())

    return inmem_file
