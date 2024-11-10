package request

type MakePaymentReq struct {
	Amount float64 `json:"amount" validate:"required"`
}
