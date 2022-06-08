package api

import (
	"fmt"
	"mayfly-go/internal/sys/api/vo"
	"mayfly-go/internal/sys/domain/entity"
	"mayfly-go/pkg/biz"
	"mayfly-go/pkg/ctx"
	"mayfly-go/pkg/ginx"
	"mayfly-go/pkg/model"
	"net/url"
)

type MockData struct {
}

func (r *MockData) MockDatas(rc *ctx.ReqCtx) {
	g := rc.GinCtx
	service := g.Query("service")
	path := g.Query("path")
	sql := "SELECT id, service, method, path, description, status, create_time, creator, update_time, modifier FROM t_mock_data WHERE 1 = 1 "
	if path != "" {
		sql = sql + " AND path LIKE '%" + path + "%'"
	}
	if service != "" {
		sql = sql + fmt.Sprintf("AND service = '%s'", service)
	}
	sql = sql + " ORDER BY update_time DESC"
	rc.ResData = model.GetPageBySql(sql, ginx.GetPageParam(g), new([]vo.MockDataListVO))
}

func (r *MockData) MockDataServices(rc *ctx.ReqCtx) {
	sql := "SELECT DISTINCT(service) service FROM t_mock_data"
	var res []entity.MockData
	model.GetListBySql2Model(sql, &res)
	rc.ResData = res
}

func (r *MockData) MockDataDetail(rc *ctx.ReqCtx) {
	g := rc.GinCtx
	MockDataDetail := &entity.MockData{}
	model.GetById(MockDataDetail, uint64(ginx.QueryInt(g, "id", 0)))
	rc.ResData = MockDataDetail
}

// 保存mock内容
func (r *MockData) SaveMockData(rc *ctx.ReqCtx) {
	g := rc.GinCtx
	mockData := &entity.MockData{}
	ginx.BindJsonAndValid(g, mockData)

	if mockData.RealUrl != "" {
		_, err := url.ParseRequestURI(mockData.RealUrl)
		biz.ErrIsNilAppendErr(err, "目标地址格式错误: %s")
	}

	mockData.SetBaseInfo(rc.LoginAccount)

	if mockData.Id == 0 {
		condition := &entity.MockData{}
		condition.Path = mockData.Path
		condition.Method = mockData.Method
		err := model.GetBy(condition)
		biz.IsTrue(err != nil, "该请求路径已配置mock数据")
		model.Insert(mockData)
	} else {
		model.UpdateById(mockData)
	}
}

// 删除mock
func (r *MockData) DelMockData(rc *ctx.ReqCtx) {
	model.DeleteById(new(entity.MockData), uint64(ginx.PathParamInt(rc.GinCtx, "id")))
}

// /mock-datas/:service/:method?username=xx
// func (r *MockData) GetMockData(rc *ctx.ReqCtx) {
// 	g := rc.GinCtx
// 	service := g.Param("service")
// 	method := g.Param("method")
// 	params := utils.MapBuilder("method", method).ToMap()
// 	params["service"] = service
// 	// 调用该mock数据的用户，若该数据指定了生效用户，则需要校验是否可访问
// 	username := g.Query("username")
// 	if username != "" {
// 		params["username"] = username
// 	}
// 	// 记录日志使用
// 	rc.ReqParam = params

// 	mockData := &entity.MockData{}
// 	mockData.Service = service
// 	mockData.Method = method
// 	err := model.GetBy(mockData)
// 	// 数据不存在或者状态为禁用
// 	biz.IsTrue(err == nil && mockData.Status == 1, "无该mock数据")

// 	eu := mockData.EffectiveUser
// 	// 如果设置的生效用户为空，则表示所有用户都生效
// 	if eu == "" {
// 		rc.ResData = mockData.Data
// 		return
// 	}
// 	biz.IsTrue(utils.StrLen(username) != 0, "该用户无法访问该mock数据")
// 	// 判断该用户是否在该数据指定的生效用户中
// 	for _, e := range strings.Split(eu, ",") {
// 		if username == e {
// 			rc.ResData = mockData.Data
// 			return
// 		}
// 	}
// 	panic(biz.NewBizErr("该用户无法访问该mock数据"))
// }
