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

)

var orderStateMap = map[int]string{
	1: "접수중",
	2: "접수완료",
	3: "오더취소",
	4: "추가주문",
	5: "조리중",
	6: "배달중",
	7: "배달완료",
}

func GetOrderStateText(orderState int) string {
	return orderStateMap[orderState]
}
