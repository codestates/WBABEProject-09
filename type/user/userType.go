package user

const (
	TypeOwner    = 1 // 피주문자
	TypeCustomer = 2 // 주문자

	TypeOwnerText    = "피주문자"
	TypeCustomerText = "주문자"
	/* [코드리뷰]
	 * 현재 코드에서는 UserType을 주석을 통해 친절하게 알려주고 있습니다.
	 * 그러나 주석과 Text가 중복되는 경우가 발생합니다.
	 * 주문자, 피주문자의 정보를 int와 str 두개의 타입으로 매핑을 하는 코드이기 때문입니다.
	 * 아래 변수들은 각각 함께 움직이는 성격으로 보여집니다. map을 활용한 key, value 형식으로 관리하면 
	 * 시스템이 복잡할 수록 발생하기 쉬운, human error를 방지할 수 있을 것입니다.
	 * 또한 코드의 가독성도 좋아질 것입니다.
	 * as-is: TypeOwner, TypeOwnerText || TypeCustomer, TypeCustomerText
	 * tp-be: UserTypeMap := map[int]string{
							"1": "피주문자",
							"2": "주문자"
							}
	 * 이후 시스템의 확장성을 고려하여 다양한 유저가 생기는 경우에도 대응하기 쉬울 것입니다.
     */
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
