import base64
import os
from pypdf import PdfReader
import logging
from tempfile import NamedTemporaryFile

def extract_text_content(file_base64):
    try:
        file_data = base64.b64decode(file_base64.decode('utf-8'))
        logging.info("Processing")
        extracted_text = ""
        # Create a temp file to store the pdf
        with NamedTemporaryFile(delete=False) as tmp:
            tmp.write(file_data)
            tmp.flush()
            pdfReader = PdfReader(tmp.name)

            # Loop through all the pages and extract the text
            for num in range(len(pdfReader.pages)):
                page = pdfReader.pages[num]
                content = page.extract_text()
                extracted_text += content
                
            os.remove(tmp.name)

        return extracted_text    
    except Exception as e:
        logging.error(f"Error extracting text from pdf: {e}")
        return ""