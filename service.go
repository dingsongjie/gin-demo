package main

import (
	"fmt"
)

func getAllInformation(paths []string) ([]FileUserInformation, error) {
	result, err := GetUserInfo(paths)
	if err != nil {
		return nil, err
	}
	err = setNewPath(paths, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func GetUserInfo(paths []string) ([]FileUserInformation, error) {
	userInfoDbClient := FileUserInfoDbClient{UserInfoConnectionString}
	result := []FileUserInformation{}
	db, err := userInfoDbClient.getDb()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	var inCodition string
	var params []interface{}
	for _, item := range paths {
		params = append(params, item)
		if inCodition != "" {
			inCodition += ", "
		}
		inCodition += "?"
	}
	sqlstr := fmt.Sprintf(`select entry.ctime, cuser.slug,entry.mtime,muser.slug,
	 entry.path from iris_name_entry as entry join iris_user as cuser on  entry.creator_uid = cuser.id  
	 join iris_user as muser on  entry.updator_uid = muser.id where entry.nsid=1 and entry.path in (%s)`, inCodition)
	rows, err := db.Query(sqlstr, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		currentInfo := FileUserInformation{}
		err := rows.Scan(&currentInfo.CreatedTime, &currentInfo.CreatedUserCode, &currentInfo.LastModifiedTime, &currentInfo.LastModifiedUserCode, &currentInfo.OriginPath)
		if err != nil {
			return nil, err
		}
		result = append(result, currentInfo)
	}
	return result, nil
}

func setNewPath(paths []string, enumerable []FileUserInformation) error {
	newPathDbClient := FileUserInfoDbClient{NewPathConnectionString}

	db, err := newPathDbClient.getDb()
	if err != nil {
		return err
	}
	defer db.Close()
	var inCodition string
	var params []interface{}
	for _, item := range paths {
		params = append(params, item)
		if inCodition != "" {
			inCodition += ", "
		}
		inCodition += "?"
	}
	sqlstr := fmt.Sprintf(`select OriginPath, NewPath FROM Path_Maps where OriginPath in (%s)`, inCodition)
	rows, err := db.Query(
		sqlstr, params...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var originPath, newPath string
		err := rows.Scan(&originPath, &newPath)
		if err != nil {
			return err
		}
		currentInfo := filter(enumerable, func(s *FileUserInformation) bool {
			return s.OriginPath == originPath
		})
		if err != nil {
			return err
		}
		if currentInfo.OriginPath != "" {
			currentInfo.NewPath = newPath
		}
	}
	return nil
}

func filter(ss []FileUserInformation, pred func(*FileUserInformation) bool) (ret *FileUserInformation) {
	for i, s := range ss {
		if pred(&s) {
			return &ss[i]
		}
	}
	return
}
