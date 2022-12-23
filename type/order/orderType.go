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
