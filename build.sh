#! /bin/sh -eu

version=`cat cmd/gaze/version`

rm -rf ./dist
make build-all

cd dist

for p in macos_amd macos_arm windows linux; do
  cp ../LICENSE ../README.md ${p}
  mv ${p} gaze_${p}_${version}
  zip -r gaze_${p}_${version}.zip ./gaze_${p}_${version}
done
