package helpers

type Payment struct {
	Amount  int64   `structs:"amount"`
	AssetId *string `structs:"assetId"`
}
type FuncArg struct {
	Type  string      `structs:"type"`
	Value interface{} `structs:"value"`
}

type FuncCall struct {
	Function ContractFunc `structs:"function"`
	Args     []FuncArg    `structs:"args"`
}

type InvokeScriptBody struct {
	Call    FuncCall  `structs:"call"`
	DApp    string    `structs:"dApp"`
	Payment []Payment `structs:"payment"`
}

func (tx *Transaction) NewInvokeScript(dapp string, funcCall FuncCall, payments []Payment, fee int) {
	tx.InvokeScriptBody = &InvokeScriptBody{
		Call:    funcCall,
		DApp:    dapp,
		Payment: payments,
	}
	tx.Fee = fee
}
