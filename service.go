package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/STRockefeller/go-linq"
	"github.com/gin-gonic/gin"
)

// @BasePath /

// getAllInformation
// @Summary getAllInformation
// @Schemes
// @Description getAllInformation
// @Tags getAllInformation
// @Accept json
// @Produce json
// @param request body requestModel true "requestUser"
// @Success 200  {array} FileUserInformation
// @Router /getUserInfomation [post]
func getAllInformation(c *gin.Context) {
	var requestUser = requestModel{}
	if err := c.BindJSON(&requestUser); err != nil {
		return
	}
	info, err := GetNewPathInfo(requestUser.Paths)
	if err != nil {
		logger.Error(err.Error())
		c.AbortWithStatus(400)
	}
	err = SetUserInfo(requestUser.Paths, info)
	if err != nil {
		logger.Error(err.Error())
		c.AbortWithStatus(400)
	} else {
		c.JSON(http.StatusOK, info)
	}
}

func SetUserInfo(paths []string, enumerable []FileUserInformation) error {
	results := linq.Linq[*FileUserInformation]{}
	for index, _ := range enumerable {
		results.Add(&enumerable[index])
	}
	userInfoDbClient := FileUserInfoDbClient{UserInfoConnectionString}
	db, err := userInfoDbClient.getDb()
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
	sqlstr := fmt.Sprintf(`select entry.path, entry.ctime, cuser.slug,entry.mtime,muser.slug,
	 entry.path from iris_name_entry as entry left join iris_user as cuser on  entry.creator_uid = cuser.id  
	 join iris_user as muser on  entry.updator_uid = muser.id where entry.nsid=1 and entry.path in (%s)`, inCodition)
	rows, err := db.Query(sqlstr, params...)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		newInfo := FileUserInformation{}
		err := rows.Scan(&newInfo.OriginPath, &newInfo.CreatedTime, &newInfo.CreatedUserCode, &newInfo.LastModifiedTime, &newInfo.LastModifiedUserCode, &newInfo.OriginPath)
		if err != nil {
			return err
		}
		currentInfo := results.FirstOrDefault(func(s *FileUserInformation) bool {
			return s.OriginPath == newInfo.OriginPath
		})
		if err != nil {
			return err
		}
		if currentInfo == nil {
			continue
		}
		if newInfo.CreatedUserCode != "" {
			currentInfo.CreatedTime = newInfo.CreatedTime
			currentInfo.LastModifiedTime = newInfo.LastModifiedTime
			currentInfo.CreatedUserCode = newInfo.CreatedUserCode
			currentInfo.LastModifiedUserCode = newInfo.LastModifiedUserCode
		}
	}

	results.ForEach(func(fui *FileUserInformation) {
		if time.Time.IsZero(fui.CreatedTime) {
			fui.CreatedTime = time.Now()
		}
		if time.Time.IsZero(fui.LastModifiedTime) {
			fui.LastModifiedTime = time.Now()
		}
	})
	return nil
}

func GetNewPathInfo(paths []string) ([]FileUserInformation, error) {
	var results = []FileUserInformation{}
	newPathDbClient := FileUserInfoDbClient{NewPathConnectionString}

	db, err := newPathDbClient.getDb()
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
	sqlstr := fmt.Sprintf(`select OriginPath, NewPath FROM Path_Maps where OriginPath in (%s)`, inCodition)
	rows, err := db.Query(
		sqlstr, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		info := FileUserInformation{}
		err := rows.Scan(&info.OriginPath, &info.NewPath)
		if err != nil {
			return nil, err
		}
		results = append(results, info)
	}
	return results, nil
}

func filter(ss []FileUserInformation, pred func(*FileUserInformation) bool) (ret *FileUserInformation) {
	for i, s := range ss {
		if pred(&s) {
			return &ss[i]
		}
	}
	return
}
