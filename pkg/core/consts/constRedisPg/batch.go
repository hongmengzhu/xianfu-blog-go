package constRedisPg

const BATCH_MGET = `
-- 1. 从 KEYS[1] 取出逗号分隔的字符串
local str = KEYS[1]

-- 2. 按逗号拆分，生成数组 arr
local arr = {}
for s in string.gmatch(str, '[^,]+') do
    table.insert(arr, s)
end

local result = redis.call('MGET', unpack(arr))
return cjson.encode(result)
`
