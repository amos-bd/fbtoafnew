package main

import (
	"errors"
	"github.com/shopspring/decimal"
)

// BalanceData represents different types of balances
type BalanceData struct {
	MainBalance    decimal.Decimal
	BonusBalance   decimal.Decimal
	WithdrawLimit  decimal.Decimal
	DailyLimit     decimal.Decimal
	MonthlyLimit   decimal.Decimal
	CreditBalance  float64 // for testing float64 conversion
	PointsBalance  int     // for testing int conversion
}

// Balance handles balance operations with decimal precision
func Balance(amount decimal.Decimal, data *BalanceData, operation string) error {
	// 1. Replace all amount <= 0 checks with amount.Cmp(decimal.Zero) <= 0
	if amount.Cmp(decimal.Zero) <= 0 {
		return errors.New("amount must be positive")
	}

	// Calculate absolute amount using decimal operations instead of math.Abs
	absAmount := amount.Abs()

	switch operation {
	case "withdraw":
		// 2. Replace amount < 0 && data.XXX < math.Abs(amount) patterns
		// For decimal type: amount.Cmp(decimal.Zero) < 0 && data.XXX.Cmp(absAmount) < 0
		if amount.Cmp(decimal.Zero) < 0 && data.MainBalance.Cmp(absAmount) < 0 {
			return errors.New("insufficient main balance for withdrawal")
		}
		
		if amount.Cmp(decimal.Zero) < 0 && data.BonusBalance.Cmp(absAmount) < 0 {
			return errors.New("insufficient bonus balance for withdrawal")
		}

		if amount.Cmp(decimal.Zero) < 0 && data.WithdrawLimit.Cmp(absAmount) < 0 {
			return errors.New("withdrawal amount exceeds limit")
		}

		// For float64 type: amount.Cmp(decimal.Zero) < 0 && decimal.NewFromFloat(data.XXX).Cmp(absAmount) < 0
		if amount.Cmp(decimal.Zero) < 0 && decimal.NewFromFloat(data.CreditBalance).Cmp(absAmount) < 0 {
			return errors.New("insufficient credit balance for withdrawal")
		}

		// For int type: amount.Cmp(decimal.Zero) < 0 && decimal.NewFromInt(data.XXX).Cmp(absAmount) < 0
		if amount.Cmp(decimal.Zero) < 0 && decimal.NewFromInt(int64(data.PointsBalance)).Cmp(absAmount) < 0 {
			return errors.New("insufficient points balance for withdrawal")
		}

		// Update balances
		data.MainBalance = data.MainBalance.Sub(absAmount)

	case "deposit":
		if amount.Cmp(decimal.Zero) <= 0 {
			return errors.New("deposit amount must be positive")
		}
		
		// Check daily limit
		if amount.Cmp(decimal.Zero) < 0 && data.DailyLimit.Cmp(absAmount) < 0 {
			return errors.New("deposit exceeds daily limit")
		}

		// Check monthly limit  
		if amount.Cmp(decimal.Zero) < 0 && data.MonthlyLimit.Cmp(absAmount) < 0 {
			return errors.New("deposit exceeds monthly limit")
		}

		// Update balances
		data.MainBalance = data.MainBalance.Add(amount)

	case "transfer":
		// Check if source has sufficient balance
		if amount.Cmp(decimal.Zero) < 0 && data.MainBalance.Cmp(absAmount) < 0 {
			return errors.New("insufficient balance for transfer")
		}

		if amount.Cmp(decimal.Zero) < 0 && data.BonusBalance.Cmp(absAmount) < 0 {
			return errors.New("insufficient bonus balance for transfer")
		}

		// Additional validation for negative amounts
		if amount.Cmp(decimal.Zero) <= 0 {
			return errors.New("transfer amount must be positive")
		}

	case "bonus":
		// Bonus can only be positive
		if amount.Cmp(decimal.Zero) <= 0 {
			return errors.New("bonus amount must be positive")
		}

		// Update bonus balance
		data.BonusBalance = data.BonusBalance.Add(amount)

	default:
		return errors.New("unsupported operation")
	}

	return nil
}

// Helper function to check balance sufficiency with decimal precision
func CheckBalanceSufficiency(amount decimal.Decimal, data *BalanceData, balanceType string) bool {
	if amount.Cmp(decimal.Zero) <= 0 {
		return false
	}

	absAmount := amount.Abs()

	switch balanceType {
	case "main":
		// Check if main balance is sufficient
		return data.MainBalance.Cmp(absAmount) >= 0
	case "bonus":
		// Check if bonus balance is sufficient
		return data.BonusBalance.Cmp(absAmount) >= 0
	case "credit":
		// Check if credit balance is sufficient
		return decimal.NewFromFloat(data.CreditBalance).Cmp(absAmount) >= 0
	case "points":
		// Check if points balance is sufficient
		return decimal.NewFromInt(int64(data.PointsBalance)).Cmp(absAmount) >= 0
	default:
		return false
	}
}