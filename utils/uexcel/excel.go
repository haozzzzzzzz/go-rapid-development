package uexcel

const NumChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func X(iX int) (x string) {
	if iX < 26 {
		x = string(NumChars[iX])
	}
	return
}
