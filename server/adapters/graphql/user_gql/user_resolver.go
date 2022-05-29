package user_gql

import (
	"context"
	"fmt"

	"github.com/dzemildupljak/risc_monolith/server/usecase"
	"github.com/dzemildupljak/risc_monolith/server/usecase/donor_usecase"
	"github.com/dzemildupljak/risc_monolith/server/utils"
	"github.com/graphql-go/graphql"
)

type DonorResolver struct {
	du donor_usecase.DonorUsecase
	l  usecase.Logger
}

func NewDonorResolver(d donor_usecase.DonorUsecase, l usecase.Logger) *DonorResolver {
	return &DonorResolver{
		du: d,
		l:  l,
	}
}

func (dr *DonorResolver) GetDonors(params graphql.ResolveParams) (interface{}, error) {

	ctx := context.Background()

	postSearch, postSearchOk := params.Args["search"].(string)
	if !postSearchOk {
		dr.l.LogError("GetDonors = something went wrong ")
		return &utils.GenericResponse{
			Status:  false,
			Message: "Unable to get donors. Please try again later",
		}, nil

		// return helper.NewHTTPResponse(http.StatusBadRequest, "invalid params", nil), errors.New("invalid params")
	}

	fmt.Println("postSearch============", postSearch)

	donors, err := dr.du.GetAllDonors(ctx)
	if err != nil {
		dr.l.LogError("GetDonors = something went wrong ")
		return &utils.GenericResponse{
			Status:  false,
			Message: "Unable to get donors. Please try again later",
		}, err
	}

	dr.l.LogAccess("GetDonors = donors list ")
	return &utils.GenericResponse{
		Data:    donors,
		Status:  true,
		Message: "",
	}, nil

	// res, err := r.postUsecase.GetAll(postSearch)
	// if err != nil {
	// 	return helper.NewHTTPResponse(http.StatusBadRequest, err.Error(), nil), err
	// }

	// return helper.NewHTTPResponse(http.StatusOK, "get posts", res), nil
}
