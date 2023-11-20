
import io
import os
from pypdf import PdfReader
from flask import jsonify
from tempfile import NamedTemporaryFile
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
            return jsonify(status='No file found in request'), 400
        
        file = request.files['document']

        if file.filename == '' or not __allowed_file(file.filename):
            return jsonify(status='Invalid file type'), 400
        
        try:
            # Try casting the file as a FileStorage object to use its function
            # to save the file
            with NamedTemporaryFile(delete=False) as tmp:
                tmp.write(file.read())
                tmp.flush()

                # Get the file size of the temporary file
                tmp_size = os.path.getsize(tmp.name)
                if tmp_size > MAX_FILE_SIZE:
                    return jsonify(status='File too large'), 400
                
                # Extract the text from the file
                pdfReader = PdfReader(tmp.name)
                extracted_text = ""
                for page_num in range(pdfReader.pages):
                    page = pdfReader.getPage(page_num)
                    text = page.extract_text()
                    extracted_text += "\n" + text
                    print("Page {}:\n{}".format(page_num, text))

            return jsonify(status='OK', data={"text": extracted_text, "size": tmp_size,}), 200
        except Exception as e:
            logging.error(f"Error reading file: {e}")
            return jsonify(status=f'Error reading file inner {e}'), 500
        finally:
            os.remove(tmp.name)


    except Exception as e:
        print(f"Error extracting text from pdf: {e}")
        return jsonify(status=f'Error processing file {e}'), 500
   