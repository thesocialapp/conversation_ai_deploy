

from flask import Blueprint, request
from controllers import health_controller as controller

health_bp = Blueprint('health', __name__)

@health_bp.route('/health', methods=['GET'])
def health_check():
    return controller.health(request)