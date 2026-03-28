package domain

// Management は入金管理エンティティ
type Management struct {
	ID         int
	BankStatus BankStatus
}

// NewManagement は入金管理エンティティを生成する
func NewManagement(id int, bankStatus BankStatus) *Management {
	return &Management{
		ID:         id,
		BankStatus: bankStatus,
	}
}
