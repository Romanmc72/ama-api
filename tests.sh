#!/usr/bin/env bash

set -euo pipefail

# --- Configuration ---
# Set the Firebase Authentication emulator host and port.
# The default port is 9099.
EMULATOR_HOST="http://localhost:9099"
PROJECT_ID=ama-dev
PORT=8088
ID_TOKEN='TODO: Get this to work with the Firebase emulator'
USER_EMAIL='t@t.com'
USER_PASSWORD='password123'
FIREBASE_USER_ID='SET ME!'
API_KEY='fake-key'

# This is not possible outside of the admin sdk because it uses some cryptographic
# stuff that are not available in plain curl/bash, need to move this all to golang :(
get_machine_to_machine_auth() {
  # This script creates a machine-to-machine user with a custom claim
  # and stores the returned custom authentication token in a variable.

  # Set a unique ID for the user you want to create.
  # This can be any string, but it's often a meaningful identifier.
  # In a real-world scenario, this might be a service name or an API client ID.
  USER_ID="question-admin"

  # Define the custom claims to attach to the user's token.
  # Here, we are setting the 'role' to 'admin'.
  CUSTOM_CLAIMS='{"role": "question/admin"}'

  # --- Script Logic ---

  # Step 1: Create the custom authentication token using curl.
  # We send a POST request to the emulator's `createCustomToken` endpoint.
  # The request body is a JSON object containing the UID and the claims.
  # The --data-raw flag is used to send the JSON string exactly as is.
  echo "Creating custom token for user: ${USER_ID} with claims: ${CUSTOM_CLAIMS}"

  # Use curl to send the request and capture the entire JSON response.
  # TOKEN_RESPONSE=$(curl -s -X POST "${EMULATOR_HOST}/createCustomToken" \
  # TOKEN_RESPONSE=$(curl -s -X POST "${EMULATOR_HOST}/emulator/v1/projects/${PROJECT_ID}/auth/token" \
  TOKEN_RESPONSE=$(curl -s -X POST "${EMULATOR_HOST}/identitytoolkit.googleapis.com/v1/projects/${PROJECT_ID}/auth/token" \
    -H "Content-Type: application/json" \
    --data-raw '{
      "uid": "'"${USER_ID}"'",
      "claims": '"${CUSTOM_CLAIMS}"'
    }')

  # Check if the curl request was successful
  if [ $? -ne 0 ]; then
    echo "Error: curl command failed."
    exit 1
  fi

  # Step 2: Parse the JSON response to extract the token.
  # We use a tool like 'jq' to parse the JSON. If 'jq' is not installed,
  # a simple grep and cut could be used, but 'jq' is more robust.
  # The '.customToken' filter retrieves the value associated with that key.
  # We strip the quotes using 'sed' or by not quoting the assignment.
  echo "Parsing the response to extract the custom token... ${TOKEN_RESPONSE}"
  CUSTOM_TOKEN=$(echo "${TOKEN_RESPONSE}" | jq -r '.customToken')

  # Check if the token was successfully extracted.
  if [ -z "${CUSTOM_TOKEN}" ]; then
    echo "Error: Failed to extract custom token from response."
    echo "Response was: ${TOKEN_RESPONSE}"
    exit 1
  fi

  # Step 3: Store and display the token.
  # The token is now stored in the CUSTOM_TOKEN variable.
  echo ""
  echo "Successfully created and retrieved custom token."
  echo "------------------------------------------------------"
  echo "USER ID: ${USER_ID}"
  echo "CUSTOM CLAIMS: ${CUSTOM_CLAIMS}"
  echo "CUSTOM TOKEN VARIABLE: ${CUSTOM_TOKEN}"
  echo "------------------------------------------------------"

  # Example of how to use the variable later in the script:
  # echo "This is how you can use the stored token variable: ${CUSTOM_TOKEN}"
  ID_TOKEN=${CUSTOM_TOKEN}
}

get_user_auth() {
  echo "Creating new user with email and password..."

  # Step 1: Send a POST request to the `signUp` endpoint.
  # The endpoint is an emulated version of the Identity Toolkit API.
  # The payload includes the email, password, and a flag to return a secure token.
  RESPONSE=$(curl -s -X POST "${EMULATOR_HOST}/identitytoolkit.googleapis.com/v1/accounts:signUp?key=${API_KEY}" \
    -H "Content-Type: application/json" \
    --data-raw '{
      "email": "'"${USER_EMAIL}"'",
      "password": "'"${USER_PASSWORD}"'",
      "returnSecureToken": true
    }')

  # Check if the curl command was successful
  if [ $? -ne 0 ]; then
    echo "Error: curl command failed."
    exit 1
  fi

  # Step 2: Check for API errors in the response.
  # The emulator returns a JSON object with an 'error' key if something went wrong.
  if echo "${RESPONSE}" | grep -q '"error"'; then
    echo "Error from API: $(echo "${RESPONSE}" | jq -r '.error.message')"
    exit 1
  fi

  # Step 3: Parse the JSON response to get the ID token.
  # We use 'jq' to extract the 'idToken' field.
  ID_TOKEN=$(echo "${RESPONSE}" | jq -r '.idToken')
  # a. Extract the second part (the payload) from the token.
  PAYLOAD_B64=$(echo "${ID_TOKEN}" | cut -d'.' -f2)

  # b. Base64 decode the payload string to get the JSON.
  PAYLOAD_JSON=$(echo "${PAYLOAD_B64}" | base64 --decode 2>/dev/null)

  # c. Use 'jq' to extract the 'sub' field, which contains the UID.
  FIREBASE_USER_ID=$(echo "${PAYLOAD_JSON}" | jq -r '.sub')
  
  echo ""
  echo "Successfully created and retrieved user token."
  echo "------------------------------------------------------"
  echo "USER ID: ${FIREBASE_USER_ID}"
  echo "EMAIL: ${USER_EMAIL}"
  echo "PASSWORD: ${USER_PASSWORD}"
  echo "------------------------------------------------------"
}

question_tests() {
  question_id=$(curl \
    --silent \
    -X POST \
    -H 'Content-Type: application/json' \
    -H "Authorization: Bearer ${ID_TOKEN}" \
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
      -H "Authorization: Bearer ${ID_TOKEN}" \
      -d '{"prompt": "test '$i' question", "tags": ["test", "another test", "'$i'"]}' \
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
    curl --silent \
      -X DELETE \
      -H "Authorization: Bearer ${ID_TOKEN}" \
      "localhost:$PORT/question/${each_question}" | jq
  done
}

user_tests() {
  echo 'Creating a user'
  USER_ID=$(curl --silent -X POST \
    -H 'Content-Type: application/json' \
    -H "Authorization: Bearer ${ID_TOKEN}" \
    -d '{
      "name": "test",
      "email": "'${USER_EMAIL}'",
      "tier": "free",
      "firebaseId": "'${FIREBASE_USER_ID}'",
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
  list_id=$(curl --silent -H "Authorization: Bearer ${ID_TOKEN}" localhost:$PORT/user/$USER_ID | jq -r '.lists[0].id')
  list_tests() {
    echo 'Testing user lists'
    for i in {0..5}
    do
      question_id=$(get_random_question_id)
      echo "Adding question $question_id to list $list_id"
      curl \
        --silent \
        -X POST \
        -H "Authorization: Bearer ${ID_TOKEN}" \
        "localhost:$PORT/user/$USER_ID/list/$list_id/question/$question_id" \
        | jq
    done
    echo 'Fetching questions from list and removing them'
    for each_question in $(curl --silent -H "Authorization: Bearer ${ID_TOKEN}" "localhost:$PORT/user/$USER_ID/list/$list_id?limit=10" | jq -r '.questions[] | .id')
    do
      echo "Removing question $each_question from list $list_id"
      curl --silent -X DELETE -H "Authorization: Bearer ${ID_TOKEN}" "localhost:$PORT/user/$USER_ID/list/$list_id/question/$each_question" | jq
    done
  }
  # echo 'Testing user lists...'
  # list_tests
  # echo 'List tests passed!'
  # echo 'Deleting user'
  # curl --silent -X DELETE -H "Authorization: Bearer ${ID_TOKEN}" "localhost:$PORT/user/$USER_ID" | jq
}

get_random_question_id() {
  echo -n $(curl --silent -H "Authorization: Bearer ${ID_TOKEN}" "localhost:$PORT/question?limit=1&tags=test&random=true" | jq -r '.[0].id')
}

# Creates reads updates and deletes data through the API as a series of tests
main() {
  # get_machine_to_machine_auth
  # echo 'Running question tests...'
  # question_tests
  # echo 'All question tests passed!'
  echo 'Running user tests...'
  get_user_auth
  user_tests
  echo 'All user tests passed!'
  # delete_all_questions
  echo 'Done!'
}

main "$@"
