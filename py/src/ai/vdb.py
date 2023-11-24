
# This is where we initialize the AI's database.  We will be using Chroma DB
import logging
import chromadb
from langchain.embeddings.sentence_transformer import SentenceTransformerEmbeddings
from langchain.embeddings.openai import OpenAIEmbeddings
import configs
from langchain.vectorstores import chroma
from langchain.chat_models import ChatOpenAI
from langchain.prompts import ChatPromptTemplate
from langchain.schema import StrOutputParser
from langchain.chains import RetrievalQA
from langchain.callbacks.streaming_stdout import StreamingStdOutCallbackHandler
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
    _instance = None

    def __init__(self, collection_name, embedding_function_name="all-MiniLM-L6-v2"):
        self.db = self.__setup_client()
        self.collection_name = collection_name
        self.llm = ChatOpenAI(
            streaming=True,
            model="gpt-3.5-turbo",
            callbacks=[StreamingStdOutCallbackHandler()],
            temperature=0,
        )
        self.db4 = self.__create(collection_name=collection_name, embedding_function_name=embedding_function_name)

    @classmethod
    def get_instance(cls):
        print("Getting instance")
        if cls._instance is None:
            cls._instance = VectorDb(collection_name="ai_pdf")
        return cls._instance

    def get_db(self):
        return self.db
    
    def __setup_client(reset=False):
        # Read the host from the environment variables
        host = configs.chromaHost
        port = configs.chromaPort
        
        # Check whether we are in production or development
        if configs.env == 'production':
            # Use the production database
            chroma_client = chromadb.HttpClient(host=host, port=port)
        else:
            # Use the ephemeral database
            chroma_client = chromadb.PersistentClient()
    
        # Reset the DB
        if reset == True and configs.env == 'production':
            chroma_client.reset()

        return chroma_client
    

    # Create a collection using chromadb
    def create_collection(self, name):
        try:
            collection = self.db.get_or_create_collection(name, embedding_function=SentenceTransformerEmbeddings(model_name="all-MiniLM-L6-v2"))
            return collection
        except Exception as e:
            print(e)
            return None
        
    def collection_count(self, name):
        try:
            collection = self.db.get_collection(name)
            return collection.count()
        except Exception as e:
            print(e)
            logging.error(f"Error getting collection count: {e}")  
            return None
    
    def get_all(self, ids: list):
        try:
            collections = self.db.get_or_create_collection("ai_pdf")
            # Load all docs inside the collection
            # Filter by ids if they exist
            if len(ids) > 0:
                print("Getting all docs")
                docs = collections.get(ids=ids)
            else:
                print("Getting all docs else")
                docs = collections.get()
            return docs
        except Exception as e:
            print(e)
            return None
        
    # Create a document using chromadb
    def create_document(self, document, metadatas, ids):
        all_cols = self.db.list_collections()

        try:
            col = self.db.get_collection("ai_pdf")
            
            # I need to map the document to extract the page content
            embeddings = SentenceTransformerEmbeddings(model_name="all-MiniLM-L6-v2")
            embs = embeddings.embed_documents(texts=document)
            doc = col.add(documents=document, metadatas=metadatas, ids=ids, embeddings=embs)
            return doc
        except Exception as e:
            logging.error(f"Error adding document to database: {e}")
            return None
    
    # Create a vector store using chromadb
    def __create(self, collection_name, embedding_function_name):
        embedding_function = SentenceTransformerEmbeddings(model_name=embedding_function_name)
    
        db4 = chroma.Chroma(
            client=self.db,
            collection_name=collection_name,
            embedding_function=embedding_function,
        )
        return db4.as_retriever()
    
    
    def query_search(self, query: str):
        return self.db4.similarity_search(query=query, k=10)
    
    def _format_docs(self, docs):
        return "\n\n".join(d.page_content for d in docs)
    
    def query_prompt(self, query: str):
        qa = RetrievalQA.from_chain_type(
            llm=self.llm,
            retriever=self.db4,
            chain_type="stuff"
        )
        # Stream the output
        return qa.run(query=query)


     

       

    

    