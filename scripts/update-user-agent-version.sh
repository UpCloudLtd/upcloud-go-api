#!/bin/sh -x

version_re="[0-9]+\.[0-9]+\.[0-9]+"

changelog=$(grep -E -m 1 "##.*$version_re.*" CHANGELOG.md | grep -Eo "$version_re");
user_agent=$(grep -E -m 1 "Version.*$version_re.*" upcloud/client/client.go | grep -Eo "$version_re")

if [ "$changelog" = "$user_agent" ]; then
    exit 0;
fi;

latest=$(echo "$changelog\n$user_agent" | sort -V | tail -n1)

sed -Ei "s/(##.*)$changelog(.*)/\1$latest\2/" CHANGELOG.md;
sed -Ei "s/(.*Version.*\")$user_agent(\".*)/\1$latest\2/" upcloud/client/client.go;
