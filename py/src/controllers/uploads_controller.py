
import datetime
import os
import uuid

from flask import jsonify
from tempfile import NamedTemporaryFile
import ai.documents as d

import logging



ALLOWED_EXTENSIONS = {'pdf'}
MAX_FILE_SIZE = 30 * 1024 * 1024  # 30MB

def __allowed_file(filename):
    return '.' in filename and filename.rsplit('.', 1)[1].lower() in ALLOWED_EXTENSIONS


def process_upload_pdf(request):
    # Read the file and extract the text
    try:
        # Check if we have files in the request otherwise sent the warning
        if 'document' not in request.files:
            return jsonify(message='No file found in request'), 400
        
        file = request.files['document']

        # Get metadata from the request
        title = request.form.get('title')
        
        description = request.form.get('description')
        created_at = request.form.get('created_at')
        
        # Ensure we have a title and decsription and created_at which we can
        # default to the current time
        if not title or not description or not created_at:
            return jsonify(message='Missing metadata: Please fill in all fields'), 400
        
        # Check if the file is empty or not allowed
        if file.filename == '' or not __allowed_file(file.filename):
            return jsonify(message='Invalid file type'), 400
        
        try:
            # Try casting the file as a FileStorage object to use its function
            # to save the file
            with NamedTemporaryFile(delete=False) as tmp:
                file.save(tmp)
                tmp.flush()

                # Get the file size of the temporary file
                tmp_size = os.path.getsize(tmp.name)
                if tmp_size > MAX_FILE_SIZE:
                    return jsonify(message='File too large'), 400

                # Load the document
                doc = d.Document()
                load = doc.load_document(tmp.name)
                docs = doc.split_to_chunks(load)

                print(f"Number of chunks: {docs}")
                
                ids = []
               
                # Generate ids equal to the number of chunks
                for i in range(len(docs)):
                    doc_id = str(uuid.uuid4())
                    ids.append(doc_id)

                (contents, meta) = _process_documents(docs)
                # Add the document to the database
                doc.add_document(contents, meta, ids)
                print(f"Number of chunks: {len(docs)} and doc id {doc_id}")
                # Close the temporary file
                tmp.close()
            return jsonify(message='OK', data={"id": ids, "docs": len(docs), "size": tmp_size,}), 200
        except Exception as e:
            logging.error(f"Error reading file: {e}")
            return jsonify(message=f'Error reading file inner {e}'), 500
        finally:
            os.remove(tmp.name)


    except Exception as e:
        print(f"Error extracting text from pdf: {e}")
        return jsonify(message=f'Error processing file {e}'), 500

def _process_documents(documents):
        content_list = []
        meta_list = []

        for document_info in documents:
            # Extract relevant metadata
            metadata_subset = {
                'created_at': str(datetime.datetime.now()),
                'source': document_info.metadata['source'],
                'page_number': document_info.metadata['page']
            }

            # Store the page content into a list
            content_list.append(document_info.page_content)
            meta_list.append(metadata_subset)
        
        return content_list, meta_list