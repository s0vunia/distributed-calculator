package agent

import (
	"github.com/stretchr/testify/assert"
	"myproject/internal/repositories/agent/mocks"
	"testing"
)

func TestCreate(t *testing.T) {
	mockRepo := mocks.NewRepository(t)
	mockRepo.
		On("Create", "testId").
		Return(nil)
	err := mockRepo.Create("testId")
	assert.NoError(t, err)

	mockRepo.AssertCalled(t, "Create", "testId")
}
