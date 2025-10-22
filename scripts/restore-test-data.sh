#! /bin/bash

PROJECT_DIR=${HOME}/Workspaces/Projects/schluckauf

rm ${PROJECT_DIR}/data/duplicates.db

cp -r ${PROJECT_DIR}/test-photos-original/* ${PROJECT_DIR}/test-photos/.
