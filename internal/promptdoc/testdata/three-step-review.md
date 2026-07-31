---
prompt: three-step-review
status: draft
type: ui
risk: standard
model: gpt-5.6-terra
model_reason: "일반 구현 작업에서 속도와 품질의 균형이 필요해 선택했습니다."
reasoning_effort: medium
created: 2026-07-31
updated: 2026-07-31
tags: [ui, review]
---

# 로컬 검토 화면 개편

## 의도와 완료 상태

프롬프트 문서를 사실 관계, 옵션 선택, 실행 계획 순서로 검토하고 승인 결과를 새 작업에 넘깁니다.

## 범위와 비범위

- 로컬 브라우저 안에서만 동작합니다.
- 외부 계정 연결과 공유 기능은 포함하지 않습니다.

## 팩트와 근거

| ID | 상태 | 내용 | 근거·확인 방법 |
| --- | --- | --- | --- |
| F1 | 확인됨 | 기존 화면은 raw Markdown을 pre 요소에 표시합니다. | internal/review/web/index.html |
| F2 | 확인됨 | 로컬 서버는 127.0.0.1의 임의 포트에 바인딩됩니다. | internal/review/server.go |
| F3 | 조건 | Mermaid와 Markdown 자산은 바이너리에 포함되어야 합니다. | README의 로컬 전용 실행 계약 |

## 조건과 결정

- Markdown은 서버에서 안전하게 렌더링합니다.
- Mermaid는 외부 CDN 없이 로컬 번들을 사용합니다.

## 선택이 필요한 항목

### Q1. 계획 화면의 정보 밀도를 어떻게 할까요?

- 권장: 문서 중심
- 선택 방식: multiple
- 옵션: **문서 중심** — 본문 폭을 제한해 긴 계획의 가독성을 우선합니다.
- 옵션: **넓은 작업대** — 본문과 실행 설정을 한 화면에 더 많이 표시합니다.

### Q2. 수정한 사실을 승인 결과에 어떻게 반영할까요?

- 권장: 프롬프트 앞에 요약
- 옵션: **프롬프트 앞에 요약** — 실행 에이전트가 사용자 결정을 가장 먼저 읽습니다.
- 옵션: **별도 JSON만 반환** — 기존 프롬프트 내용은 유지하지만 소비자가 JSON을 해석해야 합니다.

## 실행 계획

### 화면 흐름

```mermaid
flowchart LR
  A[사실 관계] --> B[옵션 선택]
  B --> C[계획 검토]
  C --> D[승인 결과]
```

### 수식 렌더링

인라인 수식 $E = mc^2$와 블록 수식을 다크·라이트 테마에서 확인합니다.

$$
\operatorname{softmax}(x_i) = \frac{e^{x_i}}{\sum_{j=1}^{n} e^{x_j}}
$$

### 검증

1. Markdown의 제목, 표, 목록, 코드 블록을 확인합니다.
2. Mermaid 다이어그램이 SVG로 바뀌는지 확인합니다.
3. 수정한 사실과 선택한 옵션이 승인 프롬프트에 포함되는지 확인합니다.

## 위험과 열린 질문

- 잘못된 Mermaid 구문은 원문과 오류 안내를 함께 표시합니다.

## Codex 실행 설정

- 모델: gpt-5.6-terra
- 추론 강도: medium

## Codex 실행 프롬프트

<!-- codex-workflow:prompt:start -->

```text
검토된 사실과 선택을 우선 적용해 UI 개편을 구현하고 검증하세요.
```

<!-- codex-workflow:prompt:end -->

## 결과와 후속

검토 대기 중입니다.
