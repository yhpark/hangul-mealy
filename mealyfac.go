package hangulmealy

var cho, jung, jong = []rune{'ㄱ', 'ㄲ', 'ㄴ', 'ㄷ', 'ㄸ', 'ㄹ', 'ㅁ', 'ㅂ', 'ㅃ', 'ㅅ', 'ㅆ', 'ㅇ', 'ㅈ', 'ㅉ', 'ㅊ', 'ㅋ', 'ㅌ', 'ㅍ', 'ㅎ'},
	[]rune{'ㅏ', 'ㅐ', 'ㅑ', 'ㅒ', 'ㅓ', 'ㅔ', 'ㅕ', 'ㅖ', 'ㅗ', 'ㅘ', 'ㅙ', 'ㅚ', 'ㅛ', 'ㅜ', 'ㅝ', 'ㅞ', 'ㅟ', 'ㅠ', 'ㅡ', 'ㅢ', 'ㅣ'},
	[]rune{'x', 'ㄱ', 'ㄲ', 'ㄳ', 'ㄴ', 'ㄵ', 'ㄶ', 'ㄷ', 'ㄹ', 'ㄺ', 'ㄻ', 'ㄼ', 'ㄽ', 'ㄾ', 'ㄿ', 'ㅀ', 'ㅁ', 'ㅂ', 'ㅄ', 'ㅅ', 'ㅆ', 'ㅇ', 'ㅈ', 'ㅊ', 'ㅋ', 'ㅌ', 'ㅍ', 'ㅎ'}

func makeChoProirityPi() map[string]PiFunc {
	Pi := make(map[string]PiFunc)

	// <초성>
	Pi["1.1"] = func(text *HangulText, c rune) {
		text.Buf(c, c)
	}

	// <중성>
	Pi["2.1"] = func(text *HangulText, c rune) {
		text.Buf(assemble(idxOfRune(cho, text.Last()), idxOfRune(jung, c), 0), c)
	}

	// 중성->중성+<중성>
	Pi["2.2"] = func(text *HangulText, c rune) {
		choi, jungi, _ := disassemble(text.Last())
		text.Buf(assemble(choi, assembleJungi(jungi, idxOfRune(jung, c)), 0), c)
	}

	// <+초성>
	Pi["3.1"] = func(text *HangulText, c rune) {
		text.BufCho(text.Last(), c)
	}

	// +초성->종성, <+초성>
	Pi["3.2"] = func(text *HangulText, c rune) {
		choi, jungi, _ := disassemble(text.Last())
		text.BufCho(assemble(choi, jungi, idxOfRune(jong, text.cho)), c)
	}

	// 새로운 글자 <초성>
	Pi["3.3"] = func(text *HangulText, c rune) {
		text.Flush()
		text.Buf(c, c)
	}

	// +초성->새로운 글자 초성 <중성>
	Pi["4.1"] = func(text *HangulText, c rune) {
		newCho := text.RemoveLast()
		text.Flush()

		text.Buf(newCho, newCho)
		text.Buf(assemble(idxOfRune(cho, newCho), idxOfRune(jung, c), 0), c)
	}

	// +초성->종성, 새로운 글자 <+초성>
	Pi["4.2"] = func(text *HangulText, c rune) {
		choi, jungi, _ := disassemble(text.Last())
		newJong := text.RemoveLast()
		text.Buf(assemble(choi, jungi, idxOfRune(jong, newJong)), 0) // 0 -> not used
		text.Flush()

		text.Buf(c, c)
	}

	// +초성->새로운 글자 초성 <중성>
	Pi["4.3"] = func(text *HangulText, c rune) {
		newCho := text.RemoveLast()
		text.Flush()

		text.Buf(newCho, newCho)
		text.Buf(assemble(idxOfRune(cho, newCho), idxOfRune(jung, c), 0), c)
	}

	// +초성->종성+종성, 새로운 글자 <초성>
	Pi["4.4"] = func(text *HangulText, c rune) {
		choi, jungi, jongi := disassemble(text.Last())
		newJongB := text.RemoveLast()
		text.Buf(assemble(choi, jungi, assembleJongi(jongi, idxOfRune(jong, newJongB))), c)
		text.Flush()

		text.Buf(c, c)
	}

	// 4.2와 같으나 space buf 제외
	Pi["5.1"] = func(text *HangulText, c rune) {
		choi, jungi, _ := disassemble(text.Last())
		newJong := text.RemoveLast()
		text.Buf(assemble(choi, jungi, idxOfRune(jong, newJong)), 0) // 0 -> not used
		text.Flush()
	}

	// 4.4와 같으나 space buf 제외
	Pi["5.2"] = func(text *HangulText, c rune) {
		choi, jungi, jongi := disassemble(text.Last())
		newJongB := text.RemoveLast()
		text.Buf(assemble(choi, jungi, assembleJongi(jongi, idxOfRune(jong, newJongB))), c)
		text.Flush()
	}

	// Dead State -> q0
	Pi["*"] = func(text *HangulText, c rune) {
		text.Flush()
	}

	// Dead State from q0
	Pi["*0"] = func(text *HangulText, c rune) {
		text.Buf(c, c)

		text.Flush()
	}

	return Pi
}

func makeJongProirityPi() map[string]PiFunc {
	Pi := make(map[string]PiFunc)

	// 초성
	Pi["1.1"] = func(text *HangulText, c rune) {
		text.Buf(c, c)
	}

	// 중성
	Pi["2.1"] = func(text *HangulText, c rune) {
		text.Buf(assemble(idxOfRune(cho, text.Last()), idxOfRune(jung, c), 0), c)
	}

	// 중성 합성
	Pi["2.2"] = func(text *HangulText, c rune) {
		choi, jungi, _ := disassemble(text.Last())
		text.Buf(assemble(choi, assembleJungi(jungi, idxOfRune(jung, c)), 0), c)
	}

	// 종성
	Pi["3.1"] = func(text *HangulText, c rune) {
		choi, jungi, _ := disassemble(text.Last())
		text.Buf(assemble(choi, jungi, idxOfRune(jong, c)), c)
	}

	// 종성 합성
	Pi["3.2"] = func(text *HangulText, c rune) {
		choi, jungi, jongi := disassemble(text.Last())
		text.Buf(assemble(choi, jungi, assembleJongi(jongi, idxOfRune(jong, c))), c)
	}

	// 종성x, 초성
	Pi["3.3"] = func(text *HangulText, c rune) {
		text.Flush()
		text.Buf(c, c)
	}

	// 종성 -> 새로운 글자 초성 <중성>
	Pi["4.1"] = func(text *HangulText, c rune) {
		choi, jungi, jongi := disassemble(text.Last())
		text.Buf(assemble(choi, jungi, 0), 0)
		text.Flush()

		newCho := jong[jongi]

		text.Buf(newCho, newCho)
		text.Buf(assemble(idxOfRune(cho, newCho), idxOfRune(jung, c), 0), c)
	}

	// 새로운 글자
	Pi["4.2"] = func(text *HangulText, c rune) {
		text.Flush()

		text.Buf(c, c)
	}

	// 종성B -> 새로운 글자 초성 <중성>
	Pi["4.3"] = func(text *HangulText, c rune) {
		newCho := text.ins[len(text.ins)-1]

		choi, jungi, _ := disassemble(text.Last())
		jongi := idxOfRune(jong, text.ins[len(text.ins)-2])
		text.Buf(assemble(choi, jungi, jongi), 0)
		text.Flush()

		text.Buf(newCho, newCho)
		text.Buf(assemble(idxOfRune(cho, newCho), idxOfRune(jung, c), 0), c)
	}

	// Dead State -> q0
	Pi["*"] = func(text *HangulText, c rune) {
		text.Flush()
	}

	// Dead State from q0
	Pi["*0"] = func(text *HangulText, c rune) {
		text.Buf(c, c)

		text.Flush()
	}

	return Pi
}

func MakeHangulMealy(choPriority bool) (*Mealy, error) {
	var Q []string = []string{"S", "V", "O", "U", "A", "I", "K", "N", "R", "L1", "L2"}

	v := []rune{'ㅏ', 'ㅑ', 'ㅓ', 'ㅕ', 'ㅗ', 'ㅛ', 'ㅜ', 'ㅠ', 'ㅡ', 'ㅣ', 'ㅐ', 'ㅒ', 'ㅔ', 'ㅖ'}
	c := []rune{'ㄱ', 'ㄴ', 'ㄷ', 'ㄹ', 'ㅁ', 'ㅂ', 'ㅅ', 'ㅇ', 'ㅈ', 'ㅊ', 'ㅋ', 'ㅌ', 'ㅍ', 'ㅎ', 'ㄲ', 'ㄸ', 'ㅃ', 'ㅆ', 'ㅉ'}

	var Sigma []rune = append(append(v, c...), '*', ' ')

	var D [][]int = make([][]int, len(Q))

	var Lambda = make([][]string, len(Q))

	for i, _ := range Q {
		D[i] = make([]int, len(Sigma))
		Lambda[i] = make([]string, len(Sigma))
		for j, _ := range Sigma {
			D[i][j] = -1
		}
	}

	// p -cs-> q
	d := func(p string, cs []rune, q string) {
		for _, c := range cs {
			D[idxOfString(Q, p)][idxOfRune(Sigma, c)] = idxOfString(Q, q)
		}
	}

	d("S", c, "V")

	// 중성
	d("V", []rune{'ㅗ'}, "O")
	d("V", []rune{'ㅜ'}, "U")
	d("V", []rune{'ㅏ', 'ㅑ', 'ㅓ', 'ㅕ', 'ㅡ'}, "A")
	d("V", []rune{'ㅐ', 'ㅒ', 'ㅔ', 'ㅖ', 'ㅛ', 'ㅠ', 'ㅣ'}, "I")

	// 중성 합체
	d("O", []rune{'ㅏ'}, "A")
	d("O", []rune{'ㅣ', 'ㅐ'}, "I")
	d("U", []rune{'ㅓ'}, "A")
	d("U", []rune{'ㅣ', 'ㅔ'}, "I")
	d("A", []rune{'ㅣ'}, "I")

	// bus
	for _, p := range []string{"O", "U", "A", "I"} {
		d(p, []rune{'ㄱ', 'ㅂ'}, "K")
		d(p, []rune{'ㄴ'}, "N")
		d(p, []rune{'ㄹ'}, "R")
		d(p, []rune{'ㄷ', 'ㅁ', 'ㅅ', 'ㅇ', 'ㅈ', 'ㅊ', 'ㅋ', 'ㅌ', 'ㅍ', 'ㅎ', 'ㄲ', 'ㅆ'}, "L2")
		d(p, []rune{'ㄸ', 'ㅃ', 'ㅉ'}, "V")
	}

	// bus (reverse)
	for _, p := range []string{"K", "N", "R", "L1", "L2"} {
		d(p, []rune{'ㅗ'}, "O")
		d(p, []rune{'ㅜ'}, "U")
		d(p, []rune{'ㅏ', 'ㅑ', 'ㅓ', 'ㅕ', 'ㅡ'}, "A")
		d(p, []rune{'ㅐ', 'ㅒ', 'ㅔ', 'ㅖ', 'ㅛ', 'ㅠ', 'ㅣ'}, "I")
	}

	// ㄱㅅ
	d("K", []rune{'ㅅ'}, "L1")
	d("K", subtractRuneSlices(c, []rune{'ㅅ'}), "V")

	// ㄴㅈ, ㄴㅎ
	d("N", []rune{'ㅈ', 'ㅎ'}, "L1")
	d("N", subtractRuneSlices(c, []rune{'ㅈ', 'ㅎ'}), "V")

	// ㄹㄱ, ㄹㅁ, ... ㄹㅎ
	d("R", []rune{'ㄱ', 'ㅁ', 'ㅂ', 'ㅅ', 'ㅌ', 'ㅍ', 'ㅎ'}, "L1")
	d("R", subtractRuneSlices(c, []rune{'ㄱ', 'ㅁ', 'ㅂ', 'ㅅ', 'ㅌ', 'ㅍ', 'ㅎ'}), "V")

	d("L1", c, "V")
	d("L2", c, "V")

	if choPriority { // space
		d("K", []rune{' '}, "S")
		d("N", []rune{' '}, "S")
		d("R", []rune{' '}, "S")
		d("L1", []rune{' '}, "S")
		d("L2", []rune{' '}, "S")
	}

	// Dead State
	for _, p := range Q {
		d(p, []rune{'*'}, "S")
	}

	var Pi = func() map[string]PiFunc {
		if choPriority {
			return makeChoProirityPi()
		} else {
			return makeJongProirityPi()
		}
	}()

	// p -cs-> : pi
	l := func(p string, cs []rune, pi string) {
		for _, c := range cs {
			Lambda[idxOfString(Q, p)][idxOfRune(Sigma, c)] = pi
		}
	}

	l("S", c, "1.1")

	l("V", v, "2.1")

	l("O", []rune{'ㅏ', 'ㅣ', 'ㅐ'}, "2.2")
	l("U", []rune{'ㅓ', 'ㅣ', 'ㅔ'}, "2.2")
	l("A", []rune{'ㅣ'}, "2.2")

	c_jong := subtractRuneSlices(c, []rune{'ㄸ', 'ㅃ', 'ㅉ'})
	l("O", c_jong, "3.1")
	l("O", []rune{'ㄸ', 'ㅃ', 'ㅉ'}, "3.3")
	l("U", c_jong, "3.1")
	l("U", []rune{'ㄸ', 'ㅃ', 'ㅉ'}, "3.3")
	l("A", c_jong, "3.1")
	l("A", []rune{'ㄸ', 'ㅃ', 'ㅉ'}, "3.3")
	l("I", c_jong, "3.1")
	l("I", []rune{'ㄸ', 'ㅃ', 'ㅉ'}, "3.3")

	l("K", []rune{'ㅅ'}, "3.2")
	l("N", []rune{'ㅈ', 'ㅎ'}, "3.2")
	l("R", []rune{'ㄱ', 'ㅁ', 'ㅂ', 'ㅅ', 'ㅌ', 'ㅍ', 'ㅎ'}, "3.2")

	l("K", subtractRuneSlices(c, []rune{'ㅅ'}), "4.2")
	l("K", v, "4.1")

	l("N", subtractRuneSlices(c, []rune{'ㅈ', 'ㅎ'}), "4.2")
	l("N", v, "4.1")

	l("R", subtractRuneSlices(c, []rune{'ㄱ', 'ㅁ', 'ㅂ', 'ㅅ', 'ㅌ', 'ㅍ', 'ㅎ'}), "4.2")
	l("R", v, "4.1")

	if choPriority {
		l("L1", c, "4.4")
		l("L1", v, "4.3")

		// space
		l("K", []rune{' '}, "5.1")
		l("N", []rune{' '}, "5.1")
		l("R", []rune{' '}, "5.1")
		l("L1", []rune{' '}, "5.2")
		l("L2", []rune{' '}, "5.1")
	} else {
		l("L1", c, "4.2")
		l("L1", v, "4.3")
	}

	l("L2", c, "4.2")
	l("L2", v, "4.1")

	// Dead State
	for _, p := range Q {
		l(p, []rune{'*'}, "*")
	}

	l("S", []rune{'*'}, "*0")

	var q0 int = idxOfString(Q, "S")

	return NewMealy(Q, Sigma, Pi, D, Lambda, q0), nil
}

func disassemble(c rune) (int, int, int) {
	ci := int(c) - 0xAC00

	choi := ci / (21 * 28)
	jungi := (ci % (21 * 28)) / 28
	jongi := ci % 28

	return choi, jungi, jongi
}

func assemble(choi int, jungi int, jongi int) rune {
	return rune(0xAC00 + ((choi*21)+jungi)*28 + jongi)
}

func assembleJungi(jungiA int, jungiB int) int {
	jungMap := map[rune]map[rune]rune{
		'ㅗ': map[rune]rune{
			'ㅏ': 'ㅘ',
			'ㅣ': 'ㅚ',
			'ㅐ': 'ㅙ',
		},
		'ㅜ': map[rune]rune{
			'ㅓ': 'ㅝ',
			'ㅣ': 'ㅟ',
			'ㅔ': 'ㅞ',
		},
		'ㅏ': map[rune]rune{
			'ㅣ': 'ㅐ',
		},
		'ㅑ': map[rune]rune{
			'ㅣ': 'ㅒ',
		},
		'ㅓ': map[rune]rune{
			'ㅣ': 'ㅔ',
		},
		'ㅕ': map[rune]rune{
			'ㅣ': 'ㅖ',
		},
		'ㅘ': map[rune]rune{
			'ㅣ': 'ㅙ',
		},
		'ㅝ': map[rune]rune{
			'ㅣ': 'ㅞ',
		},
		'ㅡ': map[rune]rune{
			'ㅣ': 'ㅢ',
		},
	}

	return idxOfRune(jung, jungMap[jung[jungiA]][jung[jungiB]])
}

func assembleJongi(jongiA int, jongiB int) int {
	jongA := jong[jongiA]
	jongB := jong[jongiB]

	jongMap := map[rune][]rune{
		'ㄱ': []rune{'ㄱ', 'ㅅ'}, // ㄱ -> not used
		'ㅂ': []rune{'ㅅ'},
		'ㄴ': []rune{'ㅈ', 'ㅎ'},
		'ㄹ': []rune{'ㄱ', 'ㅁ', 'ㅂ', 'ㅅ', 'ㅌ', 'ㅍ', 'ㅎ'},
	}

	return jongiA + idxOfRune(jongMap[jongA], jongB) + 1
}
