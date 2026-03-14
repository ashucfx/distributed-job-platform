package database

import "gorm.io/gorm"

// Exposing the underlying *gorm.DB just for the local quick List extension
func (s *GormDBService) GetDB() *gorm.DB {
	return s.db
}
