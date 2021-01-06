package estimate

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

const float64CompareThreshold = 0.001

func TestIncorrectBestCaseBelowZero(t *testing.T) {
	_, err := NewDelphiEstimate(-0.1, 1, 2)
	assert.Equal(t, "Best case must be >= 0, provided: -0.1", err.Error())
}

func TestIncorrectMostLikelyBelowZero(t *testing.T) {
	_, err := NewDelphiEstimate(0, -0.5, 2)
	assert.Equal(t, "Most Likely must be >= 0, provided: -0.5", err.Error())
}

func TestIncorrectMostLikelyBelowBestCase(t *testing.T) {
	_, err := NewDelphiEstimate(2, 1, 3)
	assert.Equal(t, "Most Likely was smaller than Best Effort", err.Error())
}

func TestIncorrectWorstCaseBelowZero(t *testing.T) {
	_, err := NewDelphiEstimate(0, 0, -1.0)
	assert.Equal(t, "Worst Case must be >= 0, provided: -1", err.Error())
}

func TestIncorrectWorstCaseBelowMostLikely(t *testing.T) {
	_, err := NewDelphiEstimate(2, 3, 2.5)
	assert.Equal(t, "Worst Case was smaller than Most Likely", err.Error())
}

func TestCorrectCalculationOfEffort(t *testing.T) {
	ex := 15.0
	de, err := NewDelphiEstimate(10, 15, 20)
	r := de.GetEffort()

	assert.NoError(t, err)
	assert.True(t, math.Abs(ex-r) <= float64CompareThreshold)
}

func TestCorrectCalculationOfStandardDeviation(t *testing.T) {
	ex := 1.666
	de, err := NewDelphiEstimate(10, 15, 20)
	r := de.GetStandardDeviation()

	assert.NoError(t, err)
	assert.True(t, math.Abs(ex-r) <= float64CompareThreshold)
}
