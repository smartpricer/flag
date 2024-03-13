#!/bin/bash

echo Pushing $(cat VERSION)
git tag $(cat VERSION)
git push origin $(cat VERSION)