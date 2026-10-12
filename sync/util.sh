# Helpers sourced by the sync scripts. Not meant to be executed directly.

# GNU sed takes an optional suffix for -i, BSD sed requires one.
if sed --version >/dev/null 2>&1 ; then
    sed_inplace() { sed -i "$@"; }
else
    sed_inplace() { sed -i "" "$@"; }
fi
