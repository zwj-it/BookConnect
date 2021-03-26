package logic

//处理业务逻辑的代码
import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snowflake"
)

func SignUp(p *models.ParamSignUp) (err error) {
	//1.验证用户是否已存在
	if err = mysql.CheckUserIsExist(p.Username); err != nil {
		return err
	}
	//2.判断用户名是否合法。视频没有，自己想的，但是没实现
	//3.生成UID
	userID := snowflake.GenID()
	//构造User实例
	user := &models.User{
		UserID:   userID,
		Username: p.Username,
		Password: p.Password,
	}
	//4.更新数据库
	err = mysql.InsertUser(user)
	return
}
func Login(p *models.ParamLogin) (user *models.User, err error) {
	user = &models.User{
		Username: p.Username,
		Password: p.Password,
	}
	//在数据库里判断数据是否正确,下面这步传递的是指针变量，能拿到user_id
	if err = mysql.Login(user); err != nil {
		return nil, err
	}

	//生成并返回JWT的token
	token, err := jwt.GenToken(user.UserID, user.Username)
	if err != nil {
		return
	}
	user.Token = token
	return
}
