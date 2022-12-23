package user

const (
	TypeOwner    = 1 // 피주문자
	TypeCustomer = 2 // 주문자

	TypeOwnerText    = "피주문자"
	TypeCustomerText = "주문자"
)

func GetUserTypeText(userType int) string {
	returnString := ""
	switch userType {
	case TypeOwner:
		returnString = TypeOwnerText
	case TypeCustomer:
		returnString = TypeCustomerText
	}

	return returnString
}
