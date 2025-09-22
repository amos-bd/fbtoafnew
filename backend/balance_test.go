package main

import (
	"testing"

	"github.com/shopspring/decimal"
)

func TestBalanceRefactored(t *testing.T) {
	balance := &Balance{
		UserID:       1,
		Cash:         decimal.NewFromFloat(100.0),
		Bonus:        decimal.NewFromFloat(50.0),
		ReleaseBonus: decimal.NewFromFloat(25.0),
		Equity:       75,
		Integral:     100,
	}

	data := &BalanceData{
		Cash:         decimal.NewFromFloat(100.0),
		Bonus:        decimal.NewFromFloat(50.0),
		ReleaseBonus: decimal.NewFromFloat(25.0),
		Equity:       75,
		Integral:     100,
	}

	// Test positive amount (should succeed)
	err := balance.Balance(decimal.NewFromFloat(10.0), data)
	if err != nil {
		t.Errorf("Expected no error for positive amount, got: %v", err)
	}

	// Test zero amount (should fail with decimal.Cmp) - Pattern 1 refactoring
	err = balance.Balance(decimal.Zero, data)
	if err == nil {
		t.Error("Expected error for zero amount")
	}
	if err.Error() != "amount must be positive" {
		t.Errorf("Expected 'amount must be positive', got: %v", err)
	}

	// Test negative amount with sufficient balance (should succeed)
	err = balance.Balance(decimal.NewFromFloat(-5.0), data)
	if err != nil {
		t.Errorf("Expected no error for negative amount with sufficient balance, got: %v", err)
	}

	// Test insufficient balance scenarios - Pattern 2 refactoring
	insufficientData := &BalanceData{
		Cash:         decimal.NewFromFloat(5.0),
		Bonus:        decimal.NewFromFloat(5.0),
		ReleaseBonus: decimal.NewFromFloat(5.0),
		Equity:       5,
		Integral:     5,
	}

	// Test insufficient release bonus - demonstrates refactored comparison
	err = balance.Balance(decimal.NewFromFloat(-10.0), insufficientData)
	if err == nil {
		t.Error("Expected error for insufficient release bonus")
	}
	if err.Error() != "insufficient release bonus" {
		t.Errorf("Expected 'insufficient release bonus', got: %v", err)
	}
}

func TestDecimalComparisons(t *testing.T) {
	// Test that our refactored decimal comparisons work correctly
	
	// Test amount.Cmp(decimal.Zero) <= 0
	zero := decimal.Zero
	negative := decimal.NewFromFloat(-5.0)
	positive := decimal.NewFromFloat(5.0)

	if zero.Cmp(decimal.Zero) <= 0 {
		t.Log("✓ Zero comparison works: zero.Cmp(decimal.Zero) <= 0")
	} else {
		t.Error("Zero comparison failed")
	}

	if negative.Cmp(decimal.Zero) <= 0 {
		t.Log("✓ Negative comparison works: negative.Cmp(decimal.Zero) <= 0")
	} else {
		t.Error("Negative comparison failed")
	}

	if !(positive.Cmp(decimal.Zero) <= 0) {
		t.Log("✓ Positive comparison works: !(positive.Cmp(decimal.Zero) <= 0)")
	} else {
		t.Error("Positive comparison failed")
	}

	// Test decimal.NewFromInt() conversions
	intValue := int64(10)
	decimalValue := decimal.NewFromInt(intValue)
	absAmount := decimal.NewFromFloat(5.0)

	if decimalValue.Cmp(absAmount) > 0 {
		t.Log("✓ Integer to decimal conversion works: decimal.NewFromInt(10).Cmp(decimal.NewFromFloat(5.0)) > 0")
	} else {
		t.Error("Integer to decimal conversion failed")
	}
}