"""
FastAPI test application with granian ASGI server
"""
from fastapi import FastAPI

app = FastAPI(title="FastAPI Granian Test")


@app.get("/")
async def root():
    return {
        "message": "Hello from FastAPI with Granian!",
        "asgi_server": "granian",
        "framework": "FastAPI"
    }


@app.get("/health")
async def health():
    return {"status": "healthy"}
