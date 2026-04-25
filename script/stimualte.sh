#!/bin/bash

# Configuration
KONG_URL="http://127.0.0.1/v1/sale/create"
# Number of requests to make (default: 5)
NUM_REQUESTS=${1:-5}

# Delay between requests in seconds (default: 1)
DELAY=${2:-1}

create_sale_unit(){
    # Loop through and make requests
    for i in $(seq 1 $NUM_REQUESTS); do
        echo "Request $i/$NUM_REQUESTS"
        
        # Generate some variation in the data
        SALE_ID=1
        POS_ID="$(printf "%04d" $i)"
        RECIPE_ID=$((($i % 5) + 1))  # Cycle through recipe IDs 1-5
        PRICE=$((($i % 100) + 1))
        ORDER_ID=$((($i % 1000) + 1))
        curl --location --request POST "${KONG_URL}" \
            --form "PosID=\"${POS_ID}\"" \
            --form "Price=\"${PRICE}\"" \
            --form "RecipeID=\"${RECIPE_ID}\"" \
            --form "OrderID=\"${ORDER_ID}\"" \
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
}

echo "Starting simulation with $NUM_REQUESTS requests..."
echo "Delay between requests: ${DELAY}s"
echo "----------------------------------------"

create_sale_unit $NUM_REQUESTS $DELAY

echo "Simulation completed!"
