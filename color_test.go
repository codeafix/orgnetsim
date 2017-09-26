package orgnetsim

import "testing"

func TestRandomlySelectAlternateColor(t *testing.T) {
	for i := 0; i < 1000; i++ {
		currentColor := Color(i % 7)
		color := RandomlySelectAlternateColor(currentColor)
		NotEqual(t, Grey, color, "Grey randomly selected")
		NotEqual(t, currentColor, color, "Existing Color randomly selected")
	}
	return
}
