#!/bin/bash

VERSION=$(cat VERSION)

echo Clearing old tags
git tag -d $VERSION
git push --delete origin $VERSION
echo Pushing $VERSION
git tag $VERSION
git push origin $VERSION