#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo "Starting API tests..."
echo "===================="

# Test 1: Submit job with valid data
echo -e "\n${GREEN}Test 1: Submit job with valid data${NC}"
RESPONSE=$(curl -s -X POST http://localhost:8080/api/submit \
-H "Content-Type: application/json" \
-d '{
    "count": 2,
    "visits": [
        {
            "store_id": "RP00001",
            "image_url": [
                "https://www.gstatic.com/webp/gallery/2.jpg",
                "https://www.gstatic.com/webp/gallery/3.jpg"
            ],
            "visit_time": "2024-01-13T12:00:00Z"
        },
        {
            "store_id": "RP00002",
            "image_url": [
                "https://www.gstatic.com/webp/gallery/3.jpg"
            ],
            "visit_time": "2024-01-13T12:30:00Z"
        }
    ]
}')

echo "Response: $RESPONSE"
JOB_ID=$(echo $RESPONSE | jq -r '.job_id')

if [ ! -z "$JOB_ID" ]; then
    echo -e "${GREEN}✓ Test 1 passed: Got job_id: $JOB_ID${NC}"
else
    echo -e "${RED}✗ Test 1 failed: No job_id received${NC}"
fi

# Test 2: Check job status
echo -e "\n${GREEN}Test 2: Check job status${NC}"
echo "Waiting 2 seconds for job processing..."
sleep 2

STATUS_RESPONSE=$(curl -s "http://localhost:8080/api/status?jobid=$JOB_ID")
echo "Response: $STATUS_RESPONSE"
STATUS=$(echo $STATUS_RESPONSE | jq -r '.status')

if [ "$STATUS" == "completed" ] || [ "$STATUS" == "ongoing" ] || [ "$STATUS" == "failed" ]; then
    echo -e "${GREEN}✓ Test 2 passed: Got valid status: $STATUS${NC}"
else
    echo -e "${RED}✗ Test 2 failed: Invalid status${NC}"
fi

# Test 3: Submit job with invalid store ID
echo -e "\n${GREEN}Test 3: Submit job with invalid store ID${NC}"
RESPONSE=$(curl -s -X POST http://localhost:8080/api/submit \
-H "Content-Type: application/json" \
-d '{
    "count": 1,
    "visits": [
        {
            "store_id": "INVALID_STORE",
            "image_url": [
                "https://www.gstatic.com/webp/gallery/2.jpg"
            ],
            "visit_time": "2024-01-13T12:00:00Z"
        }
    ]
}')

echo "Response: $RESPONSE"
JOB_ID=$(echo $RESPONSE | jq -r '.job_id')

if [ ! -z "$JOB_ID" ]; then
    echo -e "${GREEN}✓ Test 3 passed: Got job_id: $JOB_ID${NC}"
    
    echo "Waiting 2 seconds for job processing..."
    sleep 2
    
    STATUS_RESPONSE=$(curl -s "http://localhost:8080/api/status?jobid=$JOB_ID")
    echo "Status Response: $STATUS_RESPONSE"
    
    if [[ $STATUS_RESPONSE == *"failed"* ]]; then
        echo -e "${GREEN}✓ Test 3 passed: Job failed as expected${NC}"
    else
        echo -e "${RED}✗ Test 3 failed: Job should have failed${NC}"
    fi
else
    echo -e "${RED}✗ Test 3 failed: No job_id received${NC}"
fi

# Test 4: Submit job with count mismatch
echo -e "\n${GREEN}Test 4: Submit job with count mismatch${NC}"
RESPONSE=$(curl -s -X POST http://localhost:8080/api/submit \
-H "Content-Type: application/json" \
-d '{
    "count": 2,
    "visits": [
        {
            "store_id": "RP00001",
            "image_url": [
                "https://www.gstatic.com/webp/gallery/2.jpg"
            ],
            "visit_time": "2024-01-13T12:00:00Z"
        }
    ]
}')

echo "Response: $RESPONSE"
if [[ $RESPONSE == *"count does not match"* ]]; then
    echo -e "${GREEN}✓ Test 4 passed: Got expected count mismatch error${NC}"
else
    echo -e "${RED}✗ Test 4 failed: Should have received count mismatch error${NC}"
fi

# Test 5: Check status with invalid job ID
echo -e "\n${GREEN}Test 5: Check status with invalid job ID${NC}"
RESPONSE=$(curl -s "http://localhost:8080/api/status?jobid=invalid_job_id")
echo "Response: $RESPONSE"

if [[ $RESPONSE == *"job not found"* ]]; then
    echo -e "${GREEN}✓ Test 5 passed: Got expected job not found error${NC}"
else
    echo -e "${RED}✗ Test 5 failed: Should have received job not found error${NC}"
fi

echo -e "\n${GREEN}All tests completed!${NC}"