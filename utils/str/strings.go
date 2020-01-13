package str

func SubString(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func BlurRange(s string, keep int) string {
	r := []rune(s)
	lenR := len(r)
	if lenR >= 4 && lenR > keep*2 { // 4位以上
		for i, _ := range r {
			if i >= keep && i < lenR-keep { // 前后i个
				r[i] = '*'
			}
		}
	} else {
		r[0] = '*'
	}

	return string(r)
}

func BlurNickname(udId string, nickname string) (str string) {
	if nickname != "" {
		str = nickname
	} else {
		lenUdId := len(udId)
		if lenUdId > 8 {
			str = udId[:4] + udId[lenUdId-4:]
		} else {
			str = udId
		}
	}
	str = BlurRange(str, 2)
	return
}

func IsSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\n'
}
