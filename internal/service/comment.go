package service

import (
	"errors"

	"github.com/jbliao/proj-quiz1/internal/model"
	"gorm.io/gorm"
)

type CommentService interface {
	GetComment(uuid string) (*model.Comment, error)
	EnsureComment(*model.Comment) error
	DeleteComment(uuid string) error
}

type PersistedCommentService struct {
	db *gorm.DB
}

func NewPersistedCommentService(db *gorm.DB) CommentService {
	return &PersistedCommentService{db: db}
}

func (s *PersistedCommentService) GetComment(uuid string) (*model.Comment, error) {
	resultComment := &model.Comment{}
	err := s.db.Where("uuid = ?", uuid).Take(resultComment).Error
	return resultComment, err
}

func (s *PersistedCommentService) EnsureComment(comment *model.Comment) error {

	oldcmt := &model.Comment{}
	err := s.db.Select("id").Where("uuid = ?", comment.Uuid).Take(oldcmt).Error

	if err == nil {
		comment.ID = oldcmt.ID
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return s.db.Save(comment).Error
}

func (s *PersistedCommentService) DeleteComment(uuid string) error {
	if result := s.db.Where("uuid = ?", uuid).Delete(&model.Comment{}); result.Error != nil {
		return result.Error
	} else if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
