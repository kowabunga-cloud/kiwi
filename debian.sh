#!/bin/bash

DIST="${DIST:=noble}"

export CURRENT_VERSION=$(git tag --sort=-committerdate | head -1)
export PREVIOUS_VERSION=$(git tag --sort=-committerdate | head -2 | awk '{split($0, tags, "\n")} END {print tags[1]}')
export CHANGES=$(git log --pretty="- %s" $CURRENT_VERSION...$PREVIOUS_VERSION)
export VERSION="$(echo $RELEASE | cut -d 'v' -f 2)~$DIST"

cat > debian/changelog <<EOF
kowabunga-kiwi-agent (${VERSION}) ${DIST}; urgency=medium

${CHANGES}

 -- The Kowabunga Project <maintainers@kowabunga.cloud>  $(date -R)
EOF

sed -i 's%^-%  \*%g' debian/changelog
DEB_BUILD_OPTIONS=noautodbgsym fakeroot debian/rules binary
