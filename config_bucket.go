package ratelimiter

var (
	// 默认最大允许每秒请求 1000 次
	defaultBucketCapacity = 1000
	// 默认每 1000 微秒填充一个 token ，即每毫秒填充一个 token
	// 即每秒填充 1000 个 token
	defaultBucketFillEveryMicrosecond = 1000
	// 默认令牌桶的过期时间
	defaultBucketExpireSecond = 60 * 60 // 1h
)

// BucketConfig 令牌桶 bucket 配置信息
type BucketConfig struct {
	// Capacity 重要参数 令牌桶容量大小 每秒最大请求量
	Capacity int
	// FillEveryMicrosecond 重要参数 每隔多少微秒往 bucket 中放一个 token
	FillEveryMicrosecond int
	// ExpireSecond 令牌桶过期时间（单位：秒）
	ExpireSecond int
}
