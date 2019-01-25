#!/bin/sh

REPO_DIR=$PWD
export GO111MODULE=on

verify () {
	if [[ "$1" ]]; then
		cd ${REPO_DIR}/$1 || exit 1
	fi
	VERSION=${1:-v2}
	go fmt . || { echo "pre commit hook ${VERSION}: FMT FAILURE"; exit 1; }
	go vet . || { echo "pre commit hook ${VERSION}: VET FAILURE"; exit 1; }
	if [[ "$1" ]]; then
		go mod tidy -v || { echo "pre commit hook ${VERSION}: MOD TIDY FAILURE"; exit 1; }
	fi
	echo "pre commit hook ${VERSION}: OK"
}

set -e

verify
verify v3
verify v4