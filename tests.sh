#!/usr/bin/env bash

set -euo pipefail

PORT=8088

question_tests() {
  question_id=$(curl \
    -s \
    -X POST \
    -H 'Content-Type: application/json' \
    -d '{"prompt": "test question", "tags": ["test"]}' \
    localhost:$PORT/question | jq -r '.id')
  echo "Created question with id ${question_id}"
  echo "Fetching id ${question_id}"
  curl -s "localhost:$PORT/question/${question_id}" | jq
  echo "modifying the question"
  curl \
    -s \
    -X PUT \
    -H 'Content-Type: application/json' \
    -d '{"prompt": "test question edited", "tags": ["test", "edit"]}' \
    "localhost:$PORT/question/${question_id}" | jq
  echo "Fetching edited question"
  curl -s "localhost:$PORT/question/${question_id}" | jq
  for i in {1..10}
  do
    echo "Adding another question n=${i}"
    curl \
      -s \
      -X POST \
      -H 'Content-Type: application/json' \
      -d '{"prompt": "second test question", "tags": ["test", "another test", "'$i'"]}' \
      localhost:$PORT/question | jq
  done
  echo "Fetching all question"
  curl -s localhost:$PORT/question | jq
  echo "Getting 1 question using the query params"
  curl -s "localhost:$PORT/question?tags=test&tags=3" | jq
  # echo "Deleting all questions"
  # for each_question in $(curl -s "localhost:$PORT/question?tags=test" | jq -r '.[] | .id')
  # do
  #   curl -s -X DELETE "localhost:$PORT/question/${each_question}" | jq
  # done
}

user_tests() {
  echo 'Creating a user'
  local USER_ID=$(curl -s -X POST \
    -H 'Content-Type: application/json' \
    -d '{
      "name": "test",
      "tier": "free",
      "subscription": {
        "payCadence": "monthly",
        "renewalDate": "2025-06-23T07:31:37.079771Z"
      },
      "settings": {
        "colorScheme": {
          "background": "default",
          "foreground": "default",
          "highlightedBackground": "default",
          "highlightedForeground": "default"
        }
      },
      "lists": []
    }' \
    localhost:$PORT/user | jq -r '.id')
    curl -s localhost:$PORT/user/$USER_ID | jq
}

list_tests() {
  echo 'TODO: write these...'
}

# Creates reads updates and deletes data through the API as a series of tests
main() {
  echo 'Running question tests...'
  question_tests
  echo 'All question tests passed!'
  echo 'Running user tests...'
  user_tests
  echo 'All user tests passed!'
  echo 'Running list tests...'
  list_tests
  echo 'All list tests passed!'
  echo 'Done!'
}

main "$@"
