#!/usr/bin/env bash
export PROJECT=${1?Error: no project name given}

rm -rf .git
find . -type f -exec sed -i "s/fixmarketsimulator/$PROJECT/g" {} +
mv ./src/main/java/com/ettech/fixmarketsimulator/ ./src/main/java/com/ettech/$PROJECT
mv ./src/test/java/com/ettech/fixmarketsimulator/ ./src/test/java/com/ettech/$PROJECT
mv ../fixmarketsimulator ../$PROJECT
