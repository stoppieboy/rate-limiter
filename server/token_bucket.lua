-- Token Bucket in Redis (atomic)
-- KEYS[1] = bucket key (e.g., rate:{user123})
-- ARGV[1] = capacity (max tokens, >0)
-- ARGV[2] = refill_rate (tokens per second, >0)
-- ARGV[3] = requested (tokens to consume now, >=1) [default: 1]
-- ARGV[4] = now_millis (optional; if absent, uses Redis TIME)

local key = KEYS[1]
local capacity = tonumber(ARGV[1])
local rate = tonumber(ARGV[2])
local requested = tonumber(ARGV[3]) or 1

if not key or not capacity or not rate then
  return redis.error_reply("missing arguments: key, capacity, rate")
end
if capacity <= 0 or rate <= 0 or requested <= 0 then
  return redis.error_reply("invalid arguments: capacity, rate > 0 and requested >= 1 required")
end

-- Time source
local now_millis
if ARGV[4] then
  now_millis = tonumber(ARGV[4])
else
  local t = redis.call('TIME')          -- { seconds, microseconds }
  now_millis = tonumber(t[1]) * 1000 + math.floor(tonumber(t[2]) / 1000)
end

-- Load state: tokens (float), ts (millis)
local data = redis.call('HMGET', key, 'tokens', 'ts')
local tokens = tonumber(data[1])
local ts = tonumber(data[2])

-- Initialize if new
if tokens == nil or ts == nil then
  tokens = capacity
  ts = now_millis
else
  -- Refill based on elapsed time
  local elapsed = now_millis - ts
  if elapsed > 0 then
    local refill = (elapsed * rate) / 1000.0
    tokens = math.min(capacity, tokens + refill)
  end
  ts = now_millis
end

-- Try to consume
local allowed = 0
local tokens_after = tokens
if tokens >= requested then
  allowed = 1
  tokens_after = tokens - requested
end

-- Persist state
redis.call('HMSET', key, 'tokens', tokens_after, 'ts', ts)

-- Auto-expire to clean idle keys
local ttl
if tokens_after >= capacity then
  -- bucket full: short TTL (≈ time to fill twice, but capped at >=1s)
  ttl = math.ceil((2 * capacity) / rate)
else
  -- time until full
  local needed = capacity - tokens_after
  ttl = math.ceil(needed / rate)
end
if ttl < 1 then ttl = 1 end
redis.call('EXPIRE', key, ttl)

-- When denied, compute ms until enough tokens for this request
local ms_to_next = 0
if allowed == 0 then
  local deficit = requested - tokens
  ms_to_next = math.max(1, math.ceil((deficit / rate) * 1000.0))
end

-- Return:
-- 1) allowed (1/0)
-- 2) remaining tokens (float)
-- 3) ms until next token (or until request can succeed if denied)
-- 4) server time in ms (for client-side jitter/backoff if desired)
return {allowed, tokens_after, ms_to_next, now_millis }
-- return allowed