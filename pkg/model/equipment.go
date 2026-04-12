package model

// ItemType 代表物品的类型
type ItemType string

const (
	ItemTypeWeapon          ItemType = "weapon"
	ItemTypeArmor           ItemType = "armor"
	ItemTypePotion          ItemType = "potion"
	ItemTypeRing            ItemType = "ring"
	ItemTypeRod             ItemType = "rod"
	ItemTypeScroll          ItemType = "scroll"
	ItemTypeStaff           ItemType = "staff"
	ItemTypeWand            ItemType = "wand"
	ItemTypeWondrousItem    ItemType = "wondrous_item"
	ItemTypeTreasure        ItemType = "treasure"
	ItemTypeAdventuringGear ItemType = "adventuring_gear"
	ItemTypeAmmunition      ItemType = "ammunition"
	ItemTypeTool            ItemType = "tool"
	ItemTypeMount           ItemType = "mount"
	ItemTypeOther           ItemType = "other"
)

// Rarity 代表物品的稀有度
type Rarity string

const (
	RarityCommon    Rarity = "common"
	RarityUncommon  Rarity = "uncommon"
	RarityRare      Rarity = "rare"
	RarityVeryRare  Rarity = "very_rare"
	RarityLegendary Rarity = "legendary"
	RarityArtifact  Rarity = "artifact"
)

// ArmorType 代表护甲类型
type ArmorType string

const (
	ArmorTypeLight  ArmorType = "轻型护甲" // light 轻型护甲
	ArmorTypeMedium ArmorType = "中型护甲" // medium 中型护甲
	ArmorTypeHeavy  ArmorType = "重型护甲" // heavy 重型护甲
	ArmorTypeShield ArmorType = "盾牌"   // shield 盾牌
)

// DamageType 代表伤害类型（这里先定义，后面damage.go会有更详细的定义）
type DamageType string

const (
	DamageTypeAcid        DamageType = "acid"
	DamageTypeBludgeoning DamageType = "bludgeoning"
	DamageTypeCold        DamageType = "cold"
	DamageTypeFire        DamageType = "fire"
	DamageTypeForce       DamageType = "force"
	DamageTypeLightning   DamageType = "lightning"
	DamageTypeNecrotic    DamageType = "necrotic"
	DamageTypePiercing    DamageType = "piercing"
	DamageTypePoison      DamageType = "poison"
	DamageTypePsychic     DamageType = "psychic"
	DamageTypeRadiant     DamageType = "radiant"
	DamageTypeSlashing    DamageType = "slashing"
	DamageTypeThunder     DamageType = "thunder"
)

// WeaponProperties 代表武器属性
type WeaponProperties struct {
	DamageDice string     `json:"damage_dice"`          // 伤害骰，如 "1d8"
	DamageType DamageType `json:"damage_type"`          // 伤害类型
	WeaponType string     `json:"weapon_type"`          // 武器类型（近战/远程）
	Range      int        `json:"range,omitempty"`      // 射程（远程武器）
	LongRange  int        `json:"long_range,omitempty"` // 远射程
	Light      bool       `json:"light"`                // 轻型武器
	Finesse    bool       `json:"finesse"`              // 灵巧武器
	Heavy      bool       `json:"heavy"`                // 重型武器
	TwoHanded  bool       `json:"two_handed"`           // 双手武器
	Versatile  string     `json:"versatile,omitempty"`  // 多用伤害骰
	Loading    bool       `json:"loading"`              // 装填武器
	Thrown     bool       `json:"thrown"`               // 投掷武器
	Reach      bool       `json:"reach"`                // 触及武器
	Ammunition bool       `json:"ammunition"`           // 弹药武器
	Special    []string   `json:"special,omitempty"`    // 特殊属性
}

// ArmorProperties 代表护甲属性
type ArmorProperties struct {
	BaseAC              int  `json:"base_ac"`              // 基础AC
	MaxDexModifier      *int `json:"max_dex_modifier"`     // 最大敏捷修正（nil表示无限制）
	StealthDisadvantage bool `json:"stealth_disadvantage"` // 隐匿劣势
	StrengthRequirement int  `json:"strength_requirement"` // 力量要求
}

// Item 代表一个物品
type Item struct {
	ID          ID       `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Type        ItemType `json:"type"`
	Rarity      Rarity   `json:"rarity"`
	Weight      float64  `json:"weight"` // 重量（磅）
	Quantity    int      `json:"quantity"`
	Value       int      `json:"value"`                // 价值（铜币）
	Attuned     bool     `json:"attuned"`              // 是否已调音
	Attunement  string   `json:"attunement,omitempty"` // 调音要求

	// 类型特定属性
	WeaponProps *WeaponProperties `json:"weapon_props,omitempty"`
	ArmorProps  *ArmorProperties  `json:"armor_props,omitempty"`

	// 魔法属性
	MagicBonus   int      `json:"magic_bonus"`   // 魔法加值
	MagicEffects []string `json:"magic_effects"` // 魔法效果
	Charges      int      `json:"charges"`       // 充能次数
	MaxCharges   int      `json:"max_charges"`   // 最大充能
	Recharge     string   `json:"recharge"`      // 充能恢复条件

	// 消耗品
	Consumable bool   `json:"consumable"`
	Effect     string `json:"effect,omitempty"` // 使用效果
}

// EquipmentSlot 代表装备槽位
type EquipmentSlot string

const (
	SlotMainHand  EquipmentSlot = "main_hand" // 主手
	SlotOffHand   EquipmentSlot = "off_hand"  // 副手
	SlotHead      EquipmentSlot = "head"      // 头部
	SlotNeck      EquipmentSlot = "neck"      // 颈部
	SlotBody      EquipmentSlot = "body"      // 身体
	SlotChest     EquipmentSlot = "chest"     // 胸部（护甲）
	SlotHands     EquipmentSlot = "hands"     // 手部
	SlotFinger1   EquipmentSlot = "finger_1"  // 戒指1
	SlotFinger2   EquipmentSlot = "finger_2"  // 戒指2
	SlotWaist     EquipmentSlot = "waist"     // 腰部
	SlotFeet      EquipmentSlot = "feet"      // 脚部
	SlotShoulders EquipmentSlot = "shoulders" // 肩部
	SlotBack      EquipmentSlot = "back"      // 背部
)

// Equipment 代表角色的装备配置
type Equipment struct {
	Slots map[EquipmentSlot]*Item `json:"slots"`
}

// NewEquipment 创建新的装备配置
func NewEquipment() *Equipment {
	return &Equipment{
		Slots: make(map[EquipmentSlot]*Item),
	}
}

// Currency 代表货币
type Currency struct {
	Platinum int `json:"platinum"` // 白金币
	Gold     int `json:"gold"`     // 金币
	Electrum int `json:"electrum"` // 银币
	Silver   int `json:"silver"`   // 铜币
	Copper   int `json:"copper"`   // 铜币
}

// TotalInGold 将所有货币转换为金币价值
func (c *Currency) TotalInGold() float64 {
	return float64(c.Platinum)*10 + float64(c.Gold) + float64(c.Electrum)/2 + float64(c.Silver)/10 + float64(c.Copper)/100
}

// Inventory 代表角色的库存
type Inventory struct {
	ID        ID         `json:"id"`
	OwnerID   ID         `json:"owner_id"`
	Items     []*Item    `json:"items"`
	Equipment *Equipment `json:"equipment"`
	Currency  Currency   `json:"currency"`
	MaxWeight float64    `json:"max_weight"` // 最大负重
}

// NewInventory 创建新的库存
func NewInventory(ownerID ID) *Inventory {
	return &Inventory{
		ID:        NewID(),
		OwnerID:   ownerID,
		Items:     make([]*Item, 0),
		Equipment: NewEquipment(),
	}
}
