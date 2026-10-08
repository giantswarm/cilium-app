#!/usr/bin/env bash
sed_inplace() {
    if sed --version 2>/dev/null | grep -q "GNU sed" ; then
        sed -i "$@"
    else
        sed -i "" "$@"
    fi
}
