package model

// MountState 代表骑乘状态
type MountState string

const (
	MountStateUnmounted MountState = "unmounted" // 未骑乘
	MountStateMounted   MountState = "mounted"   // 已骑乘
)

// MountTack 代表鞍具类型
type MountTack string

const (
	TackNone           MountTack = "none"            // 无鞍具
	TackRidingSaddle   MountTack = "riding_saddle"   // 骑乘鞍
	TackMilitarySaddle MountTack = "military_saddle" // 军用鞍
	TackPackSaddle     MountTack = "pack_saddle"     // 驮鞍
	TackExoticSaddle   MountTack = "exotic_saddle"   // 异种鞍
)

// Mount 代表坐骑
type Mount struct {
	ActorID   ID         `json:"actor_id"`           // 生物 ID
	State     MountState `json:"state"`              // 骑乘状态
	Tack      MountTack  `json:"tack"`               // 鞍具类型
	IsTrained bool       `json:"is_trained"`         // 是否受过骑乘训练
	RiderID   ID         `json:"rider_id,omitempty"` // 骑手 ID
}

// VehicleType 代表交通工具类型
type VehicleType string

const (
	VehicleLand  VehicleType = "land"  // 陆地
	VehicleWater VehicleType = "water" // 水上
)

// Vehicle 代表交通工具
type Vehicle struct {
	ID            ID          `json:"id"`
	Name          string      `json:"name"`
	Description   string      `json:"description"`
	Type          VehicleType `json:"type"`
	AC            int         `json:"ac"`             // 护甲等级
	HP            int         `json:"hp"`             // 生命值
	Speed         int         `json:"speed"`          // 速度（尺/小时）
	CargoCapacity float64     `json:"cargo_capacity"` // 载货量（磅）
	CrewRequired  int         `json:"crew_required"`  // 所需船员/船员
	CrewCurrent   int         `json:"crew_current"`   // 当前船员
	Weight        float64     `json:"weight"`         // 重量（磅）
	Value         int         `json:"value"`          // 价值（铜币）
}

// MountAction 代表骑乘动作
type MountAction string

const (
	MountActionMount    MountAction = "mount"    // 骑上
	MountActionDismount MountAction = "dismount" // 下来
	MountActionControl  MountAction = "control"  // 控制
	MountActionDodging  MountAction = "dodging"  // 闪避
)
