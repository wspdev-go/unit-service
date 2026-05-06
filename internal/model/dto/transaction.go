package dto

type Transaction struct {
	TrID              string `json:"msg_id"`
	TransactionLinkID int    `json:"transaction_link_id"`
	TrData            string `json:"msg_data"`
}

func ConvertTransaction(tr *Transaction) *Transaction {
	return &Transaction{
		TrID:              tr.TrID,
		TransactionLinkID: tr.TransactionLinkID,
		TrData:            tr.TrData,
	}
}
