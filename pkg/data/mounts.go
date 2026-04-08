package data

import (
	"github.com/zwh8800/dnd-core/pkg/model"
)

// MountData 坐骑数据
type MountData struct {
	ID          model.ID `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Speed       int      `json:"speed"`      // 基础速度（尺）
	CarryCap    float64  `json:"carry_cap"`  // 载重能力（磅）
	Value       int      `json:"value"`      // 价值（铜币）
	IsTrained   bool     `json:"is_trained"` // 是否受过骑乘训练
}

// MountDataList 所有坐骑数据
var MountDataList = []MountData{
	{
		ID:          "horse-riding",
		Name:        "骑乘马",
		Description: "一匹受过骑乘训练的马",
		Speed:       60,
		CarryCap:    480,
		Value:       7500, // 75 gp
		IsTrained:   true,
	},
	{
		ID:          "horse-draft",
		Name:        "挽马",
		Description: "一匹用于拉车的马",
		Speed:       40,
		CarryCap:    540,
		Value:       5000, // 50 gp
		IsTrained:   false,
	},
	{
		ID:          "pony",
		Name:        "矮种马",
		Description: "一匹小型马，适合小体型骑手",
		Speed:       40,
		CarryCap:    225,
		Value:       3000, // 30 gp
		IsTrained:   true,
	},
	{
		ID:          "mastiff",
		Name:        "獒犬",
		Description: "一匹大型犬，可骑乘",
		Speed:       40,
		CarryCap:    195,
		Value:       2500, // 25 gp
		IsTrained:   true,
	},
	{
		ID:          "camel",
		Name:        "骆驼",
		Description: "一头适应沙漠的坐骑",
		Speed:       50,
		CarryCap:    480,
		Value:       5000, // 50 gp
		IsTrained:   true,
	},
	{
		ID:          "elephant",
		Name:        "大象",
		Description: "一头巨大的坐骑",
		Speed:       40,
		CarryCap:    1320,
		Value:       20000, // 200 gp
		IsTrained:   true,
	},
}

// VehicleDataList 交通工具数据
var VehicleDataList = []model.Vehicle{
	{
		ID:            "cart",
		Name:          "双轮马车",
		Description:   "一辆双轮的轻型马车",
		Type:          model.VehicleLand,
		AC:            15,
		HP:            40,
		Speed:         20,
		CargoCapacity: 300,
		CrewRequired:  1,
		CrewCurrent:   0,
		Weight:        200,
		Value:         1500, // 15 gp
	},
	{
		ID:            "carriage",
		Name:          "四轮马车",
		Description:   "一辆四轮的重型马车",
		Type:          model.VehicleLand,
		AC:            15,
		HP:            50,
		Speed:         20,
		CargoCapacity: 500,
		CrewRequired:  1,
		CrewCurrent:   0,
		Weight:        600,
		Value:         3000, // 30 gp
	},
	{
		ID:            "wagon",
		Name:          "四轮货车",
		Description:   "一辆用于运输货物的重型货车",
		Type:          model.VehicleLand,
		AC:            15,
		HP:            60,
		Speed:         15,
		CargoCapacity: 1000,
		CrewRequired:  1,
		CrewCurrent:   0,
		Weight:        800,
		Value:         3500, // 35 gp
	},
	{
		ID:            "rowing-boat",
		Name:          "划艇",
		Description:   "一艘小型划艇",
		Type:          model.VehicleWater,
		AC:            15,
		HP:            40,
		Speed:         15,
		CargoCapacity: 500,
		CrewRequired:  1,
		CrewCurrent:   0,
		Weight:        200,
		Value:         5000, // 50 gp
	},
	{
		ID:            "sailing-ship",
		Name:          "帆船",
		Description:   "一艘大型远洋帆船",
		Type:          model.VehicleWater,
		AC:            15,
		HP:            200,
		Speed:         90,
		CargoCapacity: 25000,
		CrewRequired:  20,
		CrewCurrent:   0,
		Weight:        0,
		Value:         300000, // 3000 gp
	},
	{
		ID:            "keelboat",
		Name:          "龙骨船",
		Description:   "一艘小型帆船",
		Type:          model.VehicleWater,
		AC:            15,
		HP:            100,
		Speed:         60,
		CargoCapacity: 5000,
		CrewRequired:  8,
		CrewCurrent:   0,
		Weight:        0,
		Value:         150000, // 1500 gp
	},
}

// GetMountData 获取坐骑数据
func GetMountData(id model.ID) *MountData {
	for _, mount := range MountDataList {
		if mount.ID == id {
			return &mount
		}
	}
	return nil
}

// GetVehicleData 获取交通工具数据
func GetVehicleData(id model.ID) *model.Vehicle {
	for _, vehicle := range VehicleDataList {
		if vehicle.ID == id {
			return &vehicle
		}
	}
	return nil
}
