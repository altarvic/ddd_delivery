package order

import "testing"

func TestStatus_IsValid(t *testing.T) {
	validStatuses := []string{
		"created",
		"assigned",
		"completed",
	}

	for _, status := range validStatuses {
		s := Status(status)
		if !s.IsValid() {
			t.Fail()
		}
	}

	inValidStatuses := []string{
		"",
		"unknown",
	}

	for _, status := range inValidStatuses {
		s := Status(status)
		if s.IsValid() {
			t.Fail()
		}
	}
}

func TestStatusFromString(t *testing.T) {
	validStatuses := []string{
		"created",
		"assigned",
		"completed",
	}

	for _, status := range validStatuses {
		_, err := StatusFromString(status)
		if err != nil {
			t.Fail()
		}
	}

	inValidStatuses := []string{
		"",
		"unknown",
	}

	for _, status := range inValidStatuses {
		_, err := StatusFromString(status)
		if err == nil {
			t.Fail()
		}
	}

}
