package middleware

import "testing"

func TestGetToken(t *testing.T) {
	token := `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJwaG9uZSI6IjEyMzEyMzEyMzEyIiwidXVpZCI6IjczZjJlMmM4LTc3ODQtNDA4YS05YmI3LWY2NmUyMjkwMjIxYSIsInN1YiI6ImNhaXNodWppIiwiYXVkIjpbImNhaXNodWppIl0sImV4cCI6MTcxNDUxNDE5OCwibmJmIjoxNzEzMDc0MTk4LCJpYXQiOjE3MTMwNzQxOTgsImp0aSI6IjEifQ.jBWQKmU-juMnkDrq-qZYyffjiaSWcwgniLBHnBEK-l8`
	c, err := GetTokenClaims(token)
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(c.ExpiresAt)
}
