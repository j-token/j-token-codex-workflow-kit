---
name: technical-spec-writer
description: 기술 스펙 문서, 구현 스펙, 아키텍처 스펙, IPC/API/SDK/CLI/FFI/JNI/native shell/런타임/빌드/테스트 명세를 작성하거나 다듬을 때 사용한다. PRD가 "왜/무엇"을 다룬 뒤, 구현자가 따라야 할 "어떻게"와 계약을 문서화할 때 사용한다.
---

# Technical Spec Writer

## 목적

거친 기술 설계 노트나 PRD를 구현 가능한 기술 스펙으로 바꾼다. 기본 출력은 런타임 구조, 프로토콜, API, 플랫폼 경계, 상태 모델, 보안, 테스트 기준이 한 문서에 들어간 구현 계약서다.

## PRD와의 차이

- PRD는 "왜 만드는가 / 무엇을 해결하는가 / 어떤 범위인가"를 결정한다.
- 기술 스펙은 "어떤 타입, 상태, 프로토콜, API, 경계 조건으로 구현할 것인가"를 결정한다.
- PRD의 기능 요구사항은 스펙에서 schema, state transition, public API, platform contract, acceptance test로 내려온다.
- 기술 스펙 작성에는 사용자가 확인한 PRD 경로 또는 버전이 선행되어야 한다.
- 확정된 PRD가 없는 제품·기능 요청은 `prd-writer`로 먼저 보낸다. 사용자가 "스펙", "구현 스펙", "기술 명세"를 직접 요청해도 제품 범위가 정해진 PRD가 없으면 기술 스펙부터 작성하지 않는다.
- 이미 확정된 PRD를 기반으로 프로토콜, API 계약 또는 구현 스펙을 요청하면 이 스킬을 적용한다.

## 작성 원칙

- 스펙은 설명문이 아니라 구현 계약이다.
- 스펙 작성 직전에 확정된 PRD를 파일에서 다시 읽고, 대화의 기억이나 요약만으로 작성하지 않는다.
- PRD의 모든 `FR-*`를 스펙의 타입, 상태 전이, public API, 플랫폼 계약 또는 수용 테스트에 추적 가능하게 연결한다.
- 모든 protocol, message, API, config, error는 이름과 타입을 명시한다.
- 상태를 가진 시스템은 상태 목록과 전이 규칙을 반드시 쓴다.
- 경계를 넘는 호출은 입력, 출력, ownership, threading, failure를 함께 정의한다.
- 플랫폼별 차이는 숨기지 말고 iOS/Android, client/server, JS/Rust처럼 섹션을 나눠 쓴다.
- `v0.1`, `v0.2`, `v0.3`처럼 버전 범위와 deferred feature를 구분한다.
- 수치 목표는 초기 가정임을 명시하고 측정 방법을 함께 쓴다.
- 불확실한 기술 선택은 관련 섹션 안에 "미해결 기술 결정" 블록으로 붙인다.
- 문서 끝의 미해결 기술 결정 섹션은 전체 복사본이 아니라 섹션별 인덱스와 우선순위 요약으로만 쓴다.
- 기술 스펙을 제시한 뒤 현재 턴을 끝낸다. 이후 별도 사용자 메시지에서 스펙 경로 또는 버전을 확인하며 명시적으로 구현을 승인한 경우에만 현재 작업에서 직접 코드를 수정하지 말고 `start-implementation-thread`에 확정된 기술 스펙 경로를 전달한다. 최초 요청의 미래형 구현 의사는 승인으로 인정하지 않는다.

## 작성 흐름

1. 사용자가 확인한 PRD 경로와 버전을 확인하고 파일에서 다시 읽는다. 경로 또는 승인 상태가 불명확하면 스펙을 작성하지 않고 확인한다.
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

## 출력 위치

- 사용자가 파일 작성을 요청하면 레포 루트 기준 `.codex/temp/` 디렉터리에 Markdown 파일로 저장한다. 디렉터리가 없으면 생성한다. 사용자가 다른 위치를 명시하면 그 위치를 따른다.
- 파일명을 지정하지 않으면 제품명 기반으로 `.codex/temp/<product>-technical-spec.md` 또는 `.codex/temp/<product>-spec.md`를 사용한다.
- 사용자가 채팅 답변만 원하면 파일을 만들지 않고 본문만 출력한다.

## 작성 원칙 참조

스펙 본문은 `cognitive-writing` 스킬의 원칙(인지 부하 최소화, GitHub-flavored Markdown 규칙)을 따라 작성한다.

## 기본 문서 구조

기술 스펙은 아래 구조를 기본값으로 사용한다. 주제에 맞지 않는 섹션은 줄이되, 상태 모델, API/프로토콜 계약, 에러 모델, 테스트 스펙, 섹션별 미해결 기술 결정은 가능하면 유지한다.

```md
# 기술 스펙 문서: <제품/시스템명>

**문서 버전:** v0.1 Draft
**작성일:** YYYY-MM-DD
**가칭:** <이름>
**대상 플랫폼:** <플랫폼>
**목표:** <한 줄 구현 목표>

---

# 1. 개요

# 2. 핵심 설계 원칙

# 3. 시스템 아키텍처

# 4. 런타임 상태 모델

# 5. 세션/리소스 모델

# 6. 프로토콜 스펙

# 7. 프로토콜 제약 조건

# 8. Error model

# 9. Client SDK/API 스펙

# 10. Core/Runtime API 스펙

# 11. Dispatcher/Execution flow

# 12. Boundary 스펙

# 13. Platform A 스펙

# 14. Platform B 스펙

# 15. Asset/Resource Loading 스펙

# 16. Config 스펙

# 17. Security 스펙

# 18. Logging/Diagnostics 스펙

# 19. CLI 스펙

# 20. Build/Packaging 스펙

# 21. Deferred Feature 스펙

# 22. Testing 스펙

# 23. Performance 스펙

# 24. Versioning/Compatibility 스펙

# 25. v0.1 구현 체크리스트

# 26. 기본 예제

# 27. 결정 사항 요약

# 28. 섹션별 미해결 기술 결정 인덱스

# 29. 최종 요약
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

좋은 개요 항목:

- 런타임 구조
- IPC 프로토콜
- SDK API
- native shell 계약
- FFI/JNI boundary
- lifecycle 처리
- 보안 정책
- 테스트 기준

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

### Error model

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

### Runtime/Core API 스펙

core API는 trait, struct, builder, callback, lifecycle hook 중심으로 쓴다.

반드시 정의할 항목:

- command/function signature
- context object
- result/error type
- async/threading model
- registration 방식
- dispatch 방식

### Boundary 스펙

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

### Platform 스펙

플랫폼별 섹션은 같은 구조로 반복한다.

- 기술 구성
- 주요 클래스/모듈
- bridge injection 또는 transport setup
- lifecycle mapping
- navigation/resource policy
- platform-specific failure
- acceptance criteria

### Asset/Resource Loading

개발 모드와 릴리즈 모드를 분리한다.

반드시 정의할 항목:

- mode
- dev URL/origin
- release origin
- bundled asset path
- route fallback
- external navigation 처리
- release build에서 금지되는 값

### Config 스펙

설정 파일은 예시와 validation rule을 함께 쓴다.

```toml
[app]
name = "Example"

[runtime]
version = "0.1.0"
```

validation rule은 build 전 실패해야 하는 조건을 중심으로 쓴다.

### Security 스펙

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

### Testing 스펙

테스트는 계층별로 나눈다.

- unit test
- integration test
- SDK test
- platform E2E test
- lifecycle/recreation test
- security negative test
- performance measurement

각 테스트 항목은 "무엇을 검증하는가"가 보이게 쓴다.

### Performance 스펙

목표 수치와 측정 방법을 함께 쓴다.

```md
| 항목 | 목표 |
|---|---:|
| simple round-trip | 50ms 이하 |
```

측정 이벤트 이름과 trace field를 정의한다.

### Versioning/Compatibility

runtime version, protocol version, SDK compatibility, breaking change 정책을 쓴다.

### 구현 체크리스트

v0.1에 필요한 작업만 checkbox로 압축한다.

권장 그룹:

- Core
- SDK
- Platform A
- Platform B
- CLI/Config
- Example
- Docs

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
