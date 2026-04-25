#!/bin/bash

# Configuration
BASE_URL="http://localhost:15350/v1/sale"

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