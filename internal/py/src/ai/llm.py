from langchain.llms import OpenAI
from decouple import config

llm = OpenAI(openai_api_key=config("OPENAI_API_KEY"))