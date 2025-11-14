from django.http import JsonResponse
from django.urls import path
from django.views.decorators.csrf import csrf_exempt
from django.views.decorators.http import require_http_methods
from datetime import datetime
import json
import time

# In-memory data store
items = [
    {"id": 1, "name": "Sample Item 1", "description": "First sample item", "createdAt": datetime.now().isoformat()},
    {"id": 2, "name": "Sample Item 2", "description": "Second sample item", "createdAt": datetime.now().isoformat()}
]
next_id = [3]  # Using list to make it mutable in nested functions
request_count = [0]
start_time = time.time()

# Middleware to count requests
class RequestCounterMiddleware:
    def __init__(self, get_response):
        self.get_response = get_response

    def __call__(self, request):
        request_count[0] += 1
        response = self.get_response(request)
        return response

# GET / - Root endpoint with API info
@require_http_methods(["GET"])
def get_root(request):
    return JsonResponse({
        "message": "Django Test API",
        "version": "1.0.0",
        "requestCount": request_count[0],
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
@require_http_methods(["GET"])
def get_health(request):
    uptime = int(time.time() - start_time)
    return JsonResponse({
        "status": "healthy",
        "uptime": f"{uptime} seconds",
        "timestamp": datetime.now().isoformat()
    })

# GET /items - List all items
@require_http_methods(["GET"])
def get_items(request):
    return JsonResponse({
        "success": True,
        "count": len(items),
        "data": items
    })

# POST /items - Create new item
@csrf_exempt
@require_http_methods(["POST"])
def create_item(request):
    try:
        data = json.loads(request.body)
    except json.JSONDecodeError:
        return JsonResponse({
            "success": False,
            "error": "Invalid JSON"
        }, status=400)

    if not data.get('name') or not data.get('name').strip():
        return JsonResponse({
            "success": False,
            "error": "Name is required"
        }, status=400)

    new_item = {
        "id": next_id[0],
        "name": data.get('name'),
        "description": data.get('description', ''),
        "createdAt": datetime.now().isoformat()
    }
    next_id[0] += 1
    items.append(new_item)

    return JsonResponse({
        "success": True,
        "data": new_item
    }, status=201)

# GET /items/<id> - Get specific item
@require_http_methods(["GET"])
def get_item(request, item_id):
    item = next((i for i in items if i["id"] == item_id), None)

    if not item:
        return JsonResponse({
            "success": False,
            "error": "Item not found"
        }, status=404)

    return JsonResponse({
        "success": True,
        "data": item
    })

# DELETE /items/<id> - Delete item
@csrf_exempt
@require_http_methods(["DELETE"])
def delete_item(request, item_id):
    item = next((i for i in items if i["id"] == item_id), None)

    if not item:
        return JsonResponse({
            "success": False,
            "error": "Item not found"
        }, status=404)

    items.remove(item)

    return JsonResponse({
        "success": True,
        "message": "Item deleted successfully",
        "data": item
    })

urlpatterns = [
    path('', get_root),
    path('health', get_health),
    path('items', get_items),
    path('items/create', create_item),
    path('items/<int:item_id>', get_item),
    path('items/<int:item_id>/delete', delete_item),
]
