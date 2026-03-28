package usecase

import (
	"fmt"

	"github.com/go-state-transition/go-state-transition/internal/domain"
)

// ManagementRepository は入金管理の永続化インターフェース
type ManagementRepository interface {
	FindByID(id int) (*domain.Management, error)
	Save(m *domain.Management) error
}

// ManagementUsecase は入金管理のユースケース
type ManagementUsecase struct {
	repo ManagementRepository
}

// NewManagementUsecase は ManagementUsecase を生成する
func NewManagementUsecase(repo ManagementRepository) *ManagementUsecase {
	return &ManagementUsecase{repo: repo}
}

// StartProcessing は受付中の入金管理を処理中に遷移させる
func (u *ManagementUsecase) StartProcessing(id int) error {
	m, err := u.repo.FindByID(id)
	if err != nil {
		return err
	}
	if !m.BankStatus.CanTransitionTo(domain.BankStatusProcessing) {
		return fmt.Errorf("ステータス %s から処理中への遷移は許可されていません", m.BankStatus.Label())
	}
	m.BankStatus = domain.BankStatusProcessing
	return u.repo.Save(m)
}

// Complete は処理中の入金管理を完了に遷移させる
func (u *ManagementUsecase) Complete(id int) error {
	m, err := u.repo.FindByID(id)
	if err != nil {
		return err
	}
	if !m.BankStatus.CanTransitionTo(domain.BankStatusCompleted) {
		return fmt.Errorf("ステータス %s から完了への遷移は許可されていません", m.BankStatus.Label())
	}
	m.BankStatus = domain.BankStatusCompleted
	return u.repo.Save(m)
}

// Suspend は処理中の入金管理を停止に遷移させる
func (u *ManagementUsecase) Suspend(id int) error {
	m, err := u.repo.FindByID(id)
	if err != nil {
		return err
	}
	if !m.BankStatus.CanTransitionTo(domain.BankStatusSuspended) {
		return fmt.Errorf("ステータス %s から停止への遷移は許可されていません", m.BankStatus.Label())
	}
	m.BankStatus = domain.BankStatusSuspended
	return u.repo.Save(m)
}
