package domain

import (
	"testing"
)

// TestBankStatus_Label は各ステータスの表示名を確認する
func TestBankStatus_Label(t *testing.T) {
	tests := []struct {
		name   string
		status BankStatus
		want   string
	}{
		{"受付の表示名", BankStatusAccepted, "受付"},
		{"処理中の表示名", BankStatusProcessing, "処理中"},
		{"完了の表示名", BankStatusCompleted, "完了"},
		{"停止の表示名", BankStatusSuspended, "停止"},
		{"不明な値の表示名", BankStatus(99), "不明"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.status.Label()
			if got != tt.want {
				t.Errorf("Label() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestBankStatus_IsValid は有効・無効なステータス値を確認する
func TestBankStatus_IsValid(t *testing.T) {
	tests := []struct {
		name   string
		status BankStatus
		want   bool
	}{
		{"受付は有効", BankStatusAccepted, true},
		{"処理中は有効", BankStatusProcessing, true},
		{"完了は有効", BankStatusCompleted, true},
		{"停止は有効", BankStatusSuspended, true},
		{"0は無効", BankStatus(0), false},
		{"5は無効", BankStatus(5), false},
		{"9は無効", BankStatus(9), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.status.IsValid()
			if got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestBankStatus_CanTransitionTo は遷移の許可・禁止を確認する
func TestBankStatus_CanTransitionTo(t *testing.T) {
	tests := []struct {
		name    string
		from    BankStatus
		to      BankStatus
		allowed bool
	}{
		// 許可される遷移
		{"受付から処理中への遷移は許可される", BankStatusAccepted, BankStatusProcessing, true},
		{"処理中から完了への遷移は許可される", BankStatusProcessing, BankStatusCompleted, true},
		{"処理中から停止への遷移は許可される", BankStatusProcessing, BankStatusSuspended, true},
		// 禁止される遷移
		{"受付から完了への遷移は禁止される", BankStatusAccepted, BankStatusCompleted, false},
		{"受付から停止への遷移は禁止される", BankStatusAccepted, BankStatusSuspended, false},
		{"完了から処理中への遷移は禁止される", BankStatusCompleted, BankStatusProcessing, false},
		{"停止から受付への遷移は禁止される", BankStatusSuspended, BankStatusAccepted, false},
		{"完了から受付への遷移は禁止される", BankStatusCompleted, BankStatusAccepted, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.from.CanTransitionTo(tt.to)
			if got != tt.allowed {
				t.Errorf("CanTransitionTo() = %v, want %v", got, tt.allowed)
			}
		})
	}
}

// TestBankStatus_AllowedNextStatuses は遷移可能なステータス一覧を確認する
func TestBankStatus_AllowedNextStatuses(t *testing.T) {
	tests := []struct {
		name   string
		status BankStatus
		want   []BankStatus
	}{
		{"受付から遷移できるのは処理中のみ", BankStatusAccepted, []BankStatus{BankStatusProcessing}},
		{"処理中から遷移できるのは完了と停止", BankStatusProcessing, []BankStatus{BankStatusCompleted, BankStatusSuspended}},
		{"完了からは遷移できない", BankStatusCompleted, []BankStatus{}},
		{"停止からは遷移できない", BankStatusSuspended, []BankStatus{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.status.AllowedNextStatuses()
			if len(got) != len(tt.want) {
				t.Errorf("AllowedNextStatuses() の件数 = %d, want %d", len(got), len(tt.want))
				return
			}
			for i, s := range got {
				if s != tt.want[i] {
					t.Errorf("AllowedNextStatuses()[%d] = %v, want %v", i, s, tt.want[i])
				}
			}
		})
	}
}

// TestNewBankStatus は数値からの変換を確認する
func TestNewBankStatus(t *testing.T) {
	tests := []struct {
		name    string
		code    int
		want    BankStatus
		wantErr bool
	}{
		{"1は受付に変換される", 1, BankStatusAccepted, false},
		{"2は処理中に変換される", 2, BankStatusProcessing, false},
		{"3は完了に変換される", 3, BankStatusCompleted, false},
		{"4は停止に変換される", 4, BankStatusSuspended, false},
		{"0はエラーになる", 0, 0, true},
		{"5はエラーになる", 5, 0, true},
		{"負の値はエラーになる", -1, 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewBankStatus(tt.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewBankStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("NewBankStatus() = %v, want %v", got, tt.want)
			}
		})
	}
}
