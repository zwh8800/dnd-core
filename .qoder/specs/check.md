# 为 pkg/engine 包添加 Go 文档注释

## Context

`pkg/engine` 包是 D&D 5e 游戏引擎的核心，包含约 119 个导出函数，但目前仅有 16 个函数（约 13.4%）拥有完整的 godoc 注释。需要为剩余 103 个导出函数添加完整的中文文档注释，以提高代码可维护性和 API 可发现性。

## 文档规范

基于代码库现有模式，采用以下格式：

```go
// FunctionName 功能简述
// 参数:
//
//	ctx - 上下文
//	req - 请求参数，包含XXX等字段
//
// 返回:
//
//	*ResultType - 结果描述
//	error - 可能返回 ErrNotFound（角色不存在）等错误
```

## 需要添加注释的文件及函数清单

### 1. actor.go - 角色管理（14 个函数）
- `CreatePC` - 创建玩家角色
- `CreateNPC` - 创建非玩家角色
- `CreateEnemy` - 创建敌人/怪物
- `CreateCompanion` - 创建同伴
- `GetActor` - 获取任意角色信息
- `GetPC` - 获取玩家角色完整信息
- `UpdateActor` - 更新角色状态
- `RemoveActor` - 移除角色
- `ListActors` - 列出角色列表
- `AddExperience` - 添加经验值
- `LevelUp` - 角色升级
- `ShortRest` - 短休
- `StartLongRest` - 开始长休
- `EndLongRest` - 结束长休

### 2. check.go - 检定系统（2 个函数）
- `GetSkillAbility` - 获取技能对应属性
- `GetPassivePerception` - 获取被动察觉值

### 3. dice.go - 掷骰系统（5 个函数）
- `Roll` - 执行掷骰
- `RollAdvantage` - 优势掷骰
- `RollDisadvantage` - 劣势掷骰
- `RollAbility` - 属性掷骰（4d6取前3高）
- `RollHitDice` - 生命骰掷骰

### 4. combat.go - 战斗系统（3 个函数）
- `ExecuteAttack` - 执行攻击检定和伤害
- `ExecuteDamage` - 应用伤害
- `ExecuteHealing` - 应用治疗

### 5. spell.go - 法术系统（11 个函数）
- `CastSpell` - 施放法术
- `GetSpellSlots` - 获取法术位状态
- `PrepareSpells` - 准备法术
- `LearnSpell` - 学习法术
- `ConcentrationCheck` - 专注检定
- `EndConcentration` - 结束专注
- `CastSpellRitual` - 仪式施法
- `GetPactMagicSlots` - 获取契约魔法位
- `RestorePactMagicSlots` - 恢复契约魔法位
- `IsConcentrating` - 检查是否专注
- `GetConcentrationSpell` - 获取当前专注法术

### 6. inventory.go - 物品管理（8 个函数）
- `AddItem` - 添加物品到背包
- `RemoveItem` - 从背包移除物品
- `GetInventory` - 获取完整背包
- `EquipItem` - 装备物品
- `UnequipItem` - 卸下装备
- `GetEquipment` - 获取已装备物品
- `AttuneItem` - 调谐魔法物品
- `TransferItem` - 角色间转移物品

### 7. scene.go - 场景管理（11 个函数）
- `CreateScene` - 创建场景
- `GetScene` - 获取场景信息
- `UpdateScene` - 更新场景
- `DeleteScene` - 删除场景
- `ListScenes` - 列出所有场景
- `SetCurrentScene` - 设置当前场景
- `GetCurrentScene` - 获取当前场景
- `AddSceneConnection` - 添加场景连接
- `RemoveSceneConnection` - 移除场景连接
- `MoveActorToScene` - 移动角色到场景
- `GetSceneActors` - 获取场景中的角色

### 8. quest.go - 任务管理（10 个函数）
- `CreateQuest` - 创建任务
- `GetQuest` - 获取任务信息
- `ListQuests` - 列出任务
- `AcceptQuest` - 接受任务
- `UpdateQuestObjective` - 更新任务目标进度
- `CompleteQuest` - 完成任务
- `FailQuest` - 任务失败
- `DeleteQuest` - 删除任务
- `GetActorQuests` - 获取角色的任务
- `GetQuestGiverQuests` - 获取NPC发布的任务

### 9. feat.go - 专长系统（3 个函数）
- `SelectFeat` - 选择专长
- `ListFeats` - 列出可用专长
- `RemoveFeat` - 移除专长

### 10. social.go - 社交系统（2 个函数）
- `InteractWithNPC` - 与NPC互动
- `GetNPCAttitude` - 获取NPC态度

### 11. poison.go - 毒素系统（3 个函数）
- `ApplyPoison` - 涂抹毒素
- `ResolvePoisonEffect` - 结算毒素效果
- `RemovePoison` - 清除毒素

### 12. curse.go - 诅咒系统（3 个函数）
- `CurseActor` - 施加诅咒
- `RemoveCurse` - 移除诅咒
- `GetCurses` - 获取活跃诅咒

### 13. environment.go - 环境效果（2 个函数）
- `SetEnvironment` - 设置环境条件
- `ResolveEnvironmentalDamage` - 结算环境伤害

### 14. exploration.go - 探索旅行（4 个函数）
- `StartTravel` - 开始旅行
- `AdvanceTravel` - 推进旅行
- `Forage` - 觅食检定
- `Navigate` - 导航检定

### 15. lifestyle.go - 生活方式（2 个函数）
- `SetLifestyle` - 设置生活方式
- `AdvanceGameTime` - 推进游戏时间

### 16. crafting.go - 制作系统（3 个函数）
- `StartCrafting` - 开始制作
- `AdvanceCrafting` - 推进制作进度
- `CompleteCrafting` - 完成制作

### 17. mount.go - 坐骑系统（3 个函数）
- `MountCreature` - 骑乘生物
- `Dismount` - 下马
- `CalculateMountSpeed` - 计算骑乘速度

### 18. trap.go - 陷阱系统（4 个函数）
- `PlaceTrap` - 放置陷阱
- `DetectTrap` - 检测陷阱
- `DisarmTrap` - 解除陷阱
- `TriggerTrap` - 触发陷阱

### 19. monster.go - 怪物加载（3 个函数）
- `LoadMonster` - 从模板加载怪物
- `CreateEnemyFromStatBlock` - 从数据块创建敌人
- `GetMonsterActions` - 获取怪物可用动作

### 20. state.go - 游戏状态（2 个函数）
- `GetStateSummary` - 获取游戏状态摘要
- `GetActorSheet` - 获取角色卡

### 21. game.go - 游戏会话（4 个函数）
- `NewGame` - 创建新游戏
- `LoadGame` - 加载游戏
- `SaveGame` - 保存游戏
- `DeleteGame` - 删除游戏

### 22. phase.go - 阶段管理（1 个函数）
- `GetAllowedOperations` - 获取当前阶段允许的操作

### 23. engine.go - 引擎生命周期（1 个函数）
- `NewTestEngine` - 创建测试引擎

## 实施策略

由于文件数量多，将采用批量处理方式：

1. **第一批**：核心系统（actor.go, dice.go, combat.go, spell.go）
2. **第二批**：世界交互（scene.go, quest.go, exploration.go, environment.go）
3. **第三批**：角色增强（inventory.go, feat.go, lifestyle.go, crafting.go）
4. **第四批**：特殊系统（social.go, poison.go, curse.go, mount.go, trap.go, monster.go）
5. **第五批**：游戏管理（game.go, state.go, phase.go, engine.go）

每个文件的处理步骤：
1. 读取文件内容
2. 识别所有导出函数
3. 为每个函数添加符合规范的中文 godoc 注释
4. 验证注释格式正确

## 验证方法

1. 运行 `go doc pkg/engine` 验证文档生成
2. 运行 `go vet ./pkg/engine/...` 检查代码质量
3. 运行 `go test ./pkg/engine/...` 确保测试通过
4. 随机抽查几个函数的文档完整性
