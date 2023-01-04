package menu

const (
	StateSelling = 1 // 판매중
	StateSoldOut = 2 // 판매불가
)

var menuStateMap = map[int]string{
	1: "판매중",
	2: "판매불가",
}

func GetMenuStateText(menuState int) string {
	return menuStateMap[menuState]
}
