package redisPriorityQueue

var scriptDeleteMessage = ``

var ScriptPushOne = `
-- This script increments the score of a member in a sorted set
local key = KEYS[1]
local increment = ARGV[1]
local member = ARGV[2]

-- Use the ZINCRBY command
local new_score = redis.call('ZINCRBY', key, increment, member)

return new_score`

var ScriptBatchPush = `
local key = KEYS[1]
local results = {}
for i = 1, #ARGV, 2 do
	local increment = tonumber(ARGV[i])
	local member = ARGV[i + 1]
	local new_score = redis.call('ZINCRBY', key, increment, member)
	table.insert(results, new_score)
end
return results
`

var ScriptPopOne = `
-- This script pops the member with the lowest score from a sorted set
local member = redis.call('ZRANGEBYSCORE', KEYS[1], ARGV[1], ARGV[2], 'LIMIT', 0, 1)
if #member > 0 then
    redis.call('ZREM', KEYS[1], member[1])
end
return member
`

var ScriptBatchPop = `
-- This script pops the member with the lowest score from a sorted set
local members = redis.call('ZRANGEBYSCORE', KEYS[1], ARGV[1], ARGV[2], 'LIMIT', 0, ARGV[3])
for i, member in ipairs(members) do
    redis.call('ZREM', KEYS[1], member)
end
return members
`
