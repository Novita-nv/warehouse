package product

import (
	"bytes"
	"encoding/json"

	"io"
	"net/http"
	"time"

	"gitlab.privy.id/go_graphql/internal/appctx"
	"gitlab.privy.id/go_graphql/internal/entity"
	repositories "gitlab.privy.id/go_graphql/internal/repositories/product"
	"gitlab.privy.id/go_graphql/internal/ucase/contract"
)

type createProduct struct {
	productRepo repositories.ProductRepoInterface
}

func NewCreateProduct(productRepo repositories.ProductRepoInterface) contract.UseCase {
	return &createProduct{
		productRepo: productRepo,
	}
}

func (u *createProduct) Serve(data *appctx.Data) appctx.Response {
	request := data.Request
	ctx := request.Context()

	rawRequestData := bytes.NewBuffer(make([]byte, 0, request.ContentLength))
	_, errRawRequestData := io.Copy(rawRequestData, request.Body)
	if errRawRequestData != nil {
		return *appctx.NewResponse().WithStatus("ERROR").
			WithCode(http.StatusUnprocessableEntity).WithMessage("Create Product Failed").
			WithEntity("createProductFailed").WithState("createProductFailed").WithError(errRawRequestData)
	}

	request.Body.Close()
	request.Body = io.NopCloser(rawRequestData)

	var payload entity.ProductInput
	errRequestData := json.Unmarshal(rawRequestData.Bytes(), &payload)
	if errRequestData != nil {
		return *appctx.NewResponse().WithStatus("ERROR").
			WithCode(http.StatusUnprocessableEntity).WithMessage("Create Product Failed").
			WithEntity("createProductFailed").WithState("createProductFailed").WithError(errRawRequestData)
	}

	total := payload.ProductIn - payload.ProductOut

	input := entity.ProductCreated{
		ProductName: payload.Name,
		ProductIn:   payload.ProductIn,
		ProductOut:  payload.ProductOut,
		Total:       total,
		User:        payload.User,
		CreatedAt:   time.Now().Local(),
		UpdatedAt:   time.Now().Local(),
		DeletedAt:   nil,
	}


	err := u.productRepo.CreateProduct(ctx, &input)
	if err != nil {
		return *appctx.NewResponse().WithStatus("ERROR").
			WithCode(http.StatusInternalServerError).WithMessage("Create Product Failed").
			WithEntity("createProductFailed").WithState("createProductFailed").WithError(err)
	}

	return *appctx.NewResponse().WithData(input).
		WithStatus("SUCCESS").
		WithCode(http.StatusOK).WithMessage("Create Product Success").
		WithEntity("createProduct").WithState("createProductSuccess")
}
