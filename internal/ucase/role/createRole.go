package role

import (
	"bytes"
	"encoding/json"

	"io"
	"net/http"
	"time"

	"gitlab.privy.id/go_graphql/internal/appctx"
	"gitlab.privy.id/go_graphql/internal/entity"
	repositories "gitlab.privy.id/go_graphql/internal/repositories/role"
	"gitlab.privy.id/go_graphql/internal/ucase/contract"
)

type createRole struct {
	roleRepo repositories.RoleRepoInterface
}

func NewCreateRole(roleRepo repositories.RoleRepoInterface) contract.UseCase{
	return &createRole{
		roleRepo: roleRepo,
	}
}

func (u *createRole) Serve(data *appctx.Data) appctx.Response {
	request  := data.Request
	ctx := request.Context()

	rawRequestData := bytes.NewBuffer(make([]byte, 0, request.ContentLength))
	_, errRawRequestData := io.Copy(rawRequestData, request.Body)
	if errRawRequestData != nil {
		return *appctx.NewResponse().WithStatus("ERROR").
		WithCode(http.StatusUnprocessableEntity).WithMessage("Create Role Failed").
		WithEntity("createRoleFailed").WithState("createRoleFailed").WithError(errRawRequestData)
	}
	
	request.Body.Close()
	request.Body = io.NopCloser(rawRequestData)

	var payload entity.RoleInput
	errRequestData := json.Unmarshal(rawRequestData.Bytes(), &payload)
	if errRequestData != nil {
		return *appctx.NewResponse().WithStatus("ERROR").
		WithCode(http.StatusUnprocessableEntity).WithMessage("Create Role Failed").
		WithEntity("createRoleFailed").WithState("createRoleFailed").WithError(errRawRequestData)
	}


	input := entity. RoleCreated{
		RoleName: payload.Name,
		CreatedAt: time.Now().Local(),
		UpdatedAt: time.Now().Local(),
		DeletedAt: nil,
	}

	err := u.roleRepo.CreateRole(ctx, &input)
	if err != nil {
		return *appctx.NewResponse().WithStatus("ERROR").
		WithCode(http.StatusInternalServerError).WithMessage("Create Role Failed").
		WithEntity("createRoleFailed").WithState("createRoleFailed").WithError(errRawRequestData)
	}

	return *appctx.NewResponse().WithData(input).
	WithStatus("SUCCESS").
	WithCode(http.StatusOK).WithMessage("Create Role Success").
	WithEntity("createRole").WithState("createRoleSuccess")
}