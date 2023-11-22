from decouple import config

redisPort = config('REDIS_PORT', default=6379, cast=int)
redisAddress = config('REDIS_ADDRESS', default='redis:6379')
serverPort = config('PY_PORT', default=4041, cast=int)
host = config('REDIS_HOST', default='redis')
elevenlabsKey = config('ELEVEN_LABS_APIKEY')
openAiKey = config('OPENAI_API_KEY')
pyPort = config('PY_PORT', default=4041, cast=int)
env = config('ENVIRONMENT', default='development')
chromaHost = config('CHROMA_HOST', default='http://localhost:8000')
chromaPort = config('CHROMA_PORT', default=8000, cast=int)