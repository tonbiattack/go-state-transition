package persistence

import (
	"fmt"

	"github.com/go-state-transition/go-state-transition/internal/domain"
)

// managementData は永続化層で使う中間データ構造
// DBレコードを直接ドメインへ渡さず、まずこの構造へ詰め替える
type managementData struct {
	ID             int
	BankStatusCode int
}

// newManagementDataFromRecord はDBレコードを中間データへ変換する
func newManagementDataFromRecord(record managementRecord) managementData {
	return managementData{
		ID:             int(record.ID),
		BankStatusCode: record.BankStatus,
	}
}

// ToDomain は中間データをドメインエンティティへ変換する
func (d managementData) ToDomain() (*domain.Management, error) {
	status, err := domain.NewBankStatus(d.BankStatusCode)
	if err != nil {
		return nil, fmt.Errorf("ステータス変換エラー: %w", err)
	}

	return domain.NewManagement(d.ID, status), nil
}
