#!/bin/sh
set -eu

repository="j-token/j-token-codex-workflow-kit"
version="latest"
install_dir="${CODEX_WORKFLOW_INSTALL_DIR:-${HOME}/.local/bin}"

while [ "$#" -gt 0 ]; do
  case "$1" in
    --version)
      [ "$#" -ge 2 ] || { echo "--version 값이 필요합니다" >&2; exit 2; }
      version="$2"
      shift 2
      ;;
    --install-dir)
      [ "$#" -ge 2 ] || { echo "--install-dir 값이 필요합니다" >&2; exit 2; }
      install_dir="$2"
      shift 2
      ;;
    *)
      echo "알 수 없는 옵션: $1" >&2
      exit 2
      ;;
  esac
done

case "$(uname -s)" in
  Darwin) os="darwin" ;;
  Linux) os="linux" ;;
  *) echo "지원하지 않는 운영체제입니다: $(uname -s)" >&2; exit 1 ;;
esac

case "$(uname -m)" in
  x86_64|amd64) arch="x64" ;;
  arm64|aarch64) arch="arm64" ;;
  *) echo "지원하지 않는 아키텍처입니다: $(uname -m)" >&2; exit 1 ;;
esac

asset="codex-workflow-${os}-${arch}"
if [ "$version" = "latest" ]; then
  base_url="https://github.com/${repository}/releases/latest/download"
else
  case "$version" in v*) ;; *) version="v${version}" ;; esac
  base_url="https://github.com/${repository}/releases/download/${version}"
fi

temp_dir="$(mktemp -d)"
staged_path=""
trap 'rm -rf "$temp_dir"; [ -z "$staged_path" ] || rm -f "$staged_path"' EXIT INT TERM

curl -fsSL "${base_url}/${asset}" -o "${temp_dir}/${asset}"
curl -fsSL "${base_url}/${asset}.sha256" -o "${temp_dir}/${asset}.sha256"

expected="$(awk '{print $1}' "${temp_dir}/${asset}.sha256")"
if command -v sha256sum >/dev/null 2>&1; then
  actual="$(sha256sum "${temp_dir}/${asset}" | awk '{print $1}')"
else
  actual="$(shasum -a 256 "${temp_dir}/${asset}" | awk '{print $1}')"
fi

[ "$expected" = "$actual" ] || { echo "SHA-256 검증에 실패했습니다" >&2; exit 1; }

mkdir -p "$install_dir"
staged_path="${install_dir}/.codex-workflow.new.$$"
cp "${temp_dir}/${asset}" "$staged_path"
chmod 0755 "$staged_path"
mv -f "$staged_path" "${install_dir}/codex-workflow"
staged_path=""
echo "설치 완료: ${install_dir}/codex-workflow"

case ":${PATH}:" in
  *":${install_dir}:"*) ;;
  *) echo "PATH에 다음 경로를 추가하세요: ${install_dir}" ;;
esac
