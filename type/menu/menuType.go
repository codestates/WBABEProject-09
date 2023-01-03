package menu

const (
	StateSelling = 1 // 판매중
	StateSoldOut = 2 // 판매불가

	StateSellingText = "판매중"
	StateSoldOutText = "판매불가"
	/* [코드리뷰]
	 * 간결하고 보기 편한 코드입니다.
	 * 하지만 위의 내용에서는 주석과 Text가 중복되고 있습니다.
	 * orderType.go 와 userType.go에 관련한 내용을 달아놓았으니 확인 부탁드립니다.
	 */
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

func GetMenuTypeText() map[int]string {
	return map[int]string{
		1: "판매중",
		2: "판매불가",
	}
}
