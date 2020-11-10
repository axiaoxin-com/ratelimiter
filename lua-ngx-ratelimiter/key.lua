--  业务相关的限流 key 相关方法
local json = require "cjson"
local utils = require "utils"
local blocked_keys = require "blocked_keys_config"
local limited_apis = require "limited_api_config"
local whitelist = require "caller_whitelist_config"

local _M = {
    separator = ":",
    prefix = "ratelimiter",
}

-- 生成限流 key
-- ratelimiter:region:api:caller:ip
function _M:new()
    ngx.log(ngx.DEBUG, "ratelimiter: new key module")
    local o = {}
    setmetatable(o, {__index = self})

    local host = utils.str_split(ngx.var.host, ".")[1]
    local client_ip = utils.get_client_ip() or ""

    -- 获取请求参数
    local caller = ""
    local args = utils.get_req_args()
    if args ~= nil then
        caller = args["caller"] or ""
    end

    local api_name = ngx.var.uri
    if ngx.var.request_uri ~= nil and ngx.var.request_uri ~= "" then
        local r = utils.str_split(ngx.var.request_uri, "?")
        if #r >= 1 then
           api_name = r[1]
        end
    end

    o.api_name = api_name
    o.caller = caller
    o.client_ip = client_ip
    o.key = o.prefix .. o.separator .. host .. o.separator .. api_name .. o.separator .. caller .. o.separator .. client_ip
    ngx.log(ngx.INFO, "ratelimiter: gen a new limit key -> " .. o.key)
	return o
end

-- 判断黑名单
function _M:is_in_blacklist()
    ngx.log(ngx.INFO, "ratelimiter: check blacklist for key:" .. self.key)

    -- key 在 blocked_keys 中表示被拉黑
    if blocked_keys[self.key] ~= nil then
        ngx.log(ngx.ERR, "ratelimiter: key=" .. self.key .. " hit blocked key")
        return true
    end
    return false
end

-- 判断是否是需要限流的接口
function _M:is_limited_api()
    ngx.log(ngx.INFO, "ratelimiter: check is limited api for key=" .. self.key)

    -- api_name 在 limited_apis 中表示需要进行限流检查
    if limited_apis[self.api_name] ~= nil then
        ngx.log(ngx.WARN, "ratelimiter: key=" .. self.key .. " hit limited api")
        return true
    end
    ngx.log(ngx.INFO, "ratelimiter: api=" .. self.api_name .. " does not need to be limited.")
    return false
end

-- 判断白名单
function _M:is_in_whitelist()
    ngx.log(ngx.INFO, "ratelimiter: check whitelist for key=" .. self.key)

    local caller = ''
    local args = utils.get_req_args()
    if args ~= nil then
        caller =  args['caller'] or ''
    end

    -- 请求参数中 caller 值在白名单中的全部放行
    if whitelist[caller] ~= nil then
        ngx.log(ngx.WARN, "ratelimiter: caller=" .. caller .. " hit whitelist")
        return true
    end
    return false
end

return _M
