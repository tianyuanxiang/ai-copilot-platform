# Python ASGI AI Engine

This service is the first runnable AI engine skeleton.

## Start

```powershell
cd D:\GoProject\project_new\ai-copilot-platform\service\ai\engine
python -m venv .venv
.\.venv\Scripts\activate
pip install -r requirements.txt
copy config.example.yaml config.yaml
uvicorn app.main:app --reload --host 0.0.0.0 --port 8001
```

## Smoke test

```powershell
python scripts\smoke_chat_stream.py
```

The `/v1/chat/stream` endpoint is an ASGI SSE stream. It uses mock tokens by default, so no LLM key is needed.