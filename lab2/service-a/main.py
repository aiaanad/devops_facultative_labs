import os
import httpx
import logging
from fastapi import FastAPI

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger("frontend")

app = FastAPI()

BACKEND_URL = os.getenv("BACKEND_URL")

if not BACKEND_URL:
    print("CRITICAL ERROR: BACKEND_URL environment variable is not set!")
    sys.exit(1)

@app.get("/health")
def health():
    return {"status": "ok"}

@app.get("/")
async def proxy_request():
    logger.info(f"Forwarding request to: {BACKEND_URL}/data")
    
    async with httpx.AsyncClient() as client:
        try:
            response = await client.get(f"{BACKEND_URL}/data", timeout=5.0)
            data = response.json()
            return {
                "proxy": "Service A",
                "backend_status": response.status_code,
                "received_from_backend": data
            }
        except Exception as e:
            logger.error(f"Error connecting to backend: {e}")
            return {"error": str(e)}