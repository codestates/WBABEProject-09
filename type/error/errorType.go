package error

const (
	// 유효성 에러
	Validation = 1000 // 유효성 오류
	// 유저 유효성 오류
	UserValidation     = 1101 // 유저 유효성 오류
	UserValidationID   = 1102 // 유저 ID 유효성 오류
	UserValidationType = 1103 // 유저 Type 유효성 오류
	// 메뉴 유효성 오류
	MenuValidation     = 1201 // 메뉴 유효성 오류
	MenuValidationID   = 1202 // 메뉴 ID 유효성 오류
	MenuValidationType = 1203 // 메뉴 Type 유효성 오류
	// 오더 유효성 오류
	OrderValidation     = 1301 // 오더 유효성 오류
	OrderValidationID   = 1302 // 오더 ID 유효성 오류
	OrderValidationType = 1303 // 오더 Type 유효성 오류
	// 리뷰 유효성 오류
	ReviewValidation     = 1401 // 리뷰 유효성 오류
	ReviewValidationID   = 1402 // 리뷰 ID 유효성 오류
	ReviewValidationType = 1403 // 리뷰 Type 유효성 오류

	// controller 오류
	Controller = 2000 // Controller 오류
	// 유저 컨트롤러 오류
	UserController       = 2100 // 유저 컨트롤러 오류
	UserControllerSearch = 2101 // 유저 조회 오류
	UserControllerInsert = 2102 // 유저 추가 오류
	UserControllerUpdate = 2103 // 유저 수정 오류
	UserControllerDelete = 2104 // 유저 삭제 오류
	// 메뉴 컨트롤러 오류
	MenuController       = 2200 // 메뉴 컨트롤러 오류
	MenuControllerSearch = 2201 // 메뉴 조회 오류
	MenuControllerInsert = 2202 // 메뉴 추가 오류
	MenuControllerUpdate = 2203 // 메뉴 수정 오류
	MenuControllerDelete = 2204 // 메뉴 삭제 오류
	// 오더 컨트롤러 오류
	OrderController       = 2300 // 오더 조회 오류
	OrderControllerSearch = 2301 // 오더 조회 오류
	OrderControllerInsert = 2302 // 오더 추가 오류
	OrderControllerUpdate = 2303 // 오더 수정 오류
	OrderControllerDelete = 2304 // 오더 삭제 오류
	// 리뷰 컨트롤러 오류
	ReviewController       = 2400 // 리뷰 컨트롤러 오류
	ReviewControllerSearch = 2401 // 리뷰 조회 오류
	ReviewControllerInsert = 2402 // 리뷰 추가 오류
	ReviewControllerUpdate = 2403 // 리뷰 수정 오류
	ReviewControllerDelete = 2404 // 리뷰 삭제 오류

)

var errorMessageMap = map[int]string{
	Validation:         "유효성 에러가 발생했습니다!",
	UserValidation:     "유저 정보가 유효하지 않습니다!",
	UserValidationID:   "유저 ID가 유효하지 않습니다!",
	UserValidationType: "유저 Type이 일치하지 않습니다!",

	MenuValidation:     "메뉴 정보가 유효하지 않습니다!",
	MenuValidationID:   "메뉴 ID가 유효하지 않습니다!",
	MenuValidationType: "메뉴 Type이 일치하지 않습니다!",

	OrderValidation:     "메뉴 정보가 유효하지 않습니다!",
	OrderValidationID:   "메뉴 ID가 유효하지 않습니다!",
	OrderValidationType: "메뉴 Type이 일치하지 않습니다!",

	ReviewValidation:     "리뷰 정보가 유효하지 않습니다!",
	ReviewValidationID:   "리뷰 ID가 유효하지 않습니다!",
	ReviewValidationType: "리뷰 Type이 일치하지 않습니다!",

	Controller: "컨트롤러 오류가 발생했습니다!",

	UserController:       "유저 컨트롤러 오류가 발생했습니다!",
	UserControllerSearch: "유저를 조회하지 못했습니다!",
	UserControllerInsert: "유저를 추가하지 못했습니다!",
	UserControllerUpdate: "유저를 수정하지 못했습니다!",
	UserControllerDelete: "유저를 삭제하지 못했습니다!",

	MenuController:       "메뉴 컨트롤러 오류가 발생했습니다!",
	MenuControllerSearch: "메뉴를 조회하지 발생했습니다!",
	MenuControllerInsert: "메뉴를 추가하지 못했습니다!",
	MenuControllerUpdate: "메뉴를 수정하지 못했습니다!",
	MenuControllerDelete: "메뉴를 삭제하지 못했습니다!",

	OrderController:       "오더 컨트롤러 오류가 발생했습니다!",
	OrderControllerSearch: "오더를 조회하지 못했습니다!",
	OrderControllerInsert: "오더를 추가하지 못했습니다!",
	OrderControllerUpdate: "오더를 수정하지 못했습니다!",
	OrderControllerDelete: "오더를 삭제하지 못했습니다!",

	ReviewController:       "리뷰 컨트롤러 오류가 발생했습니다!",
	ReviewControllerSearch: "리뷰를 조회하지 못했습니다!",
	ReviewControllerInsert: "리뷰를 추가하지 못했습니다!",
	ReviewControllerUpdate: "리뷰를 수정하지 못했습니다!",
	ReviewControllerDelete: "리뷰를 삭제하지 못했습니다!",
}

func GetErrorMessageText(menuState int) string {
	return errorMessageMap[menuState]
}
