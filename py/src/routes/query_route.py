from flask import Blueprint, request
from controllers import query_controller as controller

query_bp = Blueprint('query', __name__)

@query_bp.route('/all', methods=['POST'])
def load_all():
    return controller.load_docs(request)

@query_bp.route('/ask', methods=['POST'])
def ask():
    return controller.ask(request)