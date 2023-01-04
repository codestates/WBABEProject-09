package user

const (
	TypeOwner    = 1 // 피주문자
	TypeCustomer = 2 // 주문자
)

var userTypeMap = map[int]string{
	1: "피주문자",
	2: "주문자",
}

func GetUserTypeText(userType int) string {
	return userTypeMap[userType]
}
