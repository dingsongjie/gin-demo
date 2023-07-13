package controllers

import (
	"lenovo-drive-mi-api/db"
	"lenovo-drive-mi-api/entities"
	"lenovo-drive-mi-api/log"
	"lenovo-drive-mi-api/models"
	model "lenovo-drive-mi-api/models"
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
// @param request body models.RequestModel true "requestUser"
// @Success 200  {array} models.FileUserInformation
// @Router /getUserInfomation [post]
func GetAllInformation(c *gin.Context) {
	var requestUser = model.RequestModel{}
	if err := c.BindJSON(&requestUser); err != nil {
		return
	}
	info, err := getNewPathInfo(requestUser.Paths)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatus(400)
	}
	err = setUserInfo(requestUser.Paths, info)
	if err != nil {
		log.Logger.Error(err.Error())
		c.AbortWithStatus(400)
	} else {
		c.JSON(http.StatusOK, info)
	}
}

func setUserInfo(paths []string, enumerable linq.Linq[*model.FileUserInformation]) error {
	db := db.FileUserInfoDb
	rows, err := db.Table("iris_name_entry as entry").
		Select("entry.path, entry.ctime, cuser.slug,entry.mtime,muser.slug, entry.path").
		Joins("left join iris_user as cuser on  entry.creator_uid = cuser.id").
		Joins("left join iris_user as muser on  entry.updator_uid = muser.id").
		Where("entry.nsid=1 and entry.path in ?", paths).Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		newInfo := models.FileUserInformation{}
		err := rows.Scan(&newInfo.OriginPath, &newInfo.CreatedTime, &newInfo.CreatedUserCode, &newInfo.LastModifiedTime, &newInfo.LastModifiedUserCode, &newInfo.OriginPath)
		if err != nil {
			return err
		}
		currentInfo := enumerable.FirstOrDefault(func(s *models.FileUserInformation) bool {
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
	setNowIfTimeIsZero(enumerable)
	return nil
}

func setNowIfTimeIsZero(enumerable linq.Linq[*models.FileUserInformation]) {
	enumerable.ForEach(func(fui *models.FileUserInformation) {
		if time.Time.IsZero(fui.CreatedTime) {
			fui.CreatedTime = time.Now()
		}
		if time.Time.IsZero(fui.LastModifiedTime) {
			fui.LastModifiedTime = time.Now()
		}
	})
}

func getNewPathInfo(paths []string) (result linq.Linq[*models.FileUserInformation], err error) {
	db := db.FileNewPathDb
	rows, err := db.Model(&entities.PathMap{}).Select("OriginPath,NewPath").Where("OriginPath in ?", paths).Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		info := models.FileUserInformation{}
		err := rows.Scan(&info.OriginPath, &info.NewPath)
		if err != nil {
			return nil, err
		}
		result.Add(&info)
	}
	return
}
