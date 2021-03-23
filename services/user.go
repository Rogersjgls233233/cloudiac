package services

import (
	"fmt"
	"strings"

	//"errors"
	"cloudiac/consts/e"
	"cloudiac/libs/db"
	"cloudiac/models"
	"cloudiac/utils"
	"cloudiac/consts"
)

func CreateUser(tx *db.Session, user models.User) (*models.User, e.Error) {
	if err := models.Create(tx, &user); err != nil {
		if e.IsDuplicate(err) {
			return nil, e.New(e.UserAlreadyExists, err)
		}
		return nil, e.New(e.DBError, err)
	}

	return &user, nil
}

func UpdateUser(tx *db.Session, id uint, attrs models.Attrs) (user *models.User, re e.Error) {
	user = &models.User{}
	if _, err := models.UpdateAttr(tx.Where("id = ?", id), &models.User{}, attrs); err != nil {
		if e.IsDuplicate(err) {
			return nil, e.New(e.UserEmailDuplicate)
		}
		return nil, e.New(e.DBError, fmt.Errorf("update user error: %v", err))
	}
	if err := tx.Where("id = ?", id).First(user); err != nil {
		return nil, e.New(e.DBError, fmt.Errorf("query user error: %v", err))
	}
	return
}

func DeleteUser(tx *db.Session, id uint) e.Error {
	if _, err := tx.Where("id = ?", id).Delete(&models.User{}); err != nil {
		return e.New(e.DBError, fmt.Errorf("delete user error: %v", err))
	}
	return nil
}

func GetUserById(tx *db.Session, id uint) (*models.User, e.Error) {
	u := models.User{}
	if err := tx.Where("id = ?", id).First(&u); err != nil {
		if e.IsRecordNotFound(err) {
			return nil, e.New(e.UserNotExists, err)
		}
		return nil, e.New(e.DBError, err)
	}
	return &u, nil
}

func GetUserByName(tx *db.Session, name string) (*models.User, error) {
	u := models.User{}
	if err := tx.Where("name = ?", name).First(&u); err != nil {
		return nil, err
	}
	return &u, nil
}

func GetUserByEmail(tx *db.Session, email string) (*models.User, error) {
	u := models.User{}
	if err := tx.Where("email = ?", email).First(&u); err != nil {
		return nil, err
	}
	return &u, nil
}

func FindUsers(query *db.Session) (users []*models.User, err error) {
	err = query.Find(&users)
	return
}

func QueryUser(query *db.Session) *db.Session {
	return query.Model(&models.User{})
}

func HashPassword(password string) (string, e.Error) {
	if er := CheckPasswordFormat(password); er != nil {
		return "", er
	}

	hashed, err := utils.HashPassword(password)
	if err != nil {
		return "", e.New(e.InternalError, err)
	}
	return hashed, nil
}

func CheckPasswordFormat(password string) e.Error {
	//密码校验规则：数字、大写字母、小写字母两种及两种以上组合，且长度在6~30之间
	if len(password) < 6 || len(password) > 30 {
		return e.New(e.InvalidPasswordFormat, fmt.Errorf("error: password length"))
	}

	typeCount := 0
	for _, chars := range []string{consts.LowerCaseLetter, consts.UpperCaseLetter, consts.DigitChars} {
		if strings.ContainsAny(password, chars) {
			typeCount += 1
		}
	}
	if typeCount < 2 {
		return e.New(e.InvalidPasswordFormat, fmt.Errorf("error: password strength"))
	}

	return nil
}

func CreateUserOrgMap(tx *db.Session, userOrgMap models.UserOrgMap) (*models.UserOrgMap, e.Error) {
	if err := models.Create(tx, &userOrgMap); err != nil {
		if e.IsDuplicate(err) {
			return nil, e.New(e.UserAlreadyExists, err)
		}
		return nil, e.New(e.DBError, err)
	}

	return &userOrgMap, nil
}

func DeleteUserOrgMap(tx *db.Session, userId uint, orgId uint) e.Error {
	if _, err := tx.Where("user_id = ? AND org_id = ?", userId, orgId).Debug().Delete(&models.UserOrgMap{}); err != nil {
		return e.New(e.DBError, fmt.Errorf("delete user %d for org %d error: %v", userId, orgId, err))
	}
	return nil
}

func FindUsersOrgMap(query *db.Session, userId uint, orgId uint) (userOrgMap []*models.UserOrgMap, err error) {
	if err := query.Where("user_id = ? AND org_id = ?", userId, orgId).Find(&userOrgMap); err != nil {
		return nil, e.AutoNew(err, e.DBError)
	}
	return
}