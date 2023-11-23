from flask import jsonify
from ai.documents import Document


def load_docs(request):
    """Health check for the server"""
    # Get the id from the request body
    json_data = request.get_json(force=True)
    
    # Ensure we don't have an empty request
    if not json_data:
        return jsonify(message='No data found in request'), 400
    
    uuids = json_data.get('ids')

    if not uuids:
        return jsonify(message='Missing ids'), 400
    
    print(f"Loading documents with ids: {uuids}")

    doc = Document()
    docs = doc.find_all_docs(uuids)
    return jsonify(success=True, data=docs), 200

def ask(request):
    """Health check for the server"""
    # Get the id from the request body
    json_data = request.get_json(force=True)
    
    # Ensure we don't have an empty request
    if not json_data:
        return jsonify(message='No data found in request'), 400
    
    question = json_data.get('question')

    if not question:
        return jsonify(message='Missing question'), 400
    
    print(f"Question: {question}")

    doc = Document()
    docs = doc.query(question)
    return jsonify(success=True, data=docs), 200
    