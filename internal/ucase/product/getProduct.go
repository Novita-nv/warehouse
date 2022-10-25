package product

import (
	"fmt"
	"net/http"

	"gitlab.privy.id/go_graphql/internal/appctx"
	repositories "gitlab.privy.id/go_graphql/internal/repositories/product"
	"gitlab.privy.id/go_graphql/internal/ucase/contract"
)

type getProduct struct {
	productRepo repositories.ProductRepoInterface
}

func NewGetProducts(productRepo repositories.ProductRepoInterface) contract.UseCase {
	return &getProduct{
		productRepo: productRepo,
	}
}

func (u *getProduct) Serve(data *appctx.Data) appctx.Response {
	request := data.Request
	ctx := request.Context()

	products, err := u.productRepo.GetProducts(ctx)
	if err != nil {
		fmt.Println("============", err)
		return *appctx.NewResponse().WithStatus("ERROR").
			WithCode(http.StatusInternalServerError).WithMessage("get Product Failed").
			WithEntity("getProductFailed").WithState("getProductFailed").WithError(err)
	}

	return *appctx.NewResponse().WithData(products).
		WithStatus("SUCCESS").
		WithCode(http.StatusOK).WithMessage("Get Product Success").
		WithEntity("getProduct").WithState("getProductSuccess")
}
