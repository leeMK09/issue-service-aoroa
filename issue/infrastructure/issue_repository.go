package infrastructure

import (
	"time"
	issueModel "issue-service-aoroa/issue/model"
)

type IssueRepository interface {
	Create(issue issueModel.Issue) issueModel.Issue
	GetAll() []issueModel.Issue
	GetByID(id uint) (*issueModel.Issue, error)
	Update(id uint, issue issueModel.Issue) (*issueModel.Issue, error)
	GetByStatus(status string) []issueModel.Issue
}

type issueRepository struct {
	issues []issueModel.Issue
	lastID uint
}

func NewIssueRepository() IssueRepository {
	return &issueRepository{
		issues: []issueModel.Issue{},
		lastID: 0,
	}
}

func (r *issueRepository) Create(issue issueModel.Issue) issueModel.Issue {
	r.lastID++
	issue.ID = r.lastID
	issue.CreatedAt = time.Now()
	issue.UpdatedAt = time.Now()
	r.issues = append(r.issues, issue)
	return issue
}

func (r *issueRepository) GetAll() []issueModel.Issue {
	return r.issues
}

func (r *issueRepository) GetByID(id uint) (*issueModel.Issue, error) {
	for _, issue := range r.issues {
		if issue.ID == id {
			return &issue, nil
		}
	}
	return nil, nil
}

func (r *issueRepository) Update(id uint, updatedIssue issueModel.Issue) (*issueModel.Issue, error) {
	for i, issue := range r.issues {
		if issue.ID == id {
			updatedIssue.ID = id
			updatedIssue.CreatedAt = issue.CreatedAt
			updatedIssue.UpdatedAt = time.Now()
			r.issues[i] = updatedIssue
			return &r.issues[i], nil
		}
	}
	return nil, nil
}

func (r *issueRepository) GetByStatus(status string) []issueModel.Issue {
	var filtered []issueModel.Issue
	for _, issue := range r.issues {
		if issue.Status == status {
			filtered = append(filtered, issue)
		}
	}
	return filtered
}