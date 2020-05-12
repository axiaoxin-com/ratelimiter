-- 工具方法集合
local json = require "cjson"

local _M = {}

function _M.new(self)
	return self
end

-- 获取客户端 IP
function _M.get_client_ip(self)
    local client_ip = nil
    -- 如果 x_forwarded_for 中有值则使用其中的第一 IP 作为客户端 IP
    if ngx.var.http_x_forwarded_for ~= nil then
        client_ip = string.match(ngx.var.http_x_forwarded_for, "%d+.%d+.%d+.%d+", "1");
    end

    -- 获取失败则使用 headers 中的 X-Real-IP 作为客户端 IP
    if client_ip == nil then
        client_ip = ngx.req.get_headers()["X-Real-IP"]
    end

    -- 仍然失败则使用内置 remote_addr
    if client_ip == nil then
        client_ip = ngx.var.remote_addr
    end

    return client_ip
end

-- 健壮版 json decode
function _M.json_decode(data)
    if data == nil then
        ngx.log(ngx.WARN, "data is nil")
        return
    end

    if type(data) == "string" then
        local ok, r = pcall(json.decode, data)
        if not ok then
            ngx.log(ngx.ERR, "decode failed for data:"..data)
            return
        end
        return r
    end

    return data
end

-- 获取请求参数
function _M.get_req_args()
    -- GET 方法从 URL 中获取
    local args = nil
    if ngx.var.request_method == "GET" then
        args = ngx.req.get_uri_args()
    -- POST 方法从请求体中获取
    elseif ngx.var.request_method == "POST" then
        ngx.req.read_body() -- 解析 body 参数之前一定要先读取 body
        local data = ngx.req.get_body_data()
        args = _M.json_decode(data)
    end
    return args
end

-- 字符串切割 返回数组
function _M.str_split(str, sep)
    local len_sep = #sep
    local t,c,p1,p2 = {},1,1,nil
    while true do
        p2 = str:find(sep,p1,true)
        if p2 then
            t[c] = str:sub(p1,p2-1)
            p1 = p2 + len_sep
        else
            t[c] = str:sub(p1)
            return t
        end
        c = c+1
    end
end


-- 判断字符串开头 返回 bool
function _M.str_startswith(s, prefix)
   return string.sub(s, 1, string.len(prefix)) == prefix
end


return _M
