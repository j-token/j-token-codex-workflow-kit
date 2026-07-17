"""Serve a prototype's src directory over localhost for in-app browser review."""

from __future__ import annotations

import argparse
import functools
import http.server
from pathlib import Path


def main() -> None:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("src", type=Path, help="Prototype src directory")
    parser.add_argument("--port", type=int, default=4173)
    args = parser.parse_args()

    root = args.src.resolve()
    if not root.is_dir():
        raise SystemExit(f"Prototype src directory does not exist: {root}")

    handler = functools.partial(http.server.SimpleHTTPRequestHandler, directory=str(root))
    server = http.server.ThreadingHTTPServer(("127.0.0.1", args.port), handler)
    print(f"Serving {root} at http://localhost:{args.port}")
    try:
        server.serve_forever()
    except KeyboardInterrupt:
        print("\nStopping prototype server")
    finally:
        server.server_close()


if __name__ == "__main__":
    main()
