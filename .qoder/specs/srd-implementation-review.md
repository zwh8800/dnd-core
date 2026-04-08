# D&D 5e 引擎规范实现审查计划

## Context

本项目 (`dnd-core`) 是一个 Go 语言实现的 D&D 5e 游戏引擎。规范文档 `.qoder/specs/srd-evolution-plan.md` 定义了 6 个演进阶段,共计约 22,000 行代码的完整实现要求。本次审查旨在对比规范与实际实现,找出所有缺失、错误和不一致之处,为后续修复工作提供指导。

## 审查方法

采用**四维矩阵对比法**,对每个子系统从四个维度检查:

| 维度 | 检查内容 | 对应目录 |
|------|----------|----------|
| **数据层** | 数据定义是否完整? 数量达标? | `pkg/data/*.go` |
| **模型层** | 数据结构是否定义? 接口实现? | `pkg/model/*.go` |
| **规则层** | 纯函数逻辑是否正确? | `pkg/rules/*.go` |
| **引擎层** | API 是否暴露? 是否正确调用规则? | `pkg/engine/*.go` |

## 发现的问题清单

---

### [P0] 1. 职业特性钩子实现不完整 (11/12 职业)

- **分类**: 不完整
- **阶段**: 阶段二
- **涉及文件**: 
  - `pkg/model/barbarian_hooks.go` (121行)
  - `pkg/model/rogue_hooks.go` (148行)
  - `pkg/model/monk_hooks.go` (130行)
  - `pkg/model/paladin_hooks.go` (123行)
  - `pkg/model/warlock_hooks.go` (123行)
  - `pkg/model/cleric_hooks.go` (92行)
  - `pkg/model/bard_hooks.go` (107行)
  - `pkg/model/druid_hooks.go` (95行)
  - `pkg/model/sorcerer_hooks.go` (97行)
  - `pkg/model/wizard_hooks.go` (97行)
  - `pkg/model/ranger_hooks.go` (106行)
- **规范要求**: 12个职业的特性钩子全部实现,战斗中正确触发
- **当前状态**: 
  - 所有 `*_hooks.go` 文件存在,FeatureHook 接口7个方法均已实现
  - 仅 Fighter (在 `classfeatures.go` 中) 实现完整(95%)
  - 其他职业仅标记特性"存在",无实际效果:
    - Barbarian: OnAttackRoll 为空,OnDamageCalc 仅加固定值,狂暴抗性未应用
    - Rogue: 偷袭伤害用固定值而非骰子,条件判断缺失
    - Monk: 震慑拳仅标记,轻功无实现
    - 其他职业: 多数钩子方法为空实现或仅返回标记
- **影响范围**: 战斗中 11 个职业的特性无法正确生效
- **建议修复**: 
  1. Barbarian: 实现狂暴伤害抗性(修改 damage 计算逻辑)
  2. Rogue: 偷袭伤害改为骰子表达式,添加条件判断(优势/近战/无优势等)
  3. Monk: 实现震慑拳(攻击时目标 CON 豁免)
  4. 其他职业: 逐一实现核心特性的实际效果逻辑
- **验收标准**: 
  - [ ] 5级野蛮人狂暴时获得伤害抗性
  - [ ] 3级游荡者偷袭造成 2d6 额外伤害(条件满足时)
  - [ ] 5级武僧震慑拳命中后目标需 CON 豁免
  - [ ] 单元测试覆盖各职业核心特性

---

### [P0] 2. 武器掌控效果未实际执行

- **分类**: 不完整
- **阶段**: 阶段二
- **涉及文件**: 
  - `pkg/rules/weaponmastery.go` (81行)
  - `pkg/engine/combat.go`
- **规范要求**: 8种武器掌控效果(Slow/Topple/Push/Nick/Vex/Cleave/Sap/Graze)在攻击时正确应用
- **当前状态**: 
  - `ApplyWeaponMastery()` 仅添加 Effects 标记到 AttackResult
  - 未实际执行效果(如 Topple 需要目标进行豁免检定)
  - Graze 特殊处理: 未命中时设置 GrazeDamage,但 combat.go 是否应用需验证
- **影响范围**: 武器掌控系统形同虚设,不影响实际战斗
- **建议修复**: 
  1. 在 combat.go 的 ExecuteAttack 中调用 ApplyWeaponMastery
  2. 为每种掌控实现实际效果:
     - Topple: 目标进行 STR/DEX 豁免,失败则倒地
     - Push: 推离攻击者
     - Vex: 下次攻击获得优势
     - Nick: 允许额外附赠攻击
  3. Graze: 未命中时应用一半伤害
- **验收标准**: 
  - [ ] 使用 Topple 武器攻击时,目标需进行豁免检定
  - [ ] 豁免失败的目标获得倒地状态
  - [ ] Graze 效果在未命中时造成一半伤害

---

### [P0] 3. 专长特殊能力无实际逻辑

- **分类**: 不完整
- **阶段**: 阶段二
- **涉及文件**: 
  - `pkg/rules/feats.go` (L92 注释"暂不实现")
  - `pkg/engine/feat.go`
- **规范要求**: 专长效果正确应用于战斗和检定
- **当前状态**: 
  - `ApplyFeatEffects` 中 AbilityScoreMax 标记为"暂不实现"
  - 特殊能力(如 Alert 的不被突袭、Savage Attacker 的伤害重掷)仅标记,无实际逻辑
- **影响范围**: 专长系统仅属性加成有效,特殊能力无效
- **建议修复**: 
  1. 实现 Alert 专长: 战斗中不能被突袭
  2. 实现 Savage Attacker: 伤害掷骰时可重掷一次
  3. 实现 Magic Initiate: 正确添加戏法和1环法术
  4. 在战斗中检查专长标记并应用对应效果
- **验收标准**: 
  - [ ] 拥有 Alert 专长的角色不会被突袭
  - [ ] Savage Attacker 允许重掷武器伤害骰
  - [ ] Magic Initiate 获得正确的戏法和法术位

---

### [P0] 4. 专注机制未在伤害时自动触发

- **分类**: 未集成
- **阶段**: 阶段三
- **涉及文件**: 
  - `pkg/engine/spell.go` (L478-545, 专注 API 完整)
  - `pkg/engine/combat.go` (伤害处理流程)
- **规范要求**: 专注施法者在受伤时自动进行 CON 豁免检定(DC = max(10, 伤害/2))
- **当前状态**: 
  - `ConcentrationCheck()` API 完整且 DC 计算正确
  - 但 `combat.go` 中 `applyDamageToTarget` 未自动调用专注检查
  - 专注法术不会被伤害打断
- **影响范围**: 专注施法机制失效,专注法术可以无限维持
- **建议修复**: 
  1. 在 `combat.go` 的伤害应用逻辑中,检查目标是否正在专注
  2. 如果正在专注,调用 `ConcentrationCheck()`
  3. 豁免失败则调用 `EndConcentration()` 终止法术
- **验收标准**: 
  - [ ] 专注施法者受伤时自动进行 CON 豁免
  - [ ] 豁免失败时专注法术结束
  - [ ] 大伤害(≥20)时 DC 正确提升

---

### [P0] 5. 制作系统为占位实现

- **分类**: 不完整
- **阶段**: 阶段五
- **涉及文件**: 
  - `pkg/engine/crafting.go` (120行)
  - `pkg/rules/crafting.go` (48行)
  - `pkg/model/crafting.go` (44行)
- **规范要求**: 装备制作系统可创建非魔法物品,包含配方、材料、工具、时间
- **当前状态**: 
  - `pkg/engine/crafting.go` 硬编码 `TotalDays: 7`
  - 无配方数据库,无材料验证,无工具熟练影响
  - `pkg/rules/crafting.go` 仅48行占位代码
- **影响范围**: 制作系统完全无法使用
- **建议修复**: 
  1. 创建配方数据库 (`pkg/data/crafting_recipes.go`)
  2. 实现材料验证逻辑
  3. 工具熟练影响制作时间和成功率
  4. 实现进度追踪和完成逻辑
- **验收标准**: 
  - [ ] 可查询可制作物品列表
  - [ ] 制作时验证材料是否充足
  - [ ] 工具熟练减少制作时间
  - [ ] 完成制作后物品添加到库存

---

### [P1] 6. 升环法术处理简化

- **分类**: 不完整
- **阶段**: 阶段三
- **涉及文件**: 
  - `pkg/rules/spelleffects.go` (L79 注释"简化的升环处理")
- **规范要求**: 法术根据施放环级正确计算效果(伤害/治疗/持续时间)
- **当前状态**: 
  - 升环伤害线性叠加,未处理复杂公式
  - 注释明确标注"应根据法术位等级计算"
- **影响范围**: 高环法术效果不正确
- **建议修复**: 
  1. 为每个法术定义升环公式(`SpellDamageEntry` 中添加 `PerSlotLevel`)
  2. 修改 `CalculateSpellDamage()` 使用实际环级计算
- **验收标准**: 
  - [ ] Fireball 在3环造成 8d6,4环 9d6,5环 10d6
  - [ ] 升环逻辑覆盖所有有升环效果的法术

---

### [P1] 7. 魔法物品效果仅返回消息

- **分类**: 不完整
- **阶段**: 阶段三
- **涉及文件**: 
  - `pkg/rules/magicitem.go` (L62 `UseMagicItem`)
- **规范要求**: 激活魔法物品后实际执行效果(治疗/伤害/状态/属性加成)
- **当前状态**: 
  - `UseMagicItem()` 仅返回描述消息,未实际执行效果
  - 药水/毒药仅消息,不应用治疗/伤害
  - 充能物品消耗充能✅,但效果仅消息
  - `RechargeMagicItems()` 仅处理 `recharge=="dawn"`,缺 ShortRest/LongRest/1d6
- **影响范围**: 魔法物品系统形同虚设
- **建议修复**: 
  1. `UseMagicItem()` 根据物品类型执行实际效果
  2. 药水应用治疗/伤害/状态
  3. 武器/护甲应用属性加成
  4. 充能物品实现完整充逻辑
- **验收标准**: 
  - [ ] 使用治疗药水实际恢复 HP
  - [ ] +1 武器提供攻击和伤害加值
  - [ ] 充能物品正确恢复充能

---

### [P1] 8. 多职业 Extra Attack 不完整

- **分类**: 不完整
- **阶段**: 阶段二
- **涉及文件**: 
  - `pkg/rules/multiclass.go` (`getExtraAttacksForClass` 仅实现战士)
- **规范要求**: 处理 Extra Attack 不叠加规则,支持所有有额外攻击的职业
- **当前状态**: 
  - `getExtraAttacksForClass` 仅实现 Fighter(L169-182)
  - 缺少 Barbarian(4级后狂暴时额外攻击)、Ranger(5级)、Monk(5级武术)
- **影响范围**: 多职业角色的额外攻击计算不正确
- **建议修复**: 
  1. 补充 Barbarian/Ranger/Monk 的额外攻击逻辑
  2. 验证不叠加规则(取最高,不累加)
- **验收标准**: 
  - [ ] 战士5级/游侠5级 各有1次额外攻击
  - [ ] 战士5级+游侠5级 不叠加,仍为1次
  - [ ] 战士11级 提供2次额外攻击

---

### [P1] 9. E2E 测试完全缺失

- **分类**: 缺失
- **阶段**: 阶段六
- **涉及文件**: 
  - `pkg/engine/e2e_character_creation_test.go` (不存在)
  - `pkg/engine/e2e_combat_test.go` (不存在)
  - `pkg/engine/e2e_multiclass_test.go` (不存在)
  - `pkg/engine/e2e_spellcasting_test.go` (不存在)
  - `pkg/engine/e2e_monster_test.go` (不存在)
- **规范要求**: 5个端到端测试文件验证完整流程
- **当前状态**: 所有 E2E 测试文件不存在
- **影响范围**: 无法验证系统整体正确性
- **建议修复**: 
  1. 创建 `e2e_character_creation_test.go`: 完整角色创建流程
  2. 创建 `e2e_combat_test.go`: 完整战斗流程
  3. 创建 `e2e_multiclass_test.go`: 多职业验证
  4. 创建 `e2e_spellcasting_test.go`: 完整施法流程
  5. 创建 `e2e_monster_test.go`: 怪物战斗测试
- **验收标准**: 
  - [ ] 5个 E2E 测试全部通过
  - [ ] 覆盖规范中定义的核心流程

---

### [P1] 10. rules/data/model/dice/storage 测试 0% 覆盖

- **分类**: 未测试
- **阶段**: 阶段六
- **涉及文件**: 
  - `pkg/rules/*_test.go` (不存在)
  - `pkg/data/*_test.go` (不存在)
  - `pkg/model/*_test.go` (不存在)
- **当前状态**: 仅 `pkg/engine` 有测试,覆盖率 34.8%
- **影响范围**: 代码质量风险,重构无保障
- **建议修复**: 为关键模块添加单元测试
- **验收标准**: 
  - [ ] pkg/rules 覆盖率 ≥ 70%
  - [ ] pkg/data 覆盖率 ≥ 70%
  - [ ] pkg/model 覆盖率 ≥ 70%

---

### [P2] 11. 多职业熟练项简化

- **分类**: 不完整
- **阶段**: 阶段二
- **涉及文件**: `pkg/rules/multiclass.go` (L81-98)
- **当前状态**: 仅给予轻甲熟练,未按职业给予不同熟练
- **建议修复**: 根据新增职业给予对应熟练项

---

### [P2] 12. Buff/Teleport/Utility 法术仅标记

- **分类**: 不完整
- **阶段**: 阶段三
- **涉及文件**: `pkg/rules/spelleffects.go` (L255-262)
- **当前状态**: ExecuteSpell 中这些效果类型仅标记,无实际逻辑
- **建议修复**: 实现状态施加/传送/工具类法术效果

---

### [P2] 13. 怪物数据仅 5 个

- **分类**: 不完整
- **阶段**: 阶段一/六
- **涉及文件**: `pkg/data/monsters.go`
- **规范要求**: 180 个怪物数据
- **当前状态**: 仅 5 个示例怪物
- **建议修复**: 批量导入剩余 170+ 怪物

---

### [P2] 14. 死亡治疗自动稳定仅消息

- **分类**: 不完整
- **阶段**: 阶段五
- **涉及文件**: `pkg/rules/death.go` (L77)
- **当前状态**: 治疗时仅返回消息,未实际恢复 HP 和稳定状态
- **建议修复**: 治疗时自动设置稳定状态并恢复对应 HP

---

### [P3] 15. 批量导入工具不存在

- **分类**: 缺失
- **阶段**: 阶段六
- **涉及文件**: `pkg/data/import/` 和 `cmd/import/` (不存在)
- **建议修复**: 创建 CLI 工具从 SRD Markdown 批量导入数据

---

### [P3] 16. API 文档和示例不存在

- **分类**: 缺失
- **阶段**: 阶段六
- **涉及文件**: `docs/api/` 和 `examples/` (不存在)
- **建议修复**: 补充 API 文档和使用示例

---

### [P3] 17. 缓存未集成到查询

- **分类**: 未集成
- **阶段**: 阶段六
- **涉及文件**: `pkg/data/cache.go` (95行)
- **当前状态**: LRUCache 存在但 GetRace/GetClass 等查询未使用
- **建议修复**: 在数据查询函数中集成缓存

---

## 问题统计

| 优先级 | 数量 | 分类 |
|--------|------|------|
| P0 - Critical | 5 | 不完整(3), 未集成(1), 占位(1) |
| P1 - Major | 5 | 不完整(3), 缺失(1), 未测试(1) |
| P2 - Minor | 4 | 不完整(3), 数据不足(1) |
| P3 - Nice-to-have | 3 | 缺失(2), 未集成(1) |
| **总计** | **17** | |

## 验证命令

```bash
# 1. 统计法术数量
grep -c '"id":' pkg/data/spells.go pkg/data/spells_additional.go

# 2. 统计专长数量  
grep -c '"id":' pkg/data/feats.go pkg/data/feats_additional.go

# 3. 统计武器/护甲数量
grep -c '"id":' pkg/data/weapons.go pkg/data/armors.go

# 4. 统计怪物数量
grep -c '"id":' pkg/data/monsters.go

# 5. 搜索未实现的标记
grep -rn "需要额外实现\|简化\|暂不实现\|TODO\|FIXME\|仅标记" pkg/

# 6. 检查钩子调用
grep -rn "OnAttackRoll\|OnDamageCalc\|OnACCalc" pkg/engine/

# 7. 检查专注调用
grep -rn "Concentration\|IsConcentrating" pkg/engine/combat.go

# 8. 检查 ApplyWeaponMastery 调用
grep -rn "ApplyWeaponMastery" pkg/engine/

# 9. 检查测试覆盖
go test ./pkg/rules/... -v -count=1
go test ./pkg/data/... -v -count=1

# 10. 代码质量
go vet ./...
```

## 修复优先级建议

1. **第一批** (P0): 职业特性钩子 + 专注机制 + 武器掌控
2. **第二批** (P0): 专长特殊能力 + 制作系统重写
3. **第三批** (P1): 升环法术 + 魔法物品效果 + Extra Attack
4. **第四批** (P1): E2E 测试 + 单元测试补充
5. **第五批** (P2): 多职业熟练 + Buff法术 + 怪物数据 + 死亡治疗
6. **第六批** (P3): 批量导入工具 + 文档 + 缓存集成
