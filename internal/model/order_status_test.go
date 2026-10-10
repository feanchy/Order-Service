package model

import "testing"

func TestCanTransition(t *testing.T) {
	tests := []struct {
		name string
		from string
		to   string
		want bool
	}{
		{
			"created to confirmed",
			OrderStatusCreated,
			OrderStatusConfirmed,
			true,
		},
		{
			"confirmed to delivered",
			OrderStatusConfirmed,
			OrderStatusDelivered,
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CanTransition(tt.from, tt.to)

			if got != tt.want {
				t.Errorf("expected %v, got %v", tt.want, got)
			}
		})
	}
}
