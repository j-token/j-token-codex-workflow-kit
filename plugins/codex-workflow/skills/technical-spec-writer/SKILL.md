---
name: technical-spec-writer
description: 사용자가 기술 스펙, 구현 스펙, 아키텍처 스펙 또는 IPC/API/SDK/CLI/FFI/JNI/native shell/런타임/빌드/테스트 계약 문서를 요청했을 때 작성하거나 다듬습니다. 제품·기능 워크플로우에서는 별도 메시지로 승인된 PRD를 입력으로 기술 스펙을 작성하고, 문서를 제시한 뒤 별도 구현 승인을 기다립니다.
---

# 기술 스펙 작성기

## 출력 언어

사용자가 출력 언어를 명시적으로 지정하면 해당 언어를 사용합니다. 지정하지 않으면 이 스킬이 만드는 모든 사용자 대상 출력, 문서, 프롬프트, 보고서, 계획, 스펙 및 기타 산출물을 한국어로 작성합니다. 제목, 섹션, 레이블, 표, 체크리스트, 다이어그램, 템플릿에도 같은 언어를 사용합니다. 코드, 명령어, 파일 경로, 식별자, API 이름, 모델 ID, 프로토콜 이름과 필수 고유명사는 번역하지 않습니다.

## 작성 문서 첨부

요청된 모든 문서의 작성 또는 갱신을 마친 뒤 최종 응답에 각 문서를 `[문서 이름](절대 경로)` 형식의 Markdown 링크로 반드시 첨부합니다. 링크에는 독립된 디렉터리 세그먼트로 구성된 절대 경로를 사용하고, 디렉터리 구분자를 생략하거나 경로 세그먼트를 붙여 쓰지 않습니다. Windows 경로를 문자열로 표시할 때는 `C:\\Users\\WinUser\\Documents\\폴더\\문서.md`처럼 각 역슬래시를 두 번 씁니다.

## 문서 동일성 확인

- 기본 확인은 문서의 절대 경로, 문서 버전 또는 상태, Git 변경 상태와 필요한 의미 검증으로 수행합니다.
- 줄바꿈 형식이나 공백만 달라졌다는 이유로 문서 전체를 다시 작성하지 않습니다. 먼저 의미가 바뀌었는지 확인합니다.
- 외부 시스템, 감사 규정 또는 사용자가 전체 SHA-256을 명시적으로 요구할 때만 전체 해시를 계산하고 비교합니다.

## 목적

거친 기술 설계 노트나 PRD를 구현 가능한 기술 스펙으로 바꾼다. 기본 출력은 런타임 구조, 프로토콜, API, 플랫폼 경계, 상태 모델, 보안, 테스트 기준이 한 문서에 들어간 구현 계약서다.

## PRD와의 차이

- PRD는 "왜 만드는가 / 무엇을 해결하는가 / 어떤 범위인가"를 결정한다.
- 기술 스펙은 "어떤 타입, 상태, 프로토콜, API, 경계 조건으로 구현할 것인가"를 결정한다.
- PRD의 기능 요구사항은 스펙에서 schema, state transition, public API, platform contract, acceptance test로 내려온다.
- 제품·기능 워크플로우의 기술 스펙은 별도 메시지로 승인된 PRD를 입력으로 사용한다.
- 사용자가 독립 기술 스펙 작성만 요청한 경우에는 직접 확정한 범위를 입력으로 사용할 수 있다.

## 작성 원칙

- 스펙은 설명문이 아니라 구현 계약이다.
- 스펙의 목표는 구현 전 불확실성 0건이 아니라 안전한 구현을 시작할 수 있는 최소 충분 계약이다.
- 스펙 작성 직전에 입력 PRD 또는 사용자 확정 범위를 다시 읽고, 대화의 기억이나 요약만으로 작성하지 않는다.
- 입력에 `FR-*`가 있으면 타입, 상태 전이, public API, 플랫폼 계약 또는 수용 테스트에 추적 가능하게 연결한다.
- 관련 코드, interface, data model, build·배포 경로, 테스트와 필요한 외부 공식 문서를 직접 조사한다.
- 저장소에서 검증한 사실, 설계 선택과 미해결 질문을 분리하고 추정을 저장소 사실처럼 쓰지 않는다.
- 파일이나 모듈은 직접 확인한 뒤에만 확정 경로로 적고, 확인 전에는 후보라고 표시한다.
- 공개되거나 경계를 넘는 protocol, message, API, config, error는 이름과 타입을 명시한다.
- 공개·교차 모듈 계약에는 입력, 출력, validation, error 동작, ownership과 compatibility 기대값을 정의한다.
- 상태를 가진 시스템은 상태 목록과 전이 규칙을 반드시 쓴다.
- 경계를 넘는 호출은 입력, 출력, ownership, threading, failure를 함께 정의한다.
- 플랫폼별 차이는 숨기지 말고 iOS/Android, client/server, JS/Rust처럼 섹션을 나눠 쓴다.
- `v0.1`, `v0.2`, `v0.3`처럼 버전 범위와 deferred feature를 구분한다.
- 수치 목표는 초기 가정임을 명시하고 측정 방법을 함께 쓴다.
- 변경 단계는 의존 순서대로 적고 migration과 rollback 영향을 함께 정의한다.
- `에러를 적절히 처리한다`처럼 모호하게 쓰지 말고 failure mode, 사용자에게 보이는 동작, logging과 test를 명시한다.
- 선택 개선안과 미해결 결정은 별도 항목에 두고 승인된 범위를 조용히 확대하지 않는다.
- 단일 모듈 내부 자료구조, private 함수, 중간 artifact 세부 schema, 내부 상태 이름, 외부 동작에 영향이 없는 tie-break는 구현 TODO·테스트·ADR 또는 spike로 이관한다.
- 숫자와 기본값은 수용 기준·통계적 타당성·재현성·상호 운용에 영향을 줄 때만 스펙에서 고정한다. 구현으로 안전하게 검증할 값은 설정 기본값과 테스트로 둔다.
- 아직 존재하지 않는 산출물의 모든 필드와 SHA 소유 관계를 미리 완전하게 서술하지 않는다. 외부 검증에 필요한 최소 lineage와 불변식만 정의한다.
- 문서 분량, 섹션 수, 감사 발견 0건을 품질 목표로 삼지 않는다.
- 불확실한 기술 선택은 관련 섹션 안에 "미해결 기술 결정" 블록으로 붙인다.
- 문서 끝의 미해결 기술 결정 섹션은 전체 복사본이 아니라 섹션별 인덱스와 우선순위 요약으로만 쓴다.
- 보안·데이터·호환성·배포처럼 실패 비용이 크거나 사용자가 독립 감사를 요청한 경우에만 `audit-technical-spec`을 적용한다.
- 제품·기능 기술 스펙을 제시한 뒤 현재 턴을 끝낸다. 이후 별도 사용자 메시지에서 스펙 경로 또는 버전을 확인하며 구현을 승인한 경우에만 `start-implementation-thread`를 적용한다.
- 구현 시작을 막는 미해결 기술 결정은 스펙 안에만 기록하지 않고, 기술 스펙 링크와 함께 사용자 대화에도 먼저 제시한다.

## 작성 흐름

1. 제품·기능 워크플로우에서는 사용자가 별도 메시지로 승인한 PRD 경로와 버전을 확인하고 파일에서 다시 읽는다. 독립 기술 스펙 요청은 사용자가 확정한 범위를 다시 읽는다.
2. 입력을 분류한다.
   - PRD 기반 구현 스펙
   - 아키텍처/런타임 스펙
   - IPC/API/SDK 계약 스펙
   - 플랫폼 shell 또는 native boundary 스펙
   - 기존 초안 정리/확장
3. 문서 헤더와 목표를 고정한다.
4. 시스템 레이어와 모듈 구성을 먼저 그린다.
5. 상태 모델과 session/resource model을 정의한다.
6. 프로토콜, API, config, error model을 타입 중심으로 작성한다.
7. 플랫폼 경계와 lifecycle mapping을 작성한다.
8. 보안, logging, build, test, performance, versioning을 정리한다.
9. 각 섹션의 미해결 기술 결정을 해당 섹션에 기록하고, 마지막에는 v0.1 구현 체크리스트와 섹션별 미해결 기술 결정 인덱스를 둔다.
10. 파일을 저장해 사용자에게 경로와 버전을 제시합니다. 구현 시작을 막는 미해결 기술 결정이 있으면 질문 전문·선택지·권장안·영향과 답변 방법을 함께 제시하고 현재 턴을 끝냅니다. 위험 기반 조건이 맞을 때만 `audit-technical-spec`을 적용합니다.
11. 별도 사용자 메시지에서 기술 스펙 승인과 미해결 기술 결정의 답변을 확인합니다. 사용자가 먼저 승인했지만 질문을 아직 대화에 제시하지 않았다면 승인 상태는 보존하고 질문을 먼저 제시합니다.
12. 사용자의 직접 선택 또는 권장안 일괄 승인을 기술 스펙에 반영한 뒤에만 `start-implementation-thread`에 전달합니다.

## 출력 위치

- 사용자가 파일 작성을 요청하면 `.codex/temp/<product>-technical-spec.md`에 저장한다. 사용자가 다른 위치를 명시하면 그 위치를 따른다.
- 사용자가 채팅 답변만 원하면 파일을 만들지 않고 본문만 출력한다.

## 작성 원칙 참조

스펙 본문은 `cognitive-writing` 스킬의 원칙(인지 부하 최소화, GitHub-flavored Markdown 규칙)을 따라 작성한다.

## 구현 전 열린 질문 게이트

구현을 막는 미해결 기술 결정이 없고 승인·취소만 남으면 `../../references/workflow-confirmation-ui.md`에 따라 버튼형 승인 카드를 제시합니다. 버튼이 게시한 후속 메시지는 별도 사용자 구현 승인 메시지로 처리합니다.

구현 시작을 막는 미해결 기술 결정이 하나라도 있으면 기술 스펙을 첨부하는 최종 응답에 아래 정보를 포함합니다.

1. 기술 스펙 절대 경로 링크와 문서 버전 또는 상태
2. `Q1`, `Q2`처럼 식별 가능한 질문 전문
3. 각 질문의 선택지, 현재 권장안과 선택에 따른 영향
4. 다음 두 답변 방법
   - `권장안대로 승인`: 방금 대화에 제시된 모든 권장안을 선택하고 기술 스펙을 승인
   - `Q1: ..., Q2: ...`: 질문별로 직접 선택하고 그 내용으로 기술 스펙을 갱신

사용자가 질문을 보지 못한 채 `승인`만 보낼 수 있으므로 다음 규칙을 지킵니다.

- `승인`만으로 미해결 기술 결정의 답을 추정하지 않습니다.
- 승인 메시지를 받은 시점에 미제시 또는 미응답 질문을 발견하면 기술 스펙 링크와 질문 전체를 다시 제시하고 그 턴에서 멈춥니다.
- 사용자가 이미 한 기술 스펙 승인은 유지합니다. 질문 답변을 받은 뒤 별도의 기술 스펙 재승인을 반복해서 요구하지 않습니다.
- `권장안대로 승인`은 질문을 대화에 제시한 뒤 받은 답일 때만 유효합니다. 문서 안에만 있던 권장안에 대한 사전 포괄 승인은 인정하지 않습니다.
- 질문별 직접 답변을 받으면 결정과 영향을 관련 섹션 및 결정 사항 요약에 반영하고 버전 또는 상태를 갱신합니다.
- 구현을 막는 미해결 기술 결정이 없으면 `구현 전 필수 열린 질문: 없음`이라고 명시합니다.

## 기본 문서 구조

기술 스펙은 아래의 압축 구조를 기본값으로 사용합니다. 주제에 필요한 섹션만 추가하며 빈 섹션, 동일 내용의 요약 반복과 구현 내부 세부사항을 채우기 위해 문서를 확장하지 않습니다.

```md
# 기술 스펙 문서: <제품/시스템명>

**문서 버전:** v0.1 Draft
**작성일:** YYYY-MM-DD
**가칭:** <이름>
**대상 플랫폼:** <플랫폼>
**목표:** <한 줄 구현 목표>

---

# 1. 개요

## 1.1 요구사항 출처와 범위

## 1.2 기존 시스템 조사 결과

## 1.3 목표와 비목표

# 2. 핵심 설계 원칙

# 3. 시스템 아키텍처

# 4. 공개·교차 모듈 계약

## 4.1 상태와 생명주기

## 4.2 API·프로토콜·데이터 계약

## 4.3 오류와 실패 동작

# 5. 보안·데이터·호환성 경계

# 6. 설정과 운영 계약

# 7. 구현 계획과 요구사항 추적

# 8. 검증 및 수용 기준

# 9. 미해결 결정과 구현 이관 항목

# 10. 참고 자료
```

## 섹션 작성 규칙

### 섹션별 미해결 기술 결정

미해결 기술 결정은 가장 관련 있는 섹션 안에 둔다. 독자가 상태 모델, 프로토콜, boundary, 플랫폼 lifecycle 같은 맥락을 기억한 채 문서 끝의 거대한 질문 목록으로 이동하게 만들지 않는다.

권장 블록:

```md
#### 미해결 기술 결정

- **결정 필요**: 아직 정하지 못한 기술 선택
- **영향 범위**: 이 결정이 막는 구현 영역
- **후보**: 선택지 A / 선택지 B / 선택지 C
- **현재 권장안**: 있으면 1줄로 작성
- **결정 필요 시점**: prototype 전 / v0.1 구현 전 / 릴리즈 전
```

규칙:

- 해당 섹션과 직접 관련된 결정만 둔다.
- 여러 섹션에 걸친 결정은 가장 먼저 막히는 섹션에 두고, 마지막 인덱스에서 관련 섹션을 함께 표시한다.
- 문제가 1개뿐이면 짧은 bullet 하나로 줄인다.
- 문서 끝에는 상세를 반복하지 않고 섹션 링크, 우선순위, 결정 필요 시점만 요약한다.

### 개요

PRD와의 차이를 짧게 설명하고, 이 문서가 정의하는 구현 계약을 bullet로 나열한다.

### 핵심 설계 원칙

기술 판단을 명령형 원칙으로 쓴다.

예시:

- Rust first
- WebView UI first
- native shell은 thin transport layer
- IPC가 제품의 핵심 계약
- macro-heavy 설계보다 명시적 구조 우선

### 시스템 아키텍처

레이어 다이어그램을 반드시 포함한다.

```txt
Web UI
  ↓
TypeScript SDK
  ↓
Injected Bridge
  ↓
Native Shell
  ↓
Runtime Boundary
  ↓
Core Runtime
```

필요하면 repository 구조도 포함한다.

### 상태 모델

상태가 있는 구성요소마다 다음을 쓴다.

- 상태 다이어그램
- 상태 정의 표
- 전이 규칙
- invalid transition 처리
- resource cleanup 규칙

예시 상태:

- runtime state
- WebView session state
- connection state
- command execution state
- build state

### 세션/리소스 모델

reload, recreate, reconnect, dispose가 가능한 시스템은 session id 또는 resource handle 정책을 명시한다.

반드시 정의할 항목:

- id 형식
- lifetime
- ownership
- stale message 처리
- pending request cleanup
- destroyed resource 접근 시 에러

### 프로토콜 스펙

프로토콜은 타입 정의를 중심으로 쓴다.

필수 항목:

- protocol name/version
- envelope
- request
- response
- event
- control
- log/diagnostic message
- 예시 JSON

타입 예시는 TypeScript 또는 JSON Schema 중 하나를 기본으로 한다. 구현 언어가 다르면 Rust/Kotlin/Swift 타입을 추가한다.

### 프로토콜 제약 조건

제약은 표와 규칙으로 쓴다.

예시:

- message size
- timeout
- request id 형식
- command name regex
- payload encoding
- version negotiation
- stale response 처리

### 오류 모델

모든 실패는 structured error로 정의한다.

```ts
type SpecError = {
  code: string
  message: string
  data?: unknown
}
```

error code는 union type 또는 표로 정리한다. 각 code는 발생 조건과 처리 방식을 가져야 한다.

### SDK/API 스펙

public API는 signature, type, 동작 순서, 실패 조건을 함께 쓴다.

각 API는 아래 형식을 따른다.

````md
## <API 이름>

Signature:

```ts
declare function example(input: Input): Promise<Output>
```

동작 순서:

1. ...

실패 조건:

- ...
````

### 경계 스펙

FFI/JNI/native bridge처럼 경계를 넘는 부분은 가장 엄격하게 쓴다.

필수 항목:

- boundary 방향
- 함수 signature
- memory ownership
- pointer/string lifetime
- callback lifetime
- threading rule
- main thread dispatch 여부
- dispose 규칙

### 플랫폼 스펙

플랫폼별 섹션은 같은 구조로 반복한다.

- 기술 구성
- 주요 클래스/모듈
- bridge injection 또는 transport setup
- lifecycle mapping
- navigation/resource policy
- platform-specific failure
- acceptance criteria

### 에셋/리소스 로딩

개발 모드와 릴리즈 모드를 분리한다.

반드시 정의할 항목:

- mode
- dev URL/origin
- release origin
- bundled asset path
- route fallback
- external navigation 처리
- release build에서 금지되는 값

### 설정 스펙

설정 파일은 예시와 validation rule을 함께 쓴다.

```toml
[app]
name = "Example"

[runtime]
version = "0.1.0"
```

validation rule은 build 전 실패해야 하는 조건을 중심으로 쓴다.

### 보안 스펙

다음 원칙으로 작성한다.

> 외부 또는 UI 계층에서 들어오는 모든 message는 untrusted input이다.

필수 검토 항목:

- origin allowlist
- navigation policy
- command/API allowlist
- schema validation
- payload size limit
- dev/release 보안 차이
- injection policy
- file/resource 접근 제한

### 테스트 스펙

테스트는 계층별로 나눈다.

- unit test
- integration test
- SDK test
- platform E2E test
- lifecycle/recreation test
- security negative test
- performance measurement

각 테스트 항목은 "무엇을 검증하는가"가 보이게 쓴다.

### 성능 스펙

목표 수치와 측정 방법을 함께 쓴다.

```md
| 항목 | 목표 |
|---|---:|
| simple round-trip | 50ms 이하 |
```

측정 이벤트 이름과 trace field를 정의한다.

### 버전 관리/호환성

runtime version, protocol version, SDK compatibility, breaking change 정책을 쓴다.

### 기존 시스템 조사 결과

- 확인한 진입점, interface, data model, build·배포 경로와 테스트를 근거 경로와 함께 적는다.
- 저장소 사실과 새 설계 결정을 분리한다.
- 확인하지 않은 파일이나 모듈은 확정 경로가 아니라 후보로 표시한다.

### 관측 가능성

- 주요 이벤트, log level, metric, trace field와 진단에 필요한 correlation identifier를 정의한다.
- 민감정보를 기록하지 않는 기준과 운영자가 실패를 확인하는 방법을 적는다.

### 마이그레이션과 롤백

- 기존 data, protocol, API, config와 배포 방식에 미치는 영향을 적는다.
- 호환성 경계, migration 순서, 중단 또는 rollback 조건과 복구 절차를 정의한다.

### 구현 계획과 요구사항 추적

- 각 중요 요구사항을 담당 component, contract, 구현 단계와 검증 항목에 연결한다.
- 단계별 의존성, 완료 조건과 rollback 영향을 적는다.

### 참고 자료

- 설계 판단에 사용한 승인 문서, 저장소 경로와 공식 외부 문서를 연결한다.

### 결정 사항 요약

현재 확정된 기술 결정을 표로 정리한다.

```md
| 항목 | 스펙 |
|---|---|
| IPC v0.1 | JSON request/response |
| Config | `rvm.toml` |
```

### 섹션별 미해결 기술 결정 인덱스

문서 끝의 인덱스는 상세 목록이 아니라 리뷰 순서를 잡기 위한 요약이다. 상세 내용은 관련 섹션 안의 `미해결 기술 결정` 블록에 둔다.

예시:

| 섹션 | 우선순위 | 미해결 기술 결정 | 결정 필요 시점 |
|---|---|---|---|
| `10. Core/Runtime API 스펙` | 높음 | Rust async runtime 선택 | prototype 전 |
| `12. Boundary 스펙` | 높음 | iOS Rust packaging 방식 선택 | v0.1 구현 전 |
| `15. Asset/Resource Loading 스펙` | 중간 | release asset origin 구현 방식 선택 | 릴리즈 전 |

## 품질 체크리스트

출력 전 아래를 확인한다.

- [ ] 문서가 PRD가 아니라 구현 스펙으로 읽히는가
- [ ] 구현을 시작하는 데 필요하지 않은 private 세부사항을 TODO·테스트·ADR·spike로 이관했는가
- [ ] 문서 분량과 감사 발견 0건을 품질 목표로 삼지 않았는가
- [ ] 핵심 계약이 타입, schema, signature로 표현되어 있는가
- [ ] 상태 모델과 전이 규칙이 있는가
- [ ] session/resource lifetime과 cleanup 규칙이 있는가
- [ ] 에러 코드가 발생 조건과 연결되는가
- [ ] platform-specific lifecycle과 failure가 빠지지 않았는가
- [ ] boundary ownership/threading 규칙이 명확한가
- [ ] dev/release 보안 차이가 명시되어 있는가
- [ ] v0.1 범위와 deferred feature가 구분되어 있는가
- [ ] 테스트 스펙이 unit/integration/E2E/performance로 나뉘어 있는가
- [ ] 미해결 기술 결정이 관련 섹션 안에 있고 마지막에는 인덱스만 있는가
- [ ] 승인 요구사항이 구현 책임과 검증 항목으로 추적되는가
- [ ] 기존 시스템 사실과 새 설계 선택이 구분되는가
- [ ] compatibility, migration, rollback과 observability가 필요한 범위에서 정의됐는가
- [ ] 저장된 기술 스펙을 `audit-technical-spec`에 전달할 경로와 버전 또는 상태를 확인했는가
- [ ] 외부 계약이 요구하지 않는데 전체 SHA-256이나 줄바꿈 차이를 필수 게이트로 사용하지 않았는가
- [ ] 구현을 막는 미해결 기술 결정을 기술 스펙 링크와 함께 대화에 제시했는가
- [ ] 미해결 기술 결정이 없으면 없다고 명시했는가
- [ ] 일반 `승인`을 권장안 선택으로 처리하지 않았는가
