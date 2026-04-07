package model

// SceneType 代表场景类型
type SceneType string

const (
	SceneTypeIndoor     SceneType = "indoor"     // 室内
	SceneTypeOutdoor    SceneType = "outdoor"    // 室外
	SceneTypeWilderness SceneType = "wilderness" // 荒野
	SceneTypeDungeon    SceneType = "dungeon"    // 地牢
	SceneTypeCity       SceneType = "city"       // 城市
	SceneTypeTavern     SceneType = "tavern"     // 酒馆
	SceneTypeShop       SceneType = "shop"       // 商店
	SceneTypeTemple     SceneType = "temple"     // 寺庙
	SceneTypeOther      SceneType = "other"      // 其他
)

// SceneConnection 代表场景之间的连接
type SceneConnection struct {
	TargetSceneID ID     `json:"target_scene_id"` // 目标场景ID
	Description   string `json:"description"`     // 连接描述（如"北边的门"）
	Locked        bool   `json:"locked"`          // 是否上锁
	DC            int    `json:"dc,omitempty"`    // 开锁DC
	Hidden        bool   `json:"hidden"`          // 是否隐藏（需要察觉才能发现）
}

// SceneRegion 代表场景中的一个区域
type SceneRegion struct {
	ID          ID     `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Bounds      *Rect  `json:"bounds,omitempty"` // 区域边界（可选）
}

// Rect 代表矩形区域
type Rect struct {
	X      int `json:"x"`
	Y      int `json:"y"`
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Scene 代表一个场景/地点
type Scene struct {
	ID          ID        `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        SceneType `json:"type"`
	Details     string    `json:"details"` // 详细描述

	// 连接
	Connections map[ID]*SceneConnection `json:"connections"`

	// 区域（可选）
	Regions map[ID]*SceneRegion `json:"regions,omitempty"`

	// 场景中的物品（在地面上的）
	Items []ID `json:"items"`

	// 场景特性
	IsDark     bool   `json:"is_dark"`     // 是否黑暗
	LightLevel string `json:"light_level"` // 光照等级
	Weather    string `json:"weather"`     // 天气（室外场景）
	Terrain    string `json:"terrain"`     // 地形

	// 自定义数据
	CustomData map[string]any `json:"custom_data,omitempty"`
}

// NewScene 创建新场景
func NewScene(name, description string, sceneType SceneType) *Scene {
	return &Scene{
		ID:          NewID(),
		Name:        name,
		Description: description,
		Type:        sceneType,
		Connections: make(map[ID]*SceneConnection),
		Regions:     make(map[ID]*SceneRegion),
		Items:       make([]ID, 0),
		CustomData:  make(map[string]any),
	}
}

// AddConnection 添加场景连接
func (s *Scene) AddConnection(targetSceneID ID, description string, locked bool) {
	s.Connections[targetSceneID] = &SceneConnection{
		TargetSceneID: targetSceneID,
		Description:   description,
		Locked:        locked,
	}
}

// RemoveConnection 移除场景连接
func (s *Scene) RemoveConnection(targetSceneID ID) {
	delete(s.Connections, targetSceneID)
}

// AddItem 添加物品到场景
func (s *Scene) AddItem(itemID ID) {
	s.Items = append(s.Items, itemID)
}

// RemoveItem 从场景移除物品
func (s *Scene) RemoveItem(itemID ID) {
	for i, id := range s.Items {
		if id == itemID {
			s.Items = append(s.Items[:i], s.Items[i+1:]...)
			return
		}
	}
}
