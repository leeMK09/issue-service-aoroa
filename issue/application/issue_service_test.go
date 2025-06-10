package application

import (
	"testing"
	"issue-service-aoroa/issue/model"
	userModel "issue-service-aoroa/user/model"
)

type mockIssueRepository struct {
	issues []model.Issue
	lastID uint
}

func (m *mockIssueRepository) Create(issue model.Issue) model.Issue {
	m.lastID++
	issue.ID = m.lastID
	m.issues = append(m.issues, issue)
	return issue
}

func (m *mockIssueRepository) GetAll() []model.Issue {
	return m.issues
}

func (m *mockIssueRepository) GetByID(id uint) (*model.Issue, error) {
	for _, issue := range m.issues {
		if issue.ID == id {
			return &issue, nil
		}
	}
	return nil, nil
}

func (m *mockIssueRepository) Update(id uint, updatedIssue model.Issue) (*model.Issue, error) {
	for i, issue := range m.issues {
		if issue.ID == id {
			updatedIssue.ID = id
			m.issues[i] = updatedIssue
			return &m.issues[i], nil
		}
	}
	return nil, nil
}

func (m *mockIssueRepository) GetByStatus(status string) []model.Issue {
	var filtered []model.Issue
	for _, issue := range m.issues {
		if issue.Status == status {
			filtered = append(filtered, issue)
		}
	}
	return filtered
}

type mockUserRepository struct {
	users []userModel.User
}

func (m *mockUserRepository) GetByID(id uint) (*userModel.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return &user, nil
		}
	}
	return nil, nil
}

func (m *mockUserRepository) GetAll() []userModel.User {
	return m.users
}

func setupTestService() (IssueService, *mockIssueRepository, *mockUserRepository) {
	issueRepo := &mockIssueRepository{
		issues: []model.Issue{},
		lastID: 0,
	}
	
	userRepo := &mockUserRepository{
		users: []userModel.User{
			{ID: 1, Name: "김개발"},
			{ID: 2, Name: "이디자인"},
		},
	}
	
	service := NewIssueService(issueRepo, userRepo)
	return service, issueRepo, userRepo
}

func TestCreateIssue_실패_존재하지_않는_사용자(t *testing.T) {
	service, _, _ := setupTestService()
	
	nonExistentUserID := uint(999)
	_, err := service.CreateIssue("테스트 이슈", "설명", &nonExistentUserID)
	
	if err == nil {
		t.Error("존재하지 않는 사용자로 이슈 생성 시 에러가 발생해야 함")
	}
	
	expectedError := "사용자를 찾을 수 없습니다"
	if err.Error() != expectedError {
		t.Errorf("예상 에러: %s, 실제 에러: %s", expectedError, err.Error())
	}
}

func TestCreateIssue_실패_빈_제목(t *testing.T) {
	service, _, _ := setupTestService()
	
	_, err := service.CreateIssue("", "설명", nil)
	
	if err == nil {
		t.Error("빈 제목으로 이슈 생성 시 에러가 발생해야 함")
	}
	
	expectedError := "제목은 필수입니다"
	if err.Error() != expectedError {
		t.Errorf("예상 에러: %s, 실제 에러: %s", expectedError, err.Error())
	}
}

func TestCreateIssue_성공_담당자_있음(t *testing.T) {
	service, _, _ := setupTestService()
	
	userID := uint(1)
	issue, err := service.CreateIssue("테스트 이슈", "설명", &userID)
	
	if err != nil {
		t.Errorf("에러가 발생하지 않아야 함: %v", err)
	}
	
	if issue.Status != model.StatusInProgress {
		t.Errorf("담당자가 있는 이슈는 IN_PROGRESS 상태여야 함. 실제: %s", issue.Status)
	}
	
	if issue.User == nil || issue.User.ID != userID {
		t.Error("담당자가 올바르게 할당되어야 함")
	}
}

func TestUpdateIssue_실패_존재하지_않는_이슈(t *testing.T) {
	service, _, _ := setupTestService()
	
	updates := map[string]interface{}{
		"title": "새 제목",
	}
	
	_, err := service.UpdateIssue(999, updates)
	
	if err == nil {
		t.Error("존재하지 않는 이슈 업데이트 시 에러가 발생해야 함")
	}
	
	expectedError := "이슈를 찾을 수 없습니다"
	if err.Error() != expectedError {
		t.Errorf("예상 에러: %s, 실제 에러: %s", expectedError, err.Error())
	}
}

func TestUpdateIssue_실패_존재하지_않는_사용자_할당(t *testing.T) {
	service, _, _ := setupTestService()
	
	// 먼저 이슈 생성
	issue, _ := service.CreateIssue("테스트 이슈", "설명", nil)
	
	updates := map[string]interface{}{
		"userId": float64(999), // JSON에서 숫자는 float64로 파싱됨
	}
	
	_, err := service.UpdateIssue(issue.ID, updates)
	
	if err == nil {
		t.Error("존재하지 않는 사용자 할당 시 에러가 발생해야 함")
	}
	
	expectedError := "사용자를 찾을 수 없습니다"
	if err.Error() != expectedError {
		t.Errorf("예상 에러: %s, 실제 에러: %s", expectedError, err.Error())
	}
}

func TestUpdateIssue_실패_완료된_이슈_업데이트(t *testing.T) {
	service, issueRepo, _ := setupTestService()
	
	// 완료된 이슈 생성
	issue, _ := service.CreateIssue("테스트 이슈", "설명", nil)
	// 직접 완료 상태로 변경 (테스트를 위해)
	for i, iss := range issueRepo.issues {
		if iss.ID == issue.ID {
			issueRepo.issues[i].Status = model.StatusCompleted
			break
		}
	}
	
	updates := map[string]interface{}{
		"title": "새 제목",
	}
	
	_, err := service.UpdateIssue(issue.ID, updates)
	
	if err == nil {
		t.Error("완료된 이슈 업데이트 시 에러가 발생해야 함")
	}
	
	expectedError := "완료되거나 취소된 이슈는 수정할 수 없습니다"
	if err.Error() != expectedError {
		t.Errorf("예상 에러: %s, 실제 에러: %s", expectedError, err.Error())
	}
}

func TestGetIssuesByStatus_실패_유효하지_않은_상태(t *testing.T) {
	service, _, _ := setupTestService()
	
	_, err := service.GetIssuesByStatus("INVALID_STATUS")
	
	if err == nil {
		t.Error("유효하지 않은 상태로 필터링 시 에러가 발생해야 함")
	}
	
	expectedError := "유효하지 않은 상태입니다"
	if err.Error() != expectedError {
		t.Errorf("예상 에러: %s, 실제 에러: %s", expectedError, err.Error())
	}
}

func TestUpdateIssue_성공_담당자_할당_후_상태_전환(t *testing.T) {
	service, _, _ := setupTestService()
	
	// PENDING 상태의 이슈 생성
	issue, _ := service.CreateIssue("테스트 이슈", "설명", nil)
	
	// 담당자 할당
	updates := map[string]interface{}{
		"userId": float64(1),
	}
	
	updatedIssue, err := service.UpdateIssue(issue.ID, updates)
	
	if err != nil {
		t.Errorf("에러가 발생하지 않아야 함: %v", err)
	}
	
	if updatedIssue.Status != model.StatusInProgress {
		t.Errorf("담당자 할당 후 IN_PROGRESS 상태로 전환되어야 함. 실제: %s", updatedIssue.Status)
	}
	
	if updatedIssue.User == nil || updatedIssue.User.ID != 1 {
		t.Error("담당자가 올바르게 할당되어야 함")
	}
}

func TestUpdateIssue_성공_담당자_제거_후_PENDING_전환(t *testing.T) {
	service, _, _ := setupTestService()
	
	// 담당자가 있는 이슈 생성
	userID := uint(1)
	issue, _ := service.CreateIssue("테스트 이슈", "설명", &userID)
	
	// 담당자 제거
	updates := map[string]interface{}{
		"userId": nil,
	}
	
	updatedIssue, err := service.UpdateIssue(issue.ID, updates)
	
	if err != nil {
		t.Errorf("에러가 발생하지 않아야 함: %v", err)
	}
	
	if updatedIssue.Status != model.StatusPending {
		t.Errorf("담당자 제거 후 PENDING 상태로 전환되어야 함. 실제: %s", updatedIssue.Status)
	}
	
	if updatedIssue.User != nil {
		t.Error("담당자가 제거되어야 함")
	}
}