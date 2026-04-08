package engine

import (
	"context"
	"fmt"

	"github.com/zwh8800/dnd-core/pkg/model"
	"github.com/zwh8800/dnd-core/pkg/rules"
)

// StartTravelRequest 开始旅行请求
type StartTravelRequest struct {
	GameID      model.ID          `json:"game_id"`
	Destination string            `json:"destination"`
	Pace        model.TravelPace  `json:"pace"`
	Terrain     model.TerrainType `json:"terrain"`
	Distance    float64           `json:"distance"`
}

// StartTravelResult 开始旅行结果
type StartTravelResult struct {
	TravelState *model.TravelState `json:"travel_state"`
	Message     string             `json:"message"`
}

// StartTravel 开始旅行
func (e *Engine) StartTravel(ctx context.Context, req StartTravelRequest) (*StartTravelResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	travelState := &model.TravelState{
		Destination:      req.Destination,
		Pace:             req.Pace,
		Terrain:          req.Terrain,
		DistanceTotal:    req.Distance,
		DistanceTraveled: 0,
		DaysElapsed:      0,
		HoursToday:       0,
		IsActive:         true,
	}

	game.CurrentTravel = travelState

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return &StartTravelResult{
		TravelState: travelState,
		Message:     fmt.Sprintf("开始前往 %s 的旅行，总距离 %.1f 英里", req.Destination, req.Distance),
	}, nil
}

// AdvanceTravelRequest 推进旅行请求
type AdvanceTravelRequest struct {
	GameID model.ID `json:"game_id"`
	Hours  int      `json:"hours"`
}

// AdvanceTravelResult 推进旅行结果
type AdvanceTravelResult struct {
	DistanceTraveled float64                `json:"distance_traveled"`
	DaysElapsed      int                    `json:"days_elapsed"`
	ForageResult     *model.ForageResult    `json:"forage_result,omitempty"`
	NavigationResult *model.NavigationCheck `json:"navigation_result,omitempty"`
	EncounterResult  *model.EncounterCheck  `json:"encounter_result,omitempty"`
	Message          string                 `json:"message"`
}

// AdvanceTravel 推进旅行
func (e *Engine) AdvanceTravel(ctx context.Context, req AdvanceTravelRequest) (*AdvanceTravelResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	if game.CurrentTravel == nil || !game.CurrentTravel.IsActive {
		return nil, fmt.Errorf("no active travel")
	}

	travel := game.CurrentTravel
	result := &AdvanceTravelResult{}

	speed := 30
	distancePerHour := rules.CalculateTravelDistance(speed, travel.Pace, travel.Terrain) / 8.0
	distanceTraveled := distancePerHour * float64(req.Hours)
	travel.DistanceTraveled += distanceTraveled
	travel.HoursToday += req.Hours
	result.DistanceTraveled = distanceTraveled

	if travel.HoursToday >= 8 {
		travel.DaysElapsed++
		travel.HoursToday = 0

		forageResult, err := rules.ForagingCheck(14, 2, true)
		if err == nil {
			result.ForageResult = forageResult
		}

		navResult, err := rules.NavigationCheck(14, 2, true, travel.Terrain)
		if err == nil {
			result.NavigationResult = navResult
			if navResult.Lost {
				distanceTraveled /= 2
			}
		}

		encounterResult, err := rules.EncounterCheck()
		if err == nil {
			result.EncounterResult = encounterResult
		}

		result.DaysElapsed = travel.DaysElapsed
	}

	if travel.DistanceTraveled >= travel.DistanceTotal {
		travel.IsActive = false
		result.Message = fmt.Sprintf("已到达目的地！总用时 %d 天", travel.DaysElapsed)
	} else {
		result.Message = fmt.Sprintf("今日行进 %.1f 英里，已行进 %.1f/%.1f 英里", distanceTraveled, travel.DistanceTraveled, travel.DistanceTotal)
	}

	if err := e.saveGame(ctx, game); err != nil {
		return nil, err
	}

	return result, nil
}

// ForageRequest 觅食请求
type ForageRequest struct {
	GameID model.ID `json:"game_id"`
}

// ForageResultEngine 觅食结果
type ForageResultEngine struct {
	Result *model.ForageResult `json:"result"`
}

// Forage 执行觅食
func (e *Engine) Forage(ctx context.Context, req ForageRequest) (*ForageResultEngine, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	_, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	forageResult, err := rules.ForagingCheck(14, 2, true)
	if err != nil {
		return nil, err
	}

	return &ForageResultEngine{Result: forageResult}, nil
}

// NavigateRequest 导航请求
type NavigateRequest struct {
	GameID model.ID `json:"game_id"`
}

// NavigateResult 导航结果
type NavigateResult struct {
	Result *model.NavigationCheck `json:"result"`
}

// Navigate 执行导航检定
func (e *Engine) Navigate(ctx context.Context, req NavigateRequest) (*NavigateResult, error) {
	e.mu.Lock()
	defer e.mu.Unlock()

	game, err := e.loadGame(ctx, req.GameID)
	if err != nil {
		return nil, err
	}

	terrain := model.TerrainClear
	if game.CurrentTravel != nil {
		terrain = game.CurrentTravel.Terrain
	}

	navResult, err := rules.NavigationCheck(14, 2, true, terrain)
	if err != nil {
		return nil, err
	}

	return &NavigateResult{Result: navResult}, nil
}
