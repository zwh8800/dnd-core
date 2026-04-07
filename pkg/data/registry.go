package data

import (
	"fmt"
	"sync"

	"github.com/zwh8800/dnd-core/pkg/model"
)

// DataRegistry 是统一的数据注册中心
type DataRegistry struct {
	mu sync.RWMutex

	// 种族数据
	races map[string]*RaceDefinition

	// 职业数据
	classes map[model.ClassID]*ClassDefinition

	// 背景数据
	backgrounds map[string]*model.BackgroundDefinition

	// 专长数据
	feats map[string]*model.FeatDefinition

	// 怪物模板数据
	monsters map[string]*model.MonsterStatBlock

	// 法术数据
	spells map[string]*model.Spell

	// 武器数据
	weapons map[string]*model.Item

	// 护甲数据
	armors map[string]*model.Item

	// 魔法物品数据
	magicItems map[string]*model.Item
}

// GlobalRegistry 全局数据注册中心实例
var GlobalRegistry = NewDataRegistry()

// NewDataRegistry 创建新的数据注册中心
func NewDataRegistry() *DataRegistry {
	return &DataRegistry{
		races:       make(map[string]*RaceDefinition),
		classes:     make(map[model.ClassID]*ClassDefinition),
		backgrounds: make(map[string]*model.BackgroundDefinition),
		feats:       make(map[string]*model.FeatDefinition),
		monsters:    make(map[string]*model.MonsterStatBlock),
		spells:      make(map[string]*model.Spell),
		weapons:     make(map[string]*model.Item),
		armors:      make(map[string]*model.Item),
		magicItems:  make(map[string]*model.Item),
	}
}

// RegisterRace 注册种族数据
func (r *DataRegistry) RegisterRace(race *RaceDefinition) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.races[race.Name]; exists {
		return fmt.Errorf("race already registered: %s", race.Name)
	}
	r.races[race.Name] = race
	return nil
}

// GetRace 获取种族数据
func (r *DataRegistry) GetRace(id string) (*RaceDefinition, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	race, exists := r.races[id]
	return race, exists
}

// ListRaces 列出所有种族
func (r *DataRegistry) ListRaces() []*RaceDefinition {
	r.mu.RLock()
	defer r.mu.RUnlock()
	races := make([]*RaceDefinition, 0, len(r.races))
	for _, race := range r.races {
		races = append(races, race)
	}
	return races
}

// RegisterClass 注册职业数据
func (r *DataRegistry) RegisterClass(class *ClassDefinition) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.classes[class.ID]; exists {
		return fmt.Errorf("class already registered: %s", class.ID)
	}
	r.classes[class.ID] = class
	return nil
}

// GetClass 获取职业数据
func (r *DataRegistry) GetClass(id model.ClassID) (*ClassDefinition, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	class, exists := r.classes[id]
	return class, exists
}

// ListClasses 列出所有职业
func (r *DataRegistry) ListClasses() []*ClassDefinition {
	r.mu.RLock()
	defer r.mu.RUnlock()
	classes := make([]*ClassDefinition, 0, len(r.classes))
	for _, class := range r.classes {
		classes = append(classes, class)
	}
	return classes
}

// RegisterBackground 注册背景数据
func (r *DataRegistry) RegisterBackground(bg *model.BackgroundDefinition) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.backgrounds[string(bg.ID)]; exists {
		return fmt.Errorf("background already registered: %s", bg.ID)
	}
	r.backgrounds[string(bg.ID)] = bg
	return nil
}

// GetBackground 获取背景数据
func (r *DataRegistry) GetBackground(id string) (*model.BackgroundDefinition, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	bg, exists := r.backgrounds[id]
	return bg, exists
}

// ListBackgrounds 列出所有背景
func (r *DataRegistry) ListBackgrounds() []*model.BackgroundDefinition {
	r.mu.RLock()
	defer r.mu.RUnlock()
	bgs := make([]*model.BackgroundDefinition, 0, len(r.backgrounds))
	for _, bg := range r.backgrounds {
		bgs = append(bgs, bg)
	}
	return bgs
}

// RegisterFeat 注册专长数据
func (r *DataRegistry) RegisterFeat(feat *model.FeatDefinition) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.feats[feat.ID]; exists {
		return fmt.Errorf("feat already registered: %s", feat.ID)
	}
	r.feats[feat.ID] = feat
	return nil
}

// GetFeat 获取专长数据
func (r *DataRegistry) GetFeat(id string) (*model.FeatDefinition, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	feat, exists := r.feats[id]
	return feat, exists
}

// ListFeats 列出所有专长
func (r *DataRegistry) ListFeats() []*model.FeatDefinition {
	r.mu.RLock()
	defer r.mu.RUnlock()
	feats := make([]*model.FeatDefinition, 0, len(r.feats))
	for _, feat := range r.feats {
		feats = append(feats, feat)
	}
	return feats
}

// RegisterMonster 注册怪物模板数据
func (r *DataRegistry) RegisterMonster(monster *model.MonsterStatBlock) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.monsters[monster.ID]; exists {
		return fmt.Errorf("monster already registered: %s", monster.ID)
	}
	r.monsters[monster.ID] = monster
	return nil
}

// GetMonster 获取怪物模板数据
func (r *DataRegistry) GetMonster(id string) (*model.MonsterStatBlock, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	monster, exists := r.monsters[id]
	return monster, exists
}

// ListMonsters 列出所有怪物模板
func (r *DataRegistry) ListMonsters() []*model.MonsterStatBlock {
	r.mu.RLock()
	defer r.mu.RUnlock()
	monsters := make([]*model.MonsterStatBlock, 0, len(r.monsters))
	for _, monster := range r.monsters {
		monsters = append(monsters, monster)
	}
	return monsters
}

// RegisterSpell 注册法术数据
func (r *DataRegistry) RegisterSpell(spell *model.Spell) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.spells[spell.ID]; exists {
		return fmt.Errorf("spell already registered: %s", spell.ID)
	}
	r.spells[spell.ID] = spell
	return nil
}

// GetSpell 获取法术数据
func (r *DataRegistry) GetSpell(id string) (*model.Spell, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	spell, exists := r.spells[id]
	return spell, exists
}

// ListSpells 列出所有法术
func (r *DataRegistry) ListSpells() []*model.Spell {
	r.mu.RLock()
	defer r.mu.RUnlock()
	spells := make([]*model.Spell, 0, len(r.spells))
	for _, spell := range r.spells {
		spells = append(spells, spell)
	}
	return spells
}

// RegisterWeapon 注册武器数据
func (r *DataRegistry) RegisterWeapon(weapon *model.Item) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.weapons[weapon.ID.String()]; exists {
		return fmt.Errorf("weapon already registered: %s", weapon.ID)
	}
	r.weapons[weapon.ID.String()] = weapon
	return nil
}

// GetWeapon 获取武器数据
func (r *DataRegistry) GetWeapon(id string) (*model.Item, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	weapon, exists := r.weapons[id]
	return weapon, exists
}

// ListWeapons 列出所有武器
func (r *DataRegistry) ListWeapons() []*model.Item {
	r.mu.RLock()
	defer r.mu.RUnlock()
	weapons := make([]*model.Item, 0, len(r.weapons))
	for _, weapon := range r.weapons {
		weapons = append(weapons, weapon)
	}
	return weapons
}

// RegisterArmor 注册护甲数据
func (r *DataRegistry) RegisterArmor(armor *model.Item) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.armors[armor.ID.String()]; exists {
		return fmt.Errorf("armor already registered: %s", armor.ID)
	}
	r.armors[armor.ID.String()] = armor
	return nil
}

// GetArmor 获取护甲数据
func (r *DataRegistry) GetArmor(id string) (*model.Item, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	armor, exists := r.armors[id]
	return armor, exists
}

// ListArmors 列出所有护甲
func (r *DataRegistry) ListArmors() []*model.Item {
	r.mu.RLock()
	defer r.mu.RUnlock()
	armors := make([]*model.Item, 0, len(r.armors))
	for _, armor := range r.armors {
		armors = append(armors, armor)
	}
	return armors
}

// RegisterMagicItem 注册魔法物品数据
func (r *DataRegistry) RegisterMagicItem(item *model.Item) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.magicItems[item.ID.String()]; exists {
		return fmt.Errorf("magic item already registered: %s", item.ID)
	}
	r.magicItems[item.ID.String()] = item
	return nil
}

// GetMagicItem 获取魔法物品数据
func (r *DataRegistry) GetMagicItem(id string) (*model.Item, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, exists := r.magicItems[id]
	return item, exists
}

// ListMagicItems 列出所有魔法物品
func (r *DataRegistry) ListMagicItems() []*model.Item {
	r.mu.RLock()
	defer r.mu.RUnlock()
	items := make([]*model.Item, 0, len(r.magicItems))
	for _, item := range r.magicItems {
		items = append(items, item)
	}
	return items
}
