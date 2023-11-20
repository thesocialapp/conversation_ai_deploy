import base64
import io
from pypdf import PdfReader
import logging

def extract_text_content(file_base64):
    try:
        print(f"Extracting text from file: {file_base64}")
        file_data = base64.b64decode(file_base64.decode('utf-8'))
        doc_buffer = io.BytesIO(file_data)
        print(f"Doc buffer: {doc_buffer}")
        pdfReader = PdfReader(doc_buffer)
        
        extracted_text = ""
        for page_num in range(pdfReader.pages):
            page = pdfReader.getPage(page_num)
            text = page.extractText()
            extracted_text += "\n" + text
            print("Page {}:\n{}".format(page_num, text))

        return extracted_text
    except Exception as e:
        print(f"Error extracting text from pdf: {e}")