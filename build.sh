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
	echo '| cover'
	echo '|		View the code coverage report from running tests.'
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

_echo_green() {
  echo -e '\e[0;32m'$1'\e[0m'
}

# Copies the generated code from the model repo.
_cp_gen() {
	rm -rf ./gen
	cp -r ../ama-model/gen/ama/api/gen ./gen
}

_set_env_vars() {
	export FIREBASE_AUTH_EMULATOR_HOST=127.0.0.1:9099
	export FIRESTORE_EMULATOR_HOST=127.0.0.1:8080 && \
	export PROJECT_ID=ama-dev && \
	export GO_LOG=debug
}

run() {
	db &
	DATABASE_PROCESS=$!
	sleep 1
	ps -p $DATABASE_PROCESS > /dev/null
	if [[ "$?" != "0" ]]
	then
		echo 'Failed to start database emulator.'
		exit 1
	fi
	_kill_db() {
		_echo_green 'Killing database...'
		kill -INT $DATABASE_PROCESS
	}
	trap _kill_db SIGINT
	_cp_gen
	_set_env_vars
	go run . -port 8088 || _kill_db
	_echo_green 'Shut down successfully.'
}

db() {
	firebase emulators:start --project ama-dev
}

COVERAGE_FILE=.cover.out

test() {
	go test -v -coverprofile="$COVERAGE_FILE.tmp" ./... | grep -v "gen/"
}

cover() {
	test
	cat "$COVERAGE_FILE.tmp" | grep -v "gen/" > $COVERAGE_FILE
	go tool cover -html=$COVERAGE_FILE
	rm $COVERAGE_FILE "$COVERAGE_FILE.tmp"
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
	_echo_green 'Reformat complete.'
}

_get_image_repo() {
	echo -n "${GCP_REGION}-docker.pkg.dev"
}

_get_image_uri_without_tag() {
	echo -n "$(_get_image_repo)/${GCP_PROJECT_ID}/ama-api/api-server"
}

_get_image_uri() {
	echo -n "$(_get_image_uri_without_tag):$(git rev-parse --short HEAD)"
}

_get_latest_uri() {
	echo -n "$(_get_image_uri_without_tag):latest"
}

_docker_login() {
	local IMAGE_REPO="$(_get_image_repo)"
	gcloud auth configure-docker $IMAGE_REPO
	gcloud auth print-access-token \
		--impersonate-service-account $GCP_PROJECT_ID | docker login \
		-u oauth2accesstoken \
		--password-stdin "https://${IMAGE_REPO}"
}

build() {
	local IMAGE_URI="$(_get_image_uri)"
	echo "Building Docker Image ${IMAGE_URI}"
	_cp_gen
	docker build -t $IMAGE_URI .
	_echo_green 'Done.'
}

deploy() {
	_docker_login
	local IMAGE_URI="$(_get_image_uri)"
	echo "Pushing Docker Image ${IMAGE_URI}"
	docker push $IMAGE_URI
	local LATEST_URI="$(_get_latest_uri)"
	docker tag $IMAGE_URI $LATEST_URI
	docker push $LATEST_URI
	_echo_green 'Done.'
}

integ() {
	echo 'Starting integration tests...'
	_set_env_vars
	export USER_EMAIL="$(uuidgen | cut -d'-' -f2)@t.com"
	go test -count=1 -tags="integration" ./integration_test -suite=all
	echo "User Email: '${USER_EMAIL}'"
}

integsetup() {
	echo 'Starting integration tests...'
	_set_env_vars
	export USER_EMAIL="$(uuidgen | cut -d'-' -f2)@t.com"
	go test -count=1 -tags="integration" ./integration_test -suite=setup
	echo "User Email: '${USER_EMAIL}'"
}

integteardown() {
	echo 'Starting integration tests...'
	_set_env_vars
	echo "DON'T FORGET TO SET THE EMAIL FOR TEARDOWN"
	echo "User Email: '${USER_EMAIL}'"
	go test -count=1 -tags="integration" ./integration_test -suite=teardown
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
