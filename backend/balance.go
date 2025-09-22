package main

import (
	"errors"

	"github.com/shopspring/decimal"
)

// Balance represents a user's account balance with various balance types
type Balance struct {
	UserID       int64           `json:"user_id"`
	Cash         decimal.Decimal `json:"cash"`
	Bonus        decimal.Decimal `json:"bonus"`
	ReleaseBonus decimal.Decimal `json:"release_bonus"`
	Equity       int64           `json:"equity"`        // This is an integer field
	Integral     int64           `json:"integral"`      // This is an integer field
}

// BalanceData represents balance data for operations
type BalanceData struct {
	Cash         decimal.Decimal
	Bonus        decimal.Decimal
	ReleaseBonus decimal.Decimal
	Equity       int64
	Integral     int64
}

// Balance method demonstrating refactored decimal comparison patterns
func (b *Balance) Balance(amount decimal.Decimal, data *BalanceData) error {
	// Pattern 2: Refactored from amount < 0 && data.XXX < math.Abs(amount) to use decimal.Cmp
	absAmount := amount.Abs()
	
	// Example with decimal field: data.ReleaseBonus < math.Abs(amount) → data.ReleaseBonus.Cmp(absAmount) < 0
	if amount.Cmp(decimal.Zero) < 0 && data.ReleaseBonus.Cmp(absAmount) < 0 {
		return errors.New("insufficient release bonus")
	}

	// Example with integer field: float64(data.Equity) < math.Abs(amount) → decimal.NewFromInt(data.Equity).Cmp(absAmount) < 0
	if amount.Cmp(decimal.Zero) < 0 && decimal.NewFromInt(data.Equity).Cmp(absAmount) < 0 {
		return errors.New("insufficient equity")
	}

	// Another decimal field: data.Cash < math.Abs(amount) → data.Cash.Cmp(absAmount) < 0
	if amount.Cmp(decimal.Zero) < 0 && data.Cash.Cmp(absAmount) < 0 {
		return errors.New("insufficient cash")
	}

	// Another integer field: float64(data.Integral) < math.Abs(amount) → decimal.NewFromInt(data.Integral).Cmp(absAmount) < 0
	if amount.Cmp(decimal.Zero) < 0 && decimal.NewFromInt(data.Integral).Cmp(absAmount) < 0 {
		return errors.New("insufficient integral points")
	}

	// Complex pattern with bonus field: data.Bonus < math.Abs(amount) → data.Bonus.Cmp(absAmount) < 0
	if amount.Cmp(decimal.Zero) < 0 && data.Bonus.Cmp(absAmount) < 0 {
		return errors.New("insufficient bonus balance")
	}

	// Pattern 1: Refactored from amount <= 0 to amount.Cmp(decimal.Zero) <= 0 (checking for zero amounts only here since negative amounts with sufficient balance are allowed)
	if amount.Cmp(decimal.Zero) == 0 {
		return errors.New("amount must be positive")
	}

	// Update balances (simplified logic)
	b.Cash = b.Cash.Add(amount)
	b.Bonus = b.Bonus.Add(amount)
	b.ReleaseBonus = b.ReleaseBonus.Add(amount)

	return nil
}