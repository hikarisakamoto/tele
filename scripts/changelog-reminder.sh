#!/usr/bin/env bash
# Non-blocking pre-push reminder: warn if a version tag on HEAD has no matching
# CHANGELOG.md section. Always exits 0 so it never blocks a push.
set -uo pipefail

tag=$(git tag --points-at HEAD --list 'v*' | head -n1)
[ -n "$tag" ] || exit 0

version=${tag#v}
if ! grep -qE "^## \[${version}\]" CHANGELOG.md 2>/dev/null; then
  echo "reminder: pushing $tag but CHANGELOG.md has no '## [$version]' section." >&2
  echo "reminder: scripts/release.sh keeps the changelog and tag in sync." >&2
fi
exit 0
