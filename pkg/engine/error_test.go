package engine

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zwh8800/dnd-core/pkg/model"
)

func TestEngineError(t *testing.T) {
	t.Run("error message includes operation and phase", func(t *testing.T) {
		innerErr := errors.New("internal error")
		engineErr := &EngineError{
			Op:    "testOperation",
			Err:   innerErr,
			Phase: model.PhaseExploration,
		}

		msg := engineErr.Error()
		assert.Contains(t, msg, "testOperation")
		assert.Contains(t, msg, "exploration")
		assert.Contains(t, msg, "internal error")
	})

	t.Run("error message without inner error", func(t *testing.T) {
		engineErr := &EngineError{
			Op:    "testOperation",
			Phase: model.PhaseCombat,
		}

		msg := engineErr.Error()
		assert.Contains(t, msg, "testOperation")
		assert.Contains(t, msg, "combat")
	})

	t.Run("unwrap returns inner error", func(t *testing.T) {
		innerErr := errors.New("wrapped error")
		engineErr := &EngineError{
			Op:  "test",
			Err: innerErr,
		}

		unwrapped := engineErr.Unwrap()
		assert.Equal(t, innerErr, unwrapped)
	})

	t.Run("unwrap returns nil when no inner error", func(t *testing.T) {
		engineErr := &EngineError{
			Op: "test",
		}

		unwrapped := engineErr.Unwrap()
		assert.Nil(t, unwrapped)
	})

	t.Run("error works with errors.Is", func(t *testing.T) {
		innerErr := ErrNotFound
		engineErr := &EngineError{
			Op:  "loadGame",
			Err: innerErr,
		}

		assert.True(t, errors.Is(engineErr, ErrNotFound))
	})

	t.Run("error message format is consistent", func(t *testing.T) {
		engineErr := &EngineError{
			Op:    "saveGame",
			Phase: model.PhaseRest,
		}

		expected := "engine error in saveGame (phase: rest)"
		assert.Equal(t, expected, engineErr.Error())
	})

	t.Run("error with details field", func(t *testing.T) {
		engineErr := &EngineError{
			Op:  "testOp",
			Err: fmt.Errorf("test error"),
			Details: map[string]any{
				"key": "value",
			},
		}

		assert.NotNil(t, engineErr.Details)
		assert.Equal(t, "value", engineErr.Details["key"])
	})
}

func TestEngineErrorDefinitions(t *testing.T) {
	t.Run("ErrNotFound is defined", func(t *testing.T) {
		assert.NotNil(t, ErrNotFound)
		assert.Equal(t, "entity not found", ErrNotFound.Error())
	})

	t.Run("ErrAlreadyExists is defined", func(t *testing.T) {
		assert.NotNil(t, ErrAlreadyExists)
		assert.Equal(t, "entity already exists", ErrAlreadyExists.Error())
	})

	t.Run("ErrInvalidState is defined", func(t *testing.T) {
		assert.NotNil(t, ErrInvalidState)
		assert.Equal(t, "invalid game state for this operation", ErrInvalidState.Error())
	})

	t.Run("ErrCombatNotActive is defined", func(t *testing.T) {
		assert.NotNil(t, ErrCombatNotActive)
		assert.Equal(t, "no active combat", ErrCombatNotActive.Error())
	})

	t.Run("ErrCombatAlreadyActive is defined", func(t *testing.T) {
		assert.NotNil(t, ErrCombatAlreadyActive)
		assert.Equal(t, "combat is already active", ErrCombatAlreadyActive.Error())
	})

	t.Run("ErrNotYourTurn is defined", func(t *testing.T) {
		assert.NotNil(t, ErrNotYourTurn)
		assert.Equal(t, "it is not this actor's turn", ErrNotYourTurn.Error())
	})

	t.Run("ErrActionAlreadyUsed is defined", func(t *testing.T) {
		assert.NotNil(t, ErrActionAlreadyUsed)
		assert.Equal(t, "action has already been used this turn", ErrActionAlreadyUsed.Error())
	})

	t.Run("ErrInsufficientSlots is defined", func(t *testing.T) {
		assert.NotNil(t, ErrInsufficientSlots)
		assert.Equal(t, "insufficient spell slots", ErrInsufficientSlots.Error())
	})

	t.Run("ErrInvalidTarget is defined", func(t *testing.T) {
		assert.NotNil(t, ErrInvalidTarget)
		assert.Equal(t, "invalid target for this action", ErrInvalidTarget.Error())
	})

	t.Run("ErrOutOfRange is defined", func(t *testing.T) {
		assert.NotNil(t, ErrOutOfRange)
		assert.Equal(t, "target is out of range", ErrOutOfRange.Error())
	})

	t.Run("ErrNoLineOfSight is defined", func(t *testing.T) {
		assert.NotNil(t, ErrNoLineOfSight)
		assert.Equal(t, "no line of sight to target", ErrNoLineOfSight.Error())
	})

	t.Run("ErrConcentrationBroken is defined", func(t *testing.T) {
		assert.NotNil(t, ErrConcentrationBroken)
		assert.Equal(t, "concentration check failed", ErrConcentrationBroken.Error())
	})

	t.Run("ErrActorIncapacitated is defined", func(t *testing.T) {
		assert.NotNil(t, ErrActorIncapacitated)
		assert.Equal(t, "actor is incapacitated", ErrActorIncapacitated.Error())
	})

	t.Run("ErrInvalidDiceExpression is defined", func(t *testing.T) {
		assert.NotNil(t, ErrInvalidDiceExpression)
		assert.Equal(t, "invalid dice expression", ErrInvalidDiceExpression.Error())
	})

	t.Run("ErrStorageError is defined", func(t *testing.T) {
		assert.NotNil(t, ErrStorageError)
		assert.Equal(t, "storage operation failed", ErrStorageError.Error())
	})

	t.Run("ErrValidationFailed is defined", func(t *testing.T) {
		assert.NotNil(t, ErrValidationFailed)
		assert.Equal(t, "validation failed", ErrValidationFailed.Error())
	})

	t.Run("ErrPhaseNotAllowed is defined", func(t *testing.T) {
		assert.NotNil(t, ErrPhaseNotAllowed)
		assert.Equal(t, "operation not allowed in current phase", ErrPhaseNotAllowed.Error())
	})
}
