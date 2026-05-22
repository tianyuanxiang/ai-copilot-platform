import json
import urllib.request


def main() -> None:
    payload = json.dumps({"question": "hello from smoke script"}).encode("utf-8")
    req = urllib.request.Request(
        "http://127.0.0.1:8001/v1/chat/stream",
        data=payload,
        headers={
            "Content-Type": "application/json",
            "Accept": "text/event-stream",
        },
        method="POST",
    )

    with urllib.request.urlopen(req, timeout=20) as response:
        for raw in response:
            print(raw.decode("utf-8").rstrip())


if __name__ == "__main__":
    main()