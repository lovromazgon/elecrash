package ui

var brailleRunes = [4][2]rune{
	{'\u2801', '\u2808'},
	{'\u2802', '\u2810'},
	{'\u2804', '\u2820'},
	{'\u2840', '\u2880'},
}

// Braille builds a braille rune given the pairs that represent coordinates
// (row x column). There are 4 rows and 2 columns.
func Braille(pairs ...int) rune {
	var r rune
	for i := 0; i < len(pairs); i = i + 2 {
		x := pairs[i]
		y := pairs[i+1]
		r |= brailleRunes[x][y]
	}
	return r
}
