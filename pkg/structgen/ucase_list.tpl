// Package {{.PackageName}}
// Automatic generated
package {{.PackageName}}

import (
	"fmt"

	"gitlab.privy.id/go_graphql/internal/appctx"
	"gitlab.privy.id/go_graphql/internal/common"
	"gitlab.privy.id/go_graphql/internal/consts"
	"gitlab.privy.id/go_graphql/internal/presentations"
	"gitlab.privy.id/go_graphql/internal/repositories"
	"gitlab.privy.id/go_graphql/pkg/logger"
	"gitlab.privy.id/go_graphql/pkg/tracer"

	ucase "gitlab.privy.id/go_graphql/internal/ucase/contract"
)

type {{.StructName}}List struct {
	repo repositories.{{.EntityName}}er
}

func New{{.EntityName}}List(repo repositories.{{.EntityName}}er) ucase.UseCase {
	return &{{.StructName}}List{repo: repo}
}

// Serve {{.EntityName}} list data
func (u *{{.StructName}}List) Serve(dctx *appctx.Data) appctx.Response {
	var (
		param presentations.{{.EntityName}}Query
		ctx   = tracer.SpanStart(dctx.Request.Context(), "ucase.{{.FileName}}_list")
		lf    = logger.NewFields(
			logger.EventName("{{.StructName}}List"),
		)
	)
    defer tracer.SpanFinish(ctx)

	err := dctx.Cast(&param)
	if err != nil {
		logger.WarnWithContext(ctx, fmt.Sprintf("error parsing query url: %v", err), lf...)
		return *appctx.NewResponse().WithMsgKey(consts.RespValidationError)
	}

	param.Limit = common.LimitDefaultValue(param.Limit)
	param.Page = common.PageDefaultValue(param.Page)

	p, count, err := u.repo.FindWithCount(ctx, param)
	if err != nil {
	    tracer.SpanError(ctx, err)
		logger.ErrorWithContext(ctx, fmt.Sprintf("error find data to database: %v", err), lf...)
		return *appctx.NewResponse().WithMsgKey(consts.RespError)
	}

	logger.InfoWithContext(ctx, fmt.Sprintf("success fetch {{.TableName}} to database"), lf...)
	return *appctx.NewResponse().
            WithMsgKey(consts.RespSuccess).
            WithData(p).
            WithMeta(appctx.MetaData{
                    Page:       param.Page,
                    Limit:      param.Limit,
                    TotalCount: count,
                    TotalPage:  common.PageCalculate(count, param.Limit),
            })
}