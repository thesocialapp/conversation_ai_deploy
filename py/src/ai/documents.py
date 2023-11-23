# from langchain.document_loaders import PyPDFium2Loader
from langchain.text_splitter import RecursiveCharacterTextSplitter
from langchain.document_loaders import PyPDFLoader
from .vdb import VectorDb
import logging

class Document:
    def __init__(self):
        self.db = VectorDb(collection_name="ai_pdf")

    # Load the document 
    def load_document(self, path: str):
        print(f"Loading document {path}")
        """Load a document, PDF from a file path."""
        try:
            loader = PyPDFLoader(path, extract_images=False)
            return loader.load()
        except Exception as e:
            logging.error(f"Error loading document: {e}")
            print(f"Error extracting text from pdf: {e}")
            return list()
        
    # Split the document into chunks
    def split_to_chunks(self, documents: list):
        """Split a document into chunks."""
        try:
            text_splitter = RecursiveCharacterTextSplitter(chunk_size=1000, chunk_overlap=10, add_start_index=True)
            docs = text_splitter.split_documents(documents)

            return docs
        except Exception as e:
            logging.error(f"Error splitting document into chunks: {e}")
            print(f"Error splitting document into chunks: {e}")
            return list()
       
    
    # Add the document to the database
    def add_document(self, document, metadatas, ids):
        try:
            # Create a document
            doc = self.db.create_document(document, metadatas, ids)
            return doc
        except Exception as e:
            print(e)
            logging.error(f"Error adding document to database: {e}")
            return None
    
    def count(self) -> int:
        try:
            return self.db.collection_count("ai_pdf")
        except Exception as e:
            print(e)
            return None
    
    def find_all_docs(self, ids: list):
        try:
            return self.db.get_all(ids)
        except Exception as e:
            print(e)
            return None
        
    def query(self, query):
        try:
            return self.db.query_prompt(query)
        except Exception as e:
            print(e)
            return None
        
    
    


    