package engine

import (
	"errors"
	"fmt"

	"github.com/zwh8800/dnd-core/pkg/model"
)

// 引擎错误定义
var (
	ErrNotFound              = errors.New("entity not found")
	ErrAlreadyExists         = errors.New("entity already exists")
	ErrInvalidState          = errors.New("invalid game state for this operation")
	ErrCombatNotActive       = errors.New("no active combat")
	ErrCombatAlreadyActive   = errors.New("combat is already active")
	ErrNotYourTurn           = errors.New("it is not this actor's turn")
	ErrActionAlreadyUsed     = errors.New("action has already been used this turn")
	ErrInsufficientSlots     = errors.New("insufficient spell slots")
	ErrInvalidTarget         = errors.New("invalid target for this action")
	ErrOutOfRange            = errors.New("target is out of range")
	ErrNoLineOfSight         = errors.New("no line of sight to target")
	ErrConcentrationBroken   = errors.New("concentration check failed")
	ErrActorIncapacitated    = errors.New("actor is incapacitated")
	ErrInvalidDiceExpression = errors.New("invalid dice expression")
	ErrStorageError          = errors.New("storage operation failed")
	ErrValidationFailed      = errors.New("validation failed")
	ErrPhaseNotAllowed       = errors.New("operation not allowed in current phase")
)

// EngineError 包装错误并附加上下文
type EngineError struct {
	Op      string         // 操作名称
	Err     error          // 原始错误
	Phase   model.Phase    // 当前阶段
	Details map[string]any // 额外详情
}

func (e *EngineError) Error() string {
	msg := fmt.Sprintf("engine error in %s (phase: %s)", e.Op, e.Phase)
	if e.Err != nil {
		msg += ": " + e.Err.Error()
	}
	return msg
}

func (e *EngineError) Unwrap() error {
	return e.Err
}
