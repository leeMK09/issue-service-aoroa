# 이슈 관리 API

이슈 관리 REST API

## 프로젝트 구조

```
aoroa_backend_assignment/
├── issue/                      # 이슈 도메인
│   ├── model/                  # 도메인 모델
│   │   ├── issue.go           # Issue 엔티티 및 비즈니스 로직
│   │   └── update_command.go  # 업데이트 명령 패턴
│   ├── application/           # 애플리케이션 서비스
│   │   └── issue_service.go   # 이슈 비즈니스 로직 오케스트레이션
│   ├── infrastructure/        # 인프라스트럭처
│   │   └── issue_repository.go # 데이터 저장소
│   └── presentation/          # 프레젠테이션 계층
│       └── issue_controller.go # HTTP 핸들러
├── user/                      # 사용자 도메인
│   ├── model/                 # 사용자 모델
│   │   └── user.go
│   └── infrastructure/        # 사용자 저장소
│       └── user_repository.go
├── main.go                    # 애플리케이션 진입점
├── go.mod                     # Go 모듈 설정
└── README.md                  # 프로젝트 문서
```

## 실행 방법

### 1. 의존성 설치

```bash
go mod tidy
```

### 2. 애플리케이션 실행

```bash
go run main.go
```

애플리케이션은 포트 8080에서 실행됩니다.

### 3. 빌드 (선택사항)

```bash
go build -o issue-service-aoroa
./issue-service-aoroa
```

## API 테스트 방법

### cURL 명령어를 사용한 테스트

#### 1. 이슈 생성 [POST] /issue

```bash
# 담당자가 있는 이슈 생성 (상태: IN_PROGRESS)
curl -X POST http://localhost:8080/issue \
  -H "Content-Type: application/json" \
  -d '{
    "title": "버그 수정 필요",
    "description": "로그인 페이지에서 오류 발생",
    "userId": 1
  }'

# 담당자가 없는 이슈 생성 (상태: PENDING)
curl -X POST http://localhost:8080/issue \
  -H "Content-Type: application/json" \
  -d '{
    "title": "새로운 기능 추가",
    "description": "사용자 프로필 기능"
  }'
```

#### 2. 이슈 목록 조회 [GET] /issues

```bash
# 전체 이슈 조회
curl http://localhost:8080/issues

# 상태별 필터링
curl "http://localhost:8080/issues?status=PENDING"
curl "http://localhost:8080/issues?status=IN_PROGRESS"
curl "http://localhost:8080/issues?status=COMPLETED"
curl "http://localhost:8080/issues?status=CANCELLED"
```

#### 3. 이슈 상세 조회 [GET] /issue/:id

```bash
curl http://localhost:8080/issue/1
```

#### 4. 이슈 수정 [PATCH] /issue/:id

```bash
# 제목과 상태 수정
curl -X PATCH http://localhost:8080/issue/1 \
  -H "Content-Type: application/json" \
  -d '{
    "title": "로그인 버그 수정",
    "status": "COMPLETED"
  }'

# 담당자 할당 (PENDING → IN_PROGRESS 자동 전환)
curl -X PATCH http://localhost:8080/issue/1 \
  -H "Content-Type: application/json" \
  -d '{
    "userId": 2
  }'

# 담당자 제거 (→ PENDING 자동 전환)
curl -X PATCH http://localhost:8080/issue/1 \
  -H "Content-Type: application/json" \
  -d '{
    "userId": null
  }'
```

### 에러 케이스 테스트

#### 1. 유효성 검사 에러

```bash
# 제목 없이 이슈 생성 (400 Bad Request)
curl -X POST http://localhost:8080/issue \
  -H "Content-Type: application/json" \
  -d '{
    "description": "제목이 없습니다"
  }'
```

#### 2. 존재하지 않는 사용자 (400 Bad Request)

```bash
curl -X POST http://localhost:8080/issue \
  -H "Content-Type: application/json" \
  -d '{
    "title": "테스트",
    "userId": 999
  }'
```

#### 3. 잘못된 상태값 (400 Bad Request)

```bash
curl -X GET "http://localhost:8080/issues?status=INVALID_STATUS"
```

#### 4. 존재하지 않는 이슈 (404 Not Found)

```bash
curl -X GET http://localhost:8080/issue/999
```

#### 5. 완료된 이슈 수정 시도 (400 Bad Request)

```bash
# 먼저 이슈를 완료 상태로 변경
curl -X PATCH http://localhost:8080/issue/1 \
  -H "Content-Type: application/json" \
  -d '{"status": "COMPLETED"}'

# 완료된 이슈 수정 시도 (에러 발생)
curl -X PATCH http://localhost:8080/issue/1 \
  -H "Content-Type: application/json" \
  -d '{"title": "수정 시도"}'
```

## API 명세

### 데이터 모델

#### User

```json
{
  "id": 1,
  "name": "김개발"
}
```

#### Issue

```json
{
  "id": 1,
  "title": "버그 수정 필요",
  "description": "로그인 페이지에서 오류 발생",
  "status": "IN_PROGRESS",
  "user": {
    "id": 1,
    "name": "김개발"
  },
  "createdAt": "2025-06-11T10:00:00Z",
  "updatedAt": "2025-06-11T10:00:00Z"
}
```

### 상태값

- `PENDING`: 대기중 (담당자 미정)
- `IN_PROGRESS`: 진행중
- `COMPLETED`: 완료
- `CANCELLED`: 취소

### 기본 사용자

시스템에 미리 등록된 사용자:

1. 김개발 (ID: 1)
2. 이디자인 (ID: 2)
3. 박기획 (ID: 3)

## 비즈니스 규칙

### 1. 이슈 생성 규칙

- `title`은 필수 항목
- 담당자(`userId`)가 있으면 상태를 `IN_PROGRESS`로 설정
- 담당자가 없으면 상태를 `PENDING`으로 설정
- 존재하지 않는 사용자를 담당자로 지정할 수 없음

### 2. 이슈 수정 규칙

- `COMPLETED` 또는 `CANCELLED` 상태의 이슈는 수정 불가
- 담당자 없이 `IN_PROGRESS` 또는 `COMPLETED` 상태로 변경 불가
- `PENDING` 상태에서 담당자 할당 시 자동으로 `IN_PROGRESS`로 변경
- 담당자 제거 시 자동으로 `PENDING`으로 변경
- 요청 데이터에 명시되지 않은 필드는 업데이트하지 않음

### 3. 에러 처리

- 적절한 HTTP 상태 코드와 한국어 에러 메시지 반환
- 유효하지 않은 데이터에 대한 검증
- 비즈니스 규칙 위반 시 명확한 에러 메시지 제공

## 에러 응답 형식

모든 에러는 다음 형식으로 응답됩니다:

```json
{
  "error": "에러 메시지",
  "code": 400
}
```

### 주요 에러 코드 및 메시지

- `400 Bad Request`:
  - "잘못된 요청 데이터입니다"
  - "제목은 필수입니다"
  - "사용자를 찾을 수 없습니다"
  - "유효하지 않은 상태입니다"
  - "완료되거나 취소된 이슈는 수정할 수 없습니다"
  - "담당자 없이는 진행중 또는 완료 상태로 변경할 수 없습니다"
- `404 Not Found`:
  - "이슈를 찾을 수 없습니다"
- `500 Internal Server Error`:
  - "서버 내부 오류입니다"

## 테스트 실행

```bash
# 단위 테스트 실행
go test ./...

# 상세한 테스트 결과 보기
go test -v ./issue/...

# 특정 도메인 테스트
go test ./issue/model/...
go test ./issue/application/...
```
