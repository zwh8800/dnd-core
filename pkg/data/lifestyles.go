package data

import (
	"github.com/zwh8800/dnd-core/pkg/model"
)

// LifestyleData 生活方式数据
type LifestyleData struct {
	Tier        model.LifestyleTier `json:"tier"`
	DailyCost   int                 `json:"daily_cost"`   // 每日花费（铜币）
	MonthlyCost int                 `json:"monthly_cost"` // 每月花费（铜币）
	Description string              `json:"description"`
}

// LifestyleDataList 所有生活方式数据
var LifestyleDataList = []LifestyleData{
	{
		Tier:        model.LifestyleWretched,
		DailyCost:   0,
		MonthlyCost: 0,
		Description: "你生活在悲惨 conditions 中，没有稳定的食物或住所",
	},
	{
		Tier:        model.LifestyleSqualid,
		DailyCost:   10,  // 1 sp
		MonthlyCost: 300, // 3 gp
		Description: "你住在肮脏的环境中，食物质量差",
	},
	{
		Tier:        model.LifestylePoor,
		DailyCost:   20,  // 2 sp
		MonthlyCost: 600, // 6 gp
		Description: "你能够勉强维持生活，但经常需要工作",
	},
	{
		Tier:        model.LifestyleModest,
		DailyCost:   100,  // 1 gp
		MonthlyCost: 3000, // 30 gp
		Description: "你生活在普通社区，有足够的食物和住所",
	},
	{
		Tier:        model.LifestyleComfortable,
		DailyCost:   200,  // 2 gp
		MonthlyCost: 6000, // 60 gp
		Description: "你生活在舒适的家中，有足够的食物和仆人",
	},
	{
		Tier:        model.LifestyleWealthy,
		DailyCost:   400,   // 4 gp
		MonthlyCost: 12000, // 120 gp
		Description: "你生活在奢华的环境中，有充足的仆人",
	},
	{
		Tier:        model.LifestyleAristocratic,
		DailyCost:   1000,  // 10 gp
		MonthlyCost: 30000, // 300 gp
		Description: "你生活在极度奢华的环境中，有无数的仆人",
	},
}

// GetLifestyleData 获取指定生活方式的数据
func GetLifestyleData(tier model.LifestyleTier) *LifestyleData {
	for _, data := range LifestyleDataList {
		if data.Tier == tier {
			return &data
		}
	}
	return nil
}
