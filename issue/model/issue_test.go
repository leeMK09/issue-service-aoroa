package model

import (
	"testing"
	userModel "issue-service-aoroa/user/model"
)

func TestNewIssue_실패_제목이_비어있음(t *testing.T) {
	_, err := NewIssue("", "설명", nil)
	
	if err == nil {
		t.Error("제목이 비어있을 때 에러가 발생해야 함")
	}
	
	expectedError := "제목은 필수입니다"
	if err.Error() != expectedError {
		t.Errorf("예상 에러: %s, 실제 에러: %s", expectedError, err.Error())
	}
}

func TestNewIssue_성공_담당자_없음(t *testing.T) {
	issue, err := NewIssue("테스트 이슈", "설명", nil)
	
	if err != nil {
		t.Errorf("에러가 발생하지 않아야 함: %v", err)
	}
	
	if issue.Status != StatusPending {
		t.Errorf("담당자가 없으면 PENDING 상태여야 함. 실제: %s", issue.Status)
	}
}

func TestNewIssue_성공_담당자_있음(t *testing.T) {
	user := &userModel.User{ID: 1, Name: "테스트 사용자"}
	issue, err := NewIssue("테스트 이슈", "설명", user)
	
	if err != nil {
		t.Errorf("에러가 발생하지 않아야 함: %v", err)
	}
	
	if issue.Status != StatusInProgress {
		t.Errorf("담당자가 있으면 IN_PROGRESS 상태여야 함. 실제: %s", issue.Status)
	}
}

func TestIsUpdatable_완료된_이슈는_업데이트_불가(t *testing.T) {
	issue, _ := NewIssue("테스트", "설명", nil)
	issue.Status = StatusCompleted
	
	if issue.IsUpdatable() {
		t.Error("완료된 이슈는 업데이트할 수 없어야 함")
	}
}

func TestIsUpdatable_취소된_이슈는_업데이트_불가(t *testing.T) {
	issue, _ := NewIssue("테스트", "설명", nil)
	issue.Status = StatusCancelled
	
	if issue.IsUpdatable() {
		t.Error("취소된 이슈는 업데이트할 수 없어야 함")
	}
}

func TestAssignTo_실패_완료된_이슈(t *testing.T) {
	issue, _ := NewIssue("테스트", "설명", nil)
	issue.Status = StatusCompleted
	user := &userModel.User{ID: 1, Name: "테스트"}
	
	err := issue.AssignTo(user)
	
	if err == nil {
		t.Error("완료된 이슈에 담당자 할당 시 에러가 발생해야 함")
	}
	
	expectedError := "완료되거나 취소된 이슈는 수정할 수 없습니다"
	if err.Error() != expectedError {
		t.Errorf("예상 에러: %s, 실제 에러: %s", expectedError, err.Error())
	}
}

func TestAssignTo_성공_PENDING에서_IN_PROGRESS로_전환(t *testing.T) {
	issue, _ := NewIssue("테스트", "설명", nil)
	user := &userModel.User{ID: 1, Name: "테스트"}
	
	err := issue.AssignTo(user)
	
	if err != nil {
		t.Errorf("에러가 발생하지 않아야 함: %v", err)
	}
	
	if issue.Status != StatusInProgress {
		t.Errorf("PENDING에서 담당자 할당 시 IN_PROGRESS로 전환되어야 함. 실제: %s", issue.Status)
	}
	
	if issue.User != user {
		t.Error("담당자가 올바르게 할당되어야 함")
	}
}

func TestChangeStatus_실패_담당자_없이_IN_PROGRESS(t *testing.T) {
	issue, _ := NewIssue("테스트", "설명", nil)
	
	err := issue.ChangeStatus(StatusInProgress)
	
	if err == nil {
		t.Error("담당자 없이 IN_PROGRESS로 변경 시 에러가 발생해야 함")
	}
	
	expectedError := "담당자 없이는 진행중 또는 완료 상태로 변경할 수 없습니다"
	if err.Error() != expectedError {
		t.Errorf("예상 에러: %s, 실제 에러: %s", expectedError, err.Error())
	}
}

func TestChangeStatus_실패_담당자_없이_COMPLETED(t *testing.T) {
	issue, _ := NewIssue("테스트", "설명", nil)
	
	err := issue.ChangeStatus(StatusCompleted)
	
	if err == nil {
		t.Error("담당자 없이 COMPLETED로 변경 시 에러가 발생해야 함")
	}
}

func TestChangeStatus_실패_유효하지_않은_상태(t *testing.T) {
	issue, _ := NewIssue("테스트", "설명", nil)
	
	err := issue.ChangeStatus("INVALID_STATUS")
	
	if err == nil {
		t.Error("유효하지 않은 상태로 변경 시 에러가 발생해야 함")
	}
	
	expectedError := "유효하지 않은 상태입니다"
	if err.Error() != expectedError {
		t.Errorf("예상 에러: %s, 실제 에러: %s", expectedError, err.Error())
	}
}

func TestUnassign_성공_PENDING으로_전환(t *testing.T) {
	user := &userModel.User{ID: 1, Name: "테스트"}
	issue, _ := NewIssue("테스트", "설명", user)
	
	err := issue.Unassign()
	
	if err != nil {
		t.Errorf("에러가 발생하지 않아야 함: %v", err)
	}
	
	if issue.User != nil {
		t.Error("담당자가 제거되어야 함")
	}
	
	if issue.Status != StatusPending {
		t.Errorf("담당자 제거 시 PENDING 상태가 되어야 함. 실제: %s", issue.Status)
	}
}

func TestUpdateDetails_실패_완료된_이슈(t *testing.T) {
	issue, _ := NewIssue("테스트", "설명", nil)
	issue.Status = StatusCompleted
	
	newTitle := "새 제목"
	err := issue.UpdateDetails(&newTitle, nil)
	
	if err == nil {
		t.Error("완료된 이슈의 세부사항 업데이트 시 에러가 발생해야 함")
	}
}

func TestUpdateDetails_실패_빈_제목(t *testing.T) {
	issue, _ := NewIssue("테스트", "설명", nil)
	
	emptyTitle := ""
	err := issue.UpdateDetails(&emptyTitle, nil)
	
	if err == nil {
		t.Error("빈 제목으로 업데이트 시 에러가 발생해야 함")
	}
	
	expectedError := "제목은 필수입니다"
	if err.Error() != expectedError {
		t.Errorf("예상 에러: %s, 실제 에러: %s", expectedError, err.Error())
	}
}