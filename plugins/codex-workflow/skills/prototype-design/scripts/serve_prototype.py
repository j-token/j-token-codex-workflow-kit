"""앱 내 브라우저 검토를 위해 프로토타입의 src 디렉터리를 localhost로 제공합니다."""

from __future__ import annotations

import argparse
import functools
import http.server
from pathlib import Path


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("src", type=Path, help="프로토타입 src 디렉터리")
    parser.add_argument("--port", type=int, default=4173)
    args = parser.parse_args()

    root = args.src.resolve()
    if not root.is_dir():
        raise SystemExit(f"프로토타입 src 디렉터리가 없습니다: {root}")

    handler = functools.partial(http.server.SimpleHTTPRequestHandler, directory=str(root))
    server = http.server.ThreadingHTTPServer(("127.0.0.1", args.port), handler)
    print(f"{root}을(를) http://localhost:{args.port}에서 제공합니다.")
    try:
        server.serve_forever()
    except KeyboardInterrupt:
        print("\n프로토타입 서버를 중지합니다.")
    finally:
        server.server_close()


if __name__ == "__main__":
    main()
