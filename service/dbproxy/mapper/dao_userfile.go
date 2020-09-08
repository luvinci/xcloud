package mapper

import (
	"github.com/sirupsen/logrus"
	"xcloud/db/mysql"
)

// UserFileUploadFinished: 文件上传完成，保存信息到用户文件表
func UserFileUploadFinished(username, fileHash, fileName string, fileSize int64) (res SqlResult){
	stmt, err := mysql.Conn().Prepare(`insert ignore into 
    user_file (username, file_hash, file_name, file_size, status) values (?,?,?,?,1)` )
	if err != nil {
		logrus.Error(err)
		res.Succ = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, fileHash, fileName, fileSize)
	if err != nil {
		logrus.Error(err)
		res.Succ = false
		res.Msg = "更新用户文件表信息失败"
		return
	}
	res.Succ = true
	res.Msg = "更新用户文件表信息成功"
	return
}

// GetUserFiles: 批量获取用户文件
func GetUserFiles(username string, limit int64) (res SqlResult) {
	stmt, err := mysql.Conn().Prepare(
		`select file_hash, file_name, file_size, DATE_FORMAT(upload_at, '%Y-%m-%d'), 
       DATE_FORMAT(last_update, '%Y-%m-%d') from user_file where username=? limit ?`)
	if err != nil {
		logrus.Error(err)
		res.Succ = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(username, limit)
	if err != nil {
		logrus.Error(err)
		res.Succ = false
		res.Msg = err.Error()
		return
	}

	var userFiles []UserFile
	for rows.Next() {
		userFile := UserFile{}
		err = rows.Scan(&userFile.FileHash, &userFile.FileName,
			&userFile.FileSize, &userFile.UploadAt, &userFile.LastUpdate)
		if err != nil {
			logrus.Error(err)
			res.Succ = false
			res.Msg = "批量获取用户文件失败"
			return
		}
		userFiles = append(userFiles, userFile)
	}
	res.Succ = true
	res.Msg = "批量获取用户文件成功"
	res.Data = userFiles
	return
}

// RenameUserFile: 用户文件重命名
func RenameUserFile(username string, fileHash string, newFilename string) (res SqlResult) {
	stmt, err := mysql.Conn().Prepare(
		"update user_file set file_name=? where username=? and file_hash=?")
	if err != nil {
		logrus.Error(err)
		res.Succ = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	ret, err := stmt.Exec(newFilename, username, fileHash)
	if err != nil {
		logrus.Error(err)
		res.Succ = false
		res.Msg = err.Error()
		return
	}
	rf, err := ret.RowsAffected()
	if err == nil {
		if rf <= 0 {
			logrus.Error(err)
			res.Succ = false
			res.Msg = err.Error()
			return
		}
		res.Succ = true
		res.Msg = "文件重命名成功"
		return
	}
	logrus.Error(err)
	res.Succ = false
	res.Msg = "文件重命名失败"
	return
}

// DeleteUserFile: 删除用户文件记录
func DeleteUserFile(username string, fileHash string) (res SqlResult) {
	stmt, err := mysql.Conn().Prepare(
		"delete from user_file where username=? and file_hash=?")
	if err != nil {
		logrus.Error(err)
		res.Succ = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, fileHash)
	if err != nil {
		logrus.Error(err)
		res.Succ = false
		res.Msg = err.Error()
		return
	}
	res.Succ = true
	res.Msg = "删除用户文件记录成功"
	return
}

// DeleteUserFileAndUniqueFile: 开启事务删除文件
func DeleteUserFileAndUniqueFile(username string, fileHash string) (res SqlResult) {
	tx, err := mysql.Conn().Begin()
	if err != nil {
		logrus.Error(err)
		res.Succ = false
		res.Msg = "开启事务失败"
		return
	}
	stmt, err := tx.Prepare("delete from user_file where username=? and file_hash=?")
	if err != nil {
		logrus.Error(err)
		res.Succ = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(username, fileHash)
	if err != nil {
		logrus.Error(err)
		res.Succ = false
		res.Msg = err.Error()
		return
	}

	stmt, err = tx.Prepare("delete from file where file_hash=?")
	if err != nil {
		logrus.Error(err)
		res.Succ = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()
	_, err = stmt.Exec(fileHash)
	if err != nil {
		_ = tx.Rollback()
		logrus.Error(err)
		res.Succ = false
		res.Msg = err.Error()
		return
	}
	err = tx.Commit()
	if err != nil {
		logrus.Error(err)
		res.Succ = false
		res.Msg = err.Error()
		return
	}
	res.Succ = true
	res.Msg = "删除成功"
	return
}

// GetSameFileHashCount: 计算用户文件表相同的文件hash有多少个
func GetSameFileHashCount(fileHash string) (res SqlResult) {
	stmt, err := mysql.Conn().Prepare(
		"select id from user_file where file_hash=?")
	if err != nil {
		logrus.Error(err)
		res.Succ = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(fileHash)
	if err != nil {
		logrus.Error(err)
		res.Succ = false
		res.Msg = err.Error()
		return
	}
	var count int64
	for rows.Next() {
		count++
	}
	res.Succ = true
	res.Msg = "相同的文件hash数量：" + string(count)
	res.Data = count
	return
}
