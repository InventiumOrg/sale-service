#!/bin/bash

# Configuration
BASE_URL="http://localhost:15350/v1/sale"
AUTH_TOKEN="eyJhbGciOiJSUzI1NiIsImNhdCI6ImNsX0I3ZDRQRDExMUFBQSIsImtpZCI6Imluc18yclhPVXZhbHRjWHBqSTJRQUg3WFZFTUlRNWkiLCJ0eXAiOiJKV1QifQ.eyJleHAiOjE3NjIzNDI4MTQsImZ2YSI6Wzk5OTk5LC0xXSwiaWF0IjoxNzYyMzM5MjE0LCJpc3MiOiJodHRwczovL2FjZS1sb3VzZS00Mi5jbGVyay5hY2NvdW50cy5kZXYiLCJuYmYiOjE3NjIzMzkyMDQsIm9yZ19pZCI6Im9yZ18zMGNJY090cUtIVFpvTWpOYVUxQ0h1dmRsd3kiLCJvcmdfcGVybWlzc2lvbnMiOltdLCJvcmdfcm9sZSI6Im9yZzptZW1iZXIiLCJvcmdfc2x1ZyI6InRlc3QtMTc1MzkyMDUxMiIsInNpZCI6InNlc3NfMzUzV1BxR3djV3JaNTMyRWNLUVB1UnFQcUFwIiwic3RzIjoiYWN0aXZlIiwic3ViIjoidXNlcl8zMGNIUVVIU3pYVDJ6Y3lGdU81Snc5emxoZHMifQ.pcUOv_QS9uAIgqINm3GFYruG75BRh8mlHb9P_aTbuEFab5dGkQ-zKFRqMuWUAVfJ8pWEFyZpbzoJ2zxFclXIP_wzLGIOZM5S4AyM_F2OUlg8n12hZhjS8gl3DQNBxssdrUZoxkPKPekUkUECnqML-KiGFzwDMiVFmBdkTlf29ZzXlAOXHZJNAMIadVnqZrRL4rgi3fuyafbgcMAqSNxPMzF9SPhnXOlSGA6KqNNvxRBsyIUFpgFpOl2JrgERvSKT9sJMyF0TpB6k_nBtxF27ZeUncGP98Cg_8qcbTEimzvWH_nZ8RyVQACnbHIssJKgP125zwdpUu7cMUFlUKz7zhA"

# Number of requests to make (default: 5)
NUM_REQUESTS=${1:-5}

# Delay between requests in seconds (default: 1)
DELAY=${2:-1}

echo "Starting simulation with $NUM_REQUESTS requests..."
echo "Delay between requests: ${DELAY}s"
echo "----------------------------------------"

# Loop through and make requests
for i in $(seq 1 $NUM_REQUESTS); do
    echo "Request $i/$NUM_REQUESTS"
    
    # Generate some variation in the data
    SALE_ID=1
    POS_ID="pos_$(printf "%04d" $i)"
    RECIPE_ID=$((($i % 5) + 1))  # Cycle through recipe IDs 1-5
    
    curl --location --request GET "${BASE_URL}/${SALE_ID}" \
        --header "Authorization: Bearer ${AUTH_TOKEN}" \
        --form "Name=\"Test Sale $i\"" \
        --form "PosID=\"${POS_ID}\"" \
        --form "SaleRecipeID=\"${RECIPE_ID}\"" \
        --silent --show-error
    
    echo ""
    echo "Response for request $i completed"
    
    # Add delay between requests (except for the last one)
    if [ $i -lt $NUM_REQUESTS ]; then
        echo "Waiting ${DELAY}s before next request..."
        sleep $DELAY
    fi
    
    echo "----------------------------------------"
done

echo "Simulation completed!"