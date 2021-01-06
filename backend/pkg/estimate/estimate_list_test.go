package estimate

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestNewEstimateListFailsDueToEmptyList(t *testing.T) {
	var list []Estimator
	_, err := NewEstimateList(list)
	assert.Error(t, err)
	assert.Equal(t, "Cannot create a empty list", err.Error())
}

func TestLenSuccess(t *testing.T) {
	est1, _ := NewDelphiEstimate(1.0, 2.0, 3.0)
	est2, _ := NewDelphiEstimate(2.0, 3.0, 4.0)
	list := []Estimator{est1, est2}
	res, err := NewEstimateList(list)
	assert.NoError(t, err)
	assert.Equal(t, 2, res.Len())
}

func TestLessSuccess(t *testing.T) {
	est1, _ := NewDelphiEstimate(1.0, 2.0, 3.0)
	est2, _ := NewDelphiEstimate(2.0, 3.0, 4.0)
	list := []Estimator{est1, est2}
	res, err := NewEstimateList(list)
	assert.NoError(t, err)
	assert.True(t, res.Less(1, 0))
}

func TestSwapSuccess(t *testing.T) {
	est1, _ := NewDelphiEstimate(1.0, 2.0, 3.0)
	est2, _ := NewDelphiEstimate(2.0, 3.0, 4.0)
	list := []Estimator{est1, est2}
	res, err := NewEstimateList(list)
	assert.NoError(t, err)
	res.Swap(0, 1)
	assert.Equal(t, est1, res.list[1])
	assert.Equal(t, est2, res.list[0])
}

func TestSortSuccess(t *testing.T) {
	est1, _ := NewDelphiEstimate(1.0, 2.0, 3.0)
	est2, _ := NewDelphiEstimate(2.0, 3.0, 4.0)
	est3, _ := NewDelphiEstimate(3.0, 4.0, 5.0)
	list := []Estimator{est1, est2, est3}
	res, err := NewEstimateList(list)
	assert.NoError(t, err)
	sort.Sort(res)
	assert.Equal(t, est1, res.list[2])
	assert.Equal(t, est2, res.list[1])
	assert.Equal(t, est3, res.list[0])
}
