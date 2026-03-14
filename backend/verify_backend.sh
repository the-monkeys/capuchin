#!/bin/bash

BASE_URL="http://localhost:8080"
EMAIL="testuser@example.com"
PASSWORD="password123"

echo "1. Testing Health Endpoint..."
curl -s $BASE_URL/health | grep -q "ok" && echo "  - Success: Health check passed" || echo "  - Error: Health check failed"

echo -e "\n2. Testing Signup..."
SIGNUP_RES=$(curl -s -X POST $BASE_URL/signup \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"$EMAIL\", \"password\": \"$PASSWORD\"}")
echo "  - Response: $SIGNUP_RES"

echo -e "\n3. Testing Login..."
LOGIN_RES=$(curl -s -X POST $BASE_URL/login \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"$EMAIL\", \"password\": \"$PASSWORD\"}")

TOKEN=$(echo $LOGIN_RES | grep -oP '"token":"\K[^"]+')

if [ -z "$TOKEN" ]; then
  echo "  - Error: Failed to retrieve token"
  exit 1
fi
echo "  - Success: Received JWT token"

echo -e "\n4. Testing Add Todo..."
TODO_RES=$(curl -s -X POST $BASE_URL/api/user/todo \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"item": "Try out the new backend", "completed": false}')
echo "  - Response: $TODO_RES"

TODO_ID=$(echo $TODO_RES | grep -oP '"id":"\K[^"]+')

echo -e "\n5. Testing Get Todos..."
curl -s -X GET $BASE_URL/api/user/todo -H "Authorization: Bearer $TOKEN" | grep -q "$TODO_ID" && echo "  - Success: Todo found in list" || echo "  - Error: Todo not found"

echo -e "\n6. Testing Delete Todo..."
DELETE_RES=$(curl -s -X DELETE $BASE_URL/api/user/todo/$TODO_ID -H "Authorization: Bearer $TOKEN")
echo "  - Response: $DELETE_RES"

echo -e "\n7. Testing Logout..."
LOGOUT_RES=$(curl -s -X POST $BASE_URL/api/user/logout -H "Authorization: Bearer $TOKEN")
echo "  - Response: $LOGOUT_RES"

echo -e "\n8. Testing Protected Route After Logout..."
POST_LOGOUT_RES=$(curl -s -w "\n%{http_code}" -X GET $BASE_URL/api/user/todo -H "Authorization: Bearer $TOKEN")
HTTP_CODE=$(echo "$POST_LOGOUT_RES" | tail -n1)
if [ "$HTTP_CODE" -eq 401 ]; then
  echo "  - Success: Request with logged-out token was rejected"
else
  echo "  - Error: Request with logged-out token succeeded unexpectedly (HTTP $HTTP_CODE)"
fi
