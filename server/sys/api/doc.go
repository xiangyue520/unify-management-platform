package api

import (
	"fmt"
	"mayfly-go/base/biz"
	"mayfly-go/base/ctx"
	"mayfly-go/base/ginx"
	"mayfly-go/base/model"
	"mayfly-go/server/sys/api/vo"
	"mayfly-go/server/sys/domain/entity"
	"time"
)

type Doc struct {
}

func (r *Doc) Docs(rc *ctx.ReqCtx) {
	g := rc.GinCtx
	category := g.Query("category")
	title := g.Query("title")
	sql := "SELECT id, title, category, create_time, creator, update_time, modifier FROM t_doc WHERE 1 = 1 "
	if title != "" {
		sql = sql + " AND title LIKE '%" + title + "%'"
	}
	if category != "" {
		sql = sql + fmt.Sprintf("AND category = '%s'", category)
	}
	sql = sql + " ORDER BY update_time DESC"
	rc.ResData = model.GetPageBySql(sql, ginx.GetPageParam(g), new([]vo.DocListVO))
}

func (r *Doc) DocCategory(rc *ctx.ReqCtx) {
	sql := "SELECT DISTINCT(category) category FROM t_doc"
	var res []entity.Doc
	model.GetListBySql2Model(sql, &res)
	rc.ResData = res
}

func (r *Doc) DocDetail(rc *ctx.ReqCtx) {
	g := rc.GinCtx
	docDetail := &entity.Doc{}
	model.GetById(docDetail, uint64(ginx.PathParamInt(g, "id")))
	rc.ResData = docDetail
}

// 保存文档内容
func (r *Doc) SaveDoc(rc *ctx.ReqCtx) {
	g := rc.GinCtx
	doc := &entity.Doc{}
	ginx.BindJsonAndValid(g, doc)

	doc.SetBaseInfo(rc.LoginAccount)

	if doc.Id == 0 {
		model.Insert(doc)
	} else {
		model.UpdateById(doc)
	}
}

// 上传文档图片
func (r *Doc) UploadDocImage(rc *ctx.ReqCtx) {
	g := rc.GinCtx
	fileheader, err := g.FormFile("file")
	biz.ErrIsNilAppendErr(err, "读取文件失败: %s")

	path := fmt.Sprintf("/files/doc/%d-%s", time.Now().UnixMilli(), fileheader.Filename)
	err = rc.GinCtx.SaveUploadedFile(fileheader, "."+path)
	biz.ErrIsNilAppendErr(err, "读取保存失败: %s")

	rc.ResData = path
}

// 删除文档
func (r *Doc) DelDoc(rc *ctx.ReqCtx) {
	model.DeleteById(new(entity.Doc), uint64(ginx.PathParamInt(rc.GinCtx, "id")))
}
