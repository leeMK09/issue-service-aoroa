package model

import (
	"testing"
	userModel "issue-service-aoroa/user/model"
)

func TestUpdateCommand_ApplyTo_실패_완료된_이슈(t *testing.T) {
	issue, _ := NewIssue("테스트", "설명", nil)
	issue.Status = StatusCompleted
	
	cmd := NewUpdateCommand().WithTitle("새 제목")
	
	err := cmd.ApplyTo(issue)
	
	if err == nil {
		t.Error("완료된 이슈에 업데이트 명령 적용 시 에러가 발생해야 함")
	}
	
	expectedError := "완료되거나 취소된 이슈는 수정할 수 없습니다"
	if err.Error() != expectedError {
		t.Errorf("예상 에러: %s, 실제 에러: %s", expectedError, err.Error())
	}
}

func TestUpdateCommand_ApplyTo_실패_취소된_이슈(t *testing.T) {
	issue, _ := NewIssue("테스트", "설명", nil)
	issue.Status = StatusCancelled
	
	cmd := NewUpdateCommand().WithDescription("새 설명")
	
	err := cmd.ApplyTo(issue)
	
	if err == nil {
		t.Error("취소된 이슈에 업데이트 명령 적용 시 에러가 발생해야 함")
	}
}

func TestUpdateCommand_ApplyTo_성공_체이닝_업데이트(t *testing.T) {
	issue, _ := NewIssue("원래 제목", "원래 설명", nil)
	user := &userModel.User{ID: 1, Name: "테스트 사용자"}
	
	cmd := NewUpdateCommand().
		WithTitle("새 제목").
		WithDescription("새 설명").
		WithUser(1, user)
	
	err := cmd.ApplyTo(issue)
	
	if err != nil {
		t.Errorf("에러가 발생하지 않아야 함: %v", err)
	}
	
	if issue.Title != "새 제목" {
		t.Errorf("제목이 업데이트되어야 함. 예상: 새 제목, 실제: %s", issue.Title)
	}
	
	if issue.Description != "새 설명" {
		t.Errorf("설명이 업데이트되어야 함. 예상: 새 설명, 실제: %s", issue.Description)
	}
	
	if issue.User != user {
		t.Error("사용자가 할당되어야 함")
	}
	
	if issue.Status != StatusInProgress {
		t.Errorf("PENDING에서 사용자 할당 시 IN_PROGRESS로 전환되어야 함. 실제: %s", issue.Status)
	}
}

func TestUpdateCommand_ApplyTo_성공_담당자_제거(t *testing.T) {
	user := &userModel.User{ID: 1, Name: "테스트 사용자"}
	issue, _ := NewIssue("테스트", "설명", user)
	
	cmd := NewUpdateCommand().WithoutUser()
	
	err := cmd.ApplyTo(issue)
	
	if err != nil {
		t.Errorf("에러가 발생하지 않아야 함: %v", err)
	}
	
	if issue.User != nil {
		t.Error("담당자가 제거되어야 함")
	}
	
	if issue.Status != StatusPending {
		t.Errorf("담당자 제거 시 PENDING 상태로 전환되어야 함. 실제: %s", issue.Status)
	}
}

func TestUpdateCommand_ApplyTo_실패_유효하지_않은_상태(t *testing.T) {
	issue, _ := NewIssue("테스트", "설명", nil)
	
	cmd := NewUpdateCommand().WithStatus("INVALID_STATUS")
	
	err := cmd.ApplyTo(issue)
	
	if err == nil {
		t.Error("유효하지 않은 상태로 업데이트 시 에러가 발생해야 함")
	}
	
	expectedError := "유효하지 않은 상태입니다"
	if err.Error() != expectedError {
		t.Errorf("예상 에러: %s, 실제 에러: %s", expectedError, err.Error())
	}
}

func TestUpdateCommand_ApplyTo_실패_빈_제목(t *testing.T) {
	issue, _ := NewIssue("테스트", "설명", nil)
	
	cmd := NewUpdateCommand().WithTitle("")
	
	err := cmd.ApplyTo(issue)
	
	if err == nil {
		t.Error("빈 제목으로 업데이트 시 에러가 발생해야 함")
	}
	
	expectedError := "제목은 필수입니다"
	if err.Error() != expectedError {
		t.Errorf("예상 에러: %s, 실제 에러: %s", expectedError, err.Error())
	}
}