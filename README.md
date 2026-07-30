# j-token-workflow-kit

`j-token-workflow-kit`은 사용자가 최소한의 의도만 전달해도 Codex가 저장소와 필요한 외부 자료를 조사하고, 실행 계획과 새 작업용 프롬프트까지 준비하는 한국어 워크플로우 플러그인입니다.

현재 플러그인 버전: `1.0.0`

## 1.0.0에서 달라진 점

기존 흐름은 작업마다 PRD, 기술 스펙, 감사 보고서와 임시 문서를 만들고 단계마다 별도 승인을 요구했습니다. 1.0.0은 작업의 크기와 위험을 먼저 분류합니다.

- 빠른 작업: 문서를 만들지 않고 사용자의 현재 실행 요청에 따라 같은 작업에서 구현·검증합니다.
- 계획 작업: `.codex/prompts/active/<slug>.md` 한 파일에 사실, 조건, 결정, 계획, 검증, 모델, 추론 강도와 실행 프롬프트를 모읍니다.
- 민감 작업: 프롬프트 문서는 하나로 유지하고 삭제, 배포, 외부 전송, 권한 변경과 비가역 데이터 변경 직전에만 별도 승인을 받습니다.
- PRD, 기술 스펙과 독립 감사는 사용자가 명시적으로 요구하거나 외부 형식·고위험 계약에 실제로 필요할 때만 추가합니다.

## 작동 방식

```mermaid
flowchart LR
    A["최소 사용자 의도"] --> B["저장소·공식 자료 조사"]
    B --> C{"작업 등급"}
    C -- "빠른 작업" --> D["현재 작업에서 구현·검증"]
    C -- "계획·민감 작업" --> E["단일 프롬프트 문서"]
    E --> F["계획·모델·추론·프롬프트 검토"]
    F --> G["새 Codex 작업 생성"]
    G --> H["구현·검증·보고"]
    H --> I["프롬프트 문서 보관"]
```

## Codex 실행 프롬프트 문서

```text
.codex/prompts/
├── README.md
├── active/<slug>.md
├── supporting/<slug>/<명시적으로-요청된-문서>.md
└── archive/YYYY/<slug>.md
```

`README.md`는 상태·유형·태그·갱신일로 찾는 카탈로그입니다. 한 작업의 권위 있는 실행 프롬프트는 `active/<slug>.md` 하나이며, 완료되면 내용을 복제하지 않고 상태를 갱신해 `archive/YYYY/`로 옮깁니다. `supporting/`은 사용자가 별도 PRD·기술 스펙을 요청한 경우에만 사용합니다. 최초 프롬프트 생성에서는 카탈로그와 프롬프트 문서로 실제 파일 2개가 생기고, 이후 작업마다 프롬프트 문서 1개를 추가하며 카탈로그를 갱신합니다.

프롬프트 문서에는 다음 정보가 함께 들어갑니다.

- 의도와 완료 상태
- 범위와 비범위
- 확인된 사실, 추정과 근거
- 조건, 결정과 기각한 대안
- 단계별 구현·검증 계획
- 권장 모델, 추론 강도와 선택 사유
- 새 Codex 작업에 그대로 전달할 완전한 프롬프트
- 결과와 후속

모델과 추론 강도는 사용자가 새 작업을 승인하기 전에 수정할 수 있습니다. 승인 후 Codex 앱의 새 작업 도구가 있으면 해당 설정으로 작업을 만들고, 없거나 실패하면 동일한 모델·추론 강도·프롬프트를 코드 블록으로 출력합니다.

## OS별 바이너리

`codex-workflow` 바이너리는 프롬프트 문서를 로컬 브라우저에서 열고, 사용자가 모델·추론 강도·새 작업용 프롬프트를 수정한 뒤 승인하거나 취소하게 합니다. 서버는 `127.0.0.1`의 임의 포트에만 바인딩되며 외부 스크립트나 네트워크 UI 자산을 사용하지 않습니다.

| 운영체제 | 아키텍처 | 릴리스 파일 |
| --- | --- | --- |
| macOS | Apple Silicon | `codex-workflow-darwin-arm64` |
| macOS | Intel | `codex-workflow-darwin-x64` |
| Linux | arm64 | `codex-workflow-linux-arm64` |
| Linux | x64 | `codex-workflow-linux-x64` |
| Windows | arm64 | `codex-workflow-windows-arm64.exe` |
| Windows | x64 | `codex-workflow-windows-x64.exe` |

각 릴리스 파일에는 같은 이름의 `.sha256` sidecar가 함께 생성됩니다.

macOS와 Linux:

```sh
curl -fsSL https://raw.githubusercontent.com/j-token/j-token-codex-workflow-kit/main/scripts/install.sh | sh
```

Windows PowerShell:

```powershell
irm https://raw.githubusercontent.com/j-token/j-token-codex-workflow-kit/main/scripts/install.ps1 | iex
```

특정 버전 설치:

```sh
curl -fsSL https://raw.githubusercontent.com/j-token/j-token-codex-workflow-kit/v1.0.0/scripts/install.sh | sh -s -- --version 1.0.0
```

```powershell
& ([scriptblock]::Create((irm https://raw.githubusercontent.com/j-token/j-token-codex-workflow-kit/v1.0.0/scripts/install.ps1))) -Version 1.0.0
```

위 설치 URL은 `v1.0.0` tag와 GitHub Release가 게시된 뒤부터 사용할 수 있습니다. 설치 스크립트와 바이너리 버전을 함께 고정하므로 이후 `main`의 변경에 영향을 받지 않습니다.

프롬프트 검토:

```text
codex-workflow review .codex/prompts/active/<slug>.md --json
```

승인 결과는 `status`, `path`, `model`, `reasoningEffort`, `prompt`를 가진 JSON으로 stdout에 반환됩니다. 브라우저를 자동으로 열 수 없는 환경에서는 stderr에 표시된 로컬 URL을 직접 열 수 있습니다.

## 사용 예시

간단한 수정은 그대로 요청합니다.

```text
이 오타를 고치고 관련 테스트를 실행해줘.
```

조사와 별도 실행 작업이 필요한 목표는 프롬프트로 준비합니다.

```text
$setup-codex-prompt를 사용해 이 아이디어를 저장소와 공식 자료에서 조사하고 새 Codex 작업용 프롬프트를 만들어줘.
```

계획을 본 뒤 모델과 추론 강도를 바꿀 수 있습니다.

```text
모델은 gpt-5.6-sol, 추론 강도는 high로 바꾸고 이 계획으로 새 작업 시작해줘.
```

## 포함된 스킬

| 스킬 | 역할 |
| --- | --- |
| `setup-codex-prompt` | 최소 의도에서 사실·조건·계획·모델·추론 강도·실행 프롬프트를 단일 프롬프트 문서로 만듭니다. |
| `requirements-to-spec` | 거친 요구사항을 조사해 Codex 실행 프롬프트로 수렴시키고, 명시적으로 필요한 경우에만 보조 문서를 연결합니다. |
| `start-implementation-thread` | 승인된 모델·추론 강도·실행 프롬프트로 새 Codex 작업을 만들고 실패 시 같은 내용을 출력합니다. |
| `bug-report-to-fix` | 간단한 버그는 바로 수정하고 복잡한 버그만 실행 프롬프트 준비 흐름으로 올립니다. |
| `figma-flow-to-implementation` | UI 흐름과 에셋을 조사하고 작업 등급에 따라 바로 구현하거나 프롬프트로 계획합니다. |
| `workflow-composer` | 요구사항, 버그와 UI 작업을 조합하되 산출물을 단일 Codex 실행 프롬프트로 수렴시킵니다. |
| `prd-writer` | 명시적으로 요청된 장기 제품 범위 문서를 작성합니다. |
| `technical-spec-writer` | 명시적으로 요청된 장기 공개 기술 계약을 작성합니다. |
| `audit-technical-spec` | 명시 요청 또는 고위험 계약에만 독립 감사를 적용합니다. |
| `orchestrate-subagents` | 독립 병렬화나 검토 이점이 있는 작업만 하위 에이전트에 배정합니다. |
| `prototype-design` | 확정 전 아이디어를 Mermaid나 프로토타입으로 보여줍니다. |
| `cognitive-writing` | 사용자 대상 문서의 인지 부하와 Markdown 오류를 줄입니다. |
| `branch-rule`, `commit-rule`, `git-push-safety`, `pr-rule` | Git 브랜치, 커밋, push와 PR의 안전 규칙을 제공합니다. |

## 문서 관리 근거

1.0.0은 MediaWiki의 카탈로그·랜딩 페이지, OCLC FAST의 낮은 비용 패싯 분류, Diátaxis의 독자 요구 분리, ADR의 작은 장기 결정 기록과 Plannotator `setup goal`의 검증 가능한 팩트 구조를 비교해 설계했습니다.

자세한 비교와 채택·기각 사유는 [Codex 프롬프트 문서 관리 조사](plugins/codex-workflow/references/prompt-document-management.md)에 있습니다.

## 저장소 구조

```text
.agents/plugins/marketplace.json
plugins/codex-workflow/.codex-plugin/plugin.json
plugins/codex-workflow/skills/
plugins/codex-workflow/references/
cmd/codex-workflow/
internal/promptdoc/
internal/review/
scripts/install.sh
scripts/install.ps1
.github/workflows/release-binaries.yml
```

## 로컬 개발

```powershell
go test ./...
go build -o dist\codex-workflow.exe .\cmd\codex-workflow
```

Git tag `v*`를 push하면 GitHub Actions가 6개 OS·아키텍처 바이너리를 교차 컴파일하고 Linux x64 smoke test, SHA-256 생성과 GitHub Release 업로드를 수행합니다. 릴리스에는 `LICENSE`와 `THIRD_PARTY_NOTICES.md`도 포함됩니다. 릴리스는 자산 업로드가 끝날 때까지 초안으로 유지되며 같은 tag의 작업을 다시 실행해도 자산을 교체할 수 있습니다.

## 라이선스

[Apache License 2.0](LICENSE)

Plannotator를 포함한 참고 프로젝트의 저작권과 라이선스는 [서드파티 고지](THIRD_PARTY_NOTICES.md)에서 확인할 수 있습니다.
