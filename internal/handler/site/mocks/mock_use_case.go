package mocks

import (
	"context"

	domain "gitlab.com/Dokuchaevvn/site-monitor/internal/domain/site"
	siteusecase "gitlab.com/Dokuchaevvn/site-monitor/internal/usecase/site"
)

type Controller struct {
	callsCounter map[string]int
	returnArgs   map[string][]any
}

func NewController() *Controller {
	return &Controller{
		callsCounter: make(map[string]int),
		returnArgs:   make(map[string][]any),
	}
}

func (c *Controller) Call(funcName string) []any {
	res, ok := c.callsCounter[funcName]
	if ok {
		c.callsCounter[funcName] = res + 1
	} else {
		c.callsCounter[funcName] = 1
	}

	args, _ := c.returnArgs[funcName]
	return args
}

func (c *Controller) GetCallCount(funcName string) int {
	res, ok := c.callsCounter[funcName]
	if !ok {
		return 0
	}

	return res
}

func (c *Controller) Return(funcName string, args ...any) {
	c.returnArgs[funcName] = args
}

type MockGetAllUseCase struct {
	ctrl *Controller
}

func NewMockGetAllUseCase(ctrl *Controller) MockGetAllUseCase {
	return MockGetAllUseCase{
		ctrl: ctrl,
	}
}

func (m *MockGetAllUseCase) CallCount(funcName string) int {
	return m.ctrl.GetCallCount(funcName)
}

func (m *MockGetAllUseCase) Execute(ctx context.Context) ([]domain.Site, error) {
	res := m.ctrl.Call("Execute")

	ret0, _ := res[0].([]domain.Site)
	ret1, _ := res[1].(error)

	return ret0, ret1
}

type MockAddUseCase struct {
	ctrl *Controller
}

func NewMockAddUseCase(ctrl *Controller) MockAddUseCase {
	return MockAddUseCase{
		ctrl: ctrl,
	}
}

func (m *MockAddUseCase) CallCount(funcName string) int {
	return m.ctrl.GetCallCount(funcName)
}

func (m *MockAddUseCase) Execute(ctx context.Context, command siteusecase.AddCommand) (domain.Site, error) {
	res := m.ctrl.Call("Execute")

	ret0, _ := res[0].(domain.Site)
	ret1, _ := res[1].(error)

	return ret0, ret1
}

type MockDeleteByIDUseCase struct {
	ctrl *Controller
}

func NewMockDeleteByIDUseCase(ctrl *Controller) MockDeleteByIDUseCase {
	return MockDeleteByIDUseCase{
		ctrl: ctrl,
	}
}

func (m *MockDeleteByIDUseCase) Execute(ctx context.Context, command siteusecase.DeleteCommand) error {
	res := m.ctrl.Call("Execute")

	ret0, _ := res[0].(error)

	return ret0
}

type MockGetStatusByIDUseCase struct {
	ctrl *Controller
}

func NewMockGetStatusByIDUseCase(ctrl *Controller) MockGetStatusByIDUseCase {
	return MockGetStatusByIDUseCase{
		ctrl: ctrl,
	}
}

func (m *MockGetStatusByIDUseCase) Execute(ctx context.Context, command siteusecase.GetStatusCommand) (domain.Site, error) {
	res := m.ctrl.Call("Execute")

	ret0, _ := res[0].(domain.Site)
	ret1, _ := res[1].(error)

	return ret0, ret1
}

type MockGetHistoryByIDUseCase struct {
	ctrl *Controller
}

func NewMockGetHistoryByIDUseCase(ctrl *Controller) MockGetHistoryByIDUseCase {
	return MockGetHistoryByIDUseCase{
		ctrl: ctrl,
	}
}

func (m *MockGetHistoryByIDUseCase) Execute(ctx context.Context, command siteusecase.GetHistoryCommand) (siteusecase.GetHistoryResult, error) {
	res := m.ctrl.Call("Execute")

	ret0, _ := res[0].(siteusecase.GetHistoryResult)
	ret1, _ := res[1].(error)

	return ret0, ret1
}
