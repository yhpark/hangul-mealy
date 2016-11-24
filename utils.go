package main

func idxOfRune(slice []rune, val rune) int {
	for i, e := range slice {
		if e == val {
			return i
		}
	}
	return -1
}

func idxOfString(slice []string, val string) int {
	for i, e := range slice {
		if e == val {
			return i
		}
	}
	return -1
}

func subtractRuneSlices(a []rune, b []rune) []rune {
	c := make([]rune, 0)

loop:
	for _, ac := range a {
		for _, bc := range b {
			if ac == bc {
				continue loop
			}
		}
		c = append(c, ac)
	}

	return c
}

func engByteToKorRune(b byte) rune {
	r := rune(b)
	engKorMap := map[rune]rune{
		'a': 'ㅁ', 'A': 'ㅁ',
		'b': 'ㅠ', 'B': 'ㅠ',
		'c': 'ㅊ', 'C': 'ㅊ',
		'd': 'ㅇ', 'D': 'ㅇ',
		'e': 'ㄷ', 'E': 'ㄸ',
		'f': 'ㄹ', 'F': 'ㄹ',
		'g': 'ㅎ', 'G': 'ㅎ',
		'h': 'ㅗ', 'H': 'ㅗ',
		'i': 'ㅑ', 'I': 'ㅑ',
		'j': 'ㅓ', 'J': 'ㅓ',
		'k': 'ㅏ', 'K': 'ㅏ',
		'l': 'ㅣ', 'L': 'ㅣ',
		'm': 'ㅡ', 'M': 'ㅡ',
		'n': 'ㅜ', 'N': 'ㅜ',
		'o': 'ㅐ', 'O': 'ㅒ',
		'p': 'ㅔ', 'P': 'ㅖ',
		'q': 'ㅂ', 'Q': 'ㅃ',
		'r': 'ㄱ', 'R': 'ㄲ',
		's': 'ㄴ', 'S': 'ㄴ',
		't': 'ㅅ', 'T': 'ㅆ',
		'u': 'ㅕ', 'U': 'ㅕ',
		'v': 'ㅍ', 'V': 'ㅍ',
		'w': 'ㅈ', 'W': 'ㅉ',
		'x': 'ㅌ', 'X': 'ㅌ',
		'y': 'ㅛ', 'Y': 'ㅛ',
		'z': 'ㅋ', 'Z': 'ㅋ',
	}

	if engKorMap[r] == 0 {
		return r
	}
	return engKorMap[r]
}
