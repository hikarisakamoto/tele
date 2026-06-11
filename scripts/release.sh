#!/usr/bin/env bash
# Local release driver: validate, update CHANGELOG, commit, and tag.
# Usage: scripts/release.sh <version|patch|minor|major>
# Does NOT push; prints the push command to run after review.
set -euo pipefail

cd "$(git rev-parse --show-toplevel)"
# shellcheck source=scripts/changelog-lib.sh
source scripts/changelog-lib.sh

CHANGELOG="CHANGELOG.md"
die() { echo "error: $*" >&2; exit 1; }

[ $# -eq 1 ] || die "usage: scripts/release.sh <version|patch|minor|major>"

# Preconditions.
[ -z "$(git status --porcelain)" ] || die "working tree is not clean"
[ "$(git rev-parse --abbrev-ref HEAD)" = "main" ] || die "not on main branch"

latest=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
IFS=. read -r major minor patch <<<"${latest#v}"

case "$1" in
  major) version="$((major + 1)).0.0" ;;
  minor) version="${major}.$((minor + 1)).0" ;;
  patch) version="${major}.${minor}.$((patch + 1))" ;;
  *)     version="${1#v}" ;;
esac
tag="v$version"

git rev-parse "$tag" >/dev/null 2>&1 && die "tag $tag already exists"
changelog_unreleased_has_content "$CHANGELOG" \
  || die "nothing to release: [Unreleased] in $CHANGELOG is empty"

changelog_promote "$CHANGELOG" "$version" "$(date +%F)"
body=$(changelog_extract_body "$CHANGELOG" "$version")

git add "$CHANGELOG"
git commit -m "chore: release $tag"
printf '%s\n\n%s\n' "$tag" "$body" | git tag -a "$tag" -F -

echo "Tagged $tag."
echo "Review the commit and tag, then run: git push origin main --follow-tags"
