

from flask import jsonify


def health(request):
    """Health check for the server"""
    return jsonify(status='OK')