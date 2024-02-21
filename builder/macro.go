package builder

func isMacro(str string) bool {
	//todo
	return false
}

func PaserMacro(str string) string {
	//todo
	if !isMacro(str) {
		return str
	}
	return ""
}
