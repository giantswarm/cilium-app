#!/usr/bin/env bash
# Make the vendored chart pass the generic pre-commit hooks: no trailing
# whitespace, a final newline, and executable scripts with a shebang. The
# vendored Makefile.quiet keeps its trailing whitespace until the build ran.
set -o errexit
set -o nounset
set -o pipefail

for d in "$@"; do
  find "$d" -type f ! -name '*.tgz' ! -path 'vendor/*/Makefile.quiet' -print0 | while IFS= read -r -d '' f; do
    grep -Iq . "$f" || continue
    sed -i -e 's/[[:space:]]*$//' -e :a -e '/^\n*$/{$d;N;ba' -e '}' "$f"
    if [ -n "$(tail -c1 "$f")" ]; then echo >> "$f"; fi
    if [ "$(head -c2 "$f")" = '#!' ]; then chmod +x "$f"; fi
  done
done
