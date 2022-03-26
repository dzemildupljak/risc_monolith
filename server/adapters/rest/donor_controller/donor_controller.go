package donor_rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dzemildupljak/risc_monolith/server/domain"
	"github.com/dzemildupljak/risc_monolith/server/usecase"
	"github.com/dzemildupljak/risc_monolith/server/usecase/donor_usecase"
	"github.com/dzemildupljak/risc_monolith/server/utils"
	"github.com/gorilla/mux"
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
func (dc *DonorController) CreateNewDonor(w http.ResponseWriter, r *http.Request) {
	newDonor := &domain.CreateDonorParamsDTO{}

	err := json.NewDecoder(r.Body).Decode(newDonor)
	if err != nil {
		dc.log.LogError("deserialization of donor json failed", "error", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to create donor. Please try again later",
			})
		return
	}

	err = dc.di.CreateNewDonor(r.Context(), *newDonor)
	if err != nil {
		dc.log.LogError("creating donor failed", "error", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to create donor. Please try again later",
			})
		return
	}

	dc.log.LogError("successfully created donor", "error", err)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		&utils.GenericResponse{
			Status:  true,
			Message: "Donor created successfully",
		})

}

func (dc *DonorController) ListDonors(w http.ResponseWriter, r *http.Request) {
	donors, err := dc.di.GetAllDonors(r.Context())
	if err != nil {
		dc.log.LogError("DonorController-ListDonors: %s", err)
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
			Status:  true,
			Message: "",
			Data:    donors,
		})
}

func (dc *DonorController) DonorById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	donorId, err := strconv.ParseInt(params["donor_id"], 10, 64)
	if err != nil {
		dc.log.LogError("UserById = donor Id validation failed", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to get donor. Please try again later",
			})
		return
	}

	donor, err := dc.di.GetDonorById(r.Context(), donorId)

	if err != nil {
		dc.log.LogError("DonorController-DonorById: %s", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		&utils.GenericResponse{
			Status:  true,
			Message: "",
			Data:    donor,
		})
}

func (dc *DonorController) DonorByBloodType(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	bloodType, err := strconv.ParseInt(params["blood_type"], 10, 16)
	if err != nil {
		dc.log.LogError("DonorController - DonorByBloodType - donor Id validation failed", "error", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to get donor. Please try again later",
			})
		return
	}

	donor, err := dc.di.GetDonorsByBloodType(r.Context(), int16(bloodType))

	if err != nil {
		dc.log.LogError("DonorController-DonorByBloodType: %s", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		&utils.GenericResponse{
			Status:  true,
			Message: "",
			Data:    donor,
		})
}

func (dc *DonorController) DonorByUniqueNumber(w http.ResponseWriter, r *http.Request) {
	bloodType := &domain.DonorUniqueNumber{}

	err := json.NewDecoder(r.Body).Decode(bloodType)
	if err != nil {
		dc.log.LogError("deserialization of donor json failed", "error", err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to get donor. Please try again later",
			})
		return
	}

	donor, err := dc.di.GetDonorByIdentificationNum(r.Context(), bloodType.UniqueNumber)

	if err != nil {
		dc.log.LogError("DonorController-DonorByBloodType: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(
			&utils.GenericResponse{
				Status:  false,
				Message: "Unable to get donor. Please try again later",
			})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(
		&utils.GenericResponse{
			Status:  true,
			Message: "",
			Data:    donor,
		})
}
