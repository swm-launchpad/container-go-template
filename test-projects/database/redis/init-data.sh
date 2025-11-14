#!/bin/bash
# Redis test data initialization script
# This script populates Redis with comprehensive test data for validation

echo "Initializing Redis with test data..."

# Wait for Redis to be ready
redis-cli ping > /dev/null 2>&1
while [ $? -ne 0 ]; do
    echo "Waiting for Redis to be ready..."
    sleep 1
    redis-cli ping > /dev/null 2>&1
done

echo "Redis is ready. Loading test data..."

# String data types - Simple key-value pairs
redis-cli SET "app:name" "Redis Test Application"
redis-cli SET "app:version" "1.0.0"
redis-cli SET "app:environment" "test"
redis-cli SETEX "session:user123" 3600 "john_doe"
redis-cli SETEX "session:user456" 3600 "jane_smith"

# Hash data types - User profiles
redis-cli HSET "user:1" username "john_doe" email "john@example.com" age 28 balance 1500.50 status "active"
redis-cli HSET "user:2" username "jane_smith" email "jane@example.com" age 34 balance 2750.00 status "active"
redis-cli HSET "user:3" username "bob_wilson" email "bob@example.com" age 45 balance 500.25 status "inactive"
redis-cli HSET "user:4" username "alice_brown" email "alice@example.com" age 22 balance 3200.75 status "active"
redis-cli HSET "user:5" username "charlie_davis" email "charlie@example.com" age 31 balance 0.00 status "active"

# List data types - Activity logs and queues
redis-cli RPUSH "logs:login" "2024-11-13 10:00:00 - user:1 logged in"
redis-cli RPUSH "logs:login" "2024-11-13 10:15:00 - user:2 logged in"
redis-cli RPUSH "logs:login" "2024-11-13 10:30:00 - user:4 logged in"
redis-cli RPUSH "queue:tasks" "process-order-1001"
redis-cli RPUSH "queue:tasks" "send-email-user123"
redis-cli RPUSH "queue:tasks" "generate-report-daily"

# Set data types - Tags and unique collections
redis-cli SADD "tags:electronics" "laptop" "smartphone" "mouse" "keyboard"
redis-cli SADD "tags:sports" "shoes" "ball" "racket" "bottle"
redis-cli SADD "tags:kitchen" "coffee-maker" "blender" "toaster"
redis-cli SADD "active:users" "user:1" "user:2" "user:4" "user:5"
redis-cli SADD "premium:users" "user:1" "user:2" "user:4"

# Sorted Set data types - Leaderboards and rankings
redis-cli ZADD "leaderboard:points" 1500 "john_doe" 2750 "jane_smith" 500 "bob_wilson" 3200 "alice_brown" 0 "charlie_davis"
redis-cli ZADD "products:ratings" 4.5 "laptop" 4.2 "mouse" 4.0 "coffee-maker" 4.7 "running-shoes" 4.3 "desk-lamp"
redis-cli ZADD "orders:by_date" 1707667200 "order:1001" 1707753600 "order:1002" 1707840000 "order:1003"

# Counter patterns - Statistics and metrics
redis-cli SET "counter:page_views" 12456
redis-cli SET "counter:api_calls" 8934
redis-cli SET "counter:user_registrations" 247
redis-cli INCR "counter:active_sessions"
redis-cli INCR "counter:active_sessions"
redis-cli INCR "counter:active_sessions"

# Bitmap data types - User activity tracking
redis-cli SETBIT "daily:active:2024-11-13" 1 1
redis-cli SETBIT "daily:active:2024-11-13" 2 1
redis-cli SETBIT "daily:active:2024-11-13" 4 1
redis-cli SETBIT "daily:active:2024-11-13" 5 1

# HyperLogLog - Unique visitor counting
redis-cli PFADD "unique:visitors:daily" "192.168.1.1" "192.168.1.2" "192.168.1.3" "10.0.0.1" "10.0.0.2"
redis-cli PFADD "unique:visitors:weekly" "192.168.1.1" "192.168.1.2" "192.168.1.3" "10.0.0.1" "10.0.0.2" "172.16.0.1" "172.16.0.2"

# Product catalog as hashes
redis-cli HSET "product:1" name "Laptop Pro 15" price 1299.99 stock 25 category "Electronics" rating 4.5
redis-cli HSET "product:2" name "Wireless Mouse" price 29.99 stock 150 category "Electronics" rating 4.2
redis-cli HSET "product:3" name "Coffee Maker" price 79.99 stock 50 category "Home & Kitchen" rating 4.0
redis-cli HSET "product:4" name "Running Shoes" price 89.99 stock 75 category "Sports" rating 4.7
redis-cli HSET "product:5" name "Desk Lamp" price 45.99 stock 100 category "Home & Kitchen" rating 4.3

# Cache expiration examples (TTL in seconds)
redis-cli SETEX "cache:trending:products" 3600 "laptop,mouse,shoes"
redis-cli SETEX "cache:user:1:profile" 1800 '{"username":"john_doe","email":"john@example.com"}'

# Geospatial data - Store locations
redis-cli GEOADD "stores:locations" -73.9857 40.7484 "NY-Store-1"
redis-cli GEOADD "stores:locations" -118.2437 34.0522 "LA-Store-1"
redis-cli GEOADD "stores:locations" -87.6298 41.8781 "CHI-Store-1"

# JSON-like data stored as strings
redis-cli SET "config:app" '{"maxUsers":1000,"timeout":30,"debug":false,"features":["analytics","notifications"]}'
redis-cli SET "config:database" '{"host":"localhost","port":5432,"poolSize":10}'

echo "-----------------------------------"
echo "Redis test data loaded successfully!"
echo "-----------------------------------"
echo "Data Summary:"
echo "- String keys: $(redis-cli KEYS 'app:*' 'session:*' 'counter:*' 'cache:*' 'config:*' | wc -l)"
echo "- Hash keys: $(redis-cli KEYS 'user:*' 'product:*' | wc -l)"
echo "- List keys: $(redis-cli KEYS 'logs:*' 'queue:*' | wc -l)"
echo "- Set keys: $(redis-cli KEYS 'tags:*' 'active:*' 'premium:*' | wc -l)"
echo "- Sorted set keys: $(redis-cli KEYS 'leaderboard:*' 'products:*' 'orders:*' | wc -l)"
echo "- Total keys: $(redis-cli DBSIZE)"
echo "-----------------------------------"
echo ""
echo "Test commands:"
echo "  redis-cli GET app:name"
echo "  redis-cli HGETALL user:1"
echo "  redis-cli LRANGE logs:login 0 -1"
echo "  redis-cli SMEMBERS tags:electronics"
echo "  redis-cli ZRANGE leaderboard:points 0 -1 WITHSCORES"
echo "-----------------------------------"
