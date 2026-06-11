#!/usr/bin/env bash
# Shared CHANGELOG.md (Keep a Changelog) parsing/rewriting helpers.
# Pure functions over a file path; no git operations live here.
# The optional human title is an em-dash (" — ") suffix on a section heading.

# changelog_unreleased_has_content <file>
# Returns 0 if the [Unreleased] section has at least one list item, else 1.
changelog_unreleased_has_content() {
  awk '
    /^## \[Unreleased\]/ { in_sec=1; next }
    in_sec && /^## \[/   { exit }
    in_sec && /^[-*] /   { found=1; exit }
    END { exit (found ? 0 : 1) }
  ' "$1"
}

# changelog_extract_title <file> <key>
# <key> is a version like "1.2.0" or the literal "Unreleased".
# Prints the human title (text after " — " on the heading), or nothing.
changelog_extract_title() {
  local file=$1 key=$2 heading title
  heading=$(grep -m1 -E "^## \[${key}\]" "$file") || return 0
  title=${heading#*" — "}
  [ "$title" = "$heading" ] && return 0   # no separator present
  printf '%s\n' "$title"
}

# changelog_extract_body <file> <key>
# Prints the section body (lines between the heading and the next "## [..."),
# with leading and trailing blank lines trimmed.
changelog_extract_body() {
  awk -v key="$2" '
    index($0, "## [" key "]") == 1 { in_sec=1; next }
    in_sec && /^## \[/            { in_sec=0 }
    in_sec                         { buf[n++]=$0 }
    END {
      s=0;   while (s<n   && buf[s] ~ /^[[:space:]]*$/) s++
      e=n-1; while (e>=s  && buf[e] ~ /^[[:space:]]*$/) e--
      for (i=s; i<=e; i++) print buf[i]
    }
  ' "$1"
}

# changelog_promote <file> <version> <date>
# Renames the [Unreleased] section to a versioned section (carrying any title)
# and inserts a fresh empty [Unreleased] above it. Rewrites the file in place.
changelog_promote() {
  local file=$1 version=$2 date=$3 title heading
  title=$(changelog_extract_title "$file" "Unreleased")
  heading="## [$version] - $date"
  [ -n "$title" ] && heading="$heading — $title"
  awk -v nh="$heading" '
    !done && /^## \[Unreleased\]/ {
      print "## [Unreleased]"
      print ""
      print nh
      done=1
      next
    }
    { print }
  ' "$file" > "$file.tmp" && mv "$file.tmp" "$file"
}
