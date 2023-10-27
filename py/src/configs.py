from decouple import config

redisPort = config('REDIS_PORT', default=6379, cast=int)
serverPort = config('PY_PORT', default=4401, cast=int)
host = config('REDIS_HOST', default='redis')
elevenlabsKey = config('ELEVEN_LABS_APIKEY')
openAiKey = config('OPENAI_API_KEY')
env = config('ENVIRONMENT', default='development')