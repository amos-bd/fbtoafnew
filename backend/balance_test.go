package main

import (
	"testing"
	"github.com/shopspring/decimal"
)

func TestBalanceDecimalComparisons(t *testing.T) {
	// Initialize test data
	data := &BalanceData{
		MainBalance:    decimal.NewFromFloat(1000.50),
		BonusBalance:   decimal.NewFromFloat(250.25),
		WithdrawLimit:  decimal.NewFromFloat(500.00),
		DailyLimit:     decimal.NewFromFloat(10000.00),
		MonthlyLimit:   decimal.NewFromFloat(100000.00),
		CreditBalance:  750.75,
		PointsBalance:  1500,
	}

	// Test 1: amount <= 0 check using amount.Cmp(decimal.Zero) <= 0
	t.Run("Zero amount check", func(t *testing.T) {
		zeroAmount := decimal.Zero
		err := Balance(zeroAmount, data, "deposit")
		if err == nil {
			t.Error("Expected error for zero amount, got nil")
		}
		
		negativeAmount := decimal.NewFromFloat(-100.0)
		err = Balance(negativeAmount, data, "deposit")
		if err == nil {
			t.Error("Expected error for negative amount, got nil")
		}
	})

	// Test 2: Positive amount operations
	t.Run("Positive deposit", func(t *testing.T) {
		originalBalance := data.MainBalance
		depositAmount := decimal.NewFromFloat(100.0)
		err := Balance(depositAmount, data, "deposit")
		if err != nil {
			t.Errorf("Unexpected error for positive deposit: %v", err)
		}
		expectedBalance := originalBalance.Add(depositAmount)
		if !data.MainBalance.Equal(expectedBalance) {
			t.Errorf("Expected balance %v, got %v", expectedBalance, data.MainBalance)
		}
	})

	// Test 3: Bonus operation with decimal comparisons
	t.Run("Bonus addition", func(t *testing.T) {
		originalBonus := data.BonusBalance
		bonusAmount := decimal.NewFromFloat(50.0)
		err := Balance(bonusAmount, data, "bonus")
		if err != nil {
			t.Errorf("Unexpected error for bonus: %v", err)
		}
		expectedBonus := originalBonus.Add(bonusAmount)
		if !data.BonusBalance.Equal(expectedBonus) {
			t.Errorf("Expected bonus balance %v, got %v", expectedBonus, data.BonusBalance)
		}
	})

	// Test 4: Balance sufficiency check
	t.Run("Balance sufficiency checks", func(t *testing.T) {
		// Test with sufficient balance
		amount := decimal.NewFromFloat(500.0)
		sufficient := CheckBalanceSufficiency(amount, data, "main")
		if !sufficient {
			t.Error("Expected sufficient balance for main balance check")
		}

		// Test with insufficient balance (using large amount)
		largeAmount := decimal.NewFromFloat(99999.0)
		sufficient = CheckBalanceSufficiency(largeAmount, data, "main")
		if sufficient {
			t.Error("Expected insufficient balance for large amount")
		}

		// Test with zero amount
		zeroAmount := decimal.Zero
		sufficient = CheckBalanceSufficiency(zeroAmount, data, "main")
		if sufficient {
			t.Error("Expected insufficient balance for zero amount")
		}
	})

	// Test 5: Verify decimal operations maintain precision
	t.Run("Decimal precision", func(t *testing.T) {
		// Reset data for precision test
		data.MainBalance, _ = decimal.NewFromString("1000.123456789")
		preciseAmount, _ := decimal.NewFromString("0.000000001")
		
		err := Balance(preciseAmount, data, "deposit")
		if err != nil {
			t.Errorf("Unexpected error for precise deposit: %v", err)
		}
		
		expected, _ := decimal.NewFromString("1000.123456790")
		if !data.MainBalance.Equal(expected) {
			t.Errorf("Precision lost: expected %v, got %v", expected, data.MainBalance)
		}
	})
}

func TestDecimalComparisonPatterns(t *testing.T) {
	data := &BalanceData{
		MainBalance:    decimal.NewFromFloat(100.0),
		BonusBalance:   decimal.NewFromFloat(50.0),
		WithdrawLimit:  decimal.NewFromFloat(200.0),
		CreditBalance:  75.0,
		PointsBalance:  80,
	}

	// Test the specific patterns mentioned in requirements:
	// amount.Cmp(decimal.Zero) < 0 && data.XXX.Cmp(absAmount) < 0

	t.Run("Decimal type comparison pattern", func(t *testing.T) {
		// This should trigger: amount.Cmp(decimal.Zero) < 0 && data.MainBalance.Cmp(absAmount) < 0
		negativeAmount := decimal.NewFromFloat(-150.0) // More than MainBalance (100.0)
		err := Balance(negativeAmount, data, "withdraw")
		if err == nil {
			t.Error("Expected error for insufficient main balance")
		}
	})

	t.Run("Float64 type comparison pattern", func(t *testing.T) {
		// This should trigger: amount.Cmp(decimal.Zero) < 0 && decimal.NewFromFloat(data.CreditBalance).Cmp(absAmount) < 0
		negativeAmount := decimal.NewFromFloat(-100.0) // More than CreditBalance (75.0)
		err := Balance(negativeAmount, data, "withdraw")
		if err == nil {
			t.Error("Expected error for insufficient credit balance")
		}
	})

	t.Run("Int type comparison pattern", func(t *testing.T) {
		// This should trigger: amount.Cmp(decimal.Zero) < 0 && decimal.NewFromInt(data.PointsBalance).Cmp(absAmount) < 0
		negativeAmount := decimal.NewFromFloat(-100.0) // More than PointsBalance (80)
		err := Balance(negativeAmount, data, "withdraw")
		if err == nil {
			t.Error("Expected error for insufficient points balance")
		}
	})
}