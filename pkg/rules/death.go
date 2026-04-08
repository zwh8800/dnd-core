package rules

import (
	"fmt"

	"github.com/zwh8800/dnd-core/pkg/dice"
)

// MakeDeathSave 进行死亡豁免检定
func MakeDeathSave() (*DeathSaveResult, error) {
	roller := dice.New(0)
	result, err := roller.Roll("1d20")
	if err != nil {
		return nil, err
	}

	roll := result.Total
	deathResult := &DeathSaveResult{
		Roll: roll,
	}

	// 天然20：立即恢复1HP
	if roll == 20 {
		deathResult.Success = true
		deathResult.IsCritical = true
		return deathResult, nil
	}

	// 天然1：算2次失败
	if roll == 1 {
		deathResult.Success = false
		deathResult.IsCriticalFail = true
		return deathResult, nil
	}

	// DC 10
	if roll >= 10 {
		deathResult.Success = true
	} else {
		deathResult.Success = false
	}

	return deathResult, nil
}

// CheckDeathStatus 检查死亡状态
func CheckDeathStatus(successes, failures int) (isStabilized bool, isDead bool, message string) {
	if failures >= 3 {
		return false, true, "3次失败 - 角色死亡"
	}
	if successes >= 3 {
		return true, false, "3次成功 - 角色稳定"
	}
	return false, false, fmt.Sprintf("成功: %d/3, 失败: %d/3", successes, failures)
}

// HandleDamageAtZeroHP 处理0HP时受到伤害的规则
func HandleDamageAtZeroHP(damage int, maxHP int, currentDeathFails int) (newDeathFails int, isDead bool, message string) {
	// 在0HP时受到任何伤害都算1次失败
	newDeathFails = currentDeathFails + 1
	message = fmt.Sprintf("在0HP时受到伤害，死亡豁免失败+1（现在%d次）", newDeathFails)

	// 如果伤害 >= HP最大值，立即死亡
	if damage >= maxHP {
		return newDeathFails, true, "受到致命伤害（伤害 >= HP最大值），立即死亡"
	}

	return newDeathFails, false, message
}

// StabilizeCreature 稳定生物
func StabilizeCreature() string {
	return "生物已稳定，不再进行死亡豁免，但仍处于0HP和无意识状态"
}

// ApplyHealingAtZeroHP 在0HP时应用治疗
func ApplyHealingAtZeroHP(healing int) string {
	return fmt.Sprintf("恢复%d点HP，生物恢复意识并站立（如果 healed HP > 0）", healing)
}
