package order

// 오더 흐름은 접수중 -> 접수 완료 or 오더 취소 -> 조리중 or 추가 주문 -> 배달중 -> 배달완료
//
const (
	StateReceiving  = 1 // 접수중
	StateReceived   = 2 // 접수완료
	StateCancel     = 3 // 오더취소
	StateAddOrder   = 4 // 추가주문
	StateCooking    = 5 // 조리중
	StateInDelivery = 6 // 배달중
	StateDelivered  = 7 // 배달완료

	StateReceivingText  = "접수중"
	StateReceivedText   = "접수완료"
	StateCancelText     = "오더취소"
	StateAddOrderText   = "추가주문"
	StateCookingText    = "조리중"
	StateInDeliveryText = "배달중"
	StateDeliveredText  = "배달완료"
	/* [코드리뷰]
	 * 한 눈에 보기 편한 코드입니다. 잘 짜주셨습니다.
	 * 위의 내용에서는 type을 선언하는 방식에서 주석과 Text가 중복되고 있습니다.
	 * 변수가 7가지의 상황을 표현해주고 있네요. 시스템을 개발할 때, state를 관리하는 것은 
	 * 많은 개발자들을 힘들게 합니다. 점점 많은 state가 생성될 수 있기 때문이죠.
	 * 시스템 담당자가 아니면 의미를 파악하기 어려운 모호한 int 형에 친절하게 주석을 달아주셔서 감사합니다.
	 * 그러나 int를 str과 매핑을 일일히 해야하는 번거로움이 발생합니다.
	 * 개발자의 관리포인트가 증가하는 것이죠.
	 * userType.go에서도 comment를 달아놓았는데, Map 자료형을 사용해보시는 것은 어떠실까요?
	 */
)

func GetOrderStateText(orderState int) string {
	returnString := ""
	switch orderState {
	case StateReceiving:
		returnString = StateReceivingText
	case StateReceived:
		returnString = StateReceivedText
	case StateCancel:
		returnString = StateCancelText
	case StateAddOrder:
		returnString = StateAddOrderText
	case StateCooking:
		returnString = StateCookingText
	case StateInDelivery:
		returnString = StateInDeliveryText
	case StateDelivered:
		returnString = StateDeliveredText
	}

	return returnString
}
