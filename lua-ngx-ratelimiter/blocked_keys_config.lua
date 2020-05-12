-- is_in_blacklist 中对限流 key 的黑名单设置
-- _M table 中 key 为需要设置为黑名单的 key ，值为非 nil 的任意值
local _M = {}

-- 示例：将需要拉黑的 key ratelimiter:gz:/_dev_api/blocked:axiaoxin:127.0.0.1 放入 table
_M["ratelimiter:gz:/_dev_api/blocked:_dev_axiaoxin:127.0.0.1"] = 1

return _M
