package service

import (
	"context"
	"errors"
	"testing"

	"go-gin-realworld-api/internal/services"
	"go-gin-realworld-api/test/mocks"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func setupTagServiceTest(t *testing.T) (context.Context, *services.TagService, *mocks.MockTagRepository, sqlmock.Sqlmock) {
	mockTagRepo := new(mocks.MockTagRepository)
	gormDB, sqlMock := CreateMockDB(t)
	tagService := services.NewTagService(gormDB, mockTagRepo)
	ctxForTest := context.Background()

	return ctxForTest, tagService, mockTagRepo, sqlMock
}

func TestTagService_GetAllTags_Success(t *testing.T) {
	ctxForTest, tagService, mockTagRepo, _ := setupTagServiceTest(t)
	expectedTags := []string{"tag1", "tag2", "tag3"}

	mockTagRepo.On("GetAllTags", mock.Anything).Return(expectedTags, nil)

	tags, err := tagService.GetAllTags(ctxForTest)

	assert.NoError(t, err)
	assert.Equal(t, expectedTags, tags)
	mockTagRepo.AssertExpectations(t)
}

func TestTagService_GetAllTags_Error(t *testing.T) {
	ctxForTest, tagService, mockTagRepo, _ := setupTagServiceTest(t)
	expectedError := errors.New("db error")

	mockTagRepo.On("GetAllTags", mock.Anything).Return(nil, expectedError)

	tags, err := tagService.GetAllTags(ctxForTest)

	assert.Error(t, err)
	assert.Nil(t, tags)
	assert.Equal(t, expectedError, err)
	mockTagRepo.AssertExpectations(t)
}
