-- redis 相关配置
local _M = {
    host = "127.0.0.1",           -- redis 服务 IP（本地必须使用IP，不能使用localhost）
    port = 6379,                  -- redis 端口号
    password = nil,               -- redis 密码，没有密码设置为 nil
    db_index = 1,                 -- 使用的 redis db 索引
    connect_timeout = 1000,       -- 连接 redis 超时时间（毫秒）
    pool_size = 100,              -- 连接池大小
    pool_max_idle_time = 120000,  -- 连接池中每个连接 keepalive 时长（毫秒）
}

return _M
