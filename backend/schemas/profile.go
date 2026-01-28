package schemas

type BalanceFull struct {
	Id string `json:"id"`
	Balance int64 `json:"balance"`
	TotalTransactions int64 `json:"total_transactions"`
}