#!/usr/bin/env bash

set -euo pipefail

# Creates reads updates and deletes data through the API as a series of tests
main() {
  question_id=$(curl \
    -s \
    -X POST \
    -H 'Content-Type: application/json' \
    -d '{"prompt": "test question", "tags": ["test"]}' \
    localhost:8081/questions | jq -r '.id')
  echo "Created question with id ${question_id}"
  echo "Fetching id ${question_id}"
  curl -s "localhost:8081/questions/${question_id}" | jq
  echo "modifying the question"
  curl \
    -s \
    -X PUT \
    -H 'Content-Type: application/json' \
    -d '{"prompt": "test question edited", "tags": ["test", "edit"]}' \
    "localhost:8081/questions/${question_id}" | jq
  echo "Fetching edited question"
  curl -s "localhost:8081/questions/${question_id}" | jq
  for i in {1..10}
  do
    echo "Adding another question n=${i}"
    curl \
      -s \
      -X POST \
      -H 'Content-Type: application/json' \
      -d '{"prompt": "second test question", "tags": ["test", "another test", "'$i'"]}' \
      localhost:8081/questions | jq
  done
  echo "Fetching all questions"
  curl -s localhost:8081/questions | jq
  echo "Getting 1 question using the query params"
  curl -s 'localhost:8081/questions?tags=test%7C3' | jq
  echo "Deleting all questions"
  for each_question in $(curl -s 'localhost:8081/questions?tags=test' | jq -r '.[] | .id')
  do
    curl -s -X DELETE "localhost:8081/questions/${each_question}" | jq
  done
  echo "Done!"
}

main "$@"
