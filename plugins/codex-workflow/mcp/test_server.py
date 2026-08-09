import importlib.util
import unittest
from pathlib import Path


SERVER_PATH = Path(__file__).with_name("server.py")
SPEC = importlib.util.spec_from_file_location("workflow_confirmation_server", SERVER_PATH)
server = importlib.util.module_from_spec(SPEC)
assert SPEC.loader is not None
SPEC.loader.exec_module(server)


class WorkflowConfirmationServerTest(unittest.TestCase):
    def setUp(self):
        server._requests.clear()

    def test_initialize_and_tool_contract(self):
        initialized = server.handle_request({
            "jsonrpc": "2.0",
            "id": 1,
            "method": "initialize",
            "params": {"protocolVersion": "2025-11-25"},
        })
        self.assertEqual(initialized["result"]["serverInfo"]["name"], server.SERVER_NAME)
        self.assertIn("모든 파일 수정", initialized["result"]["instructions"])

        tools = server.handle_request({"jsonrpc": "2.0", "id": 2, "method": "tools/list"})
        names = {tool["name"] for tool in tools["result"]["tools"]}
        self.assertEqual(names, {"show_workflow_confirmation", "submit_workflow_decision"})

    def test_approval_is_recorded_and_retry_is_idempotent(self):
        shown = server.show_workflow_confirmation({
            "requestId": "spec-v1",
            "title": "기술 스펙 승인",
            "message": "이 스펙으로 구현을 시작할까요?",
        })
        self.assertEqual(shown["structuredContent"]["status"], "pending")

        decided = server.submit_workflow_decision({"requestId": "spec-v1", "decision": "approved"})
        retried = server.submit_workflow_decision({"requestId": "spec-v1", "decision": "approved"})
        self.assertEqual(decided["structuredContent"]["status"], "approved")
        self.assertEqual(retried["structuredContent"]["decidedAt"], decided["structuredContent"]["decidedAt"])

    def test_conflicting_retry_fails(self):
        server.show_workflow_confirmation({
            "requestId": "plan-v1",
            "title": "실험 계획 승인",
            "message": "실험을 실행할까요?",
        })
        server.submit_workflow_decision({"requestId": "plan-v1", "decision": "approved"})
        conflict = server.submit_workflow_decision({"requestId": "plan-v1", "decision": "rejected"})
        self.assertTrue(conflict["isError"])
        self.assertEqual(conflict["structuredContent"]["status"], "approved")

    def test_request_id_cannot_be_reused_for_different_content(self):
        server.show_workflow_confirmation({
            "requestId": "same-id",
            "title": "PRD 승인",
            "message": "PRD v1을 승인할까요?",
        })
        with self.assertRaisesRegex(ValueError, "다른 승인 내용"):
            server.show_workflow_confirmation({
                "requestId": "same-id",
                "title": "기술 스펙 승인",
                "message": "기술 스펙 v1을 승인할까요?",
            })

    def test_widget_resource_is_utf8_mcp_app_html(self):
        result = server.handle_request({
            "jsonrpc": "2.0",
            "id": 3,
            "method": "resources/read",
            "params": {"uri": server.WIDGET_URI},
        })
        content = result["result"]["contents"][0]
        self.assertEqual(content["mimeType"], "text/html;profile=mcp-app")
        self.assertIn("승인", content["text"])
        self.assertIn("tools/call", content["text"])


if __name__ == "__main__":
    unittest.main()
