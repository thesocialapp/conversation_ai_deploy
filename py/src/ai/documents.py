from langchain.document_loaders import PyPDFium2Loader
from langchain.text_splitter import RecursiveCharacterTextSplitter
from langchain.document_loaders import PyPDFLoader

def load_document(path: str):
    """Load a document, PDF from a file path."""
    loader = PyPDFLoader(path, extract_images=True)
    return loader.load_and_split()

def split_to_chunks(documents: list):
    """Split a document into chunks."""
    text_splitter = RecursiveCharacterTextSplitter(chunk_size=1000, chunk_overlap=100, add_start_index=True)
    docs = text_splitter.split_documents(documents)

    return docs


    