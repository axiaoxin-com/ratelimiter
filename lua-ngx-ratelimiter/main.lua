-- 限流入口文件，在 nginx 配置文件中使用 access_by_lua_file 进行加载

local function main()
    local json = require "cjson"
    local utils = require "utils"
    local key = require "key"
    local token_bucket = require "token_bucket"

    -- 生成限流 key
    local limit_key = key:new()
    -- key 创建失败则直接放行
    if limit_key.key == nil then
        ngx.log(ngx.ERR, "ratelimiter: failed to create key, so this request is not limited.")
        return
    end

    -- key 在黑名单中直接返回错误 JSON
    if limit_key:is_in_blacklist() then
        ngx.header["Content-Type"] = "application/json"
        ngx.status = ngx.HTTP_FORBIDDEN
        local args = utils.get_req_args() or {}
        local rsp = {
            seqId = args["seqId"],
            code = ngx.HTTP_FORBIDDEN,
            msg = "当前请求命中黑名单",
        }
        ngx.say(json.encode(rsp))
        ngx.exit(ngx.HTTP_FORBIDDEN)
    end

    -- 检查 caller 参数是否在白名单中，在则直接放行
    if limit_key:is_in_whitelist() then
        return
    end

    -- 检查 key 是否需要限流，key 不在 limited keys 中也直接放行
    if not limit_key:is_limited_api() then
        return
    end

    -- 获取令牌桶检测限流状态
    local token_bucket = token_bucket:new(limit_key)
    -- 被限流直接返回错误 JSON
    if token_bucket:is_limited() then
        ngx.header["Content-Type"] = "application/json"
        ngx.status = ngx.HTTP_TOO_MANY_REQUESTS
        local args = utils.get_req_args() or {}
        local rsp = {
            seqId = args["seqId"],
            code = ngx.HTTP_TOO_MANY_REQUESTS,
            msg = "请求的次数超过了频率限制",
        }
        ngx.say(json.encode(rsp))
        ngx.exit(ngx.HTTP_TOO_MANY_REQUESTS)
    end
end


local ok, r = xpcall(main, debug.traceback)
if not ok then
    ngx.log(ngx.ERR, "ratelimiter: main pcall failed", r)
end
