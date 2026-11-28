#!/usr/bin/env bash
# Keeps every Argo CD version touchpoint in lockstep with the argo-cd Go module.
#
#   argocd_version.sh check          Fail if any touchpoint drifts from the pinned argo-cd version.
#   argocd_version.sh bump <vX.Y.Z>  Move argo-cd and everything that mirrors it to <vX.Y.Z>.
#
# Touchpoints:
#   - go.mod: github.com/argoproj/argo-cd/v3 (require)
#   - go.mod: the replace block, mirrored from argo-cd's own go.mod at that version. Replace
#     directives aren't inherited by consumers, so we copy them; Renovate leaves them alone.
#   - .github/workflows/ci.yaml: ARGOCD_VERSION (the argocd CLI used by tests)
#   - .devcontainer/Dockerfile: quay.io/argoproj/argocd image (the argocd CLI in the devcontainer)
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
GOMOD="${ROOT}/go.mod"
CI="${ROOT}/.github/workflows/ci.yaml"
DOCKERFILE="${ROOT}/.devcontainer/Dockerfile"
MODULE="github.com/argoproj/argo-cd/v3"

WORK="$(mktemp -d)"
trap 'rm -rf "${WORK}"' EXIT

die() { echo "error: $*" >&2; exit 1; }

pinned_version() {
  go mod edit -json "${GOMOD}" | jq -r --arg m "${MODULE}" '.Require[] | select(.Path == $m) | .Version'
}

fetch_upstream() {
  local version="$1"
  curl -sSfL -o "${WORK}/upstream.go.mod" \
    "https://raw.githubusercontent.com/argoproj/argo-cd/${version}/go.mod" \
    || die "could not fetch argo-cd ${version}'s go.mod"
}

# Replace directives as sorted "old[@v] => new[@v]" lines. Local-path replaces (argo-cd's in-repo
# modules) are dropped: they only resolve inside argo-cd's own checkout.
replaces() {
  go mod edit -json "$1" | jq -r '
    (.Replace // [])[]
    | select(.New.Path | startswith(".") | not)
    | "\(.Old.Path)\(if .Old.Version then "@" + .Old.Version else "" end) => \(.New.Path)\(if .New.Version then "@" + .New.Version else "" end)"
  ' | sort
}

check() {
  local version status=0
  version="$(pinned_version)"
  [[ -n "${version}" ]] || die "${MODULE} not found in go.mod"
  fetch_upstream "${version}"

  if ! diff -u --label "argo-cd ${version} replace block" --label "go.mod replace block" \
      <(replaces "${WORK}/upstream.go.mod") <(replaces "${GOMOD}") >"${WORK}/replace.diff"; then
    echo "go.mod's replace block has drifted from argo-cd ${version}'s:"
    cat "${WORK}/replace.diff"
    status=1
  fi

  local ci_version docker_version
  ci_version="$(sed -nE 's/^\s*ARGOCD_VERSION: (\S+).*/\1/p' "${CI}")"
  docker_version="$(sed -nE 's#.*quay\.io/argoproj/argocd:(\S+).*#\1#p' "${DOCKERFILE}")"
  if [[ "${ci_version}" != "${version}" ]]; then
    echo "ci.yaml ARGOCD_VERSION is ${ci_version}, go.mod pins argo-cd ${version}"
    status=1
  fi
  if [[ "${docker_version}" != "${version}" ]]; then
    echo ".devcontainer/Dockerfile argocd image is ${docker_version}, go.mod pins argo-cd ${version}"
    status=1
  fi

  if [[ ${status} -ne 0 ]]; then
    echo
    echo "Run 'task argocd:bump VERSION=${version}' to resync, or bump to a new release the same way."
  else
    echo "All Argo CD touchpoints match argo-cd ${version}."
  fi
  return ${status}
}

# Rewrite go.mod's replace block from the upstream one, keeping upstream's comments.
sync_replace_block() {
  local version="$1"
  awk -v version="${version}" '
    FNR == NR {
      # First file: upstream go.mod. Capture its replace block, minus local-path replaces.
      if ($0 ~ /^replace \($/) { inblock = 1; next }
      if (inblock && $0 ~ /^\)$/) { inblock = 0; next }
      if (inblock && $0 !~ /=> \.\.?\//) { block = block $0 "\n" }
      next
    }
    # Second file: our go.mod. Drop our header comment and replace block, then append fresh ones.
    /^\/\/ Mirrored from argo-cd / { skipcomment = 1 }
    skipcomment && /^\/\// { next }
    { skipcomment = 0 }
    /^replace \($/ { skipping = 1; next }
    skipping && /^\)$/ { skipping = 0; next }
    skipping { next }
    { print }
    END {
      print "// Mirrored from argo-cd " version "'"'"'s own go.mod replace block by scripts/argocd_version.sh."
      print "// Replace directives aren'"'"'t inherited by consumers, and k8s.io/kubernetes requires its staging"
      print "// modules at unresolvable v0.0.0 placeholders without them. Do not edit by hand: run"
      print "// '"'"'task argocd:bump VERSION=<version>'"'"' (Renovate is configured to leave these lines alone)."
      printf "replace (\n%s)\n", block
    }
  ' "${WORK}/upstream.go.mod" "${GOMOD}" >"${WORK}/go.mod"
  # Collapse the blank lines left where the old block was removed.
  cat -s "${WORK}/go.mod" >"${GOMOD}"
  go mod edit -fmt "${GOMOD}"
}

bump() {
  local version="${1:-}"
  [[ "${version}" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]] || die "usage: $0 bump vX.Y.Z"
  fetch_upstream "${version}"

  cd "${ROOT}"
  sync_replace_block "${version}"
  go get "${MODULE}@${version}"
  go mod tidy
  sed -i -E "s/^(\s*ARGOCD_VERSION: )\S+/\1${version}/" "${CI}"
  sed -i -E "s#(quay\.io/argoproj/argocd:)\S+#\1${version}#" "${DOCKERFILE}"
  check
}

case "${1:-}" in
  check) check ;;
  bump) bump "${2:-}" ;;
  *) die "usage: $0 {check|bump vX.Y.Z}" ;;
esac
