package application

import (
	"errors"
	"issue-service-aoroa/issue/infrastructure"
	"issue-service-aoroa/issue/model"
	userInfra "issue-service-aoroa/user/infrastructure"
	userModel "issue-service-aoroa/user/model"
)

type IssueService interface {
	CreateIssue(title, description string, userID *uint) (*model.Issue, error)
	GetAllIssues() []model.Issue
	GetIssueByID(id uint) (*model.Issue, error)
	UpdateIssue(id uint, updates map[string]interface{}) (*model.Issue, error)
	GetIssuesByStatus(status string) ([]model.Issue, error)
}

type issueService struct {
	issueRepo infrastructure.IssueRepository
	userRepo  userInfra.UserRepository
}

func NewIssueService(issueRepo infrastructure.IssueRepository, userRepo userInfra.UserRepository) IssueService {
	return &issueService{
		issueRepo: issueRepo,
		userRepo:  userRepo,
	}
}

func (s *issueService) CreateIssue(title, description string, userID *uint) (*model.Issue, error) {
	var assignee *userModel.User

	if userID != nil {
		user, err := s.findUserByID(*userID)
		if err != nil {
			return nil, err
		}
		assignee = user
	}

	issue, err := model.NewIssue(title, description, assignee)
	if err != nil {
		return nil, err
	}

	createdIssue := s.issueRepo.Create(*issue)
	return &createdIssue, nil
}

func (s *issueService) GetAllIssues() []model.Issue {
	return s.issueRepo.GetAll()
}

func (s *issueService) GetIssueByID(id uint) (*model.Issue, error) {
	return s.issueRepo.GetByID(id)
}

func (s *issueService) UpdateIssue(id uint, updates map[string]interface{}) (*model.Issue, error) {
	existingIssue, err := s.findIssueByID(id)
	if err != nil {
		return nil, err
	}

	updateCommand, err := s.buildUpdateCommand(updates)
	if err != nil {
		return nil, err
	}

	if err := updateCommand.ApplyTo(existingIssue); err != nil {
		return nil, err
	}

	result, err := s.issueRepo.Update(id, *existingIssue)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *issueService) GetIssuesByStatus(status string) ([]model.Issue, error) {
	if !model.IsValidStatus(status) {
		return nil, errors.New("유효하지 않은 상태입니다")
	}
	return s.issueRepo.GetByStatus(status), nil
}

func (s *issueService) findIssueByID(id uint) (*model.Issue, error) {
	issue, err := s.issueRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if issue == nil {
		return nil, errors.New("이슈를 찾을 수 없습니다")
	}
	return issue, nil
}

func (s *issueService) findUserByID(userID uint) (*userModel.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("사용자를 찾을 수 없습니다")
	}
	return user, nil
}

func (s *issueService) buildUpdateCommand(updates map[string]interface{}) (*model.UpdateCommand, error) {
	cmd := model.NewUpdateCommand()

	if title, ok := updates["title"]; ok {
		if titleStr, ok := title.(string); ok {
			cmd.WithTitle(titleStr)
		}
	}

	if description, ok := updates["description"]; ok {
		if descStr, ok := description.(string); ok {
			cmd.WithDescription(descStr)
		}
	}

	if status, ok := updates["status"]; ok {
		if statusStr, ok := status.(string); ok {
			cmd.WithStatus(statusStr)
		}
	}

	if userID, ok := updates["userId"]; ok {
		if userID == nil {
			cmd.WithoutUser()
		} else if userIDFloat, ok := userID.(float64); ok {
			userIDUint := uint(userIDFloat)
			user, err := s.findUserByID(userIDUint)
			if err != nil {
				return nil, err
			}
			cmd.WithUser(userIDUint, user)
		}
	}

	return cmd, nil
}