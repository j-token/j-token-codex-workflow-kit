#!/usr/bin/env python
"""워크플로우의 이진 승인 카드를 제공하는 무의존성 MCP stdio 서버."""

from __future__ import annotations

import json
import sys
import threading
import uuid
from datetime import datetime, timezone
from pathlib import Path
from typing import Any


SERVER_NAME = "codex-workflow-confirmation"
SERVER_VERSION = "1.0.0"
PROTOCOL_VERSION = "2025-11-25"
WIDGET_URI = "ui://codex-workflow/confirmation-v1.html"
WIDGET_PATH = Path(__file__).with_name("approval.html")
MAX_TEXT_LENGTH = 4_000

_lock = threading.Lock()
_requests: dict[str, dict[str, Any]] = {}


def _text(value: Any, field: str, *, default: str | None = None) -> str:
    if value is None:
        if default is not None:
            return default
        raise ValueError(f"{field} 값이 필요합니다.")
    if not isinstance(value, str):
        raise ValueError(f"{field} 값은 문자열이어야 합니다.")
    normalized = value.strip()
    if not normalized:
        raise ValueError(f"{field} 값은 비어 있을 수 없습니다.")
    if len(normalized) > MAX_TEXT_LENGTH:
        raise ValueError(f"{field} 값은 {MAX_TEXT_LENGTH}자 이하여야 합니다.")
    return normalized


def _snapshot(record: dict[str, Any]) -> dict[str, Any]:
    return {
        "kind": "workflow_confirmation",
        "requestId": record["requestId"],
        "title": record["title"],
        "message": record["message"],
        "approveLabel": record["approveLabel"],
        "rejectLabel": record["rejectLabel"],
        "approvalPrompt": record["approvalPrompt"],
        "rejectionPrompt": record["rejectionPrompt"],
        "status": record["status"],
        "decision": record.get("decision"),
        "decidedAt": record.get("decidedAt"),
    }


def _tool_result(data: dict[str, Any], message: str, *, is_error: bool = False) -> dict[str, Any]:
    return {
        "content": [{"type": "text", "text": message}],
        "structuredContent": data,
        "isError": is_error,
    }


def show_workflow_confirmation(arguments: dict[str, Any]) -> dict[str, Any]:
    title = _text(arguments.get("title"), "title")
    message = _text(arguments.get("message"), "message")
    approve_label = _text(arguments.get("approveLabel"), "approveLabel", default="승인")
    reject_label = _text(arguments.get("rejectLabel"), "rejectLabel", default="취소")
    approval_prompt = f"{title}을 승인합니다. 승인 대상과 영향: {message}"
    rejection_prompt = f"{title} 승인을 취소합니다. 현재 단계에서 멈춰 주세요."
    request_id = arguments.get("requestId")
    if request_id is None:
        request_id = str(uuid.uuid4())
    request_id = _text(request_id, "requestId")

    with _lock:
        record = _requests.get(request_id)
        if record is None:
            record = {
                "requestId": request_id,
                "title": title,
                "message": message,
                "approveLabel": approve_label,
                "rejectLabel": reject_label,
                "approvalPrompt": approval_prompt,
                "rejectionPrompt": rejection_prompt,
                "status": "pending",
                "createdAt": datetime.now(timezone.utc).isoformat(),
            }
            _requests[request_id] = record
        elif any(
            record[key] != expected
            for key, expected in {
                "title": title,
                "message": message,
                "approveLabel": approve_label,
                "rejectLabel": reject_label,
            }.items()
        ):
            raise ValueError("같은 requestId가 다른 승인 내용에 이미 사용되었습니다.")
        data = _snapshot(record)

    fallback = (
        f"{title}\n\n{message}\n\n"
        f"UI가 표시되지 않으면 채팅에서 '{approve_label}' 또는 '{reject_label}'라고 답해 주세요."
    )
    return _tool_result(data, fallback)


def submit_workflow_decision(arguments: dict[str, Any]) -> dict[str, Any]:
    request_id = _text(arguments.get("requestId"), "requestId")
    decision = arguments.get("decision")
    if decision not in {"approved", "rejected"}:
        raise ValueError("decision 값은 approved 또는 rejected여야 합니다.")

    with _lock:
        record = _requests.get(request_id)
        if record is None:
            return _tool_result(
                {
                    "kind": "workflow_confirmation",
                    "requestId": request_id,
                    "status": "unknown",
                    "decision": None,
                },
                "승인 요청을 찾을 수 없습니다. 새 승인 카드를 요청해 주세요.",
                is_error=True,
            )
        previous = record.get("decision")
        if previous is not None and previous != decision:
            return _tool_result(
                _snapshot(record),
                "이미 반대 결정으로 처리된 승인 요청입니다.",
                is_error=True,
            )
        record["decision"] = decision
        record["status"] = "approved" if decision == "approved" else "rejected"
        record["decidedAt"] = record.get("decidedAt") or datetime.now(timezone.utc).isoformat()
        data = _snapshot(record)

    message = "승인되었습니다." if decision == "approved" else "취소되었습니다."
    return _tool_result(data, message)


TOOLS = [
    {
        "name": "show_workflow_confirmation",
        "title": "워크플로우 승인 확인",
        "description": (
            "사용자가 파일 수정, 상태 변경, 구현·검증 명령, 문서, 계획, 범위 확장 또는 다음 단계의 진행 여부를 "
            "승인하거나 취소해야 하고 미해결 열린 질문이 없을 때 사용합니다. 직접적인 구현·수정 요청도 실제 "
            "변경 전 이 도구로 별도 승인받습니다. 두 선택지만 있는 이진 확인에만 사용합니다."
        ),
        "inputSchema": {
            "type": "object",
            "properties": {
                "requestId": {"type": "string", "description": "재시도 때 재사용할 선택적 안정 ID"},
                "title": {"type": "string", "maxLength": MAX_TEXT_LENGTH},
                "message": {"type": "string", "maxLength": MAX_TEXT_LENGTH},
                "approveLabel": {"type": "string", "default": "승인"},
                "rejectLabel": {"type": "string", "default": "취소"},
            },
            "required": ["title", "message"],
            "additionalProperties": False,
        },
        "outputSchema": {
            "type": "object",
            "properties": {
                "kind": {"type": "string"},
                "requestId": {"type": "string"},
                "title": {"type": "string"},
                "message": {"type": "string"},
                "approveLabel": {"type": "string"},
                "rejectLabel": {"type": "string"},
                "approvalPrompt": {"type": "string"},
                "rejectionPrompt": {"type": "string"},
                "status": {"type": "string"},
                "decision": {"type": ["string", "null"]},
                "decidedAt": {"type": ["string", "null"]},
            },
            "required": ["kind", "requestId", "title", "message", "approveLabel", "rejectLabel", "approvalPrompt", "rejectionPrompt", "status"],
        },
        "annotations": {
            "readOnlyHint": False,
            "destructiveHint": False,
            "openWorldHint": False,
            "idempotentHint": False,
        },
        "_meta": {
            "ui": {"resourceUri": WIDGET_URI},
            "openai/outputTemplate": WIDGET_URI,
            "openai/toolInvocation/invoking": "승인 카드를 준비하는 중…",
            "openai/toolInvocation/invoked": "승인 카드가 준비되었습니다.",
        },
    },
    {
        "name": "submit_workflow_decision",
        "title": "워크플로우 승인 결정 제출",
        "description": "승인 카드에서 사용자가 누른 승인 또는 취소 결정을 기록할 때만 사용합니다.",
        "inputSchema": {
            "type": "object",
            "properties": {
                "requestId": {"type": "string"},
                "decision": {"type": "string", "enum": ["approved", "rejected"]},
            },
            "required": ["requestId", "decision"],
            "additionalProperties": False,
        },
        "outputSchema": {
            "type": "object",
            "properties": {
                "kind": {"type": "string"},
                "requestId": {"type": "string"},
                "status": {"type": "string"},
                "decision": {"type": ["string", "null"]},
            },
            "required": ["kind", "requestId", "status", "decision"],
        },
        "annotations": {
            "readOnlyHint": False,
            "destructiveHint": False,
            "openWorldHint": False,
            "idempotentHint": True,
        },
        "_meta": {
            "ui": {"visibility": ["app"]},
            "openai/visibility": "private",
            "openai/toolInvocation/invoking": "결정을 기록하는 중…",
            "openai/toolInvocation/invoked": "결정이 기록되었습니다.",
        },
    },
]


def handle_request(request: dict[str, Any]) -> dict[str, Any] | None:
    method = request.get("method")
    request_id = request.get("id")

    if request_id is None:
        return None
    try:
        if method == "initialize":
            requested = request.get("params", {}).get("protocolVersion")
            result = {
                "protocolVersion": requested or PROTOCOL_VERSION,
                "capabilities": {
                    "tools": {"listChanged": False},
                    "resources": {"subscribe": False, "listChanged": False},
                },
                "serverInfo": {"name": SERVER_NAME, "version": SERVER_VERSION},
                "instructions": (
                    "모든 파일 수정, 상태 변경, 구현 명령과 검증 명령 전에 변경 범위와 검증을 제시하고 "
                    "show_workflow_confirmation으로 별도 승인받은 뒤 다음 턴부터 실행하세요. 직접적인 구현·수정 "
                    "요청도 아직 제시하지 않은 범위의 실행 승인으로 확대하지 마세요. "
                    "여러 선택이 필요한 질문과 위험 행동의 네이티브 권한 승인은 이 도구로 대체하지 마세요."
                ),
            }
        elif method == "ping":
            result = {}
        elif method == "tools/list":
            result = {"tools": TOOLS}
        elif method == "tools/call":
            params = request.get("params") or {}
            name = params.get("name")
            arguments = params.get("arguments") or {}
            if not isinstance(arguments, dict):
                raise ValueError("arguments 값은 객체여야 합니다.")
            if name == "show_workflow_confirmation":
                result = show_workflow_confirmation(arguments)
            elif name == "submit_workflow_decision":
                result = submit_workflow_decision(arguments)
            else:
                raise ValueError(f"알 수 없는 도구입니다: {name}")
        elif method == "resources/list":
            result = {
                "resources": [
                    {
                        "uri": WIDGET_URI,
                        "name": "workflow-confirmation-card",
                        "title": "워크플로우 승인 카드",
                        "description": "승인과 취소 버튼을 표시하는 인라인 카드",
                        "mimeType": "text/html;profile=mcp-app",
                    }
                ]
            }
        elif method == "resources/templates/list":
            result = {"resourceTemplates": []}
        elif method == "resources/read":
            uri = (request.get("params") or {}).get("uri")
            if uri != WIDGET_URI:
                raise ValueError(f"알 수 없는 리소스입니다: {uri}")
            result = {
                "contents": [
                    {
                        "uri": WIDGET_URI,
                        "mimeType": "text/html;profile=mcp-app",
                        "text": WIDGET_PATH.read_text(encoding="utf-8"),
                        "_meta": {
                            "ui": {
                                "prefersBorder": True,
                                "csp": {"connectDomains": [], "resourceDomains": []},
                            },
                            "openai/widgetDescription": "사용자가 워크플로우 진행을 승인하거나 취소하는 카드입니다.",
                            "openai/widgetPrefersBorder": True,
                        },
                    }
                ]
            }
        elif method == "prompts/list":
            result = {"prompts": []}
        elif method == "logging/setLevel":
            result = {}
        else:
            return {
                "jsonrpc": "2.0",
                "id": request_id,
                "error": {"code": -32601, "message": f"지원하지 않는 메서드입니다: {method}"},
            }
        return {"jsonrpc": "2.0", "id": request_id, "result": result}
    except (ValueError, OSError) as error:
        return {
            "jsonrpc": "2.0",
            "id": request_id,
            "error": {"code": -32602, "message": str(error)},
        }
    except Exception as error:  # pragma: no cover - 최종 프로토콜 방어선
        print(f"서버 오류: {error}", file=sys.stderr, flush=True)
        return {
            "jsonrpc": "2.0",
            "id": request_id,
            "error": {"code": -32603, "message": "내부 서버 오류가 발생했습니다."},
        }


def main() -> None:
    for line in sys.stdin:
        if not line.strip():
            continue
        try:
            request = json.loads(line)
            if not isinstance(request, dict):
                raise ValueError("요청은 JSON 객체여야 합니다.")
            response = handle_request(request)
        except (json.JSONDecodeError, ValueError) as error:
            response = {
                "jsonrpc": "2.0",
                "id": None,
                "error": {"code": -32700, "message": str(error)},
            }
        if response is not None:
            print(json.dumps(response, ensure_ascii=False, separators=(",", ":")), flush=True)


if __name__ == "__main__":
    main()
