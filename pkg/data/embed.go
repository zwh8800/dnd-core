package data

// InitDefaultData 初始化所有默认数据到注册中心
func InitDefaultData() {
	// 初始化种族
	for _, race := range Races {
		GlobalRegistry.RegisterRace(race)
	}

	// 初始化背景
	for _, bg := range Backgrounds {
		GlobalRegistry.RegisterBackground(bg)
	}

	// 初始化专长
	for _, feat := range Feats {
		GlobalRegistry.RegisterFeat(feat)
	}

	// 初始化怪物
	for _, monster := range MonsterStatBlocks {
		GlobalRegistry.RegisterMonster(monster)
	}
}
