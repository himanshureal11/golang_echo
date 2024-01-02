package helper

func KeyName(key float64) string {
	var keyName string

	switch key {
	case 0.5:
		keyName = "pointfive"
	case 1:
		keyName = "one"
	case 1.5:
		keyName = "onefive"
	case 2:
		keyName = "two"
	case 2.5:
		keyName = "twofive"
	case 3:
		keyName = "three"
	case 3.5:
		keyName = "threefive"
	case 4:
		keyName = "four"
	case 4.5:
		keyName = "fourfive"
	case 5:
		keyName = "five"
	case 5.5:
		keyName = "fivefive"
	case 6:
		keyName = "six"
	case 6.5:
		keyName = "sixfive"
	case 7:
		keyName = "seven"
	case 7.5:
		keyName = "sevenfive"
	case 8:
		keyName = "eight"
	case 8.5:
		keyName = "eightfive"
	case 9:
		keyName = "nine"
	case 9.5:
		keyName = "ninefive"
	default:
		keyName = ""
	}

	return keyName
}
