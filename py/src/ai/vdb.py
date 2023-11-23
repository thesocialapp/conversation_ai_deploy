
# This is where we initialize the AI's database.  We will be using Chroma DB
import chromadb
from chromadb.config import Settings
from langchain.embeddings.sentence_transformer import SentenceTransformerEmbeddings
from langchain.embeddings.openai import OpenAIEmbeddings
from chromadb import Collection
import configs
from langchain.vectorstores import chroma
from langchain.chat_models import ChatOpenAI
from langchain.prompts import ChatPromptTemplate
from langchain.schema import StrOutputParser
from langchain.schema.runnable import RunnablePassthrough
# Create a VectorDb singleton that initializes chromadb and allows it to be accessinle across the application

# Trying to understand a working flow
# 1. The user uploads a document
# 2. The document is split into chunks
# 3. We create a collection
# 4. We create a document and add the split document to the collection
# 5. We create a vector store
# 6. We query the vector store using a prompt and return the results
class VectorDb:
    def __init__(self, collection_name, embedding_function_name="all-MiniLM-L6-v2"):
        self.db = self.__setup_client()
        self.db4 = self.__create(collection_name=collection_name, embedding_function_name=embedding_function_name)

    def get_db(self):
        return self.db
    
    def __setup_client(reset=False):
        # Read the host from the environment variables
        host = configs.chromaHost
        port = configs.chromaPort
        settings = Settings(
            allow_resets=True,
            anonymized_telemetry=False,
        )
        
        # Check whether we are in production or development
        if configs.env == 'production':
            # Use the production database
            chroma_client = chromadb.Client(host=host, port=port, settings=settings)
        else:
            # Use the ephemeral database
            chroma_client = chromadb.EphemeralClient(settings=settings)
        # Reset the DB
        if reset:
            chroma_client.reset()

        return chroma_client
    

    # Create a collection using chromadb
    def create_collection(self, name):
        try:
            collection = self.db.create_collection(name)
            return collection
        except chromadb.exceptions.CollectionAlreadyExists:
            return self.db.get_collection(name)
        except Exception as e:
            print(e)
            return None
        
    # Create a document using chromadb
    def create_document(self, collection: Collection, document, metadatas, ids):
        try:
            doc = collection.add(documents=document, metadatas=metadatas, ids=ids)
            return doc
        except Exception as e:
            print(e)
            return None
        
    def find_document(self, collection: Collection, id):
        try:
            doc = collection.find(id=id)
            return doc
        except Exception as e:
            print(e)
            return None
    
    # Create a vector store using chromadb
    def __create(self, collection_name, embedding_function_name):
        embedding_function = SentenceTransformerEmbeddings(model_name=embedding_function_name)
        embedding = OpenAIEmbeddings()
        db4 = chroma(
            client=self.db,
            embeddings=embedding,
            collection_name=collection_name,
            embedding_function=embedding_function,
        )
        return db4.as_retriever()
    
    def query_search(self, query: str):
        return self.db4.similarity_search(query=query, k=10)
    
    def __format_docs(docs):
        return "\n\n".join(d.page_content for d in docs)
    
    def query_prompt(self, query: str):
        """Query the database using a prompt."""
        template = """
        Answer the question based solely on the following context:
         {context}
         Question: {question}
         """
        prompt = ChatPromptTemplate.from_template(template)
        model = ChatOpenAI(model="gpt-3.5-turbo", )
        
        retriever = self.db4.as_retriever()
        chain = (
            {"context": retriever | self.__format_docs, "question" : RunnablePassthrough()}
            | prompt
            | model
            | StrOutputParser()
            ,
        )
        # Stream the output
        return chain.stream(query)


     

       

    

    