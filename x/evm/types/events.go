package types

// Evm module events
const (
	EventTypeEthermint  = TypeMsgEthermint
	EventTypeEthereumTx = TypeMsgEthereumTx

	EventTypeLog = "log"

	AttributeKeyContractAddress = "contract"
	AttributeKeyRecipient       = "recipient"
	AttributeValueCategory      = ModuleName
)
