package constBlogPg

// 逻辑：1. 获取集合中的所有键名 2. 批量获取这些键的值 3. 返回键和值的组合
const ArticleCategoryLua = `
		-- 第一步：获取集合 blog:category:ids 中的所有键名
		local keys = redis.call('SMEMBERS', KEYS[1])
		-- 第二步：批量获取这些键对应的字符串值（MGET 比循环 GET 效率高）
		local values = redis.call('MGET', unpack(keys))
		-- 第三步：组装键值对返回（格式：[key1, value1, key2, value2, ...]）
		local result = {}
		for i = 1, #keys do
			table.insert(result, keys[i])
			table.insert(result, values[i])
		end
		return result
	`
