package domain

import "fmt"

// BankStatus は入金管理の内部ステータスを表す型
type BankStatus int

const (
	BankStatusAccepted   BankStatus = 1 // 受付
	BankStatusProcessing BankStatus = 2 // 処理中
	BankStatusCompleted  BankStatus = 3 // 完了
	BankStatusSuspended  BankStatus = 4 // 停止
)

// 許可された遷移の一覧
// キー：現在のステータス、値：遷移先として許可されたステータスの集合
var allowedTransitions = map[BankStatus][]BankStatus{
	BankStatusAccepted:   {BankStatusProcessing},
	BankStatusProcessing: {BankStatusCompleted, BankStatusSuspended},
	BankStatusCompleted:  {}, // 終端：遷移不可
	BankStatusSuspended:  {}, // 終端：遷移不可
}

// Label はステータスの表示名を返す
func (s BankStatus) Label() string {
	switch s {
	case BankStatusAccepted:
		return "受付"
	case BankStatusProcessing:
		return "処理中"
	case BankStatusCompleted:
		return "完了"
	case BankStatusSuspended:
		return "停止"
	default:
		return "不明"
	}
}

// IsValid はステータスが有効な値かどうかを返す
func (s BankStatus) IsValid() bool {
	switch s {
	case BankStatusAccepted, BankStatusProcessing,
		BankStatusCompleted, BankStatusSuspended:
		return true
	}
	return false
}

// CanTransitionTo は次のステータスへの遷移が許可されているかを返す
func (s BankStatus) CanTransitionTo(next BankStatus) bool {
	for _, allowed := range allowedTransitions[s] {
		if allowed == next {
			return true
		}
	}
	return false
}

// AllowedNextStatuses は現在のステータスから遷移できるステータスの一覧を返す
func (s BankStatus) AllowedNextStatuses() []BankStatus {
	return allowedTransitions[s]
}

// NewBankStatus はDBの数値から BankStatus に変換する
// 不正な値は error を返す
func NewBankStatus(code int) (BankStatus, error) {
	s := BankStatus(code)
	if !s.IsValid() {
		return 0, fmt.Errorf("不正なステータス値: %d", code)
	}
	return s, nil
}
