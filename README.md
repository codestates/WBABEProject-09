# WBABEProject-09

**프로젝트 배경**

> 언택트 시대에 급증하고 있는 온라인 주문 시스템은 이미 생활전반에 그 영향을 끼치고 있는 상황에, 가깝게는 배달어플, 매장에는 키오스크, 식당에는 패드를 이용한 메뉴 주문까지 그 사용범위가 점점 확대되어 가고 있습니다. 이런 시대에 해당 시스템을 이해, 경험하고 각 단계별 프로세스를 이해하여 구현함으로써 서비스 구축에 경험을 쌓고, golang의 이해를 돕습니다.

1. 학습자는 주문자/피주문자의 역할에서 필수적인 기능을 도출, 구현합니다.
2. 학습자는 해당 시스템에 대해 요구사항을 접수하고 주문자와 피주문자 입장에서 필요한 기능을 도출하여, 기능을 확장하고 주문 서비스를 원할하게 지원할수 있는 기능을 구현합니다.
3. 주문자는 신뢰있는 주문과 배달까지를 원합니다. 또, 피주문자는 주문내역을 관리하고 원할한 서비스가 제공되어야 합니다.

## Database

![Go-order](https://user-images.githubusercontent.com/25821827/209467784-7131dc4c-2acc-43fd-9d8e-466469fe4a36.png)

tUser: 유저 정보 collection
tMenu: 메뉴 정보 collection
tOrder: 오더 정보 collection

- tOrder에 menu필드는 menuId와 name으로 이루어진 []Object으로 tMenu와 연결됨
  tOrderSave: 초기 오더 정보 및 완료 오더 정보 저장 collection
- 초기 insert 된 정보가 복사됨
- order state가 3(취소), 7(배달완료) 된 경우 tOrder에서 tOrderSave로 이동

tReview: 리뷰 정보 collection

- tOrderSave에 state: 7(배달완료)인 document와 연결됨
- orderDay와 orderId를 통해서 tOrderSave와 연결
  tReviewSave: 초기 리뷰 및 리뷰 변경사항에 대해서 보관하기 위한 collection

sID: 각 collection에 존재하는 id(objectId가 아님)를 생성하기 위한 시퀀스
sOrderCount: Daily로 orderId를 생성하기 위한 시퀀스

## 기술스택

```
go: 서버 프로그래밍 언어
gin-gonic: 웹 프레임워크
toml: 서버 config설정을 위한 toml파일 파싱
mongo-driver: go언어용 mongoDB 지원 드라이버
swaggo: 개발 문서(API 명세) 생성 지원 라이브러리
mongoDB: 도큐먼트 지향 데이터베이스 시스템
```

## 프로젝트 구조

```bash
./oos
├── config // 내부적으로 쓰이는 공통 값에 대한 config 저장
│   ├── config.go
│   └── config.toml
├── controller // 요청에 따른 로직 핸들링
│   ├── controller_test.go	// 초기 실행시 기본 data를 구축
│   └── controller.go
├── docs // controller에 명시된 swagger에 따라 생성된 swagger 문서 저장
├── logger // log 처리를 위한 구성
│   └── logger.go
├── logs // 동작중 발생하는 사항에 대한 log 저장
├── model // DB와 연결되어 데이터 처리를 담당
│   └── model.go
├── router // http 요청에 대한 controller 연결
│   └── router.go
├── type // 내부적으로 사용되는 state, type 값을 상수로 명시
│   ├── menu  > menu.go
│   ├── order > order.go
│   └── user  > user.go
├ .gitignore
├ go.mod
├ go.sum
├ main.go
└ README.md
```

## 설치

1. 저장소 clone

```
https://github.com/codestates/WBABEProject-09.git
```

2. 패키지 설치

```
go mod tidy
```

3. DB 생성
   접근 정보 및 DB명은
   `WBABEProject-09\config\config.toml`
   에 명시되어 있으며, 인증은 따로 하지않음
   초기 DB명: `WBABEProject-09`
4. 서버 실행
   `go run main.go`를 통해 서버를 실행
5. 초기 데이터 구축
   `WBABEProject-09\controller\controller_test.go' 경로에 초기 데이터 구축을 위한 Test가 존재

`go run main.go`를 통해 서버가 정상적으로 실행된 상황에서 test 진행시 DB에 테스트 진행을 위한 user, menu, order 초기 정보가 쌓임
(review 및 수정등에 대해서는 아직 미지원)

## 1차 회고록

1. 급할수록 돌아가라
   작업시간이 적다고 생각해 빠르게 개발을 하려 할수록 다시 고치고, 나중에 바꾸고 하는 시간이 점점 늘어남
   당장 작업하는 기능에 대한 예시만 찾으려 하니 정리가 안되고 똑같은 자료를 계속 반복해서 붙잡고 씨름하는데 시간을 낭비함
2. 급할수록 설계해라
   기능을 만든다에 초점을 맞춰버리면서 요구사항과 설계를 대충하게됨,
   설계를 경시한만큼 그 대가는 작업한 규모가 커질수록 더 큰 변경사항으로 다시 다가옴
   요구사항에 단편적으로 대응하면서 전체적인 시스템을 고려하지 못함
   개발 능력만큼 설계능력 또한 매우 중요하다는 것을 잊어버림
3. 누군가 옆에 있다고 생각하라
   혼자 개발한다고 생각하니 깔끔한 코드나 구조에 대해서 신경안쓰고 급한대로 이리저리 붙여쓴 스파게티 코드가 됨
   제출기한이 다가오면서 다시한번 바라보니 한것도 없는데 난잡해 보임
4. 시간이 된다면 공식문서를 보라
   가장 많이 시간이 낭비된 것은 mongodb aggregate에서 lookup을 통해 join하는 과정이었고, 두번째로 낭비된 것은 aggregate query를 Go에 mongo-driver 상으로 bson 구조를 맞추는 일이었음
   마음이 급한만큼 기술적인 이해대신 구글창에 `three collection lookup example in go`와 같이 당장 내게 필요한 예제를 찾아다님
   결과적으로 이해하지도 못하는 예제를 붙잡으며 시간을 낭비함
