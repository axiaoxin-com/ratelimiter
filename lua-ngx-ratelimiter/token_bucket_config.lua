-- token bucket 的相关配置项
local _M = {}

-- 默认的桶配置
_M['default'] = {
    fill_count = 500,                    -- 令牌桶每次填充数
    interval_microsecond = 1000000,      -- 每次填充间隔时间（微秒）
    bucket_capacity = 500,               -- 令牌桶容量（最大限流值/秒）
    expire_second = 60 * 10,             -- 过期时间
}

-- 对特定的 caller 和 api 设置自定义的桶配置（限额配置）
-- 示例：对 caller 为 _dev_ashin 的 /_dev_api/limited 接口的请求进行自定义限流配置（每秒只能请求1次）
_M['_dev_ashin:/_dev_api/limited'] = {
    fill_count = 1,                      -- 令牌桶每次填充数
    interval_microsecond = 1000000,      -- 每次填充间隔时间（微秒）
    bucket_capacity = 1,                 -- 令牌桶容量（最大限流值/秒）
    expire_second = 60 * 10,             -- 过期时间
}
return _M
