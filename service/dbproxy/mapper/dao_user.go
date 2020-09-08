package mapper

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"xcloud/cache/redis"
	"xcloud/db/mysql"
)

// UserSignUp: 注册
func UserSignUp(username string, password string, passwordSalt string) (res SqlResult) {
	stmt, err := mysql.Conn().Prepare(
		"insert into user (username, password, password_salt) value (?,?,?)")
	if err != nil {
		logrus.Error(err)
		res.Succ = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	ret, err := stmt.Exec(username, password, passwordSalt)
	if err != nil {
		logrus.Error(err)
		res.Succ = false
		res.Msg = err.Error()
		return
	}
	if rf, err := ret.RowsAffected(); err == nil && rf > 0 {
		res.Succ = true
		return
	}
	res.Succ = false
	res.Msg = "注册失败"
	return
}

// UserSignIn: 登陆
func UserSignIn(username string, password string) (res SqlResult) {
	stmt, err := mysql.Conn().Prepare(
		"select password from user where username=?")
	if err != nil {
		logrus.Error(err)
		res.Succ = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(username)
	if err != nil {
		logrus.Error(err)
		res.Succ = false
		res.Msg = err.Error()
		return
	} else if rows == nil {
		res.Succ = false
		res.Msg = "该用户未注册"
		return
	}

	parseRows := mysql.ParseRows(rows)
	if len(parseRows) > 0 && string(parseRows[0]["password"].([]byte)) == password {
		res.Succ = true
		res.Msg = "登陆成功"
		return
	}
	res.Succ = false
	res.Msg = "用户名/密码不一致"
	return
}

// UserExist: 判断用户是否存在
func UserExist(username string) (res SqlResult) {
	stmt, err := mysql.Conn().Prepare(
		"select username from user where username=?")
	if err != nil {
		logrus.Error(err)
		res.Succ = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	var name string
	err = stmt.QueryRow(username).Scan(&name)
	if err != nil || err == sql.ErrNoRows {
		res.Succ = false
		res.Msg = "该用户不存在"
		return
	}
	res.Succ = true
	res.Msg = "该用户已存在"
	return
}

// GetPasswordSalt: 获取用户密码盐值
func GetPasswordSalt(username string) (res SqlResult) {
	stmt, err := mysql.Conn().Prepare(
		"select password_salt from user where username=?")
	if err != nil {
		logrus.Error(err)
		res.Succ = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	var passwordSalt string
	err = stmt.QueryRow(username).Scan(&passwordSalt)
	if err != nil {
		logrus.Error(err)
		res.Succ = false
		res.Msg = "该用户不存在"
		return
	}
	res.Succ = true
	res.Msg = "获取密码盐值成功"
	res.Data = passwordSalt
	return
}

// UpdateUserToken: 保存或更新用户token
func UpdateUserToken(username string, token string) (res SqlResult) {
	stmt, err := mysql.Conn().Prepare(
		"replace into user_token (username, token) values (?,?)")
	if err != nil {
		logrus.Error(err)
		res.Succ = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, token)
	if err != nil {
		logrus.Error(err)
		res.Succ = false
		res.Msg = err.Error()
		return
	}
	res.Succ = true
	return
}

// SaveTokenToRedis: 用户登陆成功后，将token保存到redis，并设置过期时间
func SaveTokenToRedis(op, username, expireTime, token string) (res SqlResult) {
	conn := redis.Pool().Get()
	defer conn.Close()
	_, err := conn.Do(op, username, expireTime, token)
	if err != nil {
		res.Succ = false
		res.Msg = "设置token到redis失败"
		return
	}
	res.Succ = true
	res.Msg = "设置token到redis成功"
	return
}

// GetUserInfo: 获取用户信息
func GetUserInfo(username string) (res SqlResult) {
	stmt, err := mysql.Conn().Prepare(
		"select username, email, phone, DATE_FORMAT(signup_at, '%Y-%m-%d'), last_active_at, status from user where username=?")
	if err != nil {
		logrus.Error(err)
		res.Succ = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	user := UserInfo{}
	err = stmt.QueryRow(username).Scan(&user.UserName, &user.Email, &user.Phone,
		&user.SignupAt, &user.LastActiveAt, &user.Status)
	if err != nil {
		logrus.Error(err)
		res.Succ = false
		res.Msg = err.Error()
		return
	}
	res.Succ = true
	res.Msg = "获取用户信息成功"
	res.Data = user
	return
}
