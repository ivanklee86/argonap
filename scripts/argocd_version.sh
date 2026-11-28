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
#     bump rewrites every replace directive in go.mod and check flags any extras, so
#     project-specific replaces don't belong in go.mod (use a go.work file for local development).
#   - .github/workflows/ci.yaml: ARGOCD_VERSION (the argocd CLI used by tests)
#   - .devcontainer/Dockerfile: quay.io/argoproj/argocd image (the argocd CLI in the devcontainer)
#
# Portable to BSD/macOS userland: POSIX sed classes, no sed -i.
set -euo pipefail

ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
GOMOD="${ROOT}/go.mod"
CI="${ROOT}/.github/workflows/ci.yaml"
DOCKERFILE="${ROOT}/.devcontainer/Dockerfile"
MAJOR=3
MODULE="github.com/argoproj/argo-cd/v${MAJOR}"

WORK="$(mktemp -d)"
trap 'rm -rf "${WORK}"' EXIT

die() { echo "error: $*" >&2; exit 1; }

# Edit a file in place with sed, without sed -i (GNU and BSD disagree on its syntax).
sed_inplace() {
  local expr="$1" file="$2"
  sed -E "${expr}" "${file}" >"${WORK}/sed.out"
  cat "${WORK}/sed.out" >"${file}"
}

pinned_version() {
  go mod edit -json "${GOMOD}" | jq -r --arg m "${MODULE}" '.Require[] | select(.Path == $m) | .Version'
}

# argo-cd's go.mod at <version>, from the Go module proxy: it serves just the .mod file, is cached,
# and stays up when raw.githubusercontent.com rate-limits. Honors the first http(s) GOPROXY entry.
fetch_upstream() {
  local version="$1" proxy
  proxy="$(go env GOPROXY | tr ',|' '\n\n' | grep -E '^https?://' | head -n 1 || true)"
  proxy="${proxy:-https://proxy.golang.org}"
  curl -sSfL --retry 3 --retry-all-errors -o "${WORK}/upstream.go.mod" \
    "${proxy%/}/${MODULE}/@v/${version}.mod" \
    || die "could not fetch ${MODULE}@${version}'s go.mod from ${proxy}"
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
  [[ -n "${version}" ]] || die "${MODULE} not found in go.mod (argo-cd major bump? see 'bump' in this script)"
  fetch_upstream "${version}"

  if ! diff -u <(replaces "${WORK}/upstream.go.mod") <(replaces "${GOMOD}") >"${WORK}/replace.diff"; then
    echo "go.mod's replace directives have drifted from argo-cd ${version}'s (- argo-cd, + go.mod):"
    tail -n +3 "${WORK}/replace.diff"
    status=1
  fi

  local ci_version docker_version
  ci_version="$(sed -nE 's/^[[:space:]]*ARGOCD_VERSION: ([^[:space:]]+).*/\1/p' "${CI}")"
  docker_version="$(sed -nE 's#.*quay\.io/argoproj/argocd:([^[:space:]]+).*#\1#p' "${DOCKERFILE}")"
  if [[ "${ci_version}" != "${version}" ]]; then
    echo "ci.yaml ARGOCD_VERSION is '${ci_version}', go.mod pins argo-cd ${version}"
    status=1
  fi
  if [[ "${docker_version}" != "${version}" ]]; then
    echo ".devcontainer/Dockerfile argocd image is '${docker_version}', go.mod pins argo-cd ${version}"
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

# Rewrite go.mod's replace directives as one block holding upstream's, keeping upstream's comments.
sync_replace_block() {
  local version="$1"
  awk -v version="${version}" '
    FNR == NR {
      # First file: upstream go.mod. Collect its replace block plus any single-line replaces,
      # minus local-path replaces.
      if ($0 ~ /^replace \($/) { inblock = 1; next }
      if (inblock && $0 ~ /^\)$/) { inblock = 0; next }
      if (inblock && $0 !~ /=> \.\.?\//) { block = block $0 "\n" }
      if (!inblock && $0 ~ /^replace [^(]/ && $0 !~ /=> \.\.?\//) {
        line = $0; sub(/^replace /, "", line); single = single "\t" line "\n"
      }
      next
    }
    # Second file: our go.mod. Drop our header comment and every replace directive, then append
    # the fresh header and block.
    /^\/\/ Mirrored from argo-cd / { skipcomment = 1 }
    skipcomment && /^\/\// { next }
    { skipcomment = 0 }
    /^replace \($/ { skipping = 1; next }
    skipping && /^\)$/ { skipping = 0; next }
    skipping { next }
    /^replace [^(]/ { next }
    { print }
    END {
      print "// Mirrored from argo-cd " version "'"'"'s own go.mod replace directives by scripts/argocd_version.sh."
      print "// Replace directives aren'"'"'t inherited by consumers, and k8s.io/kubernetes requires its staging"
      print "// modules at unresolvable v0.0.0 placeholders without them. Do not edit by hand or add"
      print "// project-specific replaces here: run '"'"'task argocd:bump VERSION=<version>'"'"', which rewrites them."
      printf "replace (\n%s%s)\n", block, single
    }
  ' "${WORK}/upstream.go.mod" "${GOMOD}" >"${WORK}/go.mod"
  # Collapse the blank lines left where the old block was removed.
  cat -s "${WORK}/go.mod" >"${GOMOD}"
  go mod edit -fmt "${GOMOD}"
}

bump() {
  local version="${1:-}"
  [[ "${version}" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]] || die "usage: $0 bump vX.Y.Z"
  # Checked before anything is edited, so a major bump fails without touching go.mod.
  if [[ "${version}" != "v${MAJOR}."* ]]; then
    die "${version} is a major argo-cd bump: its module path is no longer ${MODULE}. Update MAJOR in \
this script and the ${MODULE} imports (and renovate.json's argocd group) first, then re-run bump."
  fi
  fetch_upstream "${version}"

  cd "${ROOT}"
  sync_replace_block "${version}"
  go get "${MODULE}@${version}"
  go mod tidy
  sed_inplace "s/^([[:space:]]*ARGOCD_VERSION: )[^[:space:]]+/\1${version}/" "${CI}"
  sed_inplace "s#(quay\.io/argoproj/argocd:)[^[:space:]]+#\1${version}#" "${DOCKERFILE}"
  check
}

case "${1:-}" in
  check) check ;;
  bump) bump "${2:-}" ;;
  *) die "usage: $0 {check|bump vX.Y.Z}" ;;
esac
