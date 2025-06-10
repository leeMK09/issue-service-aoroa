package model

import (
	"errors"
	"time"
	userModel "issue-service-aoroa/user/model"
)

type Issue struct {
	ID          uint               `json:"id"`
	Title       string             `json:"title"`
	Description string             `json:"description"`
	Status      string             `json:"status"`
	User        *userModel.User    `json:"user,omitempty"`
	CreatedAt   time.Time          `json:"createdAt"`
	UpdatedAt   time.Time          `json:"updatedAt"`
}

const (
	StatusPending     = "PENDING"
	StatusInProgress  = "IN_PROGRESS"
	StatusCompleted   = "COMPLETED"
	StatusCancelled   = "CANCELLED"
)

func NewIssue(title, description string, assignee *userModel.User) (*Issue, error) {
	if err := validateTitle(title); err != nil {
		return nil, err
	}

	issue := &Issue{
		Title:       title,
		Description: description,
		User:        assignee,
	}

	issue.setInitialStatus()
	return issue, nil
}

func (i *Issue) IsUpdatable() bool {
	return i.Status != StatusCompleted && i.Status != StatusCancelled
}

func (i *Issue) AssignTo(user *userModel.User) error {
	if !i.IsUpdatable() {
		return errors.New("완료되거나 취소된 이슈는 수정할 수 없습니다")
	}

	wasUnassigned := i.wasUnassigned()
	i.User = user

	if wasUnassigned && user != nil && i.isPending() {
		i.Status = StatusInProgress
	}

	return nil
}

func (i *Issue) Unassign() error {
	if !i.IsUpdatable() {
		return errors.New("완료되거나 취소된 이슈는 수정할 수 없습니다")
	}

	i.User = nil
	i.Status = StatusPending
	return nil
}

func (i *Issue) ChangeStatus(newStatus string) error {
	if !i.IsUpdatable() {
		return errors.New("완료되거나 취소된 이슈는 수정할 수 없습니다")
	}

	if !IsValidStatus(newStatus) {
		return errors.New("유효하지 않은 상태입니다")
	}

	if err := i.validateStatusTransition(newStatus); err != nil {
		return err
	}

	i.Status = newStatus
	return nil
}

func (i *Issue) UpdateDetails(title, description *string) error {
	if !i.IsUpdatable() {
		return errors.New("완료되거나 취소된 이슈는 수정할 수 없습니다")
	}

	if title != nil {
		if err := validateTitle(*title); err != nil {
			return err
		}
		i.Title = *title
	}

	if description != nil {
		i.Description = *description
	}

	return nil
}

func (i *Issue) setInitialStatus() {
	if i.hasAssignee() {
		i.Status = StatusInProgress
	} else {
		i.Status = StatusPending
	}
}

func (i *Issue) validateStatusTransition(newStatus string) error {
	if !i.hasAssignee() && i.requiresAssignee(newStatus) {
		return errors.New("담당자 없이는 진행중 또는 완료 상태로 변경할 수 없습니다")
	}
	return nil
}

func validateTitle(title string) error {
	if title == "" {
		return errors.New("제목은 필수입니다")
	}
	return nil
}

func IsValidStatus(status string) bool {
	return status == StatusPending || status == StatusInProgress || 
		   status == StatusCompleted || status == StatusCancelled
}

func (i *Issue) wasUnassigned() bool {
	return i.User == nil
}

func (i *Issue) hasAssignee() bool {
	return i.User != nil
}

func (i *Issue) isPending() bool {
	return i.Status == StatusPending
}

func (i *Issue) requiresAssignee(status string) bool {
	return status == StatusInProgress || status == StatusCompleted
}