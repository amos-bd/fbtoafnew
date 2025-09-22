package main

import (
	"errors"
	"math"

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

// Balance method with patterns that need to be refactored
func (b *Balance) Balance(amount decimal.Decimal, data *BalanceData) error {
	// Pattern 1: amount <= 0 should be refactored to amount.Cmp(decimal.Zero) <= 0
	amountFloat, _ := amount.Float64()
	if amountFloat <= 0 {
		return errors.New("amount must be positive")
	}

	// Pattern 2: amount < 0 && data.XXX < math.Abs(amount) patterns that need refactoring
	
	// Example with decimal field
	amountFloat2, _ := amount.Float64()
	releaseFloat, _ := data.ReleaseBonus.Float64()
	if amountFloat2 < 0 && releaseFloat < math.Abs(amountFloat2) {
		return errors.New("insufficient release bonus")
	}

	// Example with integer field converted to float64
	amountFloat3, _ := amount.Float64()
	if amountFloat3 < 0 && float64(data.Equity) < math.Abs(amountFloat3) {
		return errors.New("insufficient equity")
	}

	// Another decimal field example
	amountFloat4, _ := amount.Float64()
	cashFloat, _ := data.Cash.Float64()
	if amountFloat4 < 0 && cashFloat < math.Abs(amountFloat4) {
		return errors.New("insufficient cash")
	}

	// Another integer field example
	amountFloat5, _ := amount.Float64()
	if amountFloat5 < 0 && float64(data.Integral) < math.Abs(amountFloat5) {
		return errors.New("insufficient integral points")
	}

	// More patterns with amount <= 0
	amountFloat6, _ := amount.Float64()
	if amountFloat6 <= 0 {
		return errors.New("invalid amount for processing")
	}

	// Complex pattern with bonus field
	amountFloat7, _ := amount.Float64()
	bonusFloat, _ := data.Bonus.Float64()
	if amountFloat7 < 0 && bonusFloat < math.Abs(amountFloat7) {
		return errors.New("insufficient bonus balance")
	}

	// Update balances (simplified logic)
	b.Cash = b.Cash.Add(amount)
	b.Bonus = b.Bonus.Add(amount)
	b.ReleaseBonus = b.ReleaseBonus.Add(amount)

	return nil
}