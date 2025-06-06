from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
import uvicorn

app = FastAPI()

app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://127.0.0.1:5500"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

@app.get("/")
async def simple_json():
    return [
        {
            "id": 37918,
            "title": "Ура! Учёные придумали лекарство от рака!"
        },
        {
            "id": 38120,
            "title": "Завтра ожидается сильный дождь"
        },
    ]

if __name__ == "__main__":
    uvicorn.run("main:app", host="0.0.0.0", port=8000, reload=True)