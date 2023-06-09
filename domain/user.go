package domain

import (
	"fmt"
	"regexp"

	"github.com/pkg/errors"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	id             UserID
	name           UserName
	hashedPassword HashedPassword
	role           UserRole
}

func (u *User) ID() UserID                     { return u.id }
func (u *User) Name() UserName                 { return u.name }
func (u *User) HashedPassword() HashedPassword { return u.hashedPassword }
func (u *User) Role() UserRole                 { return u.role }

func NewUser(
	name UserName,
	hashedPassword HashedPassword,
	role UserRole,
) (*User, error) {
	return &User{
		name:           name,
		hashedPassword: hashedPassword,
		role:           role,
	}, nil
}

func NewUserFromSource(
	id uint64,
	name string,
	hashedPassword string,
	role string,
) (*User, error) {
	newID, err := NewUserID(id)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	newName, err := NewUserName(name)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	newHashedPassword, err := NewHashedPassword(hashedPassword)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	newRole, err := NewUserRole(role)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &User{
		id:             newID,
		name:           newName,
		hashedPassword: newHashedPassword,
		role:           newRole,
	}, nil
}

type UserID uint64

var ErrUserIDZero = errors.New("user id: must not be zero")

func NewUserID(v uint64) (UserID, error) {
	if v == 0 {
		return 0, ErrUserIDZero
	}

	return UserID(v), nil
}

type UserName string

const UserNameMaxLength = 32

var (
	ErrUserNameEmpty   = errors.New("user name: must not be empty")
	ErrUserNameTooLong = fmt.Errorf("user name: must not be longer than %d characters", UserNameMaxLength)
)

func NewUserName(v string) (UserName, error) {
	if v == "" {
		return "", errors.WithStack(ErrUserNameEmpty)
	}

	if len([]rune(v)) > UserNameMaxLength {
		return "", errors.WithStack(ErrUserNameTooLong)
	}

	return UserName(v), nil
}

type HashedPassword []byte

var (
	ErrHashedPasswordEmpty    = errors.New("hashed password: must not be empty")
	ErrHashedPasswordNotMatch = errors.New("hashed password: not match")
)

func NewHashedPassword(v string) (HashedPassword, error) {
	if v == "" {
		return HashedPassword(""), errors.WithStack(ErrHashedPasswordEmpty)
	}

	return HashedPassword(v), nil
}

func (v HashedPassword) Verify(pass Password) error {
	if err := bcrypt.CompareHashAndPassword(v, []byte(pass)); err != nil {
		return errors.Wrap(ErrHashedPasswordNotMatch, err.Error())
	}

	return nil
}

type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

var ErrUserRoleInvalid = errors.New("user role: invalid type")

func NewUserRole(v string) (UserRole, error) {
	switch v {
	case string(RoleUser):
		return RoleUser, nil
	case string(RoleAdmin):
		return RoleAdmin, nil
	}

	return UserRole(""), errors.WithStack(ErrUserRoleInvalid)
}

type Password string

const (
	PasswordMinLength = 8
	PasswordMaxLength = 16
	PasswordHashCost  = 10
)

var (
	ErrPasswordEmpty    = errors.New("password: must not be empty")
	ErrPasswordTooShort = fmt.Errorf(
		"password: must be at least %d characters",
		PasswordMinLength,
	)
	ErrPasswordTooLong = fmt.Errorf(
		"password: must be at shorter than %d characters",
		PasswordMaxLength,
	)
	ErrPasswordDoesNotFollowRule = errors.New("password: does not follow the rules")
	PasswordCharcters            = regexp.MustCompile("^[0-9a-zA-Z!-/:-@[-`{-~]+$")
	PasswordMustIncludes         = []*regexp.Regexp{
		regexp.MustCompile("[[:alpha:]]"),
		regexp.MustCompile("[[:digit:]]"),
		regexp.MustCompile("[[:punct:]]"),
	}
)

func NewPassword(v string) (Password, error) {
	if v == "" {
		return "", errors.WithStack(ErrPasswordEmpty)
	}

	if len([]rune(v)) < PasswordMinLength {
		return "", errors.WithStack(ErrPasswordTooShort)
	}

	if PasswordMaxLength < len([]rune(v)) {
		return "", errors.WithStack(ErrPasswordTooLong)
	}

	if !PasswordCharcters.MatchString(v) {
		return "", errors.WithStack(ErrPasswordDoesNotFollowRule)
	}
	for _, expected := range PasswordMustIncludes {
		if expected.FindString(v) == "" {
			return "", errors.WithStack(ErrPasswordDoesNotFollowRule)
		}
	}

	return Password(v), nil
}

func (v Password) Hash() HashedPassword {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(v), PasswordHashCost)

	return hashed
}

type UserMetaData struct {
	userAgent UserAgent
	clientIp  ClientIp
}

func (m *UserMetaData) UserAgent() UserAgent { return m.userAgent }
func (m *UserMetaData) ClientIp() ClientIp   { return m.clientIp }

func NewUserMetadata(userAgent UserAgent, clientIp ClientIp) (*UserMetaData, error) {
	return &UserMetaData{userAgent: userAgent, clientIp: clientIp}, nil
}
