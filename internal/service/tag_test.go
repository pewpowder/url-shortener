package service

import (
	"errors"
	"testing"

	"github.com/pewpowder/url-shortener/internal/dto"
	"github.com/pewpowder/url-shortener/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type TagRepositoryMock struct {
	mock.Mock
}

func (tr *TagRepositoryMock) CreateTag(tag entity.Tag) (entity.Tag, error) {
	args := tr.Called(tag)
	return args.Get(0).(entity.Tag), args.Error(1)
}

func (tr *TagRepositoryMock) UpdateTag(id uint, tag entity.Tag) (entity.Tag, error) {
	args := tr.Called(id, tag)
	return args.Get(0).(entity.Tag), args.Error(1)
}

func (tr *TagRepositoryMock) PatchTag(id uint, tag entity.Tag) (entity.Tag, error) {
	args := tr.Called(id, tag)
	return args.Get(0).(entity.Tag), args.Error(1)
}

func (tr *TagRepositoryMock) DeleteTag(id uint) (entity.Tag, error) {
	args := tr.Called(id)
	return args.Get(0).(entity.Tag), args.Error(1)
}

func (tr *TagRepositoryMock) GetTags(q dto.TagListQuery) ([]entity.Tag, error) {
	args := tr.Called(q)
	return args.Get(0).([]entity.Tag), args.Error(1)
}

func (tr *TagRepositoryMock) GetTagByID(id uint) (entity.Tag, error) {
	args := tr.Called(id)
	return args.Get(0).(entity.Tag), args.Error(1)
}

func TestPatchTag(t *testing.T) {
	tagName := "new tag"
	tagColor := "#333"
	emptyColor := ""
	repoErr := errors.New("repo error")

	tests := []struct {
		name        string
		id          uint
		input       dto.PatchTag
		repoArg     entity.Tag
		repoResult  entity.Tag
		repoErr     error
		want        entity.Tag
		wantError   bool
		errContains string
	}{
		{
			name: "patches name and color",
			id:   1,
			input: dto.PatchTag{
				Name:  &tagName,
				Color: &tagColor,
			},
			repoArg: entity.Tag{
				Name:  tagName,
				Color: tagColor,
			},
			repoResult: entity.Tag{
				ID:    1,
				Name:  tagName,
				Color: tagColor,
			},
			want: entity.Tag{
				ID:    1,
				Name:  tagName,
				Color: tagColor,
			},
		},
		{
			name: "defaults color when empty",
			id:   2,
			input: dto.PatchTag{
				Name:  &tagName,
				Color: &emptyColor,
			},
			repoArg: entity.Tag{
				Name:  tagName,
				Color: DEFAULT_COLOR,
			},
			repoResult: entity.Tag{
				ID:    2,
				Name:  tagName,
				Color: DEFAULT_COLOR,
			},
			want: entity.Tag{
				ID:    2,
				Name:  tagName,
				Color: DEFAULT_COLOR,
			},
		},
		{
			name: "returns repo error",
			id:   3,
			input: dto.PatchTag{
				Name: &tagName,
			},
			repoArg: entity.Tag{
				Name: tagName,
			},
			repoResult:  entity.Tag{},
			repoErr:     repoErr,
			want:        entity.Tag{},
			wantError:   true,
			errContains: "repo error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &TagRepositoryMock{}
			service := tagService{
				container: nil,
				repo:      repo,
			}

			repo.On("PatchTag", tt.id, tt.repoArg).Return(tt.repoResult, tt.repoErr)

			got, err := service.PatchTag(tt.id, tt.input)

			if tt.wantError {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.errContains)
				assert.Equal(t, tt.want, got)
				repo.AssertExpectations(t)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
			repo.AssertExpectations(t)
		})
	}
}
