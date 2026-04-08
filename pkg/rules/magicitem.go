package rules

import (
	"fmt"

	"github.com/zwh8800/dnd-core/pkg/dice"
	"github.com/zwh8800/dnd-core/pkg/model"
)

// AttuneItem 尝试调音一个魔法物品
func AttuneItem(item *model.Item, inventory *model.Inventory) error {
	if item.Attunement == "" {
		return fmt.Errorf("item %s does not require attunement", item.Name)
	}

	attunedCount := 0
	for _, invItem := range inventory.Items {
		if invItem.Attuned {
			attunedCount++
		}
	}
	if attunedCount >= 3 {
		return fmt.Errorf("character already has 3 attuned items")
	}

	if item.Attuned {
		return fmt.Errorf("item %s is already attuned to someone else", item.Name)
	}

	item.Attuned = true
	return nil
}

// UnattuneItem 取消调音
func UnattuneItem(item *model.Item) error {
	if !item.Attuned {
		return fmt.Errorf("item %s is not attuned", item.Name)
	}

	item.Attuned = false
	return nil
}

// GetAttunedItemCount 获取库存中已调音的物品数量
func GetAttunedItemCount(inventory *model.Inventory) int {
	count := 0
	for _, item := range inventory.Items {
		if item.Attuned {
			count++
		}
	}
	return count
}

// MagicItemUseResult 魔法物品使用结果
type MagicItemUseResult struct {
	ItemName string   `json:"item_name"`
	Messages []string `json:"messages"`
}

// UseMagicItem 使用魔法物品
func UseMagicItem(item *model.Item, user *model.PlayerCharacter, targets []model.ID) (*MagicItemUseResult, error) {
	result := &MagicItemUseResult{
		ItemName: item.Name,
		Messages: []string{},
	}

	if item.Consumable {
		return useConsumable(item, user, result)
	}

	if item.Charges > 0 {
		return useChargedItem(item, user, targets, result)
	}

	result.Messages = append(result.Messages, fmt.Sprintf("%s provides passive effects", item.Name))
	return result, nil
}

// useConsumable 使用消耗品
func useConsumable(item *model.Item, user *model.PlayerCharacter, result *MagicItemUseResult) (*MagicItemUseResult, error) {
	roller := dice.New(0)

	switch item.ID {
	case "potion-healing":
		rollResult, _ := roller.Roll("2d4")
		healing := rollResult.Total + 2
		result.Messages = append(result.Messages, fmt.Sprintf("Drinking %s heals %d HP", item.Name, healing))
	case "potion-greater-healing":
		rollResult, _ := roller.Roll("4d4")
		healing := rollResult.Total + 4
		result.Messages = append(result.Messages, fmt.Sprintf("Drinking %s heals %d HP", item.Name, healing))
	case "potion-superior-healing":
		rollResult, _ := roller.Roll("8d4")
		healing := rollResult.Total + 8
		result.Messages = append(result.Messages, fmt.Sprintf("Drinking %s heals %d HP", item.Name, healing))
	case "potion-supreme-healing":
		rollResult, _ := roller.Roll("10d4")
		healing := rollResult.Total + 20
		result.Messages = append(result.Messages, fmt.Sprintf("Drinking %s heals %d HP", item.Name, healing))
	case "acid-vial":
		result.Messages = append(result.Messages, fmt.Sprintf("Throwing %s deals 2d6 acid damage", item.Name))
	case "alchemists-fire-flask":
		result.Messages = append(result.Messages, fmt.Sprintf("Throwing %s deals 1d4 fire damage per turn", item.Name))
	case "antitoxin-vial":
		result.Messages = append(result.Messages, fmt.Sprintf("Drinking %s grants immunity to poisoned for 1 hour", item.Name))
	case "holy-water-flask":
		result.Messages = append(result.Messages, fmt.Sprintf("Throwing %s deals 2d6 radiant damage to fiends and undead", item.Name))
	case "poison-basic-vial":
		result.Messages = append(result.Messages, fmt.Sprintf("Applying %s to weapon deals 1d4 poison damage", item.Name))
	default:
		if item.Effect != "" {
			result.Messages = append(result.Messages, fmt.Sprintf("Using %s: %s", item.Name, item.Effect))
		} else {
			result.Messages = append(result.Messages, fmt.Sprintf("Using consumable %s", item.Name))
		}
	}

	return result, nil
}

// useChargedItem 使用充能物品
func useChargedItem(item *model.Item, user *model.PlayerCharacter, targets []model.ID, result *MagicItemUseResult) (*MagicItemUseResult, error) {
	if item.Charges <= 0 {
		return nil, fmt.Errorf("item %s has no charges remaining", item.Name)
	}

	item.Charges--

	switch item.ID {
	case "wand-of-magic-detection":
		result.Messages = append(result.Messages, fmt.Sprintf("Using %s to cast Detect Magic", item.Name))
	case "brooch-of-shielding":
		result.Messages = append(result.Messages, fmt.Sprintf("Using %s to cast Magic Missile", item.Name))
	case "boots-of-speed":
		result.Messages = append(result.Messages, fmt.Sprintf("Using %s to gain Haste for 10 minutes", item.Name))
	case "ring-of-invisibility":
		result.Messages = append(result.Messages, fmt.Sprintf("Using %s to cast Invisibility", item.Name))
	case "staff-of-fire":
		result.Messages = append(result.Messages, fmt.Sprintf("Using %s to cast a fire spell", item.Name))
	case "wand-of-fireballs":
		result.Messages = append(result.Messages, fmt.Sprintf("Using %s to cast Fireball (5th level)", item.Name))
	case "staff-of-power":
		result.Messages = append(result.Messages, fmt.Sprintf("Using %s to cast a powerful spell", item.Name))
	case "robe-of-stars":
		result.Messages = append(result.Messages, fmt.Sprintf("Using %s to fire star missiles", item.Name))
	case "ring-of-wish":
		result.Messages = append(result.Messages, fmt.Sprintf("Using %s to cast Wish", item.Name))
	case "staff-of-the-magi":
		result.Messages = append(result.Messages, fmt.Sprintf("Using %s to cast a powerful spell", item.Name))
	default:
		result.Messages = append(result.Messages, fmt.Sprintf("Using charged item %s (%d charges remaining)", item.Name, item.Charges))
	}

	return result, nil
}

// RechargeMagicItems 在黎明时恢复魔法物品的充能
func RechargeMagicItems(items []*model.Item) []string {
	messages := []string{}
	for _, item := range items {
		if item.MaxCharges > 0 && item.Recharge == "dawn" {
			if item.Charges < item.MaxCharges {
				oldCharges := item.Charges
				item.Charges = item.MaxCharges
				messages = append(messages, fmt.Sprintf("%s recharged from %d to %d charges", item.Name, oldCharges, item.Charges))
			}
		}
	}
	return messages
}

// GetMagicItemBonus 获取魔法物品的加值
func GetMagicItemBonus(item *model.Item) int {
	return item.MagicBonus
}

// HasMagicEffect 检查物品是否具有某个魔法效果
func HasMagicEffect(item *model.Item, effect string) bool {
	for _, e := range item.MagicEffects {
		if e == effect {
			return true
		}
	}
	return false
}
