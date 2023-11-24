import os
from flask import Blueprint, request
from controllers import uploads_controller as controller

files_bp = Blueprint('upload', __name__)

@files_bp.route('/upload/pdf', methods=['POST'])
def upload_pdf():
    """Upload a pdf file to the server for processing and reading"""
    # Accept a pdf file and read its contents
    # Send the contents to the AI service to summarize
    # Return the summarized text
    return controller.process_upload_pdf(request)