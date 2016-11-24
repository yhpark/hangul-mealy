package hangulmealy

type HangulText struct {
	final string
	buf   rune
	cho   rune   // for cho priority
	ins   []rune // not final inputs
}

func (text *HangulText) Last() rune {
	return text.buf
}

func (text *HangulText) Buf(buf rune, in rune) {
	text.buf = buf
	text.ins = append(text.ins, in)
}

func (text *HangulText) BufCho(buf rune, cho rune) {
	text.buf = buf
	text.ins = append(text.ins, cho)
	text.cho = cho
}

func (text *HangulText) IsBufEmpty() bool {
	return text.buf == 0
}

func (text *HangulText) RemoveLast() rune {
	if text.cho > 0 {
		text.cho = 0
		r := text.ins[len(text.ins)-1]
		text.ins = text.ins[:len(text.ins)-1]
		return r
	} else if text.buf > 0 {
		text.buf = 0
		r := text.ins[len(text.ins)-1]
		text.ins = text.ins[:len(text.ins)-1]
		return r
	} else if len(text.final) > 0 {
		rs := []rune(text.final)
		r := rs[len(rs)-1]
		rs = rs[:len(rs)-1]
		text.final = string(rs)
		return r
	}
	return 0
}

func (text *HangulText) Flush() {
	text.final = text.String()
	text.buf = 0
	text.ins = nil
	text.cho = 0
}

func (text *HangulText) ClearBuf() {
	text.buf = 0
	text.ins = nil
	text.cho = 0
}

func (text *HangulText) String() string {
	if text.Last() > 0 {
		if text.cho > 0 {
			return text.final + string(text.Last()) + string(text.cho)
		}
		return text.final + string(text.Last())
	}
	return text.final
}
