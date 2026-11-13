package logic

import "testing"

func TestTicketContactPhone(t *testing.T) {
	for _, tt := range []struct {
		name, contacts, account, want string
	}{
		{"first contact", "张三 13800138000,李四 13900139000,王五 13600136000", "13700137000", "13800138000"},
		{"first available number", "张三,李四：13900139000", "13700137000", "13900139000"},
		{"account fallback", " , ", "13700137000", "13700137000"},
		{"configured without phone", "张三", "13700137000", ""},
		{"reject longer digit sequence", "张三 138001380001", "", ""},
		{"empty", "", "", ""},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := ticketContactPhone(tt.contacts, tt.account); got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
