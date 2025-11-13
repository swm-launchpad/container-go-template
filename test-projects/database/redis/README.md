# Redis Test Project

This project contains comprehensive Redis test data for validating Redis functionality across various data types.

## Environment Variables
- REDIS_PASSWORD: Redis password (optional)
- REDIS_PORT: Redis port (default: 6379)

## Initialization

To load test data into Redis:
```bash
chmod +x init-data.sh
./init-data.sh
```

## Test Data Overview

The initialization script creates the following test data:

### String Data Types
- Application configuration (app:name, app:version, app:environment)
- Session data with TTL (session:user123, session:user456)
- Counters (page views, API calls, registrations)
- Cached data with expiration

### Hash Data Types
- User profiles (user:1 through user:5) with username, email, age, balance, status
- Product catalog (product:1 through product:5) with name, price, stock, category, rating

### List Data Types
- Login logs (logs:login)
- Task queues (queue:tasks)

### Set Data Types
- Product tags by category (tags:electronics, tags:sports, tags:kitchen)
- Active users collection (active:users)
- Premium users collection (premium:users)

### Sorted Set Data Types
- User leaderboard by points (leaderboard:points)
- Products by rating (products:ratings)
- Orders by date (orders:by_date)

### Bitmap Data Types
- Daily active users tracking (daily:active:YYYY-MM-DD)

### HyperLogLog Data Types
- Unique visitors counting (unique:visitors:daily, unique:visitors:weekly)

### Geospatial Data Types
- Store locations (stores:locations) with NYC, LA, and Chicago stores

## Verification Commands

Basic connectivity:
```bash
redis-cli PING
```

String operations:
```bash
redis-cli GET app:name
redis-cli GET counter:page_views
redis-cli INCR counter:api_calls
```

Hash operations:
```bash
redis-cli HGETALL user:1
redis-cli HGET product:1 name
redis-cli HGET product:1 price
```

List operations:
```bash
redis-cli LRANGE logs:login 0 -1
redis-cli LLEN queue:tasks
redis-cli LPOP queue:tasks
```

Set operations:
```bash
redis-cli SMEMBERS tags:electronics
redis-cli SISMEMBER active:users "user:1"
redis-cli SCARD premium:users
```

Sorted Set operations:
```bash
redis-cli ZRANGE leaderboard:points 0 -1 WITHSCORES
redis-cli ZREVRANGE leaderboard:points 0 2 WITHSCORES
redis-cli ZSCORE products:ratings "laptop"
```

Geospatial operations:
```bash
redis-cli GEOPOS stores:locations "NY-Store-1"
redis-cli GEODIST stores:locations "NY-Store-1" "LA-Store-1" km
```

Statistics:
```bash
redis-cli DBSIZE
redis-cli INFO stats
redis-cli KEYS "*"
```

## Test Scenarios

### Scenario 1: User Session Management
```bash
# Create session
redis-cli SETEX "session:newuser" 3600 "username123"

# Check session
redis-cli GET "session:newuser"

# Check TTL
redis-cli TTL "session:newuser"
```

### Scenario 2: Leaderboard Operations
```bash
# Add new score
redis-cli ZADD leaderboard:points 1800 "new_player"

# Get top 3 players
redis-cli ZREVRANGE leaderboard:points 0 2 WITHSCORES

# Get player rank
redis-cli ZREVRANK leaderboard:points "alice_brown"
```

### Scenario 3: Product Search by Category
```bash
# Get all electronics tags
redis-cli SMEMBERS tags:electronics

# Get product details
redis-cli HGETALL product:1
redis-cli HGETALL product:2
```

### Scenario 4: Activity Tracking
```bash
# Count unique visitors
redis-cli PFCOUNT unique:visitors:daily

# Check daily active users
redis-cli BITCOUNT daily:active:2024-11-13
```

## Performance Testing

Test Redis performance:
```bash
redis-benchmark -q -n 100000
```

Monitor Redis operations in real-time:
```bash
redis-cli MONITOR
```
