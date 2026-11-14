package cloudworkstation

import "github.com/yarlson/pin"

func NewSpinner(text string) *pin.Pin {
	return pin.New(text,
		pin.WithPosition(pin.PositionRight),
		pin.WithDoneSymbol('✅'),
		// pin.WithSpinnerColor(pin.ColorCyan),
		// pin.WithTextColor(pin.ColorYellow),
		// pin.WithDoneSymbolColor(pin.ColorGreen),
	)
}
