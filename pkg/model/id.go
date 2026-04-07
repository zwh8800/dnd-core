package model

import (
	"crypto/rand"
	"fmt"
	"sync"
	"time"

	"github.com/oklog/ulid/v2"
)

// ID 是游戏中所有实体的唯一标识符，基于ULID实现
type ID string

// NewID 生成一个新的ULID
func NewID() ID {
	t := time.Now()
	entropy := ulid.Monotonic(rand.Reader, 0)
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return ID(id.String())
}

// ParseID 从字符串解析ID
func ParseID(s string) (ID, error) {
	if s == "" {
		return "", fmt.Errorf("empty ID string")
	}
	_, err := ulid.ParseStrict(s)
	if err != nil {
		return "", fmt.Errorf("invalid ID %q: %w", s, err)
	}
	return ID(s), nil
}

// IsValid 检查ID是否有效
func (id ID) IsValid() bool {
	_, err := ulid.ParseStrict(string(id))
	return err == nil
}

// String 返回ID的字符串表示
func (id ID) String() string {
	return string(id)
}

// IDGenerator 提供可测试的ID生成
type IDGenerator interface {
	Next() ID
}

// DefaultIDGenerator 使用ULID的默认生成器
type DefaultIDGenerator struct {
	mu  sync.Mutex
	clk *ulid.MonotonicEntropy
}

// NewDefaultIDGenerator 创建默认ID生成器
func NewDefaultIDGenerator() *DefaultIDGenerator {
	return &DefaultIDGenerator{
		clk: ulid.Monotonic(rand.Reader, 0),
	}
}

// Next 生成下一个ID
func (g *DefaultIDGenerator) Next() ID {
	g.mu.Lock()
	defer g.mu.Unlock()
	t := time.Now()
	id := ulid.MustNew(ulid.Timestamp(t), g.clk)
	return ID(id.String())
}
