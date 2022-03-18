package donor_rest

import (
	"encoding/json"
	"net/http"

	"github.com/dzemildupljak/risc_monolith/server/usecase"
	"github.com/dzemildupljak/risc_monolith/server/usecase/donor_usecase"
	"github.com/dzemildupljak/risc_monolith/server/utils"
)

type DonorController struct {
	di  donor_usecase.DonorInteractor
	log usecase.Logger
}

func NewDonorController(
	di donor_usecase.DonorInteractor,
	l usecase.Logger) *DonorController {

	return &DonorController{
		di:  di,
		log: l,
	}
}

func (dc *DonorController) ListDonors(w http.ResponseWriter, r *http.Request) {
	donors, err := dc.di.GetAllDonors(r.Context())
	if err != nil {
		dc.log.LogError("UserController-Index: %s", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
	}

	if len(donors) == 0 {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Currently have no donors",
			})
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		&utils.GenericResponse{
			Status:  false,
			Message: "",
			Data:    donors,
		})
}
