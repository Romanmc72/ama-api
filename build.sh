#!/usr/bin/env bash

set -euo pipefail

help() {
	echo '| ./build.sh ($1)'
	echo '+====================================================================='
	echo '| Description'
	echo '| ==========='
	echo '| This is the script to use for building the AMA API. Run the script'
	echo '| with the desired command to perform the desired build step.'
	echo '|'
	echo '| Params'
	echo '| ======'
	echo '| $1 : command : default = help'
	echo '| One of the commands listed below.'
	echo '|'
	echo '| Commands'
	echo '| --------'
	echo '| run'
	echo '|		Launch the app locally (not for production).'
	echo '|'
	echo '| db'
	echo '|		Launch the database emulator locally (not for production).'
	echo '|'
	echo '| test'
	echo '|		Run the unit tests.'
	echo '|'
	echo '| lint'
	echo '|		Run the lint check and fail if it there are issues.'
	echo '|'
	echo '| format'
	echo '|		Run the linter to update the formatting in place.'
	echo '|'
	echo '| build'
	echo '|		Create the main executable program file.'
	echo '|'
	echo '| deploy'
	echo '|		Package and ship the executable to its destination.'
	echo '|'
	echo '| release'
	echo '|		Run several commands all together to lint, test, build, and deploy'
	echo ''
}

run() {
	export FIRESTORE_EMULATOR_HOST=127.0.0.1:8080 && \
	export PROJECT_ID=ama-dev && \
	export GO_LOG=debug && \
	go run .
}

db() {
	firebase emulators:start --project ama-dev
}

test() {
	go test -v ./...
}

# Internal flags used to show linting errors
_go_lint() {
	gofmt -l -d ./
}

lint() {
	if [[ $(_go_lint) ]]
	then
		printf '\U0001F480 Found some errors!'
		_go_lint
		return 1;
	else
		printf '\U0001F389 Looks Good!';
	fi
}

format() {
	echo 'Running formatter...'
	echo 'Changes:'
	_go_lint
	gofmt -w -s ./
	echo 'Reformat complete.'
}

build() {
	go build .
}

deploy() {
	echo 'HAVE NOT IMPLEMENTED THIS'
}

ci() {
	lint
	test
	build
	deploy
}

main() {
	if [ $# -eq 0 ]
	then
		help
	else
		"$@"
	fi
}

main $@
