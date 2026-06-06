#!/bin/bash

TAG=$1

# Tag with path-relative names that match module import path
git tag $TAG
git tag drivers/sqlite/$TAG
git tag drivers/postgres/$TAG
git tag drivers/mariadb/$TAG

# Push the correct tags
git push origin $TAG
git push origin drivers/sqlite/$TAG
git push origin drivers/postgres/$TAG
git push origin drivers/mariadb/$TAG
