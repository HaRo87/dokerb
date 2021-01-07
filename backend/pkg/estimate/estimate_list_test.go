package estimate

import (
	"github.com/stretchr/testify/assert"
	"sort"
	"testing"
)

func TestNewEstimateListFailsDueToEmptyList(t *testing.T) {
	var list []UserEstimate
	_, err := NewEstimateList(list)
	assert.Error(t, err)
	assert.Equal(t, "Cannot create a empty list", err.Error())
}

func TestLenSuccess(t *testing.T) {
	est1, _ := NewDelphiEstimate(1.0, 2.0, 3.0)
	list := []UserEstimate{
		UserEstimate{
			Name:     "Tigger",
			Estimate: est1,
		},
	}
	res, err := NewEstimateList(list)
	assert.NoError(t, err)
	assert.Equal(t, 1, res.Len())
	assert.Equal(t, "", res.GetLastUser())
}

func TestLessSuccess(t *testing.T) {
	est1, _ := NewDelphiEstimate(1.0, 2.0, 3.0)
	est2, _ := NewDelphiEstimate(2.0, 3.0, 4.0)
	list := []UserEstimate{
		UserEstimate{
			Name:     "Tigger",
			Estimate: est1,
		},
		UserEstimate{
			Name:     "Rabbit",
			Estimate: est2,
		},
	}
	res, err := NewEstimateList(list)
	assert.NoError(t, err)
	assert.True(t, res.Less(1, 0))
}

func TestSwapSuccess(t *testing.T) {
	est1, _ := NewDelphiEstimate(1.0, 2.0, 3.0)
	est2, _ := NewDelphiEstimate(2.0, 3.0, 4.0)
	list := []UserEstimate{
		UserEstimate{
			Name:     "Tigger",
			Estimate: est1,
		},
		UserEstimate{
			Name:     "Rabbit",
			Estimate: est2,
		},
	}
	res, err := NewEstimateList(list)
	assert.NoError(t, err)
	res.Swap(0, 1)
	assert.Equal(t, est1, res.list[1].Estimate)
	assert.Equal(t, est2, res.list[0].Estimate)
}

func TestSortSuccess(t *testing.T) {
	est1, _ := NewDelphiEstimate(1.0, 2.0, 3.0)
	est2, _ := NewDelphiEstimate(2.0, 3.0, 4.0)
	est3, _ := NewDelphiEstimate(3.0, 4.0, 5.0)
	list := []UserEstimate{
		UserEstimate{
			Name:     "Tigger",
			Estimate: est1,
		},
		UserEstimate{
			Name:     "Rabbit",
			Estimate: est2,
		},
		UserEstimate{
			Name:     "Piglet",
			Estimate: est3,
		},
	}
	res, err := NewEstimateList(list)
	assert.NoError(t, err)
	assert.Equal(t, "Tigger", res.GetFirstUser())
	assert.Equal(t, "Piglet", res.GetLastUser())
	sort.Sort(res)
	assert.Equal(t, est1, res.list[2].Estimate)
	assert.Equal(t, est2, res.list[1].Estimate)
	assert.Equal(t, est3, res.list[0].Estimate)
	assert.Equal(t, "Piglet", res.GetFirstUser())
	assert.Equal(t, "Tigger", res.GetLastUser())
}
