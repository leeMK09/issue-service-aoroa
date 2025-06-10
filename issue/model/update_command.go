package model

import (
	"errors"
	userModel "issue-service-aoroa/user/model"
)

type UpdateCommand struct {
	Title       *string
	Description *string
	Status      *string
	UserID      *uint
	User        *userModel.User
}

func (cmd *UpdateCommand) ApplyTo(issue *Issue) error {
	if !issue.IsUpdatable() {
		return errors.New("완료되거나 취소된 이슈는 수정할 수 없습니다")
	}

	if err := issue.UpdateDetails(cmd.Title, cmd.Description); err != nil {
		return err
	}

	if cmd.UserID != nil {
		if *cmd.UserID == 0 || cmd.User == nil {
			if err := issue.Unassign(); err != nil {
				return err
			}
		} else {
			if err := issue.AssignTo(cmd.User); err != nil {
				return err
			}
		}
	}

	if cmd.Status != nil {
		if err := issue.ChangeStatus(*cmd.Status); err != nil {
			return err
		}
	}

	return nil
}

func NewUpdateCommand() *UpdateCommand {
	return &UpdateCommand{}
}

func (cmd *UpdateCommand) WithTitle(title string) *UpdateCommand {
	cmd.Title = &title
	return cmd
}

func (cmd *UpdateCommand) WithDescription(description string) *UpdateCommand {
	cmd.Description = &description
	return cmd
}

func (cmd *UpdateCommand) WithStatus(status string) *UpdateCommand {
	cmd.Status = &status
	return cmd
}

func (cmd *UpdateCommand) WithUser(userID uint, user *userModel.User) *UpdateCommand {
	cmd.UserID = &userID
	cmd.User = user
	return cmd
}

func (cmd *UpdateCommand) WithoutUser() *UpdateCommand {
	userID := uint(0)
	cmd.UserID = &userID
	cmd.User = nil
	return cmd
}