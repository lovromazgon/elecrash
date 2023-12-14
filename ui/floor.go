package ui

import "strconv"

func floorToRune(f int) rune {
	if f == 0 {
		return 'G'
	}
	return rune(strconv.Itoa(f)[0])
}
