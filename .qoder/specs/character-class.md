# D&D 5e 职业系统重构计划

## Context

当前 `pkg/model/character.go` 中的 `ClassLevel` 设计过于简单,仅使用字符串存储职业名称,缺乏类型安全和扩展性。D&D 5e 的职业系统非常复杂,包含12种官方职业,每种职业有独特的生命骰、熟练项、豁免、技能选择和职业特性。需要将职业逻辑从 character.go 分离出来,建立类型安全的枚举和完整的数据定义框架,以战士(Fighter)为例实现部分职业特性跟踪结构,并设计职业系统与战斗、法术等其他系统的交互机制。

## 实现步骤

### 步骤 1: 创建 `pkg/model/class.go` - 核心类型和枚举

**内容:**
1. **ClassID 枚举** (string 类型,遵循现有 Alignment 模式)
   - 12个常量,值使用**中文名称**:
     - `ClassBarbarian ClassID = "野蛮人"`
     - `ClassBard ClassID = "吟游诗人"`
     - `ClassCleric ClassID = "牧师"`
     - `ClassDruid ClassID = "德鲁伊"`
     - `ClassFighter ClassID = "战士"`
     - `ClassMonk ClassID = "武僧"`
     - `ClassPaladin ClassID = "圣武士"`
     - `ClassRanger ClassID = "游侠"`
     - `ClassRogue ClassID = "游荡者"`
     - `ClassSorcerer ClassID = "术士"`
     - `ClassWarlock ClassID = "邪术师"`
     - `ClassWizard ClassID = "法师"`

2. **辅助函数**
   - `(id ClassID) String() string` - 返回中文名称
   - `(id ClassID) IsValid() bool` - 验证是否为有效职业
   - `AllClasses() []ClassID` - 返回所有12种职业

3. **重构 ClassLevel**
   ```go
   type ClassLevel struct {
       Class    ClassID  `json:"class"`
       Level    int      `json:"level"`
       Features []string `json:"features,omitempty"` // 获得的职业特性
   }
   ```

4. **战士专属类型**
   - `FightingStyle` 枚举(中文值): 箭术、防御、对决、巨武器战斗、守护、双武器战斗
   - `MartialArchetype` 枚举(中文值): 勇士、战斗大师、奥法骑士
   - `FighterFeatures` 结构体:
     ```go
     type FighterFeatures struct {
         SelectedFightingStyle FightingStyle      // 选择的战斗风格
         SelectedArchetype     MartialArchetype   // 选择的武术范型
         SecondWindMax       int                // 回气最大使用次数(随等级)
         SecondWindUsed      int                // 回气已使用次数
         ActionSurgeMax      int                // 动作如潮最大使用次数
         ActionSurgeUsed     int                // 动作如潮已使用次数
         ExtraAttacks        int                // 额外攻击次数
         IndomitableMax      int                // 不屈最大使用次数
         IndomitableUsed     int                // 不屈已使用次数
     }
     ```

5. **ClassState 接口** - 用于扩展其他职业的特性跟踪
   ```go
   type ClassState interface {
       ClassID() ClassID
   }
   ```

### 步骤 2: 创建 `pkg/data/classes.go` - 职业定义注册表

**内容:**
1. **ClassDefinition 结构体** (遵循 races.go 模式)
   ```go
   type ClassDefinition struct {
       ID                  model.ClassID       // 职业ID
       Name                string              // 中文名(与ClassID相同)
       HitDie              int                 // 生命骰: 6/8/10/12
       PrimaryAbilities    []model.Ability     // 主要属性(按重要性排序)
       SavingThrows        []model.Ability     // 豁免熟练(2项)
       SkillChoices        []model.Skill       // 可选技能列表
       NumberOfSkills      int                 // 1级可选技能数量
       ArmorProficiencies  []model.ArmorType   // 护甲熟练
       WeaponProficiencies []string            // 武器熟练(Simple/Martial/具体武器)
       ToolProficiencies   []string            // 工具熟练
       SpellcastingAbility model.Ability       // 施法属性(非施法职业为空)
       CasterType          CasterType          // 施法者类型(Full/Half/Third/None)
       Description         string              // 中文描述
   }
   ```

2. **CasterType 枚举** (在 class.go 中定义)
   - `CasterTypeNone` - 非施法者(战士、野蛮人等)
   - `CasterTypeFull` - 全施法者(法师、牧师、德鲁伊、术士、吟游诗人)
   - `CasterTypeHalf` - 半施法者(圣武士、游侠)
   - `CasterTypeThird` - 1/3施法者(奥法骑士、诡术师)

3. **Classes 注册表**: `map[model.ClassID]*ClassDefinition` 包含全部12种职业的完整定义

4. **辅助函数**
   - `GetClass(id model.ClassID) *ClassDefinition`
   - `GetClassID(name string) (model.ClassID, error)` - 支持中英文查找
   - `GetClassNames() []string`

5. **战士特性数据**
   - `FighterFeaturesByLevel`: map[int][]string 定义每级获得的特性名称
   - 战斗风格描述映射
   - 武术范型描述映射

### 步骤 3: 创建 `pkg/model/classfeatures.go` - 职业特性交互接口

**目的:** 定义职业特性如何与战斗、法术等系统交互的接口和数据结构

**内容:**
1. **FeatureHook 接口** - 职业特性钩子
   ```go
   type FeatureHook interface {
       // OnAttackRoll 攻击掷骰时调用,可修改攻击加值
       OnAttackRoll(ctx *AttackContext)
       
       // OnDamageCalc 伤害计算时调用,可修改伤害值
       OnDamageCalc(ctx *DamageContext)
       
       // OnACCalc 护甲等级计算时调用,可修改AC
       OnACCalc(ctx *ACContext)
       
       // OnSpellCalc 法术计算时调用,可修改法术DC/攻击加值
       OnSpellCalc(ctx *SpellContext)
       
       // GetAvailableActions 返回可用的特殊动作
       GetAvailableActions() []ActionTemplate
       
       // OnShortRest 短休时调用,恢复资源
       OnShortRest()
       
       // OnLongRest 长休时调用,恢复资源
       OnLongRest()
   }
   ```

2. **上下文结构体**
   ```go
   type AttackContext struct {
       BaseBonus      int           // 基础攻击加值
       Bonus          int           // 额外加值(可修改)
       WeaponType     WeaponType    // 武器类型
       IsRanged       bool          // 是否远程
       CriticalRange  int           // 暴击范围(默认20,勇士范型可改为19或18)
   }
   
   type DamageContext struct {
       BaseDamage     int           // 基础伤害
       Bonus          int           // 额外伤害(可修改)
       DamageType     DamageType    // 伤害类型
   }
   
   type ACContext struct {
       BaseAC         int           // 基础AC
       Bonus          int           // 额外AC(可修改)
       HasShield      bool          // 是否持盾
       HasArmor       bool          // 是否着装护甲
   }
   
   type SpellContext struct {
       SpellSaveDC    int           // 法术豁免DC
       SpellAttackBonus int         // 法术攻击加值
       SpellSlots     [][]int       // 法术位
   }
   ```

3. **ActionTemplate 结构** - 特殊动作模板
   ```go
   type ActionTemplate struct {
       Type        ActionType
       Name        string
       IsBonus     bool          // 是否附赠动作
       UsesPerRest int           // 每rest可用次数
       CurrentUses int           // 当前剩余次数
   }
   ```

4. **FighterFeatureHooks 实现** - 战士特性钩子示例
   ```go
   type FighterFeatureHooks struct {
       Features *FighterFeatures
   }
   // 实现 FeatureHook 接口:
   // - OnAttackRoll: 箭术+2远程, 对决+2单手近战伤害
   // - OnACCalc: 防御+1AC(着装时)
   // - GetAvailableActions: 回气(附赠), 动作如潮(自由)
   // - OnShortRest/OnLongRest: 恢复回气/动作如潮/不屈
   ```

### 步骤 4: 修改 `pkg/model/character.go`

- 删除 `ClassLevel` 结构体定义(已移至 class.go)
- `PlayerCharacter` 添加特性钩子字段:
  ```go
  type PlayerCharacter struct {
      // ... 现有字段 ...
      
      // 职业特性系统
      FeatureHooks map[ClassID]FeatureHook `json:"-"` // 运行时特性钩子
      FighterState *FighterFeatures        `json:"fighter_state,omitempty"`
  }
  ```

### 步骤 5: 修改 `pkg/rules/constants.go`

- 删除 `HitDiceByClass` map(生命骰信息已整合到 ClassDefinition)

### 步骤 6: 修改 `pkg/engine/actor.go`

**需要更新的位置:**
1. **ClassInfo 结构体** (约第100行): `ClassName string` → `Class model.ClassID`
2. **CreatePC 函数** (约第339行): 
   - 使用 `data.GetClassID()` 转换并验证职业名称
   - 根据职业类型初始化 FeatureHooks 和对应 State(如 FighterState)
   - 自动填充 ClassLevel.Features
3. **LevelUp 函数** (约第800行): 
   - 从 `data.GetClass()` 获取 HitDie 替代 `rules.HitDiceByClass`
   - 升级时更新 ClassLevel.Features
   - 更新职业特性状态(如 ExtraAttacks, ActionSurgeMax 等)
4. **calculateMaxHP 函数** (约第1103行): 使用 `cl.Class` (ClassID) 查询职业定义获取 HitDie
5. **playerCharacterToInfo 函数** (约第1168行): 更新 ClassInfo 映射

### 步骤 7: 修改 `pkg/rules/attack.go` - 集成职业特性

**修改 CalcAttachBonus 函数:**
```go
func CalcAttachBonus(attacker any, weaponType model.WeaponType, isRanged bool) int {
    // ... 现有逻辑获取 profBonus 和 abilityMod ...
    
    bonus := profBonus + abilityMod
    
    // 应用职业特性加值
    if pc, ok := attacker.(*model.PlayerCharacter); ok {
        if hook, exists := pc.FeatureHooks[model.ClassFighter]; exists {
            ctx := &model.AttackContext{
                BaseBonus:     bonus,
                Bonus:         0,
                WeaponType:    weaponType,
                IsRanged:      isRanged,
                CriticalRange: 20,
            }
            hook.OnAttackRoll(ctx)
            bonus = ctx.BaseBonus + ctx.Bonus
        }
    }
    
    return bonus
}
```

### 步骤 8: 修改 `pkg/rules/calculator.go` - 添加职业辅助函数

**新增函数:**
```go
// GetSpellcastingAbilityForClass 根据职业获取施法属性
func GetSpellcastingAbilityForClass(classID model.ClassID) model.Ability {
    classDef := data.GetClass(classID)
    if classDef == nil {
        return ""
    }
    return classDef.SpellcastingAbility
}

// GetCasterLevel 根据总等级和职业计算等效施法者等级
func GetCasterLevel(classes []model.ClassLevel) int {
    casterLevel := 0
    for _, cl := range classes {
        classDef := data.GetClass(cl.Class)
        if classDef == nil {
            continue
        }
        switch classDef.CasterType {
        case CasterTypeFull:
            casterLevel += cl.Level
        case CasterTypeHalf:
            casterLevel += cl.Level / 2
        case CasterTypeThird:
            casterLevel += cl.Level / 3
        }
    }
    return casterLevel
}
```

### 步骤 9: 修改 `pkg/engine/combat.go` - 集成特性钩子

**修改 ExecuteAttack 函数:**
- 在调用 `CalcAttachBonus` 后,检查是否有职业特性修改暴击范围
- 支持 Extra Attack: 检查战士/游侠/野蛮人等级,允许多次攻击

**修改 ExecuteAction 函数:**
- 添加 `ActionSecondWind` case: 调用 FighterFeatureHooks 恢复HP
- 添加 `ActionActionSurge` case: 授予额外动作

### 步骤 10: 修改 `pkg/engine/spell.go` - 集成职业施法规则

**修改法术位计算:**
- 使用 `GetCasterLevel()` 计算等效施法者等级
- 根据 CasterType 和施法者等级确定法术位数量

**修改 PrepareSpells 函数:**
- 验证职业是否为施法者
- 根据职业限制可准备的法术列表

### 步骤 11: 更新测试

- 修改 `pkg/engine/actor_test.go` 中的职业名称引用(使用中文)
- 添加 ClassID 枚举测试
- 添加职业定义验证测试
- 添加战士特性跟踪测试
- 添加职业特性钩子测试

## 关键文件

| 文件 | 操作 | 说明 |
|------|------|------|
| `pkg/model/class.go` | 新建 | ClassID枚举(中文)、ClassLevel、战士特性类型、CasterType |
| `pkg/model/classfeatures.go` | 新建 | FeatureHook接口、上下文结构体、战士钩子实现 |
| `pkg/data/classes.go` | 新建 | 12种职业定义、施法者类型、战士特性数据 |
| `pkg/model/character.go` | 修改 | 删除ClassLevel,添加FeatureHooks和FighterState |
| `pkg/rules/constants.go` | 修改 | 删除HitDiceByClass |
| `pkg/rules/attack.go` | 修改 | CalcAttachBonus集成特性钩子 |
| `pkg/rules/calculator.go` | 修改 | 添加施法属性/等级计算函数 |
| `pkg/engine/actor.go` | 修改 | 更新所有职业相关引用,初始化特性系统 |
| `pkg/engine/combat.go` | 修改 | 集成特性钩子到攻击和动作系统 |
| `pkg/engine/spell.go` | 修改 | 基于职业计算法术位 |
| `pkg/engine/actor_test.go` | 修改 | 更新测试用例 |

## 职业系统交互设计

### 与战斗系统交互
1. **攻击加值修改**: 战斗风格(箭术+2远程,对决+2单手伤害)
2. **暴击范围修改**: 勇士范型(19-20,18-20)
3. **额外攻击**: 战士5级/11级/20级(2/3/4次)
4. **特殊动作**: 回气(附赠,1d10+等级HP),动作如潮(额外动作)
5. **AC修改**: 防御风格(+1AC着装时)

### 与法术系统交互
1. **施法属性确定**: 每种职业有固定施法属性(INT/WIS/CHA)
2. **法术位计算**: 根据CasterType(Full/Half/Third)和等效施法者等级
3. **法术准备限制**: 准备施法者(牧师/法师)vs 已知法术(术士/吟游诗人)
4. **法术DC计算**: 8+熟练+施法属性修正

### 与休息系统交互
1. **短休恢复**: 回气、动作如潮、生命骰
2. **长休恢复**: 所有资源、力竭减少

## 验证方法

1. 运行 `go build ./...` 确保编译通过
2. 运行 `go test ./pkg/...` 确保所有测试通过
3. 验证 ClassID 枚举包含全部12种职业,值为中文
4. 验证战士职业定义的 HitDie 为 10, 豁免为 STR/CON
5. 验证创建战士角色时能正确设置职业特性和 FighterState
6. 验证战斗风格能正确影响攻击加值和AC
7. 验证法术位能根据职业类型正确计算
