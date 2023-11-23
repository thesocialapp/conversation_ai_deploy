from .vdb import VectorDB


# QA class that handles initializing the VectorDB and the QA model
# and provides a method to give answers back
class QA:
    def __init__(self, ):
        self.db = VectorDB(collection_name='qa')
    
    # We will use the upload PDF to handle adding the documents to the vector store
    # then we will save the id of the document in the database
    # Use the id of the document to handle the question answering
    