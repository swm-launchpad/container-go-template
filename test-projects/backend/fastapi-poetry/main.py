from fastapi import FastAPI, HTTPException, status
from pydantic import BaseModel
from typing import Optional
from datetime import datetime
import time

app = FastAPI()

# In-memory data store
items = [
    {"id": 1, "name": "Sample Item 1", "description": "First sample item", "createdAt": datetime.now().isoformat()},
    {"id": 2, "name": "Sample Item 2", "description": "Second sample item", "createdAt": datetime.now().isoformat()}
]
next_id = 3
request_count = 0
start_time = time.time()

# Request counter middleware
@app.middleware("http")
async def count_requests(request, call_next):
    global request_count
    request_count += 1
    response = await call_next(request)
    return response

# Pydantic models
class CreateItemRequest(BaseModel):
    name: str
    description: Optional[str] = ""

class Item(BaseModel):
    id: int
    name: str
    description: str
    createdAt: str

# GET / - Root endpoint with API info
@app.get("/")
def get_root():
    return {
        "message": "FastAPI Test API",
        "version": "1.0.0",
        "requestCount": request_count,
        "timestamp": datetime.now().isoformat(),
        "endpoints": [
            "GET /",
            "GET /health",
            "GET /items",
            "POST /items",
            "GET /items/{id}",
            "DELETE /items/{id}"
        ]
    }

# GET /health - Health check
@app.get("/health")
def get_health():
    uptime = int(time.time() - start_time)
    return {
        "status": "healthy",
        "uptime": f"{uptime} seconds",
        "timestamp": datetime.now().isoformat()
    }

# GET /items - List all items
@app.get("/items")
def get_items():
    return {
        "success": True,
        "count": len(items),
        "data": items
    }

# POST /items - Create new item
@app.post("/items", status_code=status.HTTP_201_CREATED)
def create_item(request: CreateItemRequest):
    global next_id

    if not request.name or not request.name.strip():
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="Name is required"
        )

    new_item = {
        "id": next_id,
        "name": request.name,
        "description": request.description or "",
        "createdAt": datetime.now().isoformat()
    }
    next_id += 1
    items.append(new_item)

    return {
        "success": True,
        "data": new_item
    }

# GET /items/{id} - Get specific item
@app.get("/items/{item_id}")
def get_item(item_id: int):
    item = next((i for i in items if i["id"] == item_id), None)

    if not item:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Item not found"
        )

    return {
        "success": True,
        "data": item
    }

# DELETE /items/{id} - Delete item
@app.delete("/items/{item_id}")
def delete_item(item_id: int):
    global items
    item = next((i for i in items if i["id"] == item_id), None)

    if not item:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Item not found"
        )

    items = [i for i in items if i["id"] != item_id]

    return {
        "success": True,
        "message": "Item deleted successfully",
        "data": item
    }
