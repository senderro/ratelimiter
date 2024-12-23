local tokens_key = KEYS[1]..":tokens"
local last_access_key = KEYS[1]..":last_access"

local capacity = tonumber(ARGV[1])
local rate = tonumber(ARGV[2])
local now = tonumber(ARGV[3])
local requested = tonumber(ARGV[4])

local last_tokens = tonumber(redis.call("get", tokens_key))
if last_tokens == nil then
    last_tokens = capacity
end

local last_access = tonumber(redis.call("get", last_access_key))
if last_access == nil then
    last_access = 0
end

local elapsed = math.max(0, now - last_access)
local add_tokens = math.floor(elapsed * rate / 1000000)
local new_tokens = math.min(capacity, last_tokens + add_tokens)

local new_access_time = last_access + math.ceil(add_tokens * 1000000 / rate)

local allowed = new_tokens >= requested
if allowed then
    new_tokens = new_tokens - requested
end

redis.call("setex", tokens_key, 60, new_tokens)
redis.call("setex", last_access_key, 60, new_access_time)

return allowed and 1 or 0
