package donation_rest

import (
	"encoding/json"
	"net/http"

	"github.com/dzemildupljak/risc_monolith/server/domain"
	"github.com/dzemildupljak/risc_monolith/server/usecase"
	donationevent_usecase "github.com/dzemildupljak/risc_monolith/server/usecase/donation_event"
	"github.com/dzemildupljak/risc_monolith/server/utils"
)

type DonationEventController struct {
	dec donationevent_usecase.DonationEventUsecase
	log usecase.Logger
}

func NewDonationEventController(
	d donationevent_usecase.DonationEventUsecase,
	l usecase.Logger) *DonationEventController {

	return &DonationEventController{
		dec: d,
		log: l,
	}
}

func (dc *DonationEventController) DonationList(w http.ResponseWriter, r *http.Request) {

	params := domain.DonationEventsListParams{
		RowOrder:  "",
		LimitSize: 10,
	}
	donations, err := dc.dec.GetDonationEventsList(r.Context(), params)

	if err != nil {
		dc.log.LogError("DonationEventController-DonationList: %s", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
	}

	if len(donations) == 0 {
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
			Status:  true,
			Message: "",
			Data:    donations,
		})

}
