# Codex 워크플로우 플러그인 1.2.1 설계 참조

저장소 루트의 `workflow-plugin-design.md`가 현재 권위 설계 문서다. `prompt-document-management.md`는 1.0.0에서 채택했던 단일 프롬프트 문서 방식의 역사적 조사 자료이며 현재 실행 계약이 아니다.

## 핵심 계약

1. 제품·기능 요청은 요구사항 조사, PRD 제시와 별도 승인, 기술 스펙 제시와 별도 승인, 새 구현 작업 순서로 진행한다.
2. PRD와 기술 스펙을 제시한 응답은 해당 단계에서 끝내며, 최초 요청의 미래형 구현 의사를 다음 단계 승인으로 간주하지 않는다.
3. 버그 수정 요청은 요청 증상과 직접 관련된 원인 및 필수 변경만 조사 후 같은 작업에서 수정·검증한다.
4. 조사 중 발견한 별도 버그, 독립 리팩터링과 새 기능은 범위 확장 승인 없이 수정하지 않는다.
5. UI 구현은 단일 UI 스펙/구현 문서를 제시한 뒤 별도 승인 후 새 구현 작업으로 인계한다.
6. 삭제, 배포, 외부 전송, 권한 변경과 비가역 데이터 변경은 해당 행동 직전에 승인받는다.
7. 문서 동일성은 경로, 버전 또는 상태, Git 변경 상태와 의미 검증으로 확인한다. 외부 계약이 요구할 때만 전체 SHA-256을 사용한다.
8. 새 작업 생성에 실패하면 승인된 문서에서 구성한 모델, 추론 강도, 사유와 완전한 프롬프트를 출력한다.

## 관련 문서

- 현재 설계: `workflow-plugin-design.md`
- 요구사항 진입점: `plugins/codex-workflow/skills/requirements-to-spec/SKILL.md`
- 버그 수정 경계: `plugins/codex-workflow/skills/bug-report-to-fix/SKILL.md`
- 구현 인계: `plugins/codex-workflow/skills/start-implementation-thread/SKILL.md`
- 역사적 조사: `plugins/codex-workflow/references/prompt-document-management.md`
