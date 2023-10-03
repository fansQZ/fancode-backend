// Code generated by MockGen. DO NOT EDIT.
// Source: dao/problem_attempt_dao.go

// Package dao is a generated GoMock package.
package dao

import (
	po "FanCode/models/po"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	gorm "gorm.io/gorm"
)

// MockProblemAttemptDao is a mock of ProblemAttemptDao interface.
type MockProblemAttemptDao struct {
	ctrl     *gomock.Controller
	recorder *MockProblemAttemptDaoMockRecorder
}

// MockProblemAttemptDaoMockRecorder is the mock recorder for MockProblemAttemptDao.
type MockProblemAttemptDaoMockRecorder struct {
	mock *MockProblemAttemptDao
}

// NewMockProblemAttemptDao creates a new mock instance.
func NewMockProblemAttemptDao(ctrl *gomock.Controller) *MockProblemAttemptDao {
	mock := &MockProblemAttemptDao{ctrl: ctrl}
	mock.recorder = &MockProblemAttemptDaoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockProblemAttemptDao) EXPECT() *MockProblemAttemptDaoMockRecorder {
	return m.recorder
}

// GetProblemAttemptByID mocks base method.
func (m *MockProblemAttemptDao) GetProblemAttemptByID(db *gorm.DB, userId, problemId uint) (*po.ProblemAttempt, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProblemAttemptByID", db, userId, problemId)
	ret0, _ := ret[0].(*po.ProblemAttempt)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProblemAttemptByID indicates an expected call of GetProblemAttemptByID.
func (mr *MockProblemAttemptDaoMockRecorder) GetProblemAttemptByID(db, userId, problemId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProblemAttemptByID", reflect.TypeOf((*MockProblemAttemptDao)(nil).GetProblemAttemptByID), db, userId, problemId)
}

// GetProblemAttemptState mocks base method.
func (m *MockProblemAttemptDao) GetProblemAttemptState(db *gorm.DB, userId, problemID uint) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProblemAttemptState", db, userId, problemID)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProblemAttemptState indicates an expected call of GetProblemAttemptState.
func (mr *MockProblemAttemptDaoMockRecorder) GetProblemAttemptState(db, userId, problemID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProblemAttemptState", reflect.TypeOf((*MockProblemAttemptDao)(nil).GetProblemAttemptState), db, userId, problemID)
}

// InsertProblemAttempt mocks base method.
func (m *MockProblemAttemptDao) InsertProblemAttempt(db *gorm.DB, problemAttempt *po.ProblemAttempt) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertProblemAttempt", db, problemAttempt)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertProblemAttempt indicates an expected call of InsertProblemAttempt.
func (mr *MockProblemAttemptDaoMockRecorder) InsertProblemAttempt(db, problemAttempt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertProblemAttempt", reflect.TypeOf((*MockProblemAttemptDao)(nil).InsertProblemAttempt), db, problemAttempt)
}

// UpdateProblemAttempt mocks base method.
func (m *MockProblemAttemptDao) UpdateProblemAttempt(db *gorm.DB, problemAttempt *po.ProblemAttempt) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateProblemAttempt", db, problemAttempt)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateProblemAttempt indicates an expected call of UpdateProblemAttempt.
func (mr *MockProblemAttemptDaoMockRecorder) UpdateProblemAttempt(db, problemAttempt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateProblemAttempt", reflect.TypeOf((*MockProblemAttemptDao)(nil).UpdateProblemAttempt), db, problemAttempt)
}
