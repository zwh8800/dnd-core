package model

// TravelPace 代表旅行步伐
type TravelPace string

const (
	TravelPaceFast   TravelPace = "fast"   // 快速
	TravelPaceNormal TravelPace = "normal" // 正常
	TravelPaceSlow   TravelPace = "slow"   // 慢速
)

// TerrainType 代表地形类型
type TerrainType string

const (
	TerrainClear     TerrainType = "clear"     // 开阔地
	TerrainGrassland TerrainType = "grassland" // 草原
	TerrainForest    TerrainType = "forest"    // 森林
	TerrainMountain  TerrainType = "mountain"  // 山地
	TerrainSwamp     TerrainType = "swamp"     // 沼泽
	TerrainDesert    TerrainType = "desert"    // 沙漠
	TerrainArctic    TerrainType = "arctic"    // 寒带
)

// TravelState 代表旅行状态
type TravelState struct {
	CurrentLocation  string      `json:"current_location"`  // 当前位置
	Destination      string      `json:"destination"`       // 目的地
	Pace             TravelPace  `json:"pace"`              // 旅行步伐
	Terrain          TerrainType `json:"terrain"`           // 当前地形
	DistanceTotal    float64     `json:"distance_total"`    // 总距离（英里）
	DistanceTraveled float64     `json:"distance_traveled"` // 已行距离（英里）
	DaysElapsed      int         `json:"days_elapsed"`      // 已用天数
	HoursToday       int         `json:"hours_today"`       // 今日已行小时数
	IsActive         bool        `json:"is_active"`         // 是否正在旅行
}

// ForageResult 代表觅食结果
type ForageResult struct {
	Success       bool   `json:"success"`        // 是否成功
	FoodObtained  int    `json:"food_obtained"`  // 获得的食物份数
	WaterObtained bool   `json:"water_obtained"` // 是否获得水源
	RollTotal     int    `json:"roll_total"`
	DC            int    `json:"dc"`
	Message       string `json:"message"`
}

// NavigationCheck 代表导航检定
type NavigationCheck struct {
	Success   bool   `json:"success"`
	RollTotal int    `json:"roll_total"`
	DC        int    `json:"dc"`
	Lost      bool   `json:"lost"` // 是否迷路
	Message   string `json:"message"`
}

// EncounterCheck 代表遭遇检定
type EncounterCheck struct {
	Encountered   bool   `json:"encountered"`
	EncounterType string `json:"encounter_type,omitempty"` // "monster", "npc", "treasure", "trap"
	DC            int    `json:"dc"`
	Roll          int    `json:"roll"`
	Message       string `json:"message"`
}
