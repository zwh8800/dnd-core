package rules

import (
	"fmt"

	"github.com/zwh8800/dnd-core/pkg/model"
)

// ============================================================================
// 借机攻击 (Opportunity Attack)
// ============================================================================

// OpportunityAttackRequest 借机攻击请求
type OpportunityAttackRequest struct {
	AttackerID model.ID `json:"attacker_id"` // 借机攻击者ID
	TargetID   model.ID `json:"target_id"`   // 触发借机攻击的目标ID
}

// OpportunityAttackResult 借机攻击结果
type OpportunityAttackResult struct {
	Triggered bool   `json:"triggered"` // 是否触发
	CanTake   bool   `json:"can_take"`  // 是否能进行借机攻击
	Message   string `json:"message"`   // 描述消息
}

// CanTakeOpportunityAttack 检查是否能进行借机攻击
// PHB Ch.9: 当一个你能看到的敌对生物离开你的触及范围时,你可以用反应对其进行一次近战攻击
func CanTakeOpportunityAttack(attacker, target any, attackerActor, targetActor *model.Actor) (bool, string) {
	// 检查是否失能
	if attackerActor.IsIncapacitated() {
		return false, "攻击者失能,无法进行借机攻击"
	}

	// 检查反应是否可用(需要外部追踪)
	// 这里假设调用方会检查反应是否已使用

	// 检查目标是否在触及范围内
	inReach := IsWithinReach(attackerActor, targetActor)
	if !inReach {
		return false, "目标不在触及范围内"
	}

	// 检查目标是否正在离开触及范围(需要外部传入移动信息)
	// 这里返回true表示满足基本条件

	return true, "满足借机攻击条件"
}

// IsWithinReach 检查目标是否在触及范围内
// 默认触及范围为5尺(1格),长触及武器为10尺(2格)
func IsWithinReach(attacker, target *model.Actor) bool {
	if attacker.Position == nil || target.Position == nil {
		// 如果没有位置信息,默认为相邻
		return true
	}

	// 计算距离(5-10-5规则简化版)
	distance := CalculateGridDistance(attacker.Position, target.Position)

	// 默认触及5尺
	reach := 5

	// TODO: 检查武器是否有reach属性,如果有则reach=10

	return distance <= reach
}

// CalculateGridDistance 计算网格距离(5-10-5规则)
func CalculateGridDistance(from, to *model.Point) int {
	if from == nil || to == nil {
		return 0
	}

	dx := to.X - from.X
	dy := to.Y - from.Y
	if dx < 0 {
		dx = -dx
	}
	if dy < 0 {
		dy = -dy
	}

	// 5-10-5规则: 相邻格5尺,对角格10尺
	// 简化计算: max(dx, dy) * 5 + min(dx, dy) * 5
	if dx == 0 && dy == 0 {
		return 0
	}

	// 简单实现: 曼哈顿距离 * 5
	return (dx + dy) * 5
}

// ============================================================================
// 双持武器 (Two-Weapon Fighting)
// ============================================================================

// TwoWeaponFightRequest 双持武器攻击请求
type TwoWeaponFightRequest struct {
	AttackerID     model.ID `json:"attacker_id"`      // 攻击者ID
	TargetID       model.ID `json:"target_id"`        // 目标ID
	MainHandWeapon string   `json:"main_hand_weapon"` // 主手武器
	OffHandWeapon  string   `json:"off_hand_weapon"`  // 副手武器
}

// TwoWeaponFightResult 双持武器攻击结果
type TwoWeaponFightResult struct {
	CanUseOffHand bool   `json:"can_use_off_hand"` // 是否能使用副手攻击
	Message       string `json:"message"`          // 描述消息
}

// CanUseTwoWeaponFighting 检查是否能使用双持武器战斗
// PHB Ch.9: 当你使用光型近战武器进行攻击动作时,可以用附赠动作攻击另一把轻型近战武器
func CanUseTwoWeaponFighting(mainHand, offHand *model.WeaponProperties) (bool, string) {
	// 两把武器都必须是轻型
	if !mainHand.Light {
		return false, "主手武器不是轻型武器"
	}
	if !offHand.Light {
		return false, "副手武器不是轻型武器"
	}

	return true, "满足双持武器条件"
}

// CalculateOffHandDamage 计算副手攻击伤害
// PHB Ch.9: 副手攻击的伤害不加上属性调整值(除非是负数)
func CalculateOffHandDamage(damageDice int, abilityMod int) int {
	damage := damageDice
	// 副手攻击不加属性调整值(除非双持战士战斗风格或特定专长)
	// 但如果属性调整值是负数,仍然要减去
	if abilityMod < 0 {
		damage += abilityMod
	}
	if damage < 0 {
		damage = 0
	}
	return damage
}

// ============================================================================
// 擒抱 (Grapple)
// ============================================================================

// GrappleRequest 擒抱请求
type GrappleRequest struct {
	GrapplerID model.ID `json:"grappler_id"` // 擒抱者ID
	TargetID   model.ID `json:"target_id"`   // 目标ID
}

// GrappleResult 擒抱结果
type GrappleResult struct {
	Success       bool   `json:"success"`        // 是否成功
	GrapplerRoll  int    `json:"grappler_roll"`  // 擒抱者掷骰
	TargetRoll    int    `json:"target_roll"`    // 目标掷骰
	GrapplerTotal int    `json:"grappler_total"` // 擒抱者总值
	TargetTotal   int    `json:"target_total"`   // 目标总值
	EscapeDC      int    `json:"escape_dc"`      // 逃脱DC
	Message       string `json:"message"`        // 描述消息
}

// PerformGrapple 执行擒抱
// PHB Ch.9: 使用力量(运动)检定对抗目标的力量(运动)或敏捷(体操)检定
func PerformGrapple(grapplerLevel int, grapplerStr int, targetStr int, targetDex int) *GrappleResult {
	// 擒抱者: 力量(运动)检定
	grappleBonus := AbilityModifier(grapplerStr)
	// 假设擒抱者有运动技能熟练(实际应检查)
	// 简化实现: 只加属性修正
	grappleRoll := RollD20()
	grappleTotal := grappleRoll + grappleBonus

	// 目标可以选择力量(运动)或敏捷(体操),取较高值
	strMod := AbilityModifier(targetStr)
	dexMod := AbilityModifier(targetDex)

	targetBonus := strMod
	if dexMod > strMod {
		targetBonus = dexMod
	}
	targetRoll := RollD20()
	targetTotal := targetRoll + targetBonus

	// 对抗检定: 擒抱者 >= 目标则成功
	success := grappleTotal >= targetTotal

	// 逃脱DC = 擒抱者的力量(运动)总值
	escapeDC := grappleTotal

	message := fmt.Sprintf("擒抱检定: %d vs %d", grappleTotal, targetTotal)
	if success {
		message += " - 擒抱成功"
	} else {
		message += " - 擒抱失败"
	}

	return &GrappleResult{
		Success:       success,
		GrapplerRoll:  grappleRoll,
		TargetRoll:    targetRoll,
		GrapplerTotal: grappleTotal,
		TargetTotal:   targetTotal,
		EscapeDC:      escapeDC,
		Message:       message,
	}
}

// CanGrapple 检查是否能擒抱目标
// PHB Ch.9: 目标体型不能超过你一级,且你必须有一只空手
func CanGrapple(grapplerSize, targetSize model.Size) (bool, string) {
	// 体型比较
	sizeLevels := map[model.Size]int{
		model.SizeTiny:       0,
		model.SizeSmall:      1,
		model.SizeMedium:     2,
		model.SizeLarge:      3,
		model.SizeHuge:       4,
		model.SizeGargantuan: 5,
	}

	grapplerLevel := sizeLevels[grapplerSize]
	targetLevel := sizeLevels[targetSize]

	if targetLevel > grapplerLevel+1 {
		return false, "目标体型太大,无法擒抱"
	}

	return true, "满足擒抱条件"
}

// EscapeGrapple 尝试逃脱擒抱
// PHB Ch.9: 被擒抱生物可以用动作进行力量(运动)或敏捷(体操)检定对抗擒抱者的逃脱DC
func EscapeGrapple(escapeeStr int, escapeeDex int, escapeDC int) (bool, int, string) {
	// 选择较高的属性
	strMod := AbilityModifier(escapeeStr)
	dexMod := AbilityModifier(escapeeDex)

	bonus := strMod
	if dexMod > strMod {
		bonus = dexMod
	}

	roll := RollD20()
	total := roll + bonus

	success := total >= escapeDC

	message := fmt.Sprintf("逃脱检定: %d vs DC %d", total, escapeDC)
	if success {
		message += " - 逃脱成功"
	} else {
		message += " - 逃脱失败"
	}

	return success, total, message
}

// ============================================================================
// 推撞 (Shove)
// ============================================================================

// ShoveResult 推撞结果
type ShoveResult struct {
	Success     bool   `json:"success"`      // 是否成功
	ShoverRoll  int    `json:"shover_roll"`  // 推撞者掷骰
	TargetRoll  int    `json:"target_roll"`  // 目标掷骰
	ShoverTotal int    `json:"shover_total"` // 推撞者总值
	TargetTotal int    `json:"target_total"` // 目标总值
	Effect      string `json:"effect"`       // 效果: knocked_prone 或 pushed_away
	Message     string `json:"message"`      // 描述消息
}

// PerformShove 执行推撞
// PHB Ch.9: 使用力量(运动)检定对抗目标的力量(运动)或敏捷(体操)检定
// 成功则目标倒地或被推开5尺
func PerformShove(shoverLevel int, shoverStr int, targetStr int, targetDex int, knockProne bool) *ShoveResult {
	// 推撞者: 力量(运动)检定
	shoveBonus := AbilityModifier(shoverStr)
	shoveRoll := RollD20()
	shoveTotal := shoveRoll + shoveBonus

	// 目标可以选择力量(运动)或敏捷(体操),取较高值
	strMod := AbilityModifier(targetStr)
	dexMod := AbilityModifier(targetDex)

	targetBonus := strMod
	if dexMod > strMod {
		targetBonus = dexMod
	}
	targetRoll := RollD20()
	targetTotal := targetRoll + targetBonus

	// 对抗检定
	success := shoveTotal >= targetTotal

	effect := ""
	if success {
		if knockProne {
			effect = "knocked_prone"
		} else {
			effect = "pushed_away"
		}
	}

	message := fmt.Sprintf("推撞检定: %d vs %d", shoveTotal, targetTotal)
	if success {
		if knockProne {
			message += " - 目标倒地"
		} else {
			message += " - 目标被推开5尺"
		}
	} else {
		message += " - 推撞失败"
	}

	return &ShoveResult{
		Success:     success,
		ShoverRoll:  shoveRoll,
		TargetRoll:  targetRoll,
		ShoverTotal: shoveTotal,
		TargetTotal: targetTotal,
		Effect:      effect,
		Message:     message,
	}
}

// CanShove 检查是否能推撞目标
func CanShove(shoverSize, targetSize model.Size) (bool, string) {
	// 推撞的目标体型不能超过你一级
	sizeLevels := map[model.Size]int{
		model.SizeTiny:       0,
		model.SizeSmall:      1,
		model.SizeMedium:     2,
		model.SizeLarge:      3,
		model.SizeHuge:       4,
		model.SizeGargantuan: 5,
	}

	shoverLevel := sizeLevels[shoverSize]
	targetLevel := sizeLevels[targetSize]

	if targetLevel > shoverLevel+1 {
		return false, "目标体型太大,无法推撞"
	}

	return true, "满足推撞条件"
}

// ============================================================================
// 掩护系统 (Cover)
// ============================================================================

// CalculateACWithCover 计算包含掩护的AC
// PHB Ch.9: 半身掩护+2 AC,四分之三掩护+5 AC
func CalculateACWithCover(baseAC int, cover model.CoverType) int {
	bonus := cover.GetCoverBonus()
	return baseAC + bonus.ACBonus
}

// CalculateDexSaveWithCover 计算包含掩护的敏捷豁免
func CalculateDexSaveWithCover(baseSave int, cover model.CoverType) int {
	bonus := cover.GetCoverBonus()
	return baseSave + bonus.DexSaveBonus
}

// CanTargetWithCover 检查是否能攻击有掩护的目标
// PHB Ch.9: 全身掩护的目标无法被直接攻击
func CanTargetWithCover(cover model.CoverType) (bool, string) {
	if cover == model.CoverTotal {
		return false, "目标有全身掩护,无法被直接攻击"
	}
	return true, "目标可以被攻击"
}

// ============================================================================
// 跳跃系统 (Jumping)
// ============================================================================

// JumpRequest 跳跃请求
type JumpRequest struct {
	JumperID        model.ID       `json:"jumper_id"`         // 跳跃者ID
	Type            model.JumpType `json:"type"`              // 跳跃类型
	HasRunningStart bool           `json:"has_running_start"` // 是否有助跑
}

// JumpResult 跳跃结果
type JumpResult struct {
	Type            model.JumpType `json:"type"`              // 跳跃类型
	Distance        int            `json:"distance"`          // 跳跃距离(尺)
	HasRunningStart bool           `json:"has_running_start"` // 是否有助跑
	Message         string         `json:"message"`           // 描述消息
}

// CalculateJumpDistance 计算跳跃距离
// PHB Ch.8:
// 跳远: 有助跑时=力量值尺数,立定跳远=力量值一半
// 跳高: 有助跑时=3+力量调整值尺数,立定跳高=一半
func CalculateJumpDistance(strength int, jumpType model.JumpType, hasRunningStart bool) int {
	strMod := AbilityModifier(strength)

	var distance int
	switch jumpType {
	case model.JumpTypeLong:
		if hasRunningStart {
			distance = strength
		} else {
			distance = strength / 2
		}
	case model.JumpTypeHigh:
		baseJump := 3 + strMod
		if hasRunningStart {
			distance = baseJump
		} else {
			distance = baseJump / 2
		}
	}

	if distance < 0 {
		distance = 0
	}

	return distance
}

// ============================================================================
// 坠落伤害 (Falling)
// ============================================================================

// CalculateFallDamage 计算坠落伤害
// PHB Ch.8: 每10尺1d6伤害,最多20d6
func CalculateFallDamage(distance int) (damage int, diceCount int, maxDamage int) {
	// 每10尺1d6
	diceCount = distance / 10

	// 最多20d6
	if diceCount > 20 {
		diceCount = 20
	}

	if diceCount == 0 {
		return 0, 0, 0
	}

	// 平均伤害(用于简化)
	// 1d6的平均值是3.5,向上取整为4
	damage = diceCount * 4
	maxDamage = 20 * 6

	return damage, diceCount, maxDamage
}

// ============================================================================
// 窒息规则 (Suffocation)
// ============================================================================

// SuffocationResult 窒息结果
type SuffocationResult struct {
	CanHoldBreath          bool   `json:"can_hold_breath"`          // 是否能继续闭气
	SecondsRemaining       int    `json:"seconds_remaining"`        // 剩余闭气时间(秒)
	RoundsUntilUnconscious int    `json:"rounds_until_unconscious"` // 失去意识前的轮数
	Message                string `json:"message"`                  // 描述消息
}

// CalculateHoldBreathTime 计算闭气时间
// PHB Ch.8: 生物可以闭气的时间 = 1 + 体质调整值(分钟),最少30秒
func CalculateHoldBreathTime(constitution int) int {
	conMod := AbilityModifier(constitution)
	minutes := 1 + conMod
	if minutes < 1 {
		minutes = 1 // 最少30秒
	}
	return minutes * 60 // 转换为秒
}

// CalculateSuffocationRounds 计算窒息轮数
// PHB Ch.8: 窒息后还能存活的轮数 = 体质调整值(最少1轮)
func CalculateSuffocationRounds(constitution int) int {
	conMod := AbilityModifier(constitution)
	rounds := conMod
	if rounds < 1 {
		rounds = 1
	}
	return rounds
}

// ============================================================================
// 专注豁免机制 (Concentration Saves)
// ============================================================================

// CalculateConcentrationDC 计算专注豁免DC
// PHB Ch.10: DC = 10 或 伤害的一半,取较高值
func CalculateConcentrationDC(damage int) int {
	halfDamage := damage / 2
	dc := 10
	if halfDamage > dc {
		dc = halfDamage
	}
	return dc
}

// PerformConcentrationSave 执行专注豁免
// PHB Ch.10: 受到伤害时必须进行体质豁免,DC = 10或伤害/2取高
func PerformConcentrationSave(casterLevel int, constitution int, damage int) (int, bool, string) {
	dc := CalculateConcentrationDC(damage)
	conMod := AbilityModifier(constitution)

	// 假设 caster 有体质豁免熟练(实际应根据职业判断)
	profBonus := ProficiencyBonus(casterLevel)
	bonus := conMod + profBonus

	roll := RollD20()
	total := roll + bonus

	success := total >= dc

	message := fmt.Sprintf("专注豁免: %d vs DC %d", total, dc)
	if success {
		message += " - 专注维持"
	} else {
		message += " - 专注打断"
	}

	return dc, success, message
}

// ============================================================================
// 团体检定 (Group Checks)
// ============================================================================

// PerformGroupCheck 执行团体检定
// PHB Ch.7: 当多个生物尝试一起完成某件事时,如果至少一半成员成功,则团体成功
func PerformGroupCheck(participants []GroupCheckParticipant, dc int) *GroupCheckResult {
	successCount := 0
	failCount := 0

	for i := range participants {
		participants[i].DC = dc
		if participants[i].Total >= dc {
			participants[i].Success = true
			successCount++
		} else {
			participants[i].Success = false
			failCount++
		}
	}

	// 至少一半成功
	overallSuccess := successCount >= (len(participants)+1)/2

	message := fmt.Sprintf("团体检定: %d 成功, %d 失败", successCount, failCount)
	if overallSuccess {
		message += " - 团体成功"
	} else {
		message += " - 团体失败"
	}

	return &GroupCheckResult{
		Participants:   participants,
		SuccessCount:   successCount,
		FailCount:      failCount,
		OverallSuccess: overallSuccess,
		Message:        message,
	}
}

// GroupCheckParticipant 团体检定参与者
type GroupCheckParticipant struct {
	ActorID   model.ID `json:"actor_id"`   // 角色ID
	ActorName string   `json:"actor_name"` // 角色名称
	Roll      int      `json:"roll"`       // 掷骰结果
	Total     int      `json:"total"`      // 总值
	DC        int      `json:"dc"`         // 难度等级
	Success   bool     `json:"success"`    // 是否成功
}

// GroupCheckResult 团体检定结果
type GroupCheckResult struct {
	Participants   []GroupCheckParticipant `json:"participants"`    // 参与者
	SuccessCount   int                     `json:"success_count"`   // 成功数量
	FailCount      int                     `json:"fail_count"`      // 失败数量
	OverallSuccess bool                    `json:"overall_success"` // 整体是否成功
	Message        string                  `json:"message"`         // 描述消息
}

// ============================================================================
// 合作检定 (Working Together)
// ============================================================================

// PerformWorkingTogether 执行合作检定
// PHB Ch.7: 当多个生物合作完成一件事时,领导者进行检定并具有优势
func PerformWorkingTogether(leaderName string, leaderBonus int, helpers []string, dc int) *WorkingTogetherResult {
	// 领导者具有优势(掷两次d20取高)
	roll1 := RollD20()
	roll2 := RollD20()

	roll := roll1
	if roll2 > roll1 {
		roll = roll2
	}

	total := roll + leaderBonus
	success := total >= dc

	message := fmt.Sprintf("%s 合作检定: %d (优势取高: %d vs %d)", leaderName, total, roll1, roll2)
	if success {
		message += " - 成功"
	} else {
		message += " - 失败"
	}

	return &WorkingTogetherResult{
		LeaderName:   leaderName,
		LeaderRoll:   roll,
		LeaderTotal:  total,
		HasAdvantage: true,
		Helpers:      helpers,
		DC:           dc,
		Success:      success,
		Message:      message,
	}
}

// WorkingTogetherResult 合作检定结果
type WorkingTogetherResult struct {
	LeaderName   string   `json:"leader_name"`   // 领导者名称
	LeaderRoll   int      `json:"leader_roll"`   // 领导者掷骰
	LeaderTotal  int      `json:"leader_total"`  // 领导者总值
	HasAdvantage bool     `json:"has_advantage"` // 是否具有优势
	Helpers      []string `json:"helpers"`       // 协助者列表
	DC           int      `json:"dc"`            // 难度等级
	Success      bool     `json:"success"`       // 是否成功
	Message      string   `json:"message"`       // 描述消息
}

// ============================================================================
// 暴击规则 (Critical Hits)
// ============================================================================

// CalculateCriticalDamage 计算暴击伤害
// PHB Ch.9: 暴击时,将所有伤害骰掷两次并相加,然后加上修正值
func CalculateCriticalDamage(baseDamageDice int, modifier int) int {
	// 暴击: 伤害骰掷两次(即基础伤害骰 * 2)
	// 修正值只加一次
	criticalDamage := baseDamageDice*2 + modifier
	return criticalDamage
}

// IsCriticalHit 检查是否是暴击(自然20)
func IsCriticalHit(attackRoll int) bool {
	return attackRoll == 20
}

// IsCriticalFumble 检查是否是大失败(自然1)
func IsCriticalFumble(attackRoll int) bool {
	return attackRoll == 1
}

// ============================================================================
// 辅助函数
// ============================================================================

// RollD20 掷d20
func RollD20() int {
	// 简单实现,实际应使用roller
	return 10 // 占位符,实际应由外部传入
}
