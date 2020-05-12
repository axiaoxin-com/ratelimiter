-- caller 参数白名单配置
-- 在需要限频的接口请求中，如果 caller 在这个配置中则直接放行
local _M = {}

-- 示例：将 caller 值为 _dev_whitelist 的请求设为白名单
_M["_dev_whitelist"] = 1

return _M
