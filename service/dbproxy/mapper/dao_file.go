package mapper

import (
	"github.com/sirupsen/logrus"
	"xcloud/db/mysql"
)

// FileUploadFinished: 文件上传完成，保存信息到唯一文件表
func FileUploadFinished(fileHash, fileName, fileAddr string, fileSize int64) (res SqlResult){
	stmt, err := mysql.Conn().Prepare(`insert ignore into 
    file (file_hash, file_name, file_size, file_addr, status) values (?,?,?,?,1)`)
	if err != nil {
		logrus.Error(err)
		res.Succ = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	ret, err := stmt.Exec(fileHash, fileName, fileSize, fileAddr)
	if err != nil {
		logrus.Error(err)
		res.Succ = false
		res.Msg = err.Error()
		return
	}
	if rf, err := ret.RowsAffected(); err == nil {
		if rf > 0 {
			res.Succ = true
			res.Msg = "插入唯一文件表信息成功"
			return
		} else if rf == 0 {
			res.Succ = true
			res.Msg = "文件哈希值相同"
			return
		}
	}
	res.Succ = false
	res.Msg = "更新唯一文件表信息失败"
	return
}

// UpdateFileAddr: 更新文件存储地址
func UpdateFileAddr(fileHash string, fileAddr string) (res SqlResult) {
	stmt, err := mysql.Conn().Prepare(
		"update file set file_addr=? where file_hash=?")
	if err != nil {
		logrus.Error(err)
		res.Succ = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	ret, err := stmt.Exec(fileAddr, fileHash)
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
		res.Msg = "更新文件存储地址成功"
		return
	}
	logrus.Error(err)
	res.Succ = false
	res.Msg = "更新文件存储地址失败"
	return
}

// DeleteFile: 删除文件表记录
func DeleteFile(fileHash string) (res SqlResult) {
	stmt, err := mysql.Conn().Prepare(
		"delete from file where file_hash=?")
	if err != nil {
		logrus.Error(err)
		res.Succ = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(fileHash)
	if err != nil {
		logrus.Error(err)
		res.Succ = false
		res.Msg = err.Error()
		return
	}
	res.Succ = true
	res.Msg = "删除文件记录成功"
	return
}

func GetFileAddr(fileHash string) (res SqlResult) {
	stmt, err := mysql.Conn().Prepare(
		"select file_addr from file where file_hash=?")
	if err != nil {
		logrus.Error(err)
		res.Succ = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	var fileAddr string
	err = stmt.QueryRow(fileHash).Scan(&fileAddr)
	if err != nil {
		logrus.Error(err)
		res.Succ = false
		res.Msg = "该文件不存在"
		return
	}
	res.Succ = true
	res.Msg = "找到该文件存储路径"
	res.Data = fileAddr
	return
}

func GetFileMeta(fileHash string) (res SqlResult) {
	stmt, err := mysql.Conn().Prepare(
		"select file_size, file_addr from file where file_hash=?")
	if err != nil {
		logrus.Error(err)
		res.Succ = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	var fileInfo FileInfo
	err = stmt.QueryRow(fileHash).Scan(&fileInfo.FileSize, &fileInfo.FileAddr)
	if err != nil {
		res.Succ = false
		res.Msg = err.Error()
		return
	}
	res.Succ = true
	res.Msg = "获取文件信息成功"
	res.Data = fileInfo
	return
}
