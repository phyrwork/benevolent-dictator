package database

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID    int    `gorm:"primaryKey;not null"`
	Name  string `gorm:"unique;not null"`
	Email string `gorm:"unique;not null"`
	Salt  []byte `gorm:"not null"`
	Key   []byte `gorm:"not null"`
	Likes []Rule `gorm:"many2many:likes"`
}

func (u User) IDRef() *int {
	return &u.ID
}

func (u User) IDAfter() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(&u).Where("id > ?", u.ID)
	}
}

func (u User) IDBeforeOrEqual() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(&u).Where("id <= ?", u.ID)
	}
}

func (u User) Item() *User {
	return &u
}

func (u User) NameLike() func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(&u).Where("name LIKE ?", u.Name)
	}
}

type Rule struct {
	ID      int `gorm:"primaryKey;not null"`
	UserID  int `gorm:"not null"` // TODO: Rename to UserID
	User    *User
	Created time.Time `gorm:"not null"`
	Summary string    `gorm:"not null"`
	Detail  *string
	Likes   []User `gorm:"many2many:likes"`
}

func (r Rule) IDRef() *int {
	return &r.ID
}

func (r Rule) IDAfter() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(&r).Where("id > ?", r.ID)
	}
}

func (r Rule) IDBeforeOrEqual() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(&r).Where("id <= ?", r.ID)
	}
}

type Like struct {
	UserID int `gorm:"primaryKey;not null"`
	User   *User
	RuleID int `gorm:"primaryKey;not null"`
	Rule   *Rule
}

type UserLike Like

func (l UserLike) TableName() string {
	return "likes"
}

func (l UserLike) IDRef() *int {
	return &l.UserID
}

func (l UserLike) IDAfter() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(&l).Where("user_id > ?", l.UserID)
	}
}

func (l UserLike) IDBeforeOrEqual() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(&l).Where("user_id <= ?", l.UserID)
	}
}

type RuleLike Like

func (l RuleLike) TableName() string {
	return "likes"
}

func (l RuleLike) IDRef() *int {
	return &l.RuleID
}

func (l RuleLike) IDAfter() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(&l).Where("rule_id > ?", l.RuleID)
	}
}

func (l RuleLike) IDBeforeOrEqual() func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Model(&l).Where("rule_id <= ?", l.RuleID)
	}
}
