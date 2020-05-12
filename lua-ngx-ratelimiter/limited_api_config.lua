-- 需要限流的 uri 名配置
-- 限流只对在该配置中出现的 uri 进行限频
local _M = {}

-- 示例：将接口 http://localhost:8080/_dev_api/limited 设置为需要进行限频判断
_M["/_dev_api/limited"] = 1

return _M
