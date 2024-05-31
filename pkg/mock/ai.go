// Code generated by MockGen. DO NOT EDIT.
// Source: ./ai.go
//
// Generated by this command:
//
//	mockgen -source ./ai.go -destination ./mock/ai.go
//

// Package mock_discogpt is a generated GoMock package.
package mock_discogpt

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockMessageGenerator is a mock of MessageGenerator interface.
type MockMessageGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockMessageGeneratorMockRecorder
}

// MockMessageGeneratorMockRecorder is the mock recorder for MockMessageGenerator.
type MockMessageGeneratorMockRecorder struct {
	mock *MockMessageGenerator
}

// NewMockMessageGenerator creates a new mock instance.
func NewMockMessageGenerator(ctrl *gomock.Controller) *MockMessageGenerator {
	mock := &MockMessageGenerator{ctrl: ctrl}
	mock.recorder = &MockMessageGeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMessageGenerator) EXPECT() *MockMessageGeneratorMockRecorder {
	return m.recorder
}

// Generate mocks base method.
func (m *MockMessageGenerator) Generate(ctx context.Context, prompt string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Generate", ctx, prompt)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Generate indicates an expected call of Generate.
func (mr *MockMessageGeneratorMockRecorder) Generate(ctx, prompt any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Generate", reflect.TypeOf((*MockMessageGenerator)(nil).Generate), ctx, prompt)
}
