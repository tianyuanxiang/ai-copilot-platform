"""Generate Python gRPC stubs from the Go AI RPC protobuf contract.

Run this script whenever service/ai/rpc/pb/ai.proto changes. The generated
grpc module uses a package-relative import so it can be imported as
app.generated.ai_pb2_grpc from the FastAPI service.
"""

from pathlib import Path

from grpc_tools import protoc


ENGINE_ROOT = Path(__file__).resolve().parent.parent
PROTO_DIR = ENGINE_ROOT.parent / "rpc" / "pb"
PROTO_FILE = PROTO_DIR / "ai.proto"
OUTPUT_DIR = ENGINE_ROOT / "app" / "generated"


def main() -> None:
    OUTPUT_DIR.mkdir(parents=True, exist_ok=True)
    result = protoc.main(
        [
            "grpc_tools.protoc",
            f"-I{PROTO_DIR}",
            f"--python_out={OUTPUT_DIR}",
            f"--grpc_python_out={OUTPUT_DIR}",
            str(PROTO_FILE),
        ]
    )
    if result != 0:
        raise SystemExit(result)

    grpc_file = OUTPUT_DIR / "ai_pb2_grpc.py"
    content = grpc_file.read_text(encoding="utf-8")
    content = content.replace("import ai_pb2 as ai__pb2", "from . import ai_pb2 as ai__pb2")
    grpc_file.write_text(content, encoding="utf-8")


if __name__ == "__main__":
    main()
