package engine

import "github.com/zwh8800/dnd-core/pkg/storage"

// Config 包含引擎的配置选项
type Config struct {
	// Storage 指定存储后端，用于游戏状态的持久化
	// 如果为nil，将使用默认的内存存储
	Storage storage.Store

	// DiceSeed 指定骰子随机数生成器的种子
	// 如果为0，将使用系统时间作为种子
	// 设置固定种子可用于测试或可重现的游戏
	DiceSeed int64

	// DataPath 指定自定义数据文件的路径
	// 用于覆盖内置的种族、职业、法术等数据
	// 如果为空，将仅使用内置数据
	DataPath string
}
