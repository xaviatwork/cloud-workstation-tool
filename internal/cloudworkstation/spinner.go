package cloudworkstation

import "github.com/yarlson/pin"

func NewSpinner(text string) *pin.Pin {
	return pin.New("Loading ...",
		pin.WithSpinnerColor(pin.ColorCyan),
		pin.WithTextColor(pin.ColorYellow),
		pin.WithDoneSymbol('✅'),
		pin.WithDoneSymbolColor(pin.ColorGreen),
		pin.WithPrefix(text),
		pin.WithPosition(pin.PositionRight),
	)
}
