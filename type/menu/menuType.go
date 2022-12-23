package menu

const (
	StateSelling = 1 // 판매중
	StateSoldOut = 2 // 판매불가

	StateSellingText = "판매중"
	StateSoldOutText = "판매불가"
)

func GetMenuStateText(menuState int) string {
	returnString := ""
	switch menuState {
	case StateSelling:
		returnString = StateSellingText
	case StateSoldOut:
		returnString = StateSoldOutText
	}

	return returnString
}
