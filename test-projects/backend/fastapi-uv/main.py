"""
FastAPI test application with uv dependency manager
"""
from fastapi import FastAPI

app = FastAPI(title="FastAPI uv Test")


@app.get("/")
async def root():
    return {
        "message": "Hello from FastAPI with uv!",
        "dependency_manager": "uv",
        "framework": "FastAPI"
    }


@app.get("/health")
async def health():
    return {"status": "healthy"}
