package rules

import (
	"fmt"

	"github.com/zwh8800/dnd-core/pkg/data"
	"github.com/zwh8800/dnd-core/pkg/model"
)

// ValidateMulticlass 验证多职业选择是否合法（13分规则）
func ValidateMulticlass(pc *model.PlayerCharacter, newClass model.ClassID) error {
	// 获取新职业的定义
	classDef, exists := getClassDefinition(newClass)
	if !exists {
		return fmt.Errorf("class definition not found: %s", newClass)
	}

	// SRD 5.2.1 多职业要求：新职业的主要属性至少 13 分
	if len(classDef.PrimaryAbilities) > 0 {
		primaryAbility := classDef.PrimaryAbilities[0]
		if pc.AbilityScores.Get(primaryAbility) < 13 {
			return fmt.Errorf("multiclass requirement not met: %s requires at least 13 in %s", newClass, primaryAbility)
		}
	}

	// 如果角色已有职业，也需要满足原职业的多职业要求
	if len(pc.Classes) > 0 {
		for _, classLevel := range pc.Classes {
			oldClassDef, exists := getClassDefinition(classLevel.Class)
			if !exists {
				continue
			}

			if len(oldClassDef.PrimaryAbilities) > 0 {
				primaryAbility := oldClassDef.PrimaryAbilities[0]
				if pc.AbilityScores.Get(primaryAbility) < 13 {
					return fmt.Errorf("existing class requirement not met: %s requires at least 13 in %s", classLevel.Class, primaryAbility)
				}
			}
		}
	}

	return nil
}

// GetMulticlassSpellSlots 计算多职业施法者的法术位表
func GetMulticlassSpellSlots(classes []model.ClassLevel) [][]int {
	// 计算有效施法者等级
	fullCasterLevels := 0
	halfCasterLevels := 0
	thirdCasterLevels := 0

	for _, classLevel := range classes {
		classDef, exists := getClassDefinition(classLevel.Class)
		if !exists {
			continue
		}

		switch classDef.CasterType {
		case model.CasterTypeFull:
			fullCasterLevels += classLevel.Level
		case model.CasterTypeHalf:
			halfCasterLevels += classLevel.Level
		case model.CasterTypeThird:
			thirdCasterLevels += classLevel.Level
		}
	}

	// 有效施法者等级 = 全施法者等级 + 半施法者等级/2 + 1/3施法者等级/3
	effectiveCasterLevel := fullCasterLevels + halfCasterLevels/2 + thirdCasterLevels/3

	if effectiveCasterLevel == 0 {
		return nil
	}

	// 返回对应等级的法术位表
	return GetSpellSlotTable(effectiveCasterLevel)
}

// GetMulticlassProficiencies 计算多职业获得的熟练项
func GetMulticlassProficiencies(oldClasses []model.ClassLevel, newClass model.ClassID) model.Proficiencies {
	prof := model.Proficiencies{}

	newClassDef, exists := getClassDefinition(newClass)
	if !exists {
		return prof
	}

	// 多职业时只获得部分熟练项（根据 SRD 5.2.1 规则）
	// 这里简化处理，实际应根据具体职业给予不同熟练项
	if len(newClassDef.ArmorProficiencies) > 0 {
		// 只有某些职业在多职业时给予护甲熟练
		// 简化：给予轻甲熟练
		if prof.ArmorProficiencies == nil {
			prof.ArmorProficiencies = make(map[model.ArmorType]bool)
		}
		prof.ArmorProficiencies[model.ArmorTypeLight] = true
	}

	return prof
}

// ValidateExtraAttack 处理 Extra Attack 不叠加规则
// SRD 5.2.1: 如果多个职业都给予 Extra Attack，它们不叠加
func ValidateExtraAttack(classes []model.ClassLevel) int {
	maxExtraAttacks := 0

	for _, classLevel := range classes {
		_, exists := getClassDefinition(classLevel.Class)
		if !exists {
			continue
		}

		// 检查该职业在当前等级是否给予 Extra Attack
		extraAttacks := getExtraAttacksForClass(classLevel.Class, classLevel.Level)
		if extraAttacks > maxExtraAttacks {
			maxExtraAttacks = extraAttacks
		}
	}

	return maxExtraAttacks
}

// ValidateUnarmoredDefense 处理无甲防御不叠加规则
// SRD 5.2.1: 如果具有多个无甲防御特性，只能选择其中一个生效
func ValidateUnarmoredDefense(classes []model.ClassLevel) bool {
	// 检查是否具有多个无甲防御来源
	hasUnarmoredDefense := false

	for _, classLevel := range classes {
		_, exists := getClassDefinition(classLevel.Class)
		if !exists {
			continue
		}

		// 某些职业在特定等级获得无甲防御
		if hasUnarmoredDefenseFeature(classLevel.Class, classLevel.Level) {
			if hasUnarmoredDefense {
				// 已有一个无甲防御，新的不叠加
				return false
			}
			hasUnarmoredDefense = true
		}
	}

	return hasUnarmoredDefense
}

// ValidateLevelUpChoice 验证升级选择是否合法
func ValidateLevelUpChoice(pc *model.PlayerCharacter, newClass model.ClassID) error {
	// 如果是新职业（多职业），验证多职业要求
	isNewClass := true
	for _, classLevel := range pc.Classes {
		if classLevel.Class == newClass {
			isNewClass = false
			break
		}
	}

	if isNewClass {
		return ValidateMulticlass(pc, newClass)
	}

	return nil
}

// 辅助函数

func getExtraAttacksForClass(classID model.ClassID, level int) int {
	// 战士: 5级1次, 11级2次, 20级3次
	if classID == model.ClassFighter {
		if level >= 20 {
			return 3
		} else if level >= 11 {
			return 2
		} else if level >= 5 {
			return 1
		}
	}

	// 野蛮人: 5级1次额外攻击(通过额外攻击特性)
	if classID == model.ClassBarbarian {
		if level >= 5 {
			return 1
		}
	}

	// 游侠: 5级1次额外攻击
	if classID == model.ClassRanger {
		if level >= 5 {
			return 1
		}
	}

	// 武僧: 5级通过武术获得额外攻击
	if classID == model.ClassMonk {
		if level >= 5 {
			return 1
		}
	}

	return 0
}

func hasUnarmoredDefenseFeature(classID model.ClassID, level int) bool {
	// 蛮族和武僧在 1 级获得无甲防御
	if classID == model.ClassBarbarian || classID == model.ClassMonk {
		return level >= 1
	}
	return false
}

// GetSpellSlotTable 获取指定施法者等级的法术位表
func GetSpellSlotTable(casterLevel int) [][]int {
	// SRD 5.2.1 标准法术位表 [环级][总数, 已用]
	// 索引 0 = 戏法，索引 1-9 = 1-9 环
	spellSlotTable := [][]int{
		{}, // 0 级（戏法无限）
		{2, 0},
		{3, 0},
		{4, 2, 0},
		{4, 3, 0},
		{4, 3, 2, 0},
		{4, 3, 3, 0},
		{4, 3, 3, 1, 0},
		{4, 3, 3, 2, 0},
		{4, 3, 3, 3, 1, 0},
		{4, 3, 3, 3, 2, 0},
		{4, 3, 3, 3, 2, 1, 0},
		{4, 3, 3, 3, 2, 2, 0},
		{4, 3, 3, 3, 2, 2, 1, 0},
		{4, 3, 3, 3, 2, 2, 2, 0},
		{4, 3, 3, 3, 2, 2, 2, 1, 0},
		{4, 3, 3, 3, 2, 2, 2, 2, 0},
		{4, 3, 3, 3, 2, 2, 2, 2, 1, 0},
		{4, 3, 3, 3, 3, 2, 2, 2, 2, 0},
		{4, 3, 3, 3, 3, 3, 2, 2, 2, 0},
		{4, 3, 3, 3, 3, 3, 3, 2, 2, 0},
	}

	if casterLevel >= 1 && casterLevel <= 20 {
		return spellSlotTable[1 : casterLevel+1]
	}

	return nil
}

func getClassDefinition(classID model.ClassID) (*data.ClassDefinition, bool) {
	return data.GlobalRegistry.GetClass(classID)
}
