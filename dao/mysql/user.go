package mysql

import (
	"bluebell/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"go.uber.org/zap"
)

const secret = "20160828" //加言的字符串？设计到密码学知识

// CheckUserIsExist 把对数据库的某项操作都分别写封装成函数
//等待logic层根据业务需求调用
func CheckUserIsExist(username string) (err error) {
	sqlStr := `select count(user_id) from user where username = ?`
	var count int64
	if err := db.Get(&count, sqlStr, username); err != nil {
		return err //查询失败
	}
	if count > 0 { //根据username在数据库中的个数返回对错
		return ErrorUserExist
	}
	return
}

// InsertUser 向数据库中插入一条新的用户记录
func InsertUser(user *models.User) (err error) {
	//对密码进行加密
	user.Password = encryptPassword(user.Password)
	//执行插入语句
	sqlStr := `insert into user(user_id,username,password) values(?,?,?)`
	_, err = db.Exec(sqlStr, user.UserID, user.Username, user.Password)
	return
}

// encryptPassword 对密码进行md5加密
func encryptPassword(pw string) string {
	h := md5.New()
	h.Write([]byte(secret))                      // secret 加言的字符串？设计到密码学知识
	return hex.EncodeToString(h.Sum([]byte(pw))) //hex.是把数据转换为16进制的字符串。h.sum就是给pw进行md5加密
}

func Login(user *models.User) (err error) {
	sqlStr := "select user_id,username,password from user where username = ?"
	opassword := user.Password
	err = db.Get(user, sqlStr, user.Username)
	//判断用户是否存在，但是在正经项目里，一般是不会返回给用户的。因为这样就提示用户该账号未被注册
	if err == sql.ErrNoRows {
		return ErrorUserNoExist
	}
	//查询数据库失败
	if err != nil {
		return err
	}
	//判断密码是否正确
	if user.Password != encryptPassword(opassword) {
		return ErrorInvalidPassword
	}
	return
}

// GetUserbyID 根据userid查询用户
func GetUserbyID(uid int64) (user *models.User, err error) {
	user = new(models.User)
	sqlStr := `select user_id,username from user where user_id = ?`
	if err := db.Get(user, sqlStr, uid); err != nil {
		if err == sql.ErrNoRows {
			zap.L().Warn("there is no user in db")
			err = nil
		}
	}
	return
}
