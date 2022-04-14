package base

import (
	accessmodel "gin-base/internal/model/access"
	model "gin-base/internal/model/common"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"reflect"
	"strings"
)

type Gin struct {
	C *gin.Context
}

func (g *Gin) EnsureSysUser() *accessmodel.SysUser {
	user, ok := g.C.Get("userInfo")
	if !ok {
		return nil
	}
	sysUser := user.(*accessmodel.SysUser)
	return sysUser
}

func (g *Gin) RespSuccess(data interface{}, msg string) {
	if msg == "" {
		msg = GetMsg(SUCCESS)
	}
	g.C.JSON(http.StatusOK, Response{
		Code: 0,
		Data: data,
		Ok:   true,
		Msg:  msg,
	})
}

func (g *Gin) RespServiceFail(failCode int, msg string) {
	if msg == "" {
		msg = GetMsg(Fail)
	}
	g.C.JSON(http.StatusOK, Response{
		Code: failCode,
		Ok:   false,
		Msg:  msg,
	})
}

func (g *Gin) RespForbidden(msg string) {
	if msg == "" {
		msg = GetMsg(FORBIDDEN)
	}
	g.C.Abort()
	g.C.JSON(http.StatusForbidden, Response{
		Code: FORBIDDEN,
		Ok:   false,
		Msg:  msg,
	})
}

func (g *Gin) RespUnauthorized(msg string) {
	if msg == "" {
		msg = GetMsg(UNAUTHORIZED)
	}
	g.C.Abort()
	g.C.JSON(http.StatusUnauthorized, Response{
		Code: UNAUTHORIZED,
		Ok:   false,
		Msg:  msg,
	})
}

func (g *Gin) RespError(err error, msg string) {
	if msg == "" {
		msg = GetMsg(ERROR)
	}
	g.C.Abort()
	g.C.JSON(http.StatusInternalServerError, Response{
		Code: ERROR,
		Err:  err.Error(),
		Ok:   false,
		Msg:  msg,
	})
}

func (g *Gin) RespNewError(httpStatus, errCode int, err error, msg string) {
	if msg == "" {
		msg = GetMsg(ERROR)
	}
	g.C.Abort()
	g.C.JSON(httpStatus, Response{
		Code: errCode,
		Err:  err.Error(),
		Ok:   false,
		Msg:  msg,
	})
}

func (g *Gin) RespBadRequest(msg string) {
	g.C.Abort()
	g.C.JSON(http.StatusBadRequest, Response{
		Code: INVALID_PARAMS,
		Ok:   false,
		Msg:  msg,
	})
}

func (g *Gin) Abort(err error) {
	g.C.Abort()
	if err == gorm.ErrRecordNotFound {
		g.RespNewError(http.StatusBadRequest, 0, err, "请求ID不存在")
		return
	} else if strings.Contains(err.Error(), "Error 1062") {
		g.RespNewError(http.StatusBadRequest, 0, err, "已存在此资源")
		return
	}
	g.RespError(err, "")
}

type Controller struct {
	db          *gorm.DB
	routerGroup *gin.RouterGroup
	model       model.Model
}

func NewController(db *gorm.DB, router *gin.RouterGroup, m model.Model) *Controller {
	return &Controller{db: db, routerGroup: router, model: m}
}

func (controller *Controller) GetRouter() *gin.RouterGroup {
	return controller.routerGroup
}

func (controller *Controller) Wrap(f func(g *Gin)) gin.HandlerFunc {
	return func(c *gin.Context) {
		f(&Gin{C: c})
	}
}

func (controller *Controller) newModel() interface{} {
	return controller.InitValue(controller.model)
}

func (controller *Controller) InitValue(data interface{}) interface{} {
	reflectVal := reflect.ValueOf(data)
	t := reflect.Indirect(reflectVal).Type()
	newObj := reflect.New(t)
	return newObj.Interface()
}

func (controller *Controller) BuildCreateApi(validateObj interface{}, fun func(data interface{}) (interface{}, error)) {
	controller.GetRouter().POST(controller.model.GetResourceName(), controller.Wrap(func(g *Gin) {
		// 验证字段
		obj := controller.InitValue(validateObj)
		if ok := g.ValidateStruct(obj); !ok {
			return
		}
		if ok := g.ValidateStruct(obj); !ok {
			return
		}

		result, err := fun(obj)
		if err != nil {
			g.Abort(err)
			return
		}

		g.RespSuccess(result, "创建成功")
	}))
}

func (controller *Controller) BuildRetrieveApi() {
	controller.GetRouter().GET(controller.model.GetResourceName()+"/:id", controller.Wrap(func(g *Gin) {
		id, ok := g.ValidateId()
		if !ok {
			return
		}
		data := controller.newModel()
		reflect.ValueOf(data).Elem().FieldByName("ID").Set(reflect.ValueOf(id))
		if err := controller.db.Where("id = ?", id).Take(data).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				g.RespNewError(http.StatusBadRequest, 0, err, "请求ID不存在")
				return
			}
			g.Abort(err)
			return
		}
		g.RespSuccess(data, "查询成功")
	}))
}

func (controller *Controller) BuildSearchApi(fun func(param *SearchParam) (interface{}, error)) {
	controller.GetRouter().POST(controller.model.GetResourceName()+"/_search", controller.Wrap(func(g *Gin) {
		var body SearchParam
		ok := g.ValidateStruct(&body)
		if !ok {
			return
		}
		param := &body
		model := controller.newModel()
		slice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(model)), 0, 0)
		s := reflect.New(slice.Type()).Interface()
		var err error
		if fun == nil {
			err = param.Search(controller.db).Find(s).Error
		} else {
			s, err = fun(param)
		}
		if err != nil {
			g.Abort(err)
			return
		}
		g.RespSuccess(param.NewPagination(s, controller.model), "查询成功")

	}))
}

func (controller *Controller) BuildUpdateApi(validateObj interface{}, fun func(id int, data interface{}) (interface{}, error)) {
	controller.GetRouter().PATCH(controller.model.GetResourceName()+"/:id", controller.Wrap(func(g *Gin) {
		id, ok := g.ValidateId()
		if !ok {
			return
		}
		// 验证字段
		obj := controller.InitValue(validateObj)
		if ok := g.ValidateStruct(obj); !ok {
			return
		}
		result, err := fun(id, obj)
		if err != nil {
			g.Abort(err)
			return
		}
		g.RespSuccess(result, "修改成功")
	}))

}

func (controller *Controller) BuildDeleteApi(fun func(id int) error) {
	controller.GetRouter().DELETE(controller.model.GetResourceName()+"/:id", controller.Wrap(func(g *Gin) {
		id, ok := g.ValidateId()
		if !ok {
			return
		}
		err := fun(id)
		if err != nil {
			g.Abort(err)
			return
		}
		g.RespSuccess(nil, "删除成功")
	}))
}
