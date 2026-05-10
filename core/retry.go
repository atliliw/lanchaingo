package core

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"time"
)

// RetryConfig 控制重试行为
type RetryConfig struct {
	MaxRetries   int           // 最大重试次数（默认 3）
	BaseDelay    time.Duration // 基础延迟（默认 1s）
	MaxDelay     time.Duration // 最大延迟（默认 30s）
	RetryableErr func(error) bool // 判断错误是否可重试
}

// DefaultRetryConfig 返回默认重试配置
// 指数退避: 1s → 2s → 4s → 8s → 16s (抖动 ±25%)
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:   3,
		BaseDelay:    time.Second,
		MaxDelay:     30 * time.Second,
		RetryableErr: DefaultRetryableErr,
	}
}

// DefaultRetryableErr 默认错误判断：所有非 nil 错误都重试
func DefaultRetryableErr(err error) bool {
	return err != nil
}

// DoWithRetry 带指数退避重试执行函数 fn
// fn 返回 (result, error)，error 为 nil 表示成功
// 如果所有重试都失败，返回最后一次错误
func DoWithRetry[T any](ctx context.Context, cfg RetryConfig, fn func(context.Context) (T, error)) (T, error) {
	var lastErr error
	var zero T

	for attempt := 0; attempt <= cfg.MaxRetries; attempt++ {
		// 检查 context 是否已取消
		if ctx.Err() != nil {
			return zero, ctx.Err()
		}

		// 执行
		result, err := fn(ctx)
		if err == nil {
			return result, nil
		}

		lastErr = err

		// 判断是否可重试
		if !cfg.RetryableErr(err) {
			return zero, err
		}

		// 最后一次尝试失败后不再等待
		if attempt >= cfg.MaxRetries {
			break
		}

		// 计算等待时间：baseDelay * 2^attempt + 随机抖动
		delay := float64(cfg.BaseDelay) * math.Pow(2, float64(attempt))
		if delay > float64(cfg.MaxDelay) {
			delay = float64(cfg.MaxDelay)
		}
		// 添加 ±25% 的抖动
		jitter := delay * (0.75 + rand.Float64()*0.5)
		wait := time.Duration(jitter)

		// 等待或 ctx 取消
		select {
		case <-ctx.Done():
			return zero, ctx.Err()
		case <-time.After(wait):
		}
	}

	return zero, fmt.Errorf("retry exhausted after %d attempts: %w", cfg.MaxRetries, lastErr)
}

// RetryableHTTP 判断 HTTP 响应状态码是否可重试
func RetryableHTTP(statusCode int) bool {
	return statusCode == 429 || statusCode >= 500
}
