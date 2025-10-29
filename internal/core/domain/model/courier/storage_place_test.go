package courier

import (
	"github.com/google/uuid"
	"testing"
)

func TestNewStoragePLace(t *testing.T) {

	tests := []struct {
		testTitle     string
		name          string
		totalVolume   int
		errIsExpected bool
	}{
		{
			testTitle:     "Empty name",
			name:          "",
			totalVolume:   100,
			errIsExpected: true,
		},
		{
			testTitle:     "Wrong volume",
			name:          "suitcase",
			totalVolume:   0,
			errIsExpected: true,
		},
		{
			testTitle:     "Correct",
			name:          "suitcase",
			totalVolume:   100,
			errIsExpected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.testTitle, func(t *testing.T) {
			_, err := NewStoragePlace(test.name, test.totalVolume)
			if test.errIsExpected {
				if err == nil {
					t.Fail()
				}
			} else {
				if err != nil {
					t.Fail()
				}
			}
		})
	}
}

func TestStoragePlace_CanStore(t *testing.T) {
	tests := []struct {
		testTitle   string
		name        string
		totalVolume int
		volume      int
		expected    bool
	}{
		{
			testTitle:   "Small bag",
			name:        "Bag",
			totalVolume: 100,
			volume:      500,
			expected:    false,
		},
		{
			testTitle:   "Large bag",
			name:        "Bag",
			totalVolume: 1000,
			volume:      500,
			expected:    true,
		},
		{
			testTitle:   "Perfect fit",
			name:        "Bag",
			totalVolume: 500,
			volume:      500,
			expected:    true,
		},
	}
	for _, test := range tests {
		t.Run(test.testTitle, func(t *testing.T) {
			s, _ := NewStoragePlace(test.name, test.totalVolume)
			if got := s.CanStore(test.volume); got != test.expected {
				t.Fail()
			}
		})
	}
}

func TestStoragePlace_Store(t *testing.T) {
	sp, _ := NewStoragePlace("Bag", 100)

	orderId := uuid.New()
	err := sp.Store(orderId, 200)
	if err == nil {
		t.Fail()
	}

	err = sp.Store(orderId, 100)
	if err != nil {
		t.Fail()
	}

	if *sp.orderID != orderId {
		t.Fail()
	}
}

func TestStoragePlace_Clear(t *testing.T) {
	sp, _ := NewStoragePlace("Bag", 100)

	_ = sp.Store(uuid.New(), 80)

	sp.Clear()
	if sp.orderID != nil {
		t.Fail()
	}
}
