package httphandler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/juhonamnam/wedding-invitation-server/env"
	"github.com/juhonamnam/wedding-invitation-server/sqldb"
)

type AdminHandler struct {
	http.Handler
}

func (h *AdminHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(loginPage))
		return
	}

	if r.Method == http.MethodPost {
		r.ParseForm()
		password := r.FormValue("password")
		action := r.FormValue("action")

		if env.AdminPassword == "" || password != env.AdminPassword {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(loginPageWithError))
			return
		}

		if action == "delete" {
			idStr := r.FormValue("id")
			id, err := strconv.Atoi(idStr)
			if err == nil {
				sqldb.DeleteAttendance(id)
			}
		}

		rows, err := sqldb.GetAllAttendance()
		if err != nil {
			http.Error(w, "DB 오류", http.StatusInternalServerError)
			return
		}

		totalCount := 0
		groomCount := 0
		brideCount := 0
		mealCount := 0

		tableRows := ""
		for _, row := range rows {
			t := time.Unix(row.Timestamp, 0).In(time.FixedZone("KST", 9*60*60))
			mealLabel := row.Meal
			switch row.Meal {
			case "expected":
				mealLabel = "식사 예정"
				mealCount += row.Count
			case "unexpected":
				mealLabel = "식사 미정"
			case "no":
				mealLabel = "불참"
			}
			sideLabel := row.Side
			if row.Side == "groom" {
				sideLabel = "신랑측"
				groomCount += row.Count
			} else if row.Side == "bride" {
				sideLabel = "신부측"
				brideCount += row.Count
			}
			totalCount += row.Count
			tableRows += fmt.Sprintf(`<tr>
				<td>%s</td>
				<td>%s</td>
				<td>%s</td>
				<td>%d명</td>
				<td>%s</td>
				<td><form method="POST" onsubmit="return confirm('삭제하시겠습니까?')">
					<input type="hidden" name="password" value="%s">
					<input type="hidden" name="action" value="delete">
					<input type="hidden" name="id" value="%d">
					<button type="submit" style="background:#e55;color:white;border:none;padding:4px 10px;border-radius:4px;cursor:pointer">삭제</button>
				</form></td>
			</tr>`, sideLabel, row.Name, mealLabel, row.Count, t.Format("01/02 15:04"), password, row.Id)
		}

		html := fmt.Sprintf(attendancePage, totalCount, groomCount, brideCount, mealCount, tableRows)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write([]byte(html))
		return
	}

	w.WriteHeader(http.StatusMethodNotAllowed)
}

const loginPage = `<!DOCTYPE html>
<html lang="ko">
<head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1">
<title>관리자 로그인</title>
<style>
  body { font-family: sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; background: #f5f5f5; }
  .box { background: white; padding: 40px; border-radius: 12px; box-shadow: 0 2px 12px rgba(0,0,0,0.1); text-align: center; }
  h2 { margin-bottom: 24px; color: #333; }
  input { padding: 10px 16px; font-size: 16px; border: 1px solid #ddd; border-radius: 6px; width: 200px; margin-bottom: 12px; }
  button { padding: 10px 32px; font-size: 16px; background: #c28080; color: white; border: none; border-radius: 6px; cursor: pointer; }
  button:hover { background: #a06060; }
</style>
</head>
<body>
<div class="box">
  <h2>관리자 로그인</h2>
  <form method="POST">
    <input type="password" name="password" placeholder="관리자 비밀번호" autofocus /><br>
    <button type="submit">확인</button>
  </form>
</div>
</body></html>`

const loginPageWithError = `<!DOCTYPE html>
<html lang="ko">
<head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1">
<title>관리자 로그인</title>
<style>
  body { font-family: sans-serif; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; background: #f5f5f5; }
  .box { background: white; padding: 40px; border-radius: 12px; box-shadow: 0 2px 12px rgba(0,0,0,0.1); text-align: center; }
  h2 { margin-bottom: 24px; color: #333; }
  input { padding: 10px 16px; font-size: 16px; border: 1px solid #ddd; border-radius: 6px; width: 200px; margin-bottom: 12px; }
  button { padding: 10px 32px; font-size: 16px; background: #c28080; color: white; border: none; border-radius: 6px; cursor: pointer; }
  button:hover { background: #a06060; }
  .error { color: #e55; margin-bottom: 12px; font-size: 14px; }
</style>
</head>
<body>
<div class="box">
  <h2>관리자 로그인</h2>
  <p class="error">비밀번호가 올바르지 않습니다.</p>
  <form method="POST">
    <input type="password" name="password" placeholder="관리자 비밀번호" autofocus /><br>
    <button type="submit">확인</button>
  </form>
</div>
</body></html>`

const attendancePage = `<!DOCTYPE html>
<html lang="ko">
<head><meta charset="utf-8"><meta name="viewport" content="width=device-width, initial-scale=1">
<title>참석 의사 현황</title>
<style>
  body { font-family: sans-serif; padding: 24px; background: #f5f5f5; }
  h1 { color: #333; margin-bottom: 24px; }
  .summary { display: flex; gap: 16px; margin-bottom: 24px; flex-wrap: wrap; }
  .card { background: white; padding: 16px 24px; border-radius: 10px; box-shadow: 0 1px 6px rgba(0,0,0,0.1); text-align: center; }
  .card .label { font-size: 13px; color: #888; margin-bottom: 4px; }
  .card .value { font-size: 28px; font-weight: bold; color: #c28080; }
  table { width: 100%%; border-collapse: collapse; background: white; border-radius: 10px; overflow: hidden; box-shadow: 0 1px 6px rgba(0,0,0,0.1); }
  th { background: #c28080; color: white; padding: 12px 16px; text-align: left; }
  td { padding: 12px 16px; border-bottom: 1px solid #f0f0f0; }
  tr:last-child td { border-bottom: none; }
  tr:hover td { background: #fafafa; }
</style>
</head>
<body>
<h1>참석 의사 현황</h1>
<div class="summary">
  <div class="card"><div class="label">총 참석 인원</div><div class="value">%d명</div></div>
  <div class="card"><div class="label">신랑측</div><div class="value">%d명</div></div>
  <div class="card"><div class="label">신부측</div><div class="value">%d명</div></div>
  <div class="card"><div class="label">식사 예정</div><div class="value">%d명</div></div>
</div>
<table>
  <thead><tr><th>구분</th><th>성함</th><th>식사</th><th>인원</th><th>제출 시간</th><th></th></tr></thead>
  <tbody>%s</tbody>
</table>
</body></html>`
