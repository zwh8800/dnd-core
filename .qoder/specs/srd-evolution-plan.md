# D&D 5e SRD 5.2.1 完整演进计划

## Context

本项目是一个 Go 语言编写的 D&D 5e 核心游戏引擎，目前已实现基础的角色创建、战斗、检定、法术框架等核心功能（约 14,500 行代码）。项目包含完整的 SRD 5.2.1 规范文档（258 个 Markdown 文件），但当前实现仅覆盖了 SRD 的部分功能。

**目标**：通过 6 个阶段的系统性演进，使本项目成为完整实现 SRD 5.2.1 所有核心规则的 D&D 5e 游戏引擎，预计新增约 22,000 行代码。

---

## 一、现状深度对比分析

### 1.1 完全缺失的系统

| 系统 | SRD 章节 | 说明 |
|------|----------|------|
| 背景系统 (Backgrounds) | 04_character-origins | 4+ 背景，包含属性加成、专长、技能/工具熟练、起始装备 |
| 专长系统 (Feats) | 05_feats | 起源/通用/战斗/史诗恩赐专长，50+ 专长定义 |
| 多职业系统 (Multiclassing) | 02_character-creation | 13 分属性要求、熟练项计算、法术位共享规则 |
| 怪物数据系统 | 11-12_monsters | 180 个怪物数据块，含特性/动作/反应/传说动作 |
| 武器掌控 (Weapon Mastery) | 06_equipment/01_weapons | SRD 5.2.1 新特性：Slow/Topple/Push/Nick/Vex/Cleave 等 |
| 法术数据库 | 07_spells | 200+ 具体法术定义和效果 |
| 武器/护甲数据 | 06_equipment/01-02 | 30+ 武器、16 护甲完整数据 |
| 魔法物品详细实现 | 10_magic-items | 激活/调音/充能/诅咒/有灵智物品，100+ 物品 |
| 陷阱/毒药/诅咒 | 09_gameplay-toolbox | 陷阱检测/解除、毒药涂抹、诅咒施加/移除 |
| 旅行探索系统 | 09_gameplay-toolbox/01 | 旅行速度、觅食、导航、随机遭遇 |
| 社交互动系统 | 01_playing-the-game/06 | NPC 态度（友好/冷漠/敌对）、影响检定 |
| 生活方式经济 | 06_equipment/06-08 | 7 级生活开销、雇员、食物住宿 |
| 坐骑和交通工具 | 06_equipment/05 | 坐骑数据、马/船/马车、载重计算 |
| 装备制作系统 | 06_equipment/11-13 | 非魔法物品制作、药水酿造、法术卷轴书写 |
| 灵感系统 (Inspiration) | 01_playing-the-game | DM 授予的优势点，使用机制 |

### 1.2 部分实现需要完善的系统

| 系统 | 现状 | 缺失内容 |
|------|------|----------|
| 职业特性 | 仅战士部分特性 | 其余 11 职业 × 20 级特性（狂暴、偷袭、变身等） |
| 种族系统 | 9 种族基础定义 | 子种族嵌套、多速度类型、感官、完整特性效果 |
| 法术位系统 | 基础结构存在 | 各职业施法机制差异、仪式施法、Pact Magic |
| 休息系统 | 短休/长休基础 | 生命骰详细恢复、各种特性恢复规则 |
| 状态效果 | 15 种定义 | 效果应用不完整（如目盲的攻击劣势） |
| 死亡豁免 | 字段存在 | 完整规则（3 成功/失败、大出血即死） |
| 护甲等级 | 基础计算 | 无甲防御（蛮族/武僧）、特定护甲效果 |

### 1.3 已完整实现的系统

- ✅ 六大属性系统（STR/DEX/CON/INT/WIS/CHA）
- ✅ 18 技能系统及其属性映射
- ✅ 熟练加值计算（1-20 级）
- ✅ 骰子系统（优势/劣势/保留高低/修饰符）
- ✅ 13 种伤害类型和伤害计算
- ✅ 基础战斗流程（先攻/回合/动作管理）
- ✅ 游戏状态管理和持久化
- ✅ 场景系统
- ✅ 任务系统
- ✅ 存储接口（内存实现）

---

## 二、分阶段演进路线图

### 阶段一：基础设施与数据骨架

**目标**：建立数据加载架构，完善核心数据模型，实现怪物系统架构和示例数据。

#### 任务清单

**1.1 数据加载器框架**
- 新增 `pkg/data/loader.go` — 统一数据加载接口，支持 JSON/Markdown 解析
- 新增 `pkg/data/embed.go` — 使用 `//go:embed` 内嵌 SRD 数据文件
- 新增 `pkg/data/registry.go` — 数据注册中心（RegisterRace/RegisterClass/RegisterMonster 等）
- 修改 `pkg/engine/config.go` — 增加 DataLoader 配置字段

**1.2 怪物系统架构**
- 新增 `pkg/model/monster.go`:
  - `MonsterStatBlock` — 完整 SRD 怪物数据块（体型、类型、AC、先攻调整、HP、速度、六属性、豁免、技能、感官、语言、CR、XP、PB、特性、动作、附赠动作、反应、传说动作）
  - `MonsterAction` — 怪物动作定义（攻击/法术/特效、命中/未命中效果、充能机制）
  - `MonsterTrait` — 被动特性（抗性/免疫/易伤、伤害免疫、状态免疫）
  - `RechargeInfo` — 充能机制（Recharge X-Y、X/Day）
- 新增 `pkg/model/creature.go` — `CreatureType` 枚举（14 种生物类型）
- 修改 `pkg/model/actor.go` — Actor 增加 `CreatureType` 和 `ChallengeRating` 字段
- 新增 `pkg/data/monsters.go` — 5 个示例怪物：
  - Goblin (CR 1/4) — 简单近战怪物
  - Ogre (CR 2) — 中型怪物，含多动作
  - Owlbear (CR 3) — 复杂怪物，含多攻击方式
  - Mimic (CR 2) — 含伪装特性
  - Young Red Dragon (CR 10) — 复杂怪物，含传说动作、多充能动作
- 新增 `pkg/engine/monster.go` — 怪物引擎 API：
  - `LoadMonster(template string) (*Enemy, error)` — 从模板创建怪物
  - `GetMonsterActions(monster *Enemy) []Action` — 获取可用动作

**1.3 背景系统数据模型**
- 新增 `pkg/model/background.go`:
  - `BackgroundDefinition` — 背景定义（属性加成选项、关联专长、技能熟练、工具熟练、起始装备选择）
  - `BackgroundChoice` — 装备 A/B 选择
  - `BackgroundID` — 背景标识常量
- 新增 `pkg/data/backgrounds.go` — 4 个 SRD 背景：Acolyte、Criminal、Sage、Soldier

**1.4 专长系统数据模型**
- 新增 `pkg/model/feat.go`:
  - `FeatDefinition` — 专长定义（名称、类型、先决条件、属性加成、效果描述）
  - `FeatType` 枚举 — Origin / General / Combat / Epic
  - `FeatPrerequisite` — 先决条件（最低属性/职业/等级）
  - `FeatEffect` — 效果（属性修正/攻击修正/AC 修正/特殊能力标记）
  - `FeatInstance` — 专长实例（FeatID、获得来源、获得等级）
- 新增 `pkg/data/feats.go` — 5 个起源专长：Alert、Magic Initiate、Savage Attacker、Skilled、Tough

**1.5 种族系统完善**
- 修改 `pkg/model/actor.go` — 新增 `SpeedTypes` 结构体（Walk/Swim/Fly/Climb/Burrow 多速度支持）
- 修改 `pkg/data/races.go` — 按 SRD 5.2.1 完善：
  - 补充子种族嵌套结构（Elf: High/Wood/Drow、Dwarf: Hill/Mountain 等）
  - 补充感官定义（Darkvision 距离）
  - 补充完整种族特性列表和效果标记

#### 验收标准
- [ ] 数据加载器能从内嵌数据解析并注册怪物/背景/专长/种族
- [ ] 5 个示例怪物可通过 `LoadMonster("goblin")` 创建为完整 Enemy 对象
- [ ] 背景、专长、种族数据可通过 API 查询
- [ ] `MonsterStatBlock` 能完整表示 Young Red Dragon 级别的复杂怪物（含传说动作）

#### 里程碑
> DM 可通过 `"goblin"` 模板名称一键生成完整的哥布林敌人，包含所有动作、特性和正确数值，直接投入战斗。

---

### 阶段二：核心游戏机制完善

**目标**：实现专长系统、多职业系统、背景集成、完整职业特性钩子、武器掌控。

#### 任务清单

**2.1 专长系统集成**
- 修改 `pkg/model/character.go` — `PlayerCharacter` 添加 `Feats []FeatInstance` 字段
- 新增 `pkg/rules/feats.go`:
  - `CheckFeatPrerequisites(character *PlayerCharacter, featID string) bool` — 验证先决条件
  - `ApplyFeatEffects(character *PlayerCharacter, featID string)` — 应用专长效果
  - `GetFeatBonuses(character *PlayerCharacter) FeatBonuses` — 汇总所有专长加值
- 新增 `pkg/engine/feat.go`:
  - `SelectFeat(actorID, featID string) error` — 角色选择专长
  - `ListFeats(filter FeatFilter) []FeatDefinition` — 列出可选专长
  - `GetFeatDetails(featID string) (*FeatDefinition, error)` — 获取专长详情
  - `RemoveFeat(actorID, featID string) error` — 移除专长
- 修改 `pkg/engine/actor.go` — `CreatePC()` 集成背景关联的起源专长自动赋予

**2.2 多职业系统完整实现**
- 新增 `pkg/rules/multiclass.go`:
  - `ValidateMulticlass(character *PlayerCharacter, newClass ClassID) error` — 验证 13 分属性要求
  - `GetMulticlassSpellSlots(classes []ClassLevel) SpellSlotTable` — 计算多职业法术位表
  - `GetMulticlassProficiencies(classes []ClassLevel) Proficiencies` — 计算多职业熟练项
  - `ValidateExtraAttack(classes []ClassLevel) int` — 处理 Extra Attack 不叠加规则
  - `ValidateUnarmoredDefense(classes []ClassLevel) bool` — 处理无甲防御不叠加规则
- 修改 `pkg/engine/actor.go`:
  - 修改 `LevelUp()` 支持选择新职业（多职业）
  - 新增 `ValidateLevelUpChoice()` 验证升级选择合法性

**2.3 职业特性钩子系统扩展（11 个职业）**
- 修改 `pkg/model/classfeatures.go` — 为每个职业创建 `*FeatureHooks` 结构体并实现 `FeatureHook` 接口：
  - `BarbarianFeatureHooks` — 狂暴（伤害抗性/加值）、无甲防御、危险感知、额外攻击、狂暴次数
  - `BardFeatureHooks` — 吟游灵感骰跟踪、法术、技能专家
  - `ClericFeatureHooks` — 引导神力、领域特性、神圣干预
  - `DruidFeatureHooks` — 德鲁伊语言、野兽形态、Timeless Body
  - `MonkFeatureHooks` — 武术、气点、Unarmored Movement、Deflect Missiles、Stunning Strike
  - `PaladinFeatureHooks` — 神圣感知、圣手、至圣斩、光环
  - `RangerFeatureHooks` — 偏爱敌人、自然探索、原始意识
  - `RogueFeatureHooks` — 偷袭、狡诈动作、专家、直觉闪避、反射闪避
  - `SorcererFeatureHooks` — 超魔、术力点、法术源泉
  - `WarlockFeatureHooks` — 契约魔法、邪术祈唤、契约恩赐
  - `WizardFeatureHooks` — 奥术恢复、法术书、奥术传统
- 修改 `pkg/model/class.go` — 新增每个职业的 `*State` 跟踪结构体
- 修改 `pkg/engine/actor.go` — `CreatePC()` 和 `LevelUp()` 初始化对应职业特性钩子

**2.4 背景系统集成**
- 修改 `pkg/model/character.go` — `PlayerCharacter` 添加 `BackgroundID string` 字段
- 新增 `pkg/rules/background.go` — `ApplyBackground()` 应用背景效果（属性、专长、技能、工具、装备）
- 修改 `pkg/engine/actor.go` — `CreatePC()` 接受背景参数并自动应用效果

**2.5 武器掌控系统 (Weapon Mastery)**
- 修改 `pkg/model/equipment.go` — `WeaponProperties` 添加 `Mastery string` 字段
- 新增 `pkg/model/weaponmastery.go`:
  - `WeaponMasteryType` 枚举 — Slow / Topple / Push / Nick / Vex / Cleave / Sap / Graze
  - `WeaponMasteryEffect` — 掌控效果定义
- 新增 `pkg/rules/weaponmastery.go` — `ApplyWeaponMastery()` 在攻击时应用武器掌控效果

#### 验收标准
- [ ] 角色可选择专长，效果正确应用于战斗和检定（如 Alert 专长先攻+PB）
- [ ] 多职业升级验证属性要求（13 分规则），正确计算共享法术位
- [ ] 12 个职业的特性钩子全部实现，战斗中正确触发（如战士 Extra Attack 在同一回合多次攻击）
- [ ] 创建角色时选择背景自动获得对应专长、技能和起始装备
- [ ] 武器掌控效果在攻击时正确应用（如 Vex 命中后下次攻击优势）

#### 里程碑
> 创建一个 5 级战士 / 3 级游荡者角色，拥有 Criminal 背景的 Alert 专长、战士战斗风格和游荡者偷袭伤害，所有特性在战斗中正确运作。

---

### 阶段三：内容填充——法术、装备、魔法物品

**目标**：填充法术数据库、武器护甲数据、魔法物品系统。

#### 任务清单

**3.1 法术数据库**
- 新增 `pkg/model/spelleffect.go`:
  - `SpellEffect` — 法术效果执行逻辑（伤害/治疗/状态施加/召唤）
  - `SpellTargetType` 枚举 — SingleTarget / Cone / Sphere / Line / Emanation / Self
  - `SpellDamageEntry` — 基础骰子、高环升级公式
- 新增 `pkg/data/spells.go` — 50 个核心法术（戏法到 5 环，每学派至少 2 个）：
  - 戏法: Fire Bolt、Light、Mage Hand、Sacred Flame、Shocking Grasp
  - 1 环: Magic Missile、Shield、Cure Wounds、Bless、Burning Hands
  - 2 环: Scorching Ray、Misty Step、Hold Person
  - 3 环: Fireball、Counterspell、Revivify、Lightning Bolt
  - 4-5 环: Polymorph、Cone of Cold、Wall of Force 等
- 新增 `pkg/rules/spelleffects.go`:
  - `ExecuteSpell()` — 根据法术定义执行效果
  - `ResolveSpellSave()` — 处理法术豁免
  - `CalculateSpellDamage()` — 计算法术伤害（含高环升级）

**3.2 武器和护甲完整数据**
- 新增 `pkg/data/weapons.go` — 30+ 武器（简易+军用），含伤害骰、属性、价格、重量、掌控
- 新增 `pkg/data/armors.go` — 16 护甲（轻/中/重+盾牌），含 AC、DEX 上限、隐匿劣势、STR 要求
- 新增 `pkg/data/adventuring_gear.go` — 40+ 冒险装备（绳索、火把、口粮、帐篷等）
- 新增 `pkg/data/tools_data.go` — 工具数据（盗贼工具、工匠工具、乐器等）

**3.3 魔法物品系统**
- 新增 `pkg/model/magicitem.go`:
  - `MagicItemDefinition` — 名称、类型、稀有度、调音要求、充能系统、激活方式、效果、诅咒标记
  - `MagicItemActivation` 枚举 — Command / Spell / Charge / Consumable / Equipped
  - `ChargesRecharge` — 充能恢复规则（Dawn / ShortRest / LongRest / 1d6 Roll）
  - `SentientItemProperties` — 有意识魔法物品（人格、沟通方式、意志）
- 新增 `pkg/data/magicitems.go` — 30 个经典魔法物品：
  - +1/+2/+3 武器和护甲
  - Ring of Protection、Cloak of Elvenkind、Boots of Speed
  - Wand of Fireballs、Staff of Power
  - Potion of Healing、Scroll of Fireball
- 新增 `pkg/rules/magicitem.go`:
  - `ActivateMagicItem()` — 激活魔法物品
  - `CheckAttunement()` — 验证调音条件
  - `ApplyMagicItemEffects()` — 应用魔法物品效果
- 修改 `pkg/engine/inventory.go`:
  - 新增 `EquipMagicItem()` / `UnequipMagicItem()` — 调音槽位管理（最多 3 件）
  - 新增 `ActivateItem()` API

**3.4 法术系统集成完善**
- 修改 `pkg/engine/spell.go`:
  - 完善 `CastSpell()` 集成法术效果执行
  - 实现 Pact Magic 特殊规则（邪术师法术位独立追踪）
  - 实现 Ritual Casting 仪式施法（+10 分钟施法时间，不消耗法术位）
  - 实现 Concentration 专注跟踪和打断机制（受伤时 CON 豁免 DC = max(10, 伤害/2)）

#### 验收标准
- [ ] 50+ 法术可通过 API 查询并施放，效果正确（如 Fireball 8d6 火焰伤害，60 尺半径）
- [ ] 所有武器和护甲数据可通过 API 查询
- [ ] 30+ 魔法物品数据完整，激活和调音机制正确（调音上限 3 件）
- [ ] 施法者能正确消耗法术位、处理专注打断、执行仪式施法

#### 里程碑
> 法师角色可以准备法术列表、施放 Fireball 对区域内敌人造成 8d6 伤害、维持 Shield 的专注，并在法术位耗尽后使用 Arcane Recovery 恢复部分法术位。

---

### 阶段四：探索、社交与经济系统

**目标**：实现探索规则、社交互动、旅行系统、生活方式经济、坐骑载具。

#### 任务清单

**4.1 旅行和探索系统**
- 新增 `pkg/model/exploration.go`:
  - `TravelState` — 当前位置、目的地、旅行速度、已行距离、时间消耗
  - `TravelPace` 枚举 — Fast / Normal / Slow
  - `ForageResult` — 觅食结果
  - `NavigationCheck` — 导航检定
- 新增 `pkg/rules/exploration.go`:
  - `CalculateTravelDistance()` — 根据速度、地形、步伐计算日行距离
  - `ForagingCheck()` — 觅食检定（WIS 生存检定）
  - `NavigationCheck()` — 导航检定（WIS 生存检定）
  - `EncounterCheck()` — 随机遭遇检定
- 新增 `pkg/engine/exploration.go`:
  - `StartTravel()` / `AdvanceTravel()` — 开始/推进旅行
  - `Forage()` / `Navigate()` — 觅食/导航

**4.2 社交互动系统**
- 新增 `pkg/model/social.go`:
  - `NPCAttitude` 枚举 — Friendly / Indifferent / Hostile
  - `SocialInteractionState` — 当前态度、倾向、已建立的印象
- 新增 `pkg/rules/social.go`:
  - `CalculateNPCReaction()` — 基于 NPC 倾向和角色表现计算反应
  - `DetermineAttitudeChange()` — 判定态度变化
- 新增 `pkg/engine/social.go`:
  - `InteractWithNPC()` — 执行社交互动（游说/欺骗/威吓检定）
  - `GetNPCAttitude()` — 获取 NPC 当前态度

**4.3 生活方式与经济系统**
- 新增 `pkg/model/lifestyle.go`:
  - `LifestyleTier` 枚举 — Wretched / Squalid / Poor / Modest / Comfortable / Wealthy / Aristocratic
  - `LifestyleCost` — 每日/每月费用
- 新增 `pkg/data/lifestyles.go` — 7 种生活方式的详细开销数据
- 新增 `pkg/rules/lifestyle.go` — `CalculateLifestyleCost()` / `DeductLifestyle()`
- 新增 `pkg/engine/lifestyle.go`:
  - `SetLifestyle()` — 设置生活方式
  - `AdvanceGameTime()` — 推进游戏时间并自动扣除开销

**4.4 坐骑和交通工具**
- 新增 `pkg/model/mount.go`:
  - `Mount` — 生物引用、骑乘状态、鞍具
  - `Vehicle` — 类型、HP、AC、速度、载货量、船员需求
- 新增 `pkg/data/mounts.go` — 坐骑（马/骑乘犬等）和交通工具（马车/船等）数据
- 新增 `pkg/engine/mount.go` — `MountCreature()` / `Dismount()` / `CalculateMountSpeed()`

#### 验收标准
- [ ] 队伍可开始旅行，根据步伐和地形正确计算日行距离
- [ ] NPC 有态度系统，社交检定（游说/欺骗/威吓）可改变态度
- [ ] 生活方式开销随游戏时间自动扣除
- [ ] 角色可骑乘坐骑，速度正确计算（载重影响）

#### 里程碑
> 队伍从城镇出发，选择快速步伐穿越森林，途中进行觅食和导航检定，遭遇 NPC 并通过游说检定获得帮助，每晚自动扣除舒适生活方式的费用（每天 1 GP）。

---

### 阶段五：高级游戏系统

**目标**：实现陷阱、毒药、诅咒、环境效果、死亡复苏详细规则、装备制作、力竭系统。

#### 任务清单

**5.1 陷阱系统**
- 新增 `pkg/model/trap.go`:
  - `TrapDefinition` — 类型（机械/魔法）、触发条件、检测 DC、解除 DC、效果
  - `TrapState` — 是否已触发、剩余次数
- 新增 `pkg/data/traps.go` — 10 种经典陷阱（毒针、落石、陷阱坑等）
- 新增 `pkg/engine/trap.go` — `PlaceTrap()` / `DetectTrap()` / `DisarmTrap()` / `TriggerTrap()`

**5.2 毒药系统**
- 新增 `pkg/model/poison.go`:
  - `PoisonDefinition` — 类型（接触/摄入/吸入/伤口）、效果、持续时间、豁免 DC、价格
- 新增 `pkg/data/poisons.go` — 10 种毒药数据（基本到强力）
- 新增 `pkg/engine/poison.go` — `ApplyPoison()` / `ResolvePoisonEffect()`

**5.3 诅咒系统**
- 新增 `pkg/model/curse.go`:
  - `CurseDefinition` — 效果、移除条件
  - `CurseInstance` — 来源、剩余持续时间
- 新增 `pkg/engine/curse.go` — `CurseActor()` / `RemoveCurse()`

**5.4 环境效果**
- 新增 `pkg/model/environment.go`:
  - `EnvironmentalEffect` — 类型（极寒/极热/高海拔/深水等）、效果、豁免
- 新增 `pkg/rules/environment.go` — `ApplyEnvironmentalEffects()`
- 新增 `pkg/engine/environment.go` — `SetEnvironment()` / `ResolveEnvironmentalDamage()`

**5.5 死亡和复苏详细规则**
- 修改 `pkg/model/actor.go` — 完善 `IsDead()` / `IsStabilized()` 逻辑，添加 `DeathSaveRoll` 记录
- 新增 `pkg/rules/death.go`:
  - 完整死亡豁免规则（3 成功=稳定，3 失败=死亡）
  - 伤害 ≥ HP 上限时立即死亡规则
  - 受到治疗时自动稳定并恢复对应 HP
- 修改 `pkg/engine/actor.go`:
  - 新增 `MakeDeathSave()` API
  - 新增 `StabilizeActor()` 稳定濒死角色
  - 完善 `ApplyDamage()` 处理 0 HP 情况

**5.6 装备制作系统**
- 新增 `pkg/model/crafting.go`:
  - `CraftingRecipe` — 所需材料、工具、时间、DC
  - `CraftingProgress` — 进度、剩余时间
- 新增 `pkg/rules/crafting.go` — `CalculateCraftingTime()` / `CalculateCraftingCost()`
- 新增 `pkg/engine/crafting.go` — `StartCrafting()` / `AdvanceCrafting()` / `CompleteCrafting()`

**5.7 力竭系统完善**
- 新增 `pkg/rules/exhaustion.go`:
  - 6 级力竭效果表（1 级检定劣势 → 6 级死亡）
  - `ApplyExhaustionEffects()` — 根据等级应用效果
- 修改 `pkg/engine/actor.go` — 力竭等级变化时自动应用/移除效果

#### 验收标准
- [ ] 场景中可放置和检测陷阱，触发时正确处理效果
- [ ] 毒药可涂抹到武器并在命中时生效（如伤口毒药：CON 豁免失败中毒）
- [ ] 诅咒可施加和移除（需要 Remove Curse 法术等）
- [ ] 环境效果正确影响角色（如极寒每 10 分钟 CON 豁免失败受 1d6 寒冷伤害）
- [ ] 死亡豁免完整实现（3 成功=稳定，天然 20=立即恢复 1HP，天然 1=2 次失败）
- [ ] 装备制作系统可创建非魔法物品

#### 里程碑
> 队伍探索地下城，察觉（察觉检定 DC 15）并解除（妙手检定 DC 15）陷阱，给武器涂抹毒药后与怪物战斗，成员陷入濒死并通过死亡豁免稳定，最后因极寒环境获得 1 级力竭。

---

### 阶段六：整合、优化与数据完善

**目标**：全面集成测试、性能优化、数据完整性提升、API 文档化。

#### 任务清单

**6.1 批量数据导入**
- 新增 `pkg/data/import/` 目录 — 批量导入工具：
  - 怪物批量导入（解析剩余 170+ 怪物 Markdown 文件）
  - 法术批量导入
  - 魔法物品批量导入
- 新增 `cmd/import/` CLI 工具 — 从 SRD Markdown 批量导入数据

**6.2 全面集成测试**
- 新增端到端测试文件：
  - `pkg/engine/e2e_character_creation_test.go` — 完整角色创建流程
  - `pkg/engine/e2e_combat_test.go` — 完整战斗流程（多角色、多职业特性触发）
  - `pkg/engine/e2e_multiclass_test.go` — 多职业验证
  - `pkg/engine/e2e_spellcasting_test.go` — 完整施法流程（含专注、仪式、高环升级）
  - `pkg/engine/e2e_monster_test.go` — 怪物战斗测试

**6.3 API 文档和示例**
- 新增 `docs/api/` 目录 — API 使用文档、角色创建示例代码、战斗流程示例代码
- 新增 `examples/` 目录 — 完整游戏循环示例、Web API 集成示例

**6.4 性能优化**
- 新增 `pkg/data/cache.go` — LRU 缓存频繁查询的数据
- 优化 `ListActors()` 等大列表查询
- 怪物模板使用指针引用而非复制

**6.5 剩余数据填充**
- 补充剩余种族数据（按 SRD 5.2.1 规范完善子种族）
- 补充更多专长数据（通用专长、战斗专长，目标 30+ 专长）
- 补充更多法术数据（6-9 环法术，目标 100+ 法术）
- 补充更多魔法物品数据（目标 60+ 魔法物品）

#### 验收标准
- [ ] 170+ 怪物数据可通过导入工具加载
- [ ] 所有端到端测试通过
- [ ] API 文档完整可用
- [ ] 查询响应时间 < 10ms（缓存命中）
- [ ] 内存使用稳定，无泄漏

#### 里程碑
> 系统可完整支持一场 D&D 5e 游戏的所有核心规则，DM 可以加载任意怪物、玩家可以使用完整角色创建和升级流程、战斗系统正确处理所有职业特性。

---

## 三、依赖关系图

```
阶段一 (基础设施)
    ├── 数据加载器 ──────┐
    ├── 怪物系统 ─────────┤
    ├── 背景系统 ─────────┤
    ├── 专长系统 (数据) ──┤
    └── 种族完善 ─────────┤
                          ▼
阶段二 (核心机制) ◄───────┘
    ├── 专长系统集成
    ├── 多职业系统
    ├── 职业特性钩子 (11 职业)
    ├── 背景集成
    └── 武器掌控
                          ▼
阶段三 (内容填充) ◄───────┘
    ├── 法术数据库 (50 法术)
    ├── 武器/护甲数据 (46 项)
    ├── 魔法物品系统 (30 物品)
    └── 法术系统完善
                          ▼
阶段四 (探索社交) ◄───────┘
    ├── 旅行探索
    ├── 社交互动
    ├── 生活方式经济
    └── 坐骑载具
                          ▼
阶段五 (高级系统) ◄───────┘
    ├── 陷阱/毒药/诅咒
    ├── 环境效果
    ├── 死亡复苏详细规则
    └── 装备制作
                          ▼
阶段六 (整合优化) ◄───────┘
    ├── 批量数据导入 (170+ 怪物)
    ├── 端到端集成测试
    ├── API 文档
    └── 性能优化
```

---

## 四、关键架构决策

### 4.1 数据驱动架构
- 使用 `//go:embed` 内嵌 SRD 数据到二进制
- 轻量 Markdown 解析器转为 Go 结构体
- 统一注册中心模式：`RegisterRace()` / `RegisterMonster()` 等
- LRU 缓存热点数据，查询响应 < 10ms

### 4.2 职业特性扩展模式
- 继续使用现有 `FeatureHook` 接口模式
- 每个职业独立 `*FeatureHooks` 结构体
- 通过 `map[ClassID]FeatureHook` 支持多职业
- 钩子在战斗/检定/法术时自动调用

### 4.3 怪物实例化模式
- `MonsterStatBlock` = 模板定义（不可变，从数据加载）
- `Enemy` = 实例化对象（可变状态，HP、状态效果等）
- 创建时深拷贝模板数据到 Enemy

### 4.4 多职业法术位计算
- 独立追踪每个职业的已知/准备法术
- 共享法术位池（根据多职业施法者表计算）
- Pact Magic（邪术师）独立追踪但可跨池使用高环位

---

## 五、估算代码量增长

| 阶段 | 新增文件 | 修改文件 | 新增行数（估） |
|------|----------|----------|----------------|
| 一 | ~12 | ~6 | ~3,500 |
| 二 | ~8 | ~10 | ~5,000 |
| 三 | ~10 | ~5 | ~4,500 |
| 四 | ~12 | ~4 | ~3,000 |
| 五 | ~15 | ~6 | ~4,000 |
| 六 | ~8 | ~15 | ~2,000 |
| **合计** | **~65** | **~46** | **~22,000** |

最终项目预计约 **36,500 行** Go 代码。

---

## 六、关键文件清单

### 需新增的核心文件（~35 个）

```
pkg/model/
├── monster.go          # 怪物数据块结构体
├── creature.go         # 生物类型枚举
├── background.go       # 背景系统
├── feat.go             # 专长系统
├── spelleffect.go      # 法术效果
├── magicitem.go        # 魔法物品
├── exploration.go      # 旅行探索
├── social.go           # 社交互动
├── lifestyle.go        # 生活方式
├── mount.go            # 坐骑载具
├── trap.go             # 陷阱
├── poison.go           # 毒药
├── curse.go            # 诅咒
├── environment.go      # 环境效果
├── crafting.go         # 装备制作
└── weaponmastery.go    # 武器掌控类型

pkg/data/
├── loader.go           # 数据加载器
├── embed.go            # 内嵌数据
├── registry.go         # 数据注册中心
├── monsters.go         # 5 个示例怪物
├── backgrounds.go      # 4 个背景
├── feats.go            # 5 个专长
├── spells.go           # 50 个法术
├── weapons.go          # 30+ 武器
├── armors.go           # 16 护甲
├── magicitems.go       # 30 魔法物品
├── mounts.go           # 坐骑数据
├── traps.go            # 10 陷阱
├── poisons.go          # 10 毒药
├── lifestyles.go       # 7 生活方式
└── import/             # 批量导入工具

pkg/rules/
├── feats.go            # 专长效果计算
├── multiclass.go       # 多职业规则
├── background.go       # 背景应用
├── weaponmastery.go    # 武器掌控应用
├── spelleffects.go     # 法术效果执行
├── magicitem.go        # 魔法物品激活
├── exploration.go      # 旅行计算
├── social.go           # NPC 态度
├── lifestyle.go        # 开销计算
├── death.go            # 死亡豁免
├── environment.go      # 环境效果
├── crafting.go         # 制作计算
└── exhaustion.go       # 力竭效果

pkg/engine/
├── monster.go          # 怪物引擎
├── feat.go             # 专长 API
├── exploration.go      # 探索 API
├── social.go           # 社交 API
├── lifestyle.go        # 生活方式 API
├── mount.go            # 坐骑 API
├── trap.go             # 陷阱 API
├── poison.go           # 毒药 API
├── curse.go            # 诅咒 API
├── environment.go      # 环境 API
├── crafting.go         # 制作 API
└── data/cache.go       # 数据缓存
```

### 需修改的核心文件（~15 个）

```
pkg/model/
├── actor.go            # 添加 CreatureType、SpeedTypes、Inspiration
├── character.go        # 添加 BackgroundID、Feats 字段
├── class.go            # 添加各职业 State 结构体
├── classfeatures.go    # 扩展 11 个职业 FeatureHooks
├── equipment.go        # 添加 Weapon Mastery 字段
├── spell.go            # 完善施法者状态

pkg/data/
├── races.go            # 完善种族子种族和特性

pkg/rules/
├── calculator.go       # 添加无甲防御计算
├── attack.go           # 集成武器掌控

pkg/engine/
├── actor.go            # 集成背景/专长/多职业/死亡豁免
├── combat.go           # 集成职业特性钩子触发
├── spell.go            # 集成法术效果/专注/仪式
├── inventory.go        # 集成魔法物品调音
├── config.go           # 添加 DataLoader 配置
```

---

## 七、验证策略

### 7.1 单元测试
- 每个新增 `pkg/rules/` 文件对应 `*_test.go`
- 每个新增 `pkg/engine/` 文件对应 `*_test.go`
- 数据验证：确保所有 SRD 数据与原始文档一致

### 7.2 端到端测试
- `e2e_character_creation_test.go`: 创建完整角色（种族+职业+背景+专长+装备）
- `e2e_combat_test.go`: 多角色完整战斗流程（先攻→回合→攻击→法术→状态→结束）
- `e2e_multiclass_test.go`: 多职业升级验证（属性要求→法术位→特性）
- `e2e_spellcasting_test.go`: 施法全流程（准备→施放→专注→高环升级→仪式）
- `e2e_monster_test.go`: 怪物加载→创建→战斗→死亡

### 7.3 验证命令
```bash
# 运行所有测试
go test ./... -v

# 运行特定阶段测试
go test ./pkg/rules/... -v
go test ./pkg/engine/... -v -run "E2E"

# 检查代码覆盖
go test ./... -cover

# 代码质量检查
go vet ./...
go fmt ./...
```

---

## 八、执行建议

1. **严格按阶段顺序执行**，每个阶段完成后运行全量测试确保无回归
2. **阶段一优先**：数据加载器和怪物系统是后续所有阶段的基础
3. **数据与逻辑分离**：先完成数据结构，再填充数据，最后实现逻辑
4. **测试驱动开发**：每个新功能先写测试用例，再实现代码
5. **渐进式集成**：每完成一个子系统，立即与现有系统整合并测试
