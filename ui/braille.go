package ui

var brailleRunes = [4][2]rune{
	{'\u2801', '\u2808'},
	{'\u2802', '\u2810'},
	{'\u2804', '\u2820'},
	{'\u2840', '\u2880'},
}

// Braille builds a braille rune given the number.
func Braille(amount int) rune {
	if amount == 0 {
		return ' '
	}
	if amount < 0 || amount > 8 {
		panic("can only draw between 0 and 8 dots")
	}
	var r rune
	for i := 0; i < amount; i++ {
		r |= brailleRunes[i/2][i%2]
	}
	return r
}
