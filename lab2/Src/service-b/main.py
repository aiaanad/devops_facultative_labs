import os
import asyncio
import logging
from fastapi import FastAPI

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger("backend")

app = FastAPI()

VERSION = os.getenv("APP_VERSION", "v1")

@app.get("/health")
def health():
    return {"status": "ok"}

@app.get("/data")
async def get_data():
    logger.info(f"Received request on Backend {VERSION}")
    
    if VERSION == "v1":
        logger.warning("Version v1: simulating network latency (2s)...")
        await asyncio.sleep(2.0)
    
    return {
        "version": VERSION,
        "message": "Hello from Service B",
        "latency_simulated": True if VERSION == "v1" else False
    }