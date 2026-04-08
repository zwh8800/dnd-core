# pkg/engine 单元测试增强计划

## Context

pkg/engine 目录包含 D&D 5e 游戏引擎的核心代码，涵盖角色管理、检定系统、战斗系统、法术系统等。当前已存在部分测试文件，但许多导出函数的测试覆盖不足（每个函数需要至少 5 个独立测试用例）。本计划旨在补充完整测试覆盖，确保每个导出函数都有充分的单元测试。

## 现有测试分析

### 已存在的测试文件
- `check_test.go` - 2个测试用例 for PerformAbilityCheck, 2个 for PerformSkillCheck, 2个 for PerformSavingThrow, 1个 for GetPassivePerception
- `actor_test.go` - 已有较多测试
- `combat_test.go` - 已有部分测试
- `game_test.go` - 已有部分测试
- `spell_test.go` - 不存在
- `engine_test.go` - 不存在

### 需要补充测试的导出函数

#### check.go (5个导出函数)
1. **PerformAbilityCheck** - 现有2个用例，需补充3个
2. **PerformSkillCheck** - 现有2个用例，需补充3个
3. **PerformSavingThrow** - 现有2个用例，需补充3个
4. **GetSkillAbility** - 无测试，需创建5个
5. **GetPassivePerception** - 现有1个用例，需补充4个

#### engine.go (3个导出函数)
1. **New** - 无测试，需创建5个
2. **DefaultConfig** - 无测试，需创建5个
3. **Close** - 无测试，需创建5个
4. **NewTestEngine** - 测试辅助函数，已有使用

#### actor.go (16个导出函数)
需检查每个函数是否满足5个测试用例要求

#### combat.go (11个导出函数)
1. StartCombat, StartCombatWithSurprise, EndCombat, GetCurrentCombat, NextTurn, GetCurrentTurn, ExecuteAction, ExecuteAttack, ExecuteDamage, ExecuteHealing, MoveActor

#### spell.go (10个导出函数)
1. CastSpell, GetSpellSlots, PrepareSpells, LearnSpell, ConcentrationCheck, EndConcentration, CastSpellRitual, GetPactMagicSlots, RestorePactMagicSlots, IsConcentrating, GetConcentrationSpell

## 实施计划

### Phase 1: 补充 check_test.go
- 补充 PerformAbilityCheck 测试：错误处理、边界条件、不同角色类型
- 补充 PerformSkillCheck 测试：熟练加值、专家熟练、不同技能
- 补充 PerformSavingThrow 测试：熟练豁免、不同DC、劣势情况
- 创建 GetSkillAbility 测试：各种技能映射
- 补充 GetPassivePerception 测试：NPC、Enemy、Companion、不同WIS值

### Phase 2: 创建 engine_test.go
- New: 成功创建、存储初始化失败、自定义存储、配置验证
- DefaultConfig: 默认值验证、存储类型、DiceSeed
- Close: 正常关闭、重复关闭、资源释放验证

### Phase 3: 创建 spell_test.go
- CastSpell: 成功施法、法术位不足、非施法者、攻击法术、豁免法术
- GetSpellSlots: 有法术位、无法术位、非PC
- PrepareSpells: 准备型施法者、已知型施法者错误
- LearnSpell: 学习新法术、重复学习
- ConcentrationCheck: 成功、失败、未专注
- EndConcentration: 正常结束、未专注
- CastSpellRitual: 仪式施法、非仪式法术
- IsConcentrating: 专注中、未专注
- GetConcentrationSpell: 有专注、无专注

### Phase 4: 补充 combat_test.go 和 actor_test.go（如需要）

## 验证步骤
1. 运行 `go test ./pkg/engine/... -v` 确保所有测试通过
2. 运行 `go test ./pkg/engine/... -cover` 检查覆盖率
3. 确保每个导出函数至少5个测试用例
