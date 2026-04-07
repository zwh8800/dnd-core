package model

// WeaponMasteryType 代表 SRD 5.2.1 的武器掌控类型
type WeaponMasteryType string

const (
	MasterySlow   WeaponMasteryType = "Slow"   // 减缓：命中时目标速度降 10 尺直到你的下个回合开始
	MasteryTopple WeaponMasteryType = "Topple" // 击倒：命中时可尝试将目标击倒
	MasteryPush   WeaponMasteryType = "Push"   // 推击：命中时可将目标推离 5 尺
	MasteryNick   WeaponMasteryType = "Nick"   // 轻捷：使用轻武器时，额外攻击可作为附赠动作
	MasteryVex    WeaponMasteryType = "Vex"    // 烦扰：命中后下次攻击具有优势
	MasteryCleave WeaponMasteryType = "Cleave" // 劈砍：击杀目标时可攻击邻近生物
	MasterySap    WeaponMasteryType = "Sap"    // 钝击：命中时目标下次攻击具有劣势
	MasteryGraze  WeaponMasteryType = "Graze"  // 擦伤：未命中时仍造成属性修正值伤害
)

// WeaponMasteryEffect 代表武器掌控的效果
type WeaponMasteryEffect struct {
	// Type 掌控类型
	Type WeaponMasteryType `json:"type"`
	// Description 效果描述
	Description string `json:"description"`
	// SaveDC 豁免 DC（如果需要）
	SaveDC int `json:"save_dc,omitempty"`
	// SaveAbility 豁免属性（如果需要）
	SaveAbility Ability `json:"save_ability,omitempty"`
}

// GetMasteryEffect 获取指定掌控类型的效果定义
func GetMasteryEffect(mastery WeaponMasteryType) WeaponMasteryEffect {
	switch mastery {
	case MasterySlow:
		return WeaponMasteryEffect{
			Type:        MasterySlow,
			Description: "命中时，目标速度降低 10 尺，直到你的下个回合开始",
		}
	case MasteryTopple:
		return WeaponMasteryEffect{
			Type:        MasteryTopple,
			Description: "命中时，可尝试将目标击倒（目标需进行力量或敏捷豁免）",
		}
	case MasteryPush:
		return WeaponMasteryEffect{
			Type:        MasteryPush,
			Description: "命中时，可将目标推离 5 尺（目标体型需不大于你）",
		}
	case MasteryNick:
		return WeaponMasteryEffect{
			Type:        MasteryNick,
			Description: "使用此武器进行的额外攻击可作为附赠动作执行",
		}
	case MasteryVex:
		return WeaponMasteryEffect{
			Type:        MasteryVex,
			Description: "命中目标后，你对该目标的下次攻击检定具有优势",
		}
	case MasteryCleave:
		return WeaponMasteryEffect{
			Type:        MasteryCleave,
			Description: "击杀一个生物后，可立即对邻近的另一个生物进行一次攻击",
		}
	case MasterySap:
		return WeaponMasteryEffect{
			Type:        MasterySap,
			Description: "命中时，目标的下次攻击检定具有劣势，直到你的下个回合开始",
		}
	case MasteryGraze:
		return WeaponMasteryEffect{
			Type:        MasteryGraze,
			Description: "未命中时，仍造成等于你属性修正值的伤害",
		}
	default:
		return WeaponMasteryEffect{
			Type:        mastery,
			Description: "未知掌控类型",
		}
	}
}
