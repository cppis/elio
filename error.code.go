package elio

// todo range 별 클라 액션을 정하는게..? ex) 100번대는 에러리턴 후 타이틀로 이동 200번대는 에러만 출력 이런식으로...
// Common
const (
	Success     int32 = 0
	CurlSuccess int32 = 200

	ResponseCodeUnpackError  int32 = 11 // packet unpack error
	ResponseCodeMarshalError int32 = 12 // message marshal error
	ResponseCodeHttpError    int32 = 13 // http 통신 에러
)

// Frontend & Relayer
const (
	ResponseCodeNotLoggedIn      int32 = 101 // 로그인 상태가 아닌 경우
	ResponseCodeAlreadyLoggedOn  int32 = 102 // 이미 로그인 상태인 경우
	ResponseCodeDuplicateLogin   int32 = 103 // 중복로그인. 클라이언트는 이 에러코드를 받으면 타이틀로 이동.
	ResponseCodeNotInTheRoom     int32 = 104 // 방안에 있는 상태가 아닌 경우
	ResponseCodeAlreadyInTheRoom int32 = 105 // 이미 방안에 있는 상태인 경우

	ResponseCodeNotWaitingRoom int32 = 104 // 방이 존재 & 유저를 기다리는 상태가 아닌 경우

	ResponseCodeFailToLogon                int32 = 201
	ResponseCodeFailToLogoff               int32 = 202
	ResponseCodeFailToReqPvPRandomMatching int32 = 203
	ResponseCodeFailToEnterRoom            int32 = 204
	ResponseCodeFailToReadyPvPGame         int32 = 205
)

// todo 추후 삭제 예정, 샘플 서버에서 사용하는 코드는 상단의 코드로 교체 필요 (사용하는 코드는 상단에만 존재)
const (
	ErrNotLoggedon    int32 = 101
	ErrFailToFindRoom int32 = 103

	ErrAlreadyLoginUser int32 = 301

	ErrNotExistInRoom int32 = 313
)
