---
name: git-push-safety
description: 의도하지 않은 브랜치나 보호 브랜치에 push하는 사고를 막기 위한 안전 워크플로우입니다. Codex가 git push 실행, 준비, 검토, 설명, force push, upstream 설정, 브랜치 게시, PR 준비를 요청받았을 때 사용합니다. 사용자가 매번 명시적으로 push를 요청하거나 승인하지 않으면 push는 금지됩니다.
---

# Git Push 안전 규칙

## 출력 언어

이 스킬이 만드는 모든 사용자 대상 출력, 문서, 프롬프트, 보고서, 계획, 스펙 및 기타 산출물은 한국어로 작성합니다. 제목, 섹션, 레이블, 표, 체크리스트, 다이어그램, 템플릿도 한국어로 작성합니다. 코드, 명령어, 파일 경로, 식별자, API 이름, 모델 ID, 프로토콜 이름과 필수 고유명사는 번역하지 않습니다.

## 작성 문서 첨부

요청된 모든 문서의 작성 또는 갱신을 마친 뒤 최종 응답에 각 문서를 `[문서 이름](절대 경로)` 형식의 Markdown 링크로 반드시 첨부합니다. 링크에는 독립된 디렉터리 세그먼트로 구성된 절대 경로를 사용하고, 디렉터리 구분자를 생략하거나 경로 세그먼트를 붙여 쓰지 않습니다. Windows 경로를 문자열로 표시할 때는 `C:\\Users\\WinUser\\Documents\\폴더\\문서.md`처럼 각 역슬래시를 두 번 씁니다.

## 개요

모든 `git push`는 로컬 브랜치와 리모트 브랜치를 명시하는 방식으로 실행합니다. 사용자가 push 개념만 물어보면 명령 패턴과 주의사항만 설명하고 실제 명령은 실행하지 않습니다.

모든 push는 매번 사용자의 명시적인 허락을 받은 뒤에만 진행합니다. 사용자가 push를 명시적으로 언급하지 않은 경우 `git push` 실행은 금지됩니다.

## Push 전 필수 확인

`git push` 실행 전 아래 항목을 모두 확인합니다.

1. 사용자가 이번 턴에서 push를 명시적으로 요청했거나, 이 push를 승인했는지 확인합니다.
2. 현재 로컬 브랜치를 확인합니다.

```bash
git branch --show-current
```

3. upstream 추적 상태를 확인합니다.

```bash
git branch -vv
```

4. push 대상 리모트 브랜치가 의도한 브랜치이며 `main`, `master`, `develop`이 아닌지 확인합니다.
5. push 명령은 반드시 `로컬:리모트` refspec으로 작성합니다.
6. force push가 필요하면 최근 커밋을 확인하고 `--force-with-lease`만 사용합니다.

```bash
git log --oneline -5
```

## Push 명령 규칙

로컬 브랜치와 리모트 브랜치가 명시되도록 항상 `로컬:리모트` 형식을 사용합니다.

```bash
git push origin feat/my-feature:feat/my-feature
git push --force-with-lease origin feat/my-feature:feat/my-feature
```

브랜치만 적는 push는 사용하지 않습니다.

```bash
git push origin feat/my-feature
git push --force-with-lease origin feat/my-feature
```

보호 브랜치에는 직접 push하지 않습니다.

- `main`
- `master`
- `develop`

보호 브랜치 반영은 PR/MR을 통해서만 진행합니다.

## Upstream 처리

push 전에는 항상 upstream 추적 상태를 확인합니다.

```bash
git branch -vv
```

현재 브랜치가 `origin/main`, `origin/master`, `origin/develop` 또는 잘못된 리모트 브랜치를 추적하고 있으면 그 상태로 push하지 않습니다.

의도한 리모트 브랜치가 이미 있고 안전한 경우에만 upstream을 재설정합니다.

```bash
git branch --set-upstream-to=origin/<current-branch>
```

의도한 리모트 브랜치가 아직 없으면 명시적인 refspec으로 게시합니다.

```bash
git push origin <current-branch>:<current-branch>
```

## 강제 Push

- `--force`는 사용하지 않고 `--force-with-lease`만 사용합니다.
- force push 대상이 보호 브랜치가 아닌지 확인합니다.
- force push 전 최근 커밋을 확인합니다.

```bash
git log --oneline -5
```

## 체크리스트

push할 때마다 아래를 확인합니다.

1. 현재 브랜치가 의도한 브랜치인가? `git branch --show-current`
2. upstream이 올바른 리모트 브랜치를 가리키는가? `git branch -vv`
3. push 명령이 명시적인 `로컬:리모트` refspec을 사용하는가?
4. 대상이 보호 브랜치가 아닌가?
5. force push라면 `--force-with-lease`를 사용하는가?

## 응답 규칙

사용자가 "git push는 어떻게 하나요?"처럼 개념 질문을 하면 명령 패턴과 안전 주의사항만 설명합니다. `git push`는 실행하지 않습니다.

사용자가 Codex에게 push를 요청하면 실행 전에 확인한 로컬 브랜치, upstream 상태, 대상 refspec, force push 여부를 보여줍니다.
