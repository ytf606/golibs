package orm

import "gorm.io/gorm"

func IgnoreEmpty(err error) error {
	if err == gorm.ErrRecordNotFound {
		return nil
	}
	return err
}

func EmptyRecord(err error) bool {
	return err == gorm.ErrRecordNotFound
}
