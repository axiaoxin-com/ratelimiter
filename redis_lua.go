package ratelimiter

import "github.com/go-redis/redis/v8"

// redis 中执行的 lua 脚本判断 key 是否应该被限频
// 返回 json 字符串： is_limited false:无需限频，true:需限频
var tokenBucketRedisLuaIsLimitedScript = redis.NewScript(`
    -- 兼容低版本 redis 手动打开允许随机写入 （执行 TIME 指令获取时间）
    -- 避免报错 Write commands not allowed after non deterministic commands. Call redis.replicate_commands() at the start of your script in order to switch to single commands         replication mode.
    -- Redis 出于数据一致性考虑，要求脚本必须是纯函数的形式，也就是说对于一段 Lua 脚本给定相同的参数，重复执行其结果都是相同的。
    -- 这个限制的原因是 Redis 不仅仅是单机版的内存数据库，它还支持主从复制和持久化，执行过的 Lua 脚本会复制给 slave 以及持久化到磁盘，如果重复执行得到结果不同，那么就会出现内存、磁盘、 slave 之间的数据不一致，在 failover 或者重启之后造成数据错乱影响业务。
    -- 如果执行过非确定性命令（也就是 TIME ，因为时间是随机的）， Redis 就不允许执行写命令，以此来保证数据一致性。
    -- 在 Redis 中 time 命令是一个随机命令（时间是变化的），在 Lua 脚本中调用了随机命令之后禁止再调用写命令， Redis 中一共有 10 个随机类命令：
    -- spop 、 srandmember 、 sscan 、 zscan 、 hscan 、 randomkey 、 scan 、 lastsave 、 pubsub 、 time
    -- 在执行 redis.replicate_commands() 之后， Redis 就不再是把整个 Lua 脚本同步给 slave 和持久化，而是只把脚本中的写命令使用 multi/exec 包裹后直接去做复制，那么 slave 和持久化只复制了写命名，而写入的也是确定的结果。
    redis.replicate_commands()

    redis.log(redis.LOG_DEBUG, "------------ ratelimiter script begin ------------")
    -- 获取参数
    local p_key = KEYS[1]
    local p_bucket_capacity = tonumber(ARGV[1])
    local p_fill_count = tonumber(ARGV[2])
    local p_interval_microsecond = tonumber(ARGV[3])
    local p_expire_second = tonumber(ARGV[4])

    -- 返回结果
    local result = {}
    result['p_key'] = p_key
    result['p_fill_count'] = p_fill_count
    result['p_bucket_capacity'] = p_bucket_capacity
    result['p_interval_microsecond'] = p_interval_microsecond
    result['p_expire_second'] = p_expire_second

    -- 每次填充 token 数为 0 或 令牌桶容量为 0 则表示限制该请求 直接返回 无需操作 redis
    if p_fill_count <= 0 or p_bucket_capacity <= 0 then
        result['msg'] = "be limited by p_fill_count or p_bucket_capacity"
        result['is_limited'] = true
        return cjson.encode(result)
    end

    -- 判断桶是否存在
    local exists = redis.call("EXISTS", p_key)
    redis.log(redis.LOG_DEBUG, "ratelimiter: key:" .. p_key .. ", exists:" .. exists)

    -- 桶不存在则在 redis 中创建桶 并消耗当前 token
    if exists == 0 then
        -- 本次填充时间戳
        local now_timestamp_array = redis.call("TIME")
        -- 微秒级时间戳
        local last_consume_timestamp = tonumber(now_timestamp_array[1]) * 1000000 + tonumber(now_timestamp_array[2])
        redis.log(redis.LOG_DEBUG, "ratelimiter: last_consume_timestamp:" .. last_consume_timestamp .. ", remain_token_count:" .. p_bucket_capacity)
        -- 首次请求 默认为满桶  消耗一个 token
        local remain_token_count = p_bucket_capacity - 1

        -- 将当前秒级时间戳和剩余 token 数保存到 redis
        redis.call("HMSET", p_key, "last_consume_timestamp", last_consume_timestamp, "remain_token_count", remain_token_count)
        -- 设置 redis 的过期时间
        redis.call("EXPIRE", p_key, p_expire_second)
        redis.log(redis.LOG_DEBUG, "ratelimiter: call HMSET for creating bucket")
        redis.log(redis.LOG_DEBUG, "------------ ratelimiter script end ------------")

        -- 保存 result 信息
        result['msg'] = "key not exists in redis"
        -- string format 避免科学计数法
        result['last_consume_timestamp'] = string.format("%18.0f", last_consume_timestamp)
        result['remain_token_count'] = remain_token_count
        result['is_limited'] = false

        return cjson.encode(result)
    end

    -- 桶存在时，重新计算填充 token
    -- 获取 redis 中保存的上次填充时间和剩余 token 数
    local array = redis.call("HMGET", p_key, "last_consume_timestamp", "remain_token_count")
    if array == nil then
        redis.log(redis.LOG_WARNING, "ratelimiter: HMGET return nil for key:" .. p_key)
        redis.log(redis.LOG_DEBUG, "------------ ratelimiter script end ------------")

        -- 保存 result 信息
        result['msg'] = "err:HMGET data return nil"
        result['is_limited'] = false

        return cjson.encode(result)
    end
    local last_consume_timestamp, remain_token_count = tonumber(array[1]), tonumber(array[2])
    redis.log(redis.LOG_DEBUG, "ratelimiter: last_consume_timestamp:" .. last_consume_timestamp .. ", remain_token_count:" .. remain_token_count)

    -- 计算当前时间距离上次填充 token 过了多少微秒
    local now_timestamp_array = redis.call("TIME")
    local now_timestamp = tonumber(now_timestamp_array[1]) * 1000000 + tonumber(now_timestamp_array[2])
    local duration_microsecond = math.max(now_timestamp - last_consume_timestamp, 0)
    -- 根据配置计算 token 的填充速率: x token/μs
    local fill_rate = p_fill_count / p_interval_microsecond
    redis.log(redis.LOG_DEBUG, "ratelimiter: now_timestamp:" .. now_timestamp .. ", duration_microsecond:" .. duration_microsecond .. ", fill_rate:" .. fill_rate)
    -- 计算在这段时间内产生了多少 token , 浮点数向下取整
    local fill_token_count = math.floor(fill_rate * duration_microsecond)
    -- 计算桶内当前时间应有的 token 总数，总数不超过桶的容量
    local now_token_count = math.min(remain_token_count + fill_token_count, p_bucket_capacity)
    redis.log(redis.LOG_DEBUG, "ratelimiter: fill_token_count:" .. fill_token_count .. ", now_token_count:" .. now_token_count)


    -- 保存 debug 信息
    result['last_consume_timestamp'] = string.format("%18.0f", last_consume_timestamp)
    result['remain_token_count'] = remain_token_count
    result['now_timestamp'] = string.format("%18.0f", now_timestamp)
    result['duration_microsecond'] = string.format("%18.0f", duration_microsecond)
    result['fill_rate'] = string.format("%18.9f", fill_rate)
    result['fill_token_count'] = fill_token_count
    result['now_token_count'] = now_token_count


    -- 无可用 token
    if now_token_count <= 0 then
        -- 更新 redis 中的数据，被限流不消耗 now_token_count
        redis.call("HMSET", p_key, "last_consume_timestamp", last_consume_timestamp, "remain_token_count", now_token_count)
        -- 设置 redis 的过期时间
        redis.call("EXPIRE", p_key, p_expire_second)
        redis.log(redis.LOG_DEBUG, "ratelimiter: call HMSET for updating bucket")
        redis.log(redis.LOG_DEBUG, "------------ ratelimiter script end ------------")
        result['msg'] = "limit"
        result['is_limited'] = true
        return cjson.encode(result)
    end

    -- 更新 redis 中的数据, 消耗一个 token
    redis.call("HMSET", p_key, "last_consume_timestamp", now_timestamp, "remain_token_count", now_token_count - 1)
    -- 设置 redis 的过期时间
    redis.call("EXPIRE", p_key, p_expire_second)
    redis.log(redis.LOG_DEBUG, "ratelimiter: call HMSET for updating bucket")
    redis.log(redis.LOG_DEBUG, "------------ ratelimiter script end ------------")
    result['msg'] = "pass"
    result['is_limited'] = false
    return cjson.encode(result)
`)
