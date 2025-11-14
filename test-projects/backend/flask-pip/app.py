from flask import Flask, request, jsonify
from datetime import datetime
import time

app = Flask(__name__)

# In-memory data store
items = [
    {"id": 1, "name": "Sample Item 1", "description": "First sample item", "createdAt": datetime.now().isoformat()},
    {"id": 2, "name": "Sample Item 2", "description": "Second sample item", "createdAt": datetime.now().isoformat()}
]
next_id = 3
request_count = 0
start_time = time.time()

# Request counter middleware
@app.before_request
def count_requests():
    global request_count
    request_count += 1

# GET / - Root endpoint with API info
@app.route('/', methods=['GET'])
def get_root():
    return jsonify({
        "message": "Flask Test API",
        "version": "1.0.0",
        "requestCount": request_count,
        "timestamp": datetime.now().isoformat(),
        "endpoints": [
            "GET /",
            "GET /health",
            "GET /items",
            "POST /items",
            "GET /items/<id>",
            "DELETE /items/<id>"
        ]
    })

# GET /health - Health check
@app.route('/health', methods=['GET'])
def get_health():
    uptime = int(time.time() - start_time)
    return jsonify({
        "status": "healthy",
        "uptime": f"{uptime} seconds",
        "timestamp": datetime.now().isoformat()
    })

# GET /items - List all items
@app.route('/items', methods=['GET'])
def get_items():
    return jsonify({
        "success": True,
        "count": len(items),
        "data": items
    })

# POST /items - Create new item
@app.route('/items', methods=['POST'])
def create_item():
    global next_id

    data = request.get_json()
    if not data or not data.get('name') or not data.get('name').strip():
        return jsonify({
            "success": False,
            "error": "Name is required"
        }), 400

    new_item = {
        "id": next_id,
        "name": data.get('name'),
        "description": data.get('description', ''),
        "createdAt": datetime.now().isoformat()
    }
    next_id += 1
    items.append(new_item)

    return jsonify({
        "success": True,
        "data": new_item
    }), 201

# GET /items/<id> - Get specific item
@app.route('/items/<int:item_id>', methods=['GET'])
def get_item(item_id):
    item = next((i for i in items if i["id"] == item_id), None)

    if not item:
        return jsonify({
            "success": False,
            "error": "Item not found"
        }), 404

    return jsonify({
        "success": True,
        "data": item
    })

# DELETE /items/<id> - Delete item
@app.route('/items/<int:item_id>', methods=['DELETE'])
def delete_item(item_id):
    global items
    item = next((i for i in items if i["id"] == item_id), None)

    if not item:
        return jsonify({
            "success": False,
            "error": "Item not found"
        }), 404

    items = [i for i in items if i["id"] != item_id]

    return jsonify({
        "success": True,
        "message": "Item deleted successfully",
        "data": item
    })

# Error handler for 404
@app.errorhandler(404)
def not_found(error):
    return jsonify({
        "success": False,
        "error": "Endpoint not found"
    }), 404

# Error handler for 500
@app.errorhandler(500)
def internal_error(error):
    return jsonify({
        "success": False,
        "error": "Internal server error"
    }), 500

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5000)
