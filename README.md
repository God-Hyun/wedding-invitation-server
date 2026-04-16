# 모바일 청첩장 서버

신호현 ❤️ 윤유진 결혼식 모바일 청첩장의 백엔드 서버입니다.

- **프론트엔드**: https://wedding-invitation-woad-nu.vercel.app
- **배포**: Render (Free tier)
- **관리자 페이지**: https://wedding-invitation-server-d373.onrender.com/admin

## 기술 스택

- Go 1.18+
- SQLite (github.com/mattn/go-sqlite3)
- CORS (github.com/rs/cors)

## 제공 기능

- 방명록 작성 / 조회 / 삭제 API
  - 관리자 비밀번호(`ADMIN_PASSWORD`)로 강제 삭제 가능
- 참석 의사 전달 API
- 관리자 페이지 (`/admin`) — 참석 명단 조회 및 삭제

## 환경변수

| 변수 | 설명 |
|------|------|
| `ALLOW_ORIGIN` | 허용할 프론트엔드 도메인 (예: `https://xxx.vercel.app`) |
| `ADMIN_PASSWORD` | 관리자 비밀번호 (방명록 강제 삭제, /admin 로그인) |
| `PORT` | 서버 포트 (기본값: 8080) |

## 로컬 실행

```bash
git clone https://github.com/God-Hyun/wedding-invitation-server.git
cd wedding-invitation-server
go mod download

# .env 파일 생성
echo "ALLOW_ORIGIN=http://localhost:3000" > .env
echo "ADMIN_PASSWORD=your_password" >> .env

go run app.go
# → http://localhost:8080
```

## API 엔드포인트

| 메서드 | 경로 | 설명 |
|--------|------|------|
| GET | `/api/guestbook?offset=0&limit=5` | 방명록 조회 |
| POST | `/api/guestbook` | 방명록 작성 |
| PUT | `/api/guestbook` | 방명록 삭제 (비밀번호 필요) |
| POST | `/api/attendance` | 참석 의사 전달 |
| GET | `/admin` | 관리자 로그인 페이지 |
| POST | `/admin` | 참석 명단 조회 / 삭제 |

## Render 배포

- Build Command: `go build -tags netgo -ldflags '-s -w' -o app`
- Start Command: `./app`
- 환경변수: `ALLOW_ORIGIN`, `ADMIN_PASSWORD`
