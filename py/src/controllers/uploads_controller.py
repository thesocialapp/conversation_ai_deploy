
import io
import os
from pypdf import PageObject, PdfReader, PdfWriter
# Analyze the layout and read the text
from pdfminer.high_level import extract_pages
from pdfminer.layout import LTTextContainer, LTChar, LTRect, LTFigure
# Extract text from tables
import pdfplumber

# Extract text from images
from PIL import Image
from pdf2image import convert_from_path
# Perform OCR to extract text from images
import pytesseract

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
                pdfRead = PdfReader(tmp.name)

                # Dict to extract text from each page
                text_per_page = {}

                for pagenum, page in enumerate(extract_pages(tmp.name)):
                    
                    pageObj = pdfRead.pages[pagenum]
                    page_text = []
                    line_format = []
                    # Text from images
                    image_text = []
                    # Text from table
                    table_text = []
                    table_num = 0
                    
                    page_content = []

                    first_element = True
                    table_extraction_flag = False
                    pdf = pdfplumber.open(tmp.name)

                    # Find the tables in the page
                    page_tables = pdf.pages[pagenum]
                    # Find the Number of tables
                    tables = page_tables.find_tables()

                    # Find all the elements
                    page_elements = [(element.y1, element) for element in page._objs]
                    # Sort the elements as they appear in the page
                    page_elements.sort(key=lambda x: x[0], reverse=True)

                    # Find the elements that compose the page
                    for i, component in enumerate(page_elements):
                        # Extract the position of the top side of the element
                        pos = component[0]
                        # Extract the element of the page layout
                        element = component[1]

                        # Check if the element is a text container
                        if isinstance(element, LTTextContainer):
                            # Check if the text appeared in a table
                            if table_extraction_flag == False:
                                # Extract each text and format for each text element
                                (line_text, format_per_line) = text_extaction(element)
                                # Append the text to the page text
                                page_text.append(line_text)
                                # Append the format of the text
                                line_format.append(format_per_line)
                                page_content.append(line_text)
                            else:
                                pass
                        
                        # Check the elements for images
                        if isinstance(element, LTFigure):
                            # Crop the image from the text
                            pdf_path = crop_image(element, pageObj)
                            # Convert the cropped pdf to image
                            image_path = convert_to_images(pdf_path)

                            # Extract the text from the image
                            text = image_to_text(image_path)
                            image_text.append(text)
                            page_content.append(text)
                            # Add a placeholder in the text and format list
                            page_text.append('image')
                            line_format.append('image')
                        
                        # Check the element for tables
                        if isinstance(element, LTRect):
                            # Check if the element is a table
                            if first_element == True and (table_num+1) <= len(tables):
                                # Find the bounding box of the table
                                lower_side = page.bbox[3] - tables[table_num].bbox[3]
                                upper_side = element.y1

                                # Extract the information from the table
                                table = extract_tables(tmp.name, pagenum, table_num)
                                # Convert the table in structured string format
                                table_str = table_converter(table)
                                # Append the table string into list
                                table_text.append(table_str)
                                page_content.append(table_str)
                                # Set the flag to true to prevent handling this twice
                                table_extraction_flag = True
                                # Make it another element
                                first_element = False
                                # Add placeholder in the text and format list
                                page_text.append('table')
                                line_format.append('table')
                            # Check if we already extracted all the tables
                            if element.y0 >= lower_side and element.y1 <= upper_side:
                                pass
                            elif not isinstance(page_elements[i+1][1], LTRect):
                                table_extraction_flag = False
                                table_num += 1
                                first_element = True

                # Create the key for the dictoinary
                dctKey = 'Page_'+str(pagenum)
                # Add the list of list as the value of the page key
                text_per_page[dctKey] = page_content
            
            tmp.close()

            # Clear out the rest of the temporary files
            os.remove(pdf_path)
            os.remove(image_path)

            # Get the resulted text
            extracted_text = '\n'.join(text_per_page['Page_0'][4])

            return jsonify(status='OK', data={"text": extracted_text, "size": tmp_size,}), 200
        except Exception as e:
            logging.error(f"Error reading file: {e}")
            return jsonify(status=f'Error reading file inner {e}'), 500
        finally:
            os.remove(tmp.name)


    except Exception as e:
        print(f"Error extracting text from pdf: {e}")
        return jsonify(status=f'Error processing file {e}'), 500
   

def text_extaction(element: LTTextContainer):
    line_text = element.get_text()

    # Find the format of the texts
    line_format = []
    for text_line in element:
        if isinstance(text_line, LTTextContainer):
            # Iterating through each character in the line of text
            for character in text_line:
                if isinstance(character, LTChar):
                    # Append the font name to the list
                    line_format.append(character.fontname)
                    # Append the font size
                    line_format.append(character.size)
    
    format_per_line = list(set(line_format))

    return (line_text, format_per_line)

# Create function to crop the images from the text
def crop_image(element: LTFigure, pageObj: PageObject):
    # Get the coordinates of the image
    [image_left, image_top, image_right, image_bottom] = [element.x0, element.y0, element.x1, element.y1]
    # Crop the image using the coordinates
    pageObj.mediabox.lower_left = (image_left, image_bottom)
    pageObj.mediabox.upper_right = (image_right, image_top)

    cropped_pdf_writer = PdfWriter()
    cropped_pdf_writer.addpage(pageObj)


    # Save the image in a temporary file
    with NamedTemporaryFile(delete=False) as tmp:
        # Save the image to a new PDF
        cropped_pdf_writer.write(tmp.name)
        tmp.flush()

        # return the path to the cropped image
        return tmp.name
    
# Convert the PDF to image
def convert_to_images(input_file: str):
    images = convert_from_path(input_file)
    image = images[0]

    # Save the image in a temporary file
    with NamedTemporaryFile(delete=False) as tmp:
        # Save the image to a new PDF
        image.save(tmp.name, 'PNG')
        tmp.flush()

        # return the path to the cropped image
        return tmp.name
    
# Function to read text from image
def image_to_text(image_path: str) -> str:
    img = Image.open(image_path)

    # Extract the image from the text
    text = pytesseract.image_to_string(img)
    return text


# Extract tables from the PDF
def extract_tables(input_file: str, page_num: int = 0, table_num: int = 0) -> list:
    with pdfplumber.open(input_file) as pdf:
        # Extract the tables from the PDF
        tables = pdf.pages[page_num]
        # Find the examined page
        page = tables.extract_tables()[table_num]
        return page

# Convert table to appropriate format
def table_converter(table: list) -> str:
    table_string = ''
    # Iterate through each row of the table
    for row_num in range(len(table)):
        row = table[row_num]
        # Remove line breaker from the wrapped text
        cleaned_row = [item.replace('\n', ' ') if item is not None and '\n' in item else 'None' if item is None else item for item in row]
        # Convert the table into a string
        table_string += ('|'+'|'.join(cleaned_row)+'|'+ '\n')

    # Remove the last line break
    table_string = table_string[:-1]
    return table_string