# Codex 워크플로우 플러그인 1.0.0 설계 참조

저장소 루트의 `workflow-plugin-design.md`가 현재 권위 설계 문서다. 이 참조 파일은 스킬이 설계 배경을 찾는 진입점이며 상세 문서 관리 조사와 출처는 `prompt-document-management.md`에 있다.

## 핵심 계약

1. 빠른 저위험 작업은 문서를 만들지 않고 현재 사용자 요청에 따라 수행한다.
2. 계획·민감 작업은 `.codex/prompts/active/<slug>.md` 한 파일에 실행 진실을 모은다.
3. 새 작업 직전에 계획·모델·추론 강도·프롬프트를 한 번 검토한다.
4. 비가역·외부 영향 행동은 그 행동 직전에 승인받는다.
5. PRD, 기술 스펙, 감사 보고서와 하위 에이전트 검토는 필요할 때만 추가한다.
6. 새 Codex 작업 생성에 실패하면 모델, 추론 강도, 사유와 완전한 프롬프트를 출력한다.
7. 계획 작업은 가능하면 OS별 `codex-workflow review <prompt.md> --json` 브라우저 UI에서 검토한다.
8. 바이너리가 없으면 채팅 검토로 fallback하며 작업을 차단하지 않는다.

## 관련 문서

- 저장소 설계: `workflow-plugin-design.md`
- 조사와 출처: `plugins/codex-workflow/references/prompt-document-management.md`
- 핵심 스킬: `plugins/codex-workflow/skills/setup-codex-prompt/SKILL.md`
- 바이너리 진입점: `cmd/codex-workflow/main.go`
- 릴리스 워크플로우: `.github/workflows/release-binaries.yml`
