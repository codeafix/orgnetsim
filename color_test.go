package orgnetsim

import "testing"

func TestRandomlySelectAlternateColor(t *testing.T) {
	for i := 0; i < 1000; i++ {
		currentColor := Color(i % 7)
		color := RandomlySelectAlternateColor(currentColor)
		if color == Grey {
			t.Errorf("Grey randomly selected")
		}
		if color == currentColor {
			t.Errorf("Existing Color randomly selected %s", currentColor.String())
		}
	}
	return
}
