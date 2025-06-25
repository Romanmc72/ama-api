#!/usr/bin/env bash

set -euo pipefail

PORT=8088

question_tests() {
  question_id=$(curl \
    --silent \
    -X POST \
    -H 'Content-Type: application/json' \
    -d '{"prompt": "test question", "tags": ["test"]}' \
    localhost:$PORT/question | jq -r '.id')
  echo "Created question with id ${question_id}"
  echo "Fetching id ${question_id}"
  curl --silent "localhost:$PORT/question/${question_id}" | jq
  echo "modifying the question"
  curl \
    --silent \
    -X PUT \
    -H 'Content-Type: application/json' \
    -d '{"prompt": "test question edited", "tags": ["test", "edit"]}' \
    "localhost:$PORT/question/${question_id}" | jq
  echo "Fetching edited question"
  curl --silent "localhost:$PORT/question/${question_id}" | jq
  for i in {1..10}
  do
    echo "Adding another question n=${i}"
    curl \
      --silent \
      -X POST \
      -H 'Content-Type: application/json' \
      -d '{"prompt": "second test question", "tags": ["test", "another test", "'$i'"]}' \
      localhost:$PORT/question | jq
  done
  echo "Fetching all question"
  curl --silent localhost:$PORT/question | jq
  echo "Getting 1 question using the query params"
  curl --silent "localhost:$PORT/question?tags=test&tags=3" | jq
}

delete_all_questions() {
  echo "Deleting all questions"
  for each_question in $(curl --silent "localhost:$PORT/question?tags=test&limit=100" | jq -r '.[] | .id')
  do
    curl --silent -X DELETE "localhost:$PORT/question/${each_question}" | jq
  done
}

user_tests() {
  echo 'Creating a user'
  USER_ID=$(curl --silent -X POST \
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
  list_id=$(curl --silent localhost:$PORT/user/$USER_ID | jq -r '.lists[0].id')
  list_tests() {
    echo 'Testing user lists'
    for i in {0..5}
    do
      question_id=$(get_random_question_id)
      echo "Adding question $question_id to list $list_id"
      curl \
        --silent \
        -X POST \
        "localhost:$PORT/user/$USER_ID/list/$list_id/question/$question_id" \
        | jq
    done
    echo 'Fetching questions from list and removing them'
    for each_question in $(curl --silent "localhost:$PORT/user/$USER_ID/list/$list_id?limit=10" | jq -r '.questions[] | .id')
    do
      echo "Removing question $each_question from list $list_id"
      curl --silent -X DELETE "localhost:$PORT/user/$USER_ID/list/$list_id/question/$each_question" | jq
    done
  }
  echo 'Testing user lists...'
  list_tests
  echo 'List tests passed!'
  echo 'Deleting user'
  curl --silent -X DELETE "localhost:$PORT/user/$USER_ID" | jq
}

get_random_question_id() {
  echo -n $(curl --silent "localhost:$PORT/question?limit=1&tags=test&random=true" | jq -r '.[0].id')
}

# Creates reads updates and deletes data through the API as a series of tests
main() {
  echo 'Running question tests...'
  question_tests
  echo 'All question tests passed!'
  echo 'Running user tests...'
  user_tests
  echo 'All user tests passed!'
  delete_all_questions
  echo 'Done!'
}

main "$@"
