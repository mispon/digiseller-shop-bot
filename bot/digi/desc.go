package digi

type (
	PPOValue struct {
		Text string `json:"text"`
	}

	ProductPurchasesOption struct {
		ID    int      `json:"id"`
		Value PPOValue `json:"value,omitempty"`
	}

	ProductPurchasesOptions struct {
		ProductID int                      `json:"product_id"`
		Options   []ProductPurchasesOption `json:"options"`
		UnitCnt   int                      `json:"unit_cnt"`
		Lang      string                   `json:"lang"`
		IP        string                   `json:"ip"`
	}

	ProductPaymentData struct {
		Retval  int    `json:"retval"`
		RetDesc string `json:"retdesc"`
		IDPo    int    `json:"id_po"`
	}
)
