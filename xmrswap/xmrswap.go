// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package xmrswap

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// XmrSwapMetaData contains all meta data concerning the XmrSwap contract.
var XmrSwapMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"feeBps_\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"feeRecipient_\",\"type\":\"address\",\"internalType\":\"addresspayable\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"BPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_FEE_BPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MIN_TIMEOUT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"cancelOffer\",\"inputs\":[{\"name\":\"offerId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claim\",\"inputs\":[{\"name\":\"swapId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"secret\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"feeBps\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"feeRecipient\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"addresspayable\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"keyUsed\",\"inputs\":[{\"name\":\"spendPub\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"offers\",\"inputs\":[{\"name\":\"offerId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"maker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"payout\",\"type\":\"address\",\"internalType\":\"addresspayable\"},{\"name\":\"asset\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"minAmount\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"maxAmount\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"xmrPerAsset\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expiry\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"timeout1Duration\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"timeout2Duration\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"makerSpendPub\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"makerViewPriv\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"active\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"postOffer\",\"inputs\":[{\"name\":\"asset\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"minAmount\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"maxAmount\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"xmrPerAsset\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expiry\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"timeout1Duration\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"timeout2Duration\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"makerSpendPub\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"makerViewPriv\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"payout\",\"type\":\"address\",\"internalType\":\"addresspayable\"}],\"outputs\":[{\"name\":\"offerId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"refund\",\"inputs\":[{\"name\":\"swapId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"secret\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"secretMatches\",\"inputs\":[{\"name\":\"secret\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setReady\",\"inputs\":[{\"name\":\"swapId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"swaps\",\"inputs\":[{\"name\":\"swapId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"taker\",\"type\":\"address\",\"internalType\":\"addresspayable\"},{\"name\":\"maker\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"payout\",\"type\":\"address\",\"internalType\":\"addresspayable\"},{\"name\":\"asset\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"makerSpendPub\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"takerSpendPub\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"timeout1\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"timeout2\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"stage\",\"type\":\"uint8\",\"internalType\":\"enumXmrSwap.Stage\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"takeOffer\",\"inputs\":[{\"name\":\"offerId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"takerSpendPub\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"takerViewPriv\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"swapId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"Claimed\",\"inputs\":[{\"name\":\"swapId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"secret\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"payout\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"fee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OfferCancelled\",\"inputs\":[{\"name\":\"offerId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OfferPosted\",\"inputs\":[{\"name\":\"offerId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"maker\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"asset\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"minAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"maxAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"xmrPerAsset\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"expiry\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Ready\",\"inputs\":[{\"name\":\"swapId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Refunded\",\"inputs\":[{\"name\":\"swapId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"secret\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SwapCreated\",\"inputs\":[{\"name\":\"swapId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"offerId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"taker\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"xmrPiconero\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"takerSpendPub\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"takerViewPriv\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"timeout1\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"timeout2\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AlreadyReady\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AmountOutOfRange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"BadAmounts\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FeeTooHigh\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSecret\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSwap\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"KeyAlreadyUsed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotMaker\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotTaker\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotTimeToRefund\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OfferExpired\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OfferNotActive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SwapCompleted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TimeoutTooShort\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TooEarlyToClaim\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TooLateToClaim\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"WrongValue\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroKey\",\"inputs\":[]}]",
	Bin: "0x60c060405234801561000f575f5ffd5b5060405161231438038061231483398101604081905261002e916100ca565b60017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055606461ffff831611156100785760405163cd4e616760e01b815260040160405180910390fd5b61ffff82161580159061009257506001600160a01b038116155b156100b05760405163d92e233d60e01b815260040160405180910390fd5b61ffff9091166080526001600160a01b031660a052610113565b5f5f604083850312156100db575f5ffd5b825161ffff811681146100ec575f5ffd5b60208401519092506001600160a01b0381168114610108575f5ffd5b809150509250929050565b60805160a0516121d26101425f395f81816101860152610d3f01525f81816101260152610cbe01526121d25ff3fe6080604052600436106100e4575f3560e01c80636bcfee1511610087578063e4cbfed511610057578063e4cbfed514610350578063eb84e7f214610363578063f952279e1461040f578063fb0e72b81461042e575f5ffd5b80636bcfee15146102dd57806384cc9dfb146102fe578063d55be8c61461031d578063e4683a7914610331575f5ffd5b806346904840116100c25780634690484014610175578063474d3ff0146101c05780634977149d14610280578063543ad1df146102af575f5ffd5b8063249d39e9146100e857806324a9d853146101155780633faff49114610148575b5f5ffd5b3480156100f3575f5ffd5b506100fd61271081565b60405161ffff90911681526020015b60405180910390f35b348015610120575f5ffd5b506100fd7f000000000000000000000000000000000000000000000000000000000000000081565b348015610153575f5ffd5b50610167610162366004611e8f565b61045c565b60405190815260200161010c565b348015610180575f5ffd5b506101a87f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b03909116815260200161010c565b3480156101cb575f5ffd5b506102686101da366004611f34565b5f602081905290815260409020805460018201546002830154600384015460048501546005860154600687015460078801546008909801546001600160a01b039788169896881697909516956001600160801b0380861696600160801b968790049091169567ffffffffffffffff808616956801000000000000000081048216959290041692909160ff168c565b60405161010c9c9b9a99989796959493929190611f4b565b34801561028b575f5ffd5b5061029f61029a366004611fe9565b61094d565b604051901515815260200161010c565b3480156102ba575f5ffd5b506102c4610e1081565b60405167ffffffffffffffff909116815260200161010c565b3480156102e8575f5ffd5b506102fc6102f7366004611f34565b610961565b005b348015610309575f5ffd5b506102fc610318366004611fe9565b610aee565b348015610328575f5ffd5b506100fd606481565b34801561033c575f5ffd5b506102fc61034b366004611fe9565b610dd7565b61016761035e366004612009565b610fff565b34801561036e575f5ffd5b506103f961037d366004611f34565b600160208190525f9182526040909120805491810154600282015460038301546004840154600585015460068601546007909601546001600160a01b0397881697958616969486169590931693919290919067ffffffffffffffff8082169168010000000000000000810490911690600160801b900460ff168a565b60405161010c9a9998979695949392919061204c565b34801561041a575f5ffd5b506102fc610429366004611f34565b611503565b348015610439575f5ffd5b5061029f610448366004611f34565b60026020525f908152604090205460ff1681565b5f6001600160a01b03821661049d576040517fd92e233d00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b8315806104a8575082155b156104c657604051637dec24e360e11b815260040160405180910390fd5b6001600160801b038a1615806104ed5750896001600160801b0316896001600160801b0316105b806104f6575087155b1561052d576040517fc6fd446400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610e1067ffffffffffffffff871610806105525750610e1067ffffffffffffffff8616105b15610589576040517fe143262500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b428767ffffffffffffffff16116105b357604051639cb1308760e01b815260040160405180910390fd5b5f8481526002602052604090205460ff16156105e2576040516312b3b8ff60e01b815260040160405180910390fd5b5f848152600260205260408120805460ff19166001179055600380544692309233929161060e83612106565b909155506040805160208101959095526bffffffffffffffffffffffff19606094851b8116918601919091529190921b1660548301526068820152608801604051602081830303815290604052805190602001209050604051806101800160405280336001600160a01b03168152602001836001600160a01b031681526020018c6001600160a01b031681526020018b6001600160801b031681526020018a6001600160801b031681526020018981526020018867ffffffffffffffff1681526020018767ffffffffffffffff1681526020018667ffffffffffffffff168152602001858152602001848152602001600115158152505f5f8381526020019081526020015f205f820151815f015f6101000a8154816001600160a01b0302191690836001600160a01b031602179055506020820151816001015f6101000a8154816001600160a01b0302191690836001600160a01b031602179055506040820151816002015f6101000a8154816001600160a01b0302191690836001600160a01b031602179055506060820151816003015f6101000a8154816001600160801b0302191690836001600160801b0316021790555060808201518160030160106101000a8154816001600160801b0302191690836001600160801b0316021790555060a0820151816004015560c0820151816005015f6101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555060e08201518160050160086101000a81548167ffffffffffffffff021916908367ffffffffffffffff1602179055506101008201518160050160106101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555061012082015181600601556101408201518160070155610160820151816008015f6101000a81548160ff0219169083151502179055509050508a6001600160a01b0316336001600160a01b0316827f5d98652ad01d0501cc21957f9e428beced1f29456b42459ff3c1b5a4c00fb8448d8d8d8d60405161093794939291906001600160801b039485168152929093166020830152604082015267ffffffffffffffff91909116606082015260800190565b60405180910390a49a9950505050505050505050565b5f6109588383611599565b90505b92915050565b5f818152600160205260408120906007820154600160801b900460ff16600381111561098f5761098f612038565b036109ad57604051631115766760e01b815260040160405180910390fd5b60036007820154600160801b900460ff1660038111156109cf576109cf612038565b036109ed5760405163066916a960e01b815260040160405180910390fd5b60026007820154600160801b900460ff166003811115610a0f57610a0f612038565b03610a46576040517f2000ee6c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80546001600160a01b03163314610a7057604051635211a07960e01b815260040160405180910390fd5b600781015467ffffffffffffffff164210610a9e5760405163497df9d160e01b815260040160405180910390fd5b60078101805460ff60801b191670020000000000000000000000000000000017905560405182907f5fc23b25552757626e08b316cc2387ad1bc70ee1594af7204db4ce0c39f5d15f905f90a25050565b610af66115db565b5f828152600160205260408120906007820154600160801b900460ff166003811115610b2457610b24612038565b03610b4257604051631115766760e01b815260040160405180910390fd5b60036007820154600160801b900460ff166003811115610b6457610b64612038565b03610b825760405163066916a960e01b815260040160405180910390fd5b60018101546001600160a01b03163314610baf5760405163b331e42160e01b815260040160405180910390fd5b600781015467ffffffffffffffff1642108015610bec575060026007820154600160801b900460ff166003811115610be957610be9612038565b14155b15610c23576040517fd71d60b500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600781015468010000000000000000900467ffffffffffffffff164210610c5d5760405163497df9d160e01b815260040160405180910390fd5b610c6b828260050154611599565b610c885760405163abab6bd760e01b815260040160405180910390fd5b60078101805460ff60801b191670030000000000000000000000000000000017905560048101545f9061271090610ce49061ffff7f0000000000000000000000000000000000000000000000000000000000000000169061211e565b610cee9190612135565b90505f818360040154610d019190612154565b60038401546002850154919250610d25916001600160a01b03918216911683611609565b8115610d64576003830154610d64906001600160a01b03167f000000000000000000000000000000000000000000000000000000000000000084611609565b604080518581526020810183905290810183905285907fd9126119fa9c516709ff9f21d9a56ddf245584c2d2ea9261586dc9ca978af6189060600160405180910390a2505050610dd360017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055565b5050565b610ddf6115db565b5f828152600160205260408120906007820154600160801b900460ff166003811115610e0d57610e0d612038565b03610e2b57604051631115766760e01b815260040160405180910390fd5b60036007820154600160801b900460ff166003811115610e4d57610e4d612038565b03610e6b5760405163066916a960e01b815260040160405180910390fd5b80546001600160a01b03163314610e9557604051635211a07960e01b815260040160405180910390fd5b600781015468010000000000000000900467ffffffffffffffff1642108015610ef55750600781015467ffffffffffffffff1642101580610ef5575060026007820154600160801b900460ff166003811115610ef357610ef3612038565b145b15610f2c576040517f65430c1e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610f3a828260060154611599565b610f575760405163abab6bd760e01b815260040160405180910390fd5b60078101805460ff60801b1916700300000000000000000000000000000000179055600381015481546004830154610f9c926001600160a01b03908116921690611609565b827e7c875846b687732a7579c19bb1dade66cd14e9f4f809565e2b2b5e76c72b4f83604051610fcd91815260200190565b60405180910390a250610dd360017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055565b5f6110086115db565b5f858152602081905260409020600881015460ff1661103a57604051634d80826b60e11b815260040160405180910390fd5b600581015467ffffffffffffffff16421061106857604051639cb1308760e01b815260040160405180910390fd5b831580611073575082155b1561109157604051637dec24e360e11b815260040160405180910390fd5b60038101546001600160801b03168510806110bf57506003810154600160801b90046001600160801b031685115b156110f6576040517fc64200e900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f8481526002602052604090205460ff1615611125576040516312b3b8ff60e01b815260040160405180910390fd5b5f848152600260208190526040909120805460ff191660011790558101546001600160a01b03166111755784341461117057604051632635240760e21b815260040160405180910390fd5b6111ae565b341561119457604051632635240760e21b815260040160405180910390fd5b60028101546111ae906001600160a01b03163330886116f4565b60088101805460ff1916905560058101545f906111e19068010000000000000000900467ffffffffffffffff1642612167565b60058301549091505f9061120690600160801b900467ffffffffffffffff1683612167565b60408051602081018b90526bffffffffffffffffffffffff193360601b16918101919091526054810188905290915060740160408051601f19818403018152919052805160209091012093505f5f85815260016020526040902060070154600160801b900460ff16600381111561127f5761127f612038565b1461129d57604051631115766760e01b815260040160405180910390fd5b604080516101408101825233815284546001600160a01b039081166020830152600180870154821693830193909352600286015416606082015260808101899052600685015460a082015260c0810188905267ffffffffffffffff80851660e08301528316610100820152906101208201525f85815260016020818152604092839020845181546001600160a01b039182167fffffffffffffffffffffffff000000000000000000000000000000000000000091821617835592860151938201805494821694841694909417909355928401516002840180549184169183169190911790556060840151600380850180549290941691909216179091556080830151600483015560a0830151600583015560c0830151600683015560e083015160078301805461010086015167ffffffffffffffff90811668010000000000000000027fffffffffffffffffffffffffffffffff00000000000000000000000000000000909216931692909217919091178082556101208501519260ff60801b1990911690600160801b90849081111561143957611439612038565b02179055509050505f670de0b6b3a764000084600401548961145b919061211e565b6114659190612135565b604080518a8152602081018390529081018990526060810188905267ffffffffffffffff8086166080830152841660a082015290915033908a9087907f6ac60c1e1d7b18793066f728ab948222b6ac5790d9af38983264967ab615155b9060c00160405180910390a4505050506114fb60017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055565b949350505050565b5f818152602081905260409020600881015460ff1661153557604051634d80826b60e11b815260040160405180910390fd5b80546001600160a01b0316331461155f5760405163b331e42160e01b815260040160405180910390fd5b60088101805460ff1916905560405182907f3f9cb69d022b6ec319f86f2df848bcce01f2fc51c9f86396779a8081cf6ca2ea905f90a25050565b5f806115a48461172a565b9050805f036115b6575f91505061095b565b5f5f6115c18361175f565b91509150846115d08383611880565b149695505050505050565b6115e36118b3565b60027f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055565b6001600160a01b0383166116db575f826001600160a01b0316826040515f6040518083038185875af1925050503d805f8114611660576040519150601f19603f3d011682016040523d82523d5f602084013e611665565b606091505b50509050806116d5576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f455448207472616e73666572206661696c65640000000000000000000000000060448201526064015b60405180910390fd5b50505050565b6116ef6001600160a01b0384168383611910565b505050565b611702848484846001611945565b6116d557604051635274afe760e01b81526001600160a01b03851660048201526024016116cc565b5f8060205b801561175857600884811c9460ff1692901b91909117908061175081612187565b91505061172f565b5092915050565b5f5f61178260405180606001604052805f81526020015f81526020015f81525090565b6117a360405180606001604052805f81526020015f81526020015f81525090565b7f216936d3cd6e53fec0a4e231fdd6dc5c692cc7609525a7b2c9562d608f25d51a82527f6666666666666666666666666666666666666666666666666666666666666658602080840191909152600160408085018290525f8452918301819052908201525b841561183d57846001166001036118265761182381836119cb565b90505b600185901c945061183682611bac565b9150611808565b5f61184b8260400151611d68565b90506013600160ff1b03825182900982526013600160ff1b038183602001510960208301819052915196919550909350505050565b5f7f800000000000000000000000000000000000000000000000000000000000000060ff84901b1682176114fb8161172a565b7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f005460020361190e576040517f3ee5aeb500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b61191d8383836001611dcb565b6116ef57604051635274afe760e01b81526001600160a01b03841660048201526024016116cc565b6040517f23b872dd000000000000000000000000000000000000000000000000000000005f8181526001600160a01b038781166004528616602452604485905291602083606481808c5af1925060015f511483166119ba5783831516156119ae573d5f823e3d81fd5b5f883b113d1516831692505b604052505f60605295945050505050565b6119ec60405180606001604052805f81526020015f81526020015f81525090565b611a2c6040518061010001604052805f81526020015f81526020015f81526020015f81526020015f81526020015f81526020015f81526020015f81525090565b6013600160ff1b03836040015185604001510981526013600160ff1b038151800960208201526013600160ff1b03835185510960408201526013600160ff1b03836020015185602001510960608201526013600160ff1b038082606001518360400151097f52036cee2b6ffe738cc740797779e89800700a4d4141d8ab75eb4dca135978a309608082018190526013600160ff1b0390611acc9082612154565b82602001510860a08201526013600160ff1b03816080015182602001510860c08201526013600160ff1b03806060830151611b0e906013600160ff1b03612154565b6013600160ff1b036040850151611b2c906013600160ff1b03612154565b6013600160ff1b038060208a01518a51086013600160ff1b0360208c01518c51080908086013600160ff1b0360a08401518451090982526013600160ff1b038082604001518360600151086013600160ff1b0360c08401518451090960208301526013600160ff1b038160c001518260a001510960408301525092915050565b611bcd60405180606001604052805f81526020015f81526020015f81525090565b611c0d6040518061010001604052805f81526020015f81526020015f81526020015f81526020015f81526020015f81526020015f81526020015f81525090565b6013600160ff1b03602084015184510881526013600160ff1b038151800960208201526013600160ff1b038351800960408201526013600160ff1b036020840151800960608201526040810151611c6b906013600160ff1b03612154565b6080820181905260608201516013600160ff1b03910860a08201526013600160ff1b036040840151800960e08201526013600160ff1b03808260e00151600209611cbc906013600160ff1b03612154565b8260a001510860c08201526013600160ff1b0360c08201516013600160ff1b036060840151611cf2906013600160ff1b03612154565b6013600160ff1b036040860151611d10906013600160ff1b03612154565b866020015108080982526013600160ff1b03806060830151611d39906013600160ff1b03612154565b8360800151088260a001510960208301526013600160ff1b038160c001518260a0015109604083015250919050565b5f80611d7c60026013600160ff1b03612154565b90505f6013600160ff1b03905060405160208152602080820152602060408201528460608201528260808201528160a082015260208160c08360055f19fa611dc2575f5ffd5b51949350505050565b6040517fa9059cbb000000000000000000000000000000000000000000000000000000005f8181526001600160a01b038616600452602485905291602083604481808b5af1925060015f51148316611e3a578383151615611e2e573d5f823e3d81fd5b5f873b113d1516831692505b60405250949350505050565b6001600160a01b0381168114611e5a575f5ffd5b50565b80356001600160801b0381168114611e73575f5ffd5b919050565b803567ffffffffffffffff81168114611e73575f5ffd5b5f5f5f5f5f5f5f5f5f5f6101408b8d031215611ea9575f5ffd5b8a35611eb481611e46565b9950611ec260208c01611e5d565b9850611ed060408c01611e5d565b975060608b01359650611ee560808c01611e78565b9550611ef360a08c01611e78565b9450611f0160c08c01611e78565b935060e08b013592506101008b013591506101208b0135611f2181611e46565b809150509295989b9194979a5092959850565b5f60208284031215611f44575f5ffd5b5035919050565b6001600160a01b038d811682528c811660208301528b1660408201526001600160801b038a811660608301528916608082015260a0810188905267ffffffffffffffff871660c0820152610180810167ffffffffffffffff871660e083015267ffffffffffffffff86166101008301528461012083015283610140830152611fd861016083018415159052565b9d9c50505050505050505050505050565b5f5f60408385031215611ffa575f5ffd5b50508035926020909101359150565b5f5f5f5f6080858703121561201c575f5ffd5b5050823594602084013594506040840135936060013592509050565b634e487b7160e01b5f52602160045260245ffd5b5f610140820190506001600160a01b038c1682526001600160a01b038b1660208301526001600160a01b038a1660408301526001600160a01b03891660608301528760808301528660a08301528560c083015267ffffffffffffffff851660e083015267ffffffffffffffff8416610100830152600483106120dc57634e487b7160e01b5f52602160045260245ffd5b826101208301529b9a5050505050505050505050565b634e487b7160e01b5f52601160045260245ffd5b5f60018201612117576121176120f2565b5060010190565b808202811582820484141761095b5761095b6120f2565b5f8261214f57634e487b7160e01b5f52601260045260245ffd5b500490565b8181038181111561095b5761095b6120f2565b67ffffffffffffffff818116838216019081111561095b5761095b6120f2565b5f81612195576121956120f2565b505f19019056fea2646970667358221220a74d049c22fd34e1c074394eb98a4d21e5f8c8187b418202a216a7db446af23e64736f6c634300081c0033",
}

// XmrSwapABI is the input ABI used to generate the binding from.
// Deprecated: Use XmrSwapMetaData.ABI instead.
var XmrSwapABI = XmrSwapMetaData.ABI

// XmrSwapBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use XmrSwapMetaData.Bin instead.
var XmrSwapBin = XmrSwapMetaData.Bin

// DeployXmrSwap deploys a new Ethereum contract, binding an instance of XmrSwap to it.
func DeployXmrSwap(auth *bind.TransactOpts, backend bind.ContractBackend, feeBps_ uint16, feeRecipient_ common.Address) (common.Address, *types.Transaction, *XmrSwap, error) {
	parsed, err := XmrSwapMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(XmrSwapBin), backend, feeBps_, feeRecipient_)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &XmrSwap{XmrSwapCaller: XmrSwapCaller{contract: contract}, XmrSwapTransactor: XmrSwapTransactor{contract: contract}, XmrSwapFilterer: XmrSwapFilterer{contract: contract}}, nil
}

// XmrSwap is an auto generated Go binding around an Ethereum contract.
type XmrSwap struct {
	XmrSwapCaller     // Read-only binding to the contract
	XmrSwapTransactor // Write-only binding to the contract
	XmrSwapFilterer   // Log filterer for contract events
}

// XmrSwapCaller is an auto generated read-only Go binding around an Ethereum contract.
type XmrSwapCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// XmrSwapTransactor is an auto generated write-only Go binding around an Ethereum contract.
type XmrSwapTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// XmrSwapFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type XmrSwapFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// XmrSwapSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type XmrSwapSession struct {
	Contract     *XmrSwap          // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// XmrSwapCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type XmrSwapCallerSession struct {
	Contract *XmrSwapCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts  // Call options to use throughout this session
}

// XmrSwapTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type XmrSwapTransactorSession struct {
	Contract     *XmrSwapTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts  // Transaction auth options to use throughout this session
}

// XmrSwapRaw is an auto generated low-level Go binding around an Ethereum contract.
type XmrSwapRaw struct {
	Contract *XmrSwap // Generic contract binding to access the raw methods on
}

// XmrSwapCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type XmrSwapCallerRaw struct {
	Contract *XmrSwapCaller // Generic read-only contract binding to access the raw methods on
}

// XmrSwapTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type XmrSwapTransactorRaw struct {
	Contract *XmrSwapTransactor // Generic write-only contract binding to access the raw methods on
}

// NewXmrSwap creates a new instance of XmrSwap, bound to a specific deployed contract.
func NewXmrSwap(address common.Address, backend bind.ContractBackend) (*XmrSwap, error) {
	contract, err := bindXmrSwap(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &XmrSwap{XmrSwapCaller: XmrSwapCaller{contract: contract}, XmrSwapTransactor: XmrSwapTransactor{contract: contract}, XmrSwapFilterer: XmrSwapFilterer{contract: contract}}, nil
}

// NewXmrSwapCaller creates a new read-only instance of XmrSwap, bound to a specific deployed contract.
func NewXmrSwapCaller(address common.Address, caller bind.ContractCaller) (*XmrSwapCaller, error) {
	contract, err := bindXmrSwap(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &XmrSwapCaller{contract: contract}, nil
}

// NewXmrSwapTransactor creates a new write-only instance of XmrSwap, bound to a specific deployed contract.
func NewXmrSwapTransactor(address common.Address, transactor bind.ContractTransactor) (*XmrSwapTransactor, error) {
	contract, err := bindXmrSwap(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &XmrSwapTransactor{contract: contract}, nil
}

// NewXmrSwapFilterer creates a new log filterer instance of XmrSwap, bound to a specific deployed contract.
func NewXmrSwapFilterer(address common.Address, filterer bind.ContractFilterer) (*XmrSwapFilterer, error) {
	contract, err := bindXmrSwap(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &XmrSwapFilterer{contract: contract}, nil
}

// bindXmrSwap binds a generic wrapper to an already deployed contract.
func bindXmrSwap(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := XmrSwapMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_XmrSwap *XmrSwapRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _XmrSwap.Contract.XmrSwapCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_XmrSwap *XmrSwapRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _XmrSwap.Contract.XmrSwapTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_XmrSwap *XmrSwapRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _XmrSwap.Contract.XmrSwapTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_XmrSwap *XmrSwapCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _XmrSwap.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_XmrSwap *XmrSwapTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _XmrSwap.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_XmrSwap *XmrSwapTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _XmrSwap.Contract.contract.Transact(opts, method, params...)
}

// BPS is a free data retrieval call binding the contract method 0x249d39e9.
//
// Solidity: function BPS() view returns(uint16)
func (_XmrSwap *XmrSwapCaller) BPS(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _XmrSwap.contract.Call(opts, &out, "BPS")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// BPS is a free data retrieval call binding the contract method 0x249d39e9.
//
// Solidity: function BPS() view returns(uint16)
func (_XmrSwap *XmrSwapSession) BPS() (uint16, error) {
	return _XmrSwap.Contract.BPS(&_XmrSwap.CallOpts)
}

// BPS is a free data retrieval call binding the contract method 0x249d39e9.
//
// Solidity: function BPS() view returns(uint16)
func (_XmrSwap *XmrSwapCallerSession) BPS() (uint16, error) {
	return _XmrSwap.Contract.BPS(&_XmrSwap.CallOpts)
}

// MAXFEEBPS is a free data retrieval call binding the contract method 0xd55be8c6.
//
// Solidity: function MAX_FEE_BPS() view returns(uint16)
func (_XmrSwap *XmrSwapCaller) MAXFEEBPS(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _XmrSwap.contract.Call(opts, &out, "MAX_FEE_BPS")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// MAXFEEBPS is a free data retrieval call binding the contract method 0xd55be8c6.
//
// Solidity: function MAX_FEE_BPS() view returns(uint16)
func (_XmrSwap *XmrSwapSession) MAXFEEBPS() (uint16, error) {
	return _XmrSwap.Contract.MAXFEEBPS(&_XmrSwap.CallOpts)
}

// MAXFEEBPS is a free data retrieval call binding the contract method 0xd55be8c6.
//
// Solidity: function MAX_FEE_BPS() view returns(uint16)
func (_XmrSwap *XmrSwapCallerSession) MAXFEEBPS() (uint16, error) {
	return _XmrSwap.Contract.MAXFEEBPS(&_XmrSwap.CallOpts)
}

// MINTIMEOUT is a free data retrieval call binding the contract method 0x543ad1df.
//
// Solidity: function MIN_TIMEOUT() view returns(uint64)
func (_XmrSwap *XmrSwapCaller) MINTIMEOUT(opts *bind.CallOpts) (uint64, error) {
	var out []interface{}
	err := _XmrSwap.contract.Call(opts, &out, "MIN_TIMEOUT")

	if err != nil {
		return *new(uint64), err
	}

	out0 := *abi.ConvertType(out[0], new(uint64)).(*uint64)

	return out0, err

}

// MINTIMEOUT is a free data retrieval call binding the contract method 0x543ad1df.
//
// Solidity: function MIN_TIMEOUT() view returns(uint64)
func (_XmrSwap *XmrSwapSession) MINTIMEOUT() (uint64, error) {
	return _XmrSwap.Contract.MINTIMEOUT(&_XmrSwap.CallOpts)
}

// MINTIMEOUT is a free data retrieval call binding the contract method 0x543ad1df.
//
// Solidity: function MIN_TIMEOUT() view returns(uint64)
func (_XmrSwap *XmrSwapCallerSession) MINTIMEOUT() (uint64, error) {
	return _XmrSwap.Contract.MINTIMEOUT(&_XmrSwap.CallOpts)
}

// FeeBps is a free data retrieval call binding the contract method 0x24a9d853.
//
// Solidity: function feeBps() view returns(uint16)
func (_XmrSwap *XmrSwapCaller) FeeBps(opts *bind.CallOpts) (uint16, error) {
	var out []interface{}
	err := _XmrSwap.contract.Call(opts, &out, "feeBps")

	if err != nil {
		return *new(uint16), err
	}

	out0 := *abi.ConvertType(out[0], new(uint16)).(*uint16)

	return out0, err

}

// FeeBps is a free data retrieval call binding the contract method 0x24a9d853.
//
// Solidity: function feeBps() view returns(uint16)
func (_XmrSwap *XmrSwapSession) FeeBps() (uint16, error) {
	return _XmrSwap.Contract.FeeBps(&_XmrSwap.CallOpts)
}

// FeeBps is a free data retrieval call binding the contract method 0x24a9d853.
//
// Solidity: function feeBps() view returns(uint16)
func (_XmrSwap *XmrSwapCallerSession) FeeBps() (uint16, error) {
	return _XmrSwap.Contract.FeeBps(&_XmrSwap.CallOpts)
}

// FeeRecipient is a free data retrieval call binding the contract method 0x46904840.
//
// Solidity: function feeRecipient() view returns(address)
func (_XmrSwap *XmrSwapCaller) FeeRecipient(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _XmrSwap.contract.Call(opts, &out, "feeRecipient")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FeeRecipient is a free data retrieval call binding the contract method 0x46904840.
//
// Solidity: function feeRecipient() view returns(address)
func (_XmrSwap *XmrSwapSession) FeeRecipient() (common.Address, error) {
	return _XmrSwap.Contract.FeeRecipient(&_XmrSwap.CallOpts)
}

// FeeRecipient is a free data retrieval call binding the contract method 0x46904840.
//
// Solidity: function feeRecipient() view returns(address)
func (_XmrSwap *XmrSwapCallerSession) FeeRecipient() (common.Address, error) {
	return _XmrSwap.Contract.FeeRecipient(&_XmrSwap.CallOpts)
}

// KeyUsed is a free data retrieval call binding the contract method 0xfb0e72b8.
//
// Solidity: function keyUsed(bytes32 spendPub) view returns(bool)
func (_XmrSwap *XmrSwapCaller) KeyUsed(opts *bind.CallOpts, spendPub [32]byte) (bool, error) {
	var out []interface{}
	err := _XmrSwap.contract.Call(opts, &out, "keyUsed", spendPub)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// KeyUsed is a free data retrieval call binding the contract method 0xfb0e72b8.
//
// Solidity: function keyUsed(bytes32 spendPub) view returns(bool)
func (_XmrSwap *XmrSwapSession) KeyUsed(spendPub [32]byte) (bool, error) {
	return _XmrSwap.Contract.KeyUsed(&_XmrSwap.CallOpts, spendPub)
}

// KeyUsed is a free data retrieval call binding the contract method 0xfb0e72b8.
//
// Solidity: function keyUsed(bytes32 spendPub) view returns(bool)
func (_XmrSwap *XmrSwapCallerSession) KeyUsed(spendPub [32]byte) (bool, error) {
	return _XmrSwap.Contract.KeyUsed(&_XmrSwap.CallOpts, spendPub)
}

// Offers is a free data retrieval call binding the contract method 0x474d3ff0.
//
// Solidity: function offers(bytes32 offerId) view returns(address maker, address payout, address asset, uint128 minAmount, uint128 maxAmount, uint256 xmrPerAsset, uint64 expiry, uint64 timeout1Duration, uint64 timeout2Duration, bytes32 makerSpendPub, bytes32 makerViewPriv, bool active)
func (_XmrSwap *XmrSwapCaller) Offers(opts *bind.CallOpts, offerId [32]byte) (struct {
	Maker            common.Address
	Payout           common.Address
	Asset            common.Address
	MinAmount        *big.Int
	MaxAmount        *big.Int
	XmrPerAsset      *big.Int
	Expiry           uint64
	Timeout1Duration uint64
	Timeout2Duration uint64
	MakerSpendPub    [32]byte
	MakerViewPriv    [32]byte
	Active           bool
}, error) {
	var out []interface{}
	err := _XmrSwap.contract.Call(opts, &out, "offers", offerId)

	outstruct := new(struct {
		Maker            common.Address
		Payout           common.Address
		Asset            common.Address
		MinAmount        *big.Int
		MaxAmount        *big.Int
		XmrPerAsset      *big.Int
		Expiry           uint64
		Timeout1Duration uint64
		Timeout2Duration uint64
		MakerSpendPub    [32]byte
		MakerViewPriv    [32]byte
		Active           bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Maker = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Payout = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.Asset = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.MinAmount = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.MaxAmount = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.XmrPerAsset = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.Expiry = *abi.ConvertType(out[6], new(uint64)).(*uint64)
	outstruct.Timeout1Duration = *abi.ConvertType(out[7], new(uint64)).(*uint64)
	outstruct.Timeout2Duration = *abi.ConvertType(out[8], new(uint64)).(*uint64)
	outstruct.MakerSpendPub = *abi.ConvertType(out[9], new([32]byte)).(*[32]byte)
	outstruct.MakerViewPriv = *abi.ConvertType(out[10], new([32]byte)).(*[32]byte)
	outstruct.Active = *abi.ConvertType(out[11], new(bool)).(*bool)

	return *outstruct, err

}

// Offers is a free data retrieval call binding the contract method 0x474d3ff0.
//
// Solidity: function offers(bytes32 offerId) view returns(address maker, address payout, address asset, uint128 minAmount, uint128 maxAmount, uint256 xmrPerAsset, uint64 expiry, uint64 timeout1Duration, uint64 timeout2Duration, bytes32 makerSpendPub, bytes32 makerViewPriv, bool active)
func (_XmrSwap *XmrSwapSession) Offers(offerId [32]byte) (struct {
	Maker            common.Address
	Payout           common.Address
	Asset            common.Address
	MinAmount        *big.Int
	MaxAmount        *big.Int
	XmrPerAsset      *big.Int
	Expiry           uint64
	Timeout1Duration uint64
	Timeout2Duration uint64
	MakerSpendPub    [32]byte
	MakerViewPriv    [32]byte
	Active           bool
}, error) {
	return _XmrSwap.Contract.Offers(&_XmrSwap.CallOpts, offerId)
}

// Offers is a free data retrieval call binding the contract method 0x474d3ff0.
//
// Solidity: function offers(bytes32 offerId) view returns(address maker, address payout, address asset, uint128 minAmount, uint128 maxAmount, uint256 xmrPerAsset, uint64 expiry, uint64 timeout1Duration, uint64 timeout2Duration, bytes32 makerSpendPub, bytes32 makerViewPriv, bool active)
func (_XmrSwap *XmrSwapCallerSession) Offers(offerId [32]byte) (struct {
	Maker            common.Address
	Payout           common.Address
	Asset            common.Address
	MinAmount        *big.Int
	MaxAmount        *big.Int
	XmrPerAsset      *big.Int
	Expiry           uint64
	Timeout1Duration uint64
	Timeout2Duration uint64
	MakerSpendPub    [32]byte
	MakerViewPriv    [32]byte
	Active           bool
}, error) {
	return _XmrSwap.Contract.Offers(&_XmrSwap.CallOpts, offerId)
}

// SecretMatches is a free data retrieval call binding the contract method 0x4977149d.
//
// Solidity: function secretMatches(bytes32 secret, bytes32 pubKey) view returns(bool)
func (_XmrSwap *XmrSwapCaller) SecretMatches(opts *bind.CallOpts, secret [32]byte, pubKey [32]byte) (bool, error) {
	var out []interface{}
	err := _XmrSwap.contract.Call(opts, &out, "secretMatches", secret, pubKey)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SecretMatches is a free data retrieval call binding the contract method 0x4977149d.
//
// Solidity: function secretMatches(bytes32 secret, bytes32 pubKey) view returns(bool)
func (_XmrSwap *XmrSwapSession) SecretMatches(secret [32]byte, pubKey [32]byte) (bool, error) {
	return _XmrSwap.Contract.SecretMatches(&_XmrSwap.CallOpts, secret, pubKey)
}

// SecretMatches is a free data retrieval call binding the contract method 0x4977149d.
//
// Solidity: function secretMatches(bytes32 secret, bytes32 pubKey) view returns(bool)
func (_XmrSwap *XmrSwapCallerSession) SecretMatches(secret [32]byte, pubKey [32]byte) (bool, error) {
	return _XmrSwap.Contract.SecretMatches(&_XmrSwap.CallOpts, secret, pubKey)
}

// Swaps is a free data retrieval call binding the contract method 0xeb84e7f2.
//
// Solidity: function swaps(bytes32 swapId) view returns(address taker, address maker, address payout, address asset, uint256 value, bytes32 makerSpendPub, bytes32 takerSpendPub, uint64 timeout1, uint64 timeout2, uint8 stage)
func (_XmrSwap *XmrSwapCaller) Swaps(opts *bind.CallOpts, swapId [32]byte) (struct {
	Taker         common.Address
	Maker         common.Address
	Payout        common.Address
	Asset         common.Address
	Value         *big.Int
	MakerSpendPub [32]byte
	TakerSpendPub [32]byte
	Timeout1      uint64
	Timeout2      uint64
	Stage         uint8
}, error) {
	var out []interface{}
	err := _XmrSwap.contract.Call(opts, &out, "swaps", swapId)

	outstruct := new(struct {
		Taker         common.Address
		Maker         common.Address
		Payout        common.Address
		Asset         common.Address
		Value         *big.Int
		MakerSpendPub [32]byte
		TakerSpendPub [32]byte
		Timeout1      uint64
		Timeout2      uint64
		Stage         uint8
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Taker = *abi.ConvertType(out[0], new(common.Address)).(*common.Address)
	outstruct.Maker = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.Payout = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.Asset = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)
	outstruct.Value = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.MakerSpendPub = *abi.ConvertType(out[5], new([32]byte)).(*[32]byte)
	outstruct.TakerSpendPub = *abi.ConvertType(out[6], new([32]byte)).(*[32]byte)
	outstruct.Timeout1 = *abi.ConvertType(out[7], new(uint64)).(*uint64)
	outstruct.Timeout2 = *abi.ConvertType(out[8], new(uint64)).(*uint64)
	outstruct.Stage = *abi.ConvertType(out[9], new(uint8)).(*uint8)

	return *outstruct, err

}

// Swaps is a free data retrieval call binding the contract method 0xeb84e7f2.
//
// Solidity: function swaps(bytes32 swapId) view returns(address taker, address maker, address payout, address asset, uint256 value, bytes32 makerSpendPub, bytes32 takerSpendPub, uint64 timeout1, uint64 timeout2, uint8 stage)
func (_XmrSwap *XmrSwapSession) Swaps(swapId [32]byte) (struct {
	Taker         common.Address
	Maker         common.Address
	Payout        common.Address
	Asset         common.Address
	Value         *big.Int
	MakerSpendPub [32]byte
	TakerSpendPub [32]byte
	Timeout1      uint64
	Timeout2      uint64
	Stage         uint8
}, error) {
	return _XmrSwap.Contract.Swaps(&_XmrSwap.CallOpts, swapId)
}

// Swaps is a free data retrieval call binding the contract method 0xeb84e7f2.
//
// Solidity: function swaps(bytes32 swapId) view returns(address taker, address maker, address payout, address asset, uint256 value, bytes32 makerSpendPub, bytes32 takerSpendPub, uint64 timeout1, uint64 timeout2, uint8 stage)
func (_XmrSwap *XmrSwapCallerSession) Swaps(swapId [32]byte) (struct {
	Taker         common.Address
	Maker         common.Address
	Payout        common.Address
	Asset         common.Address
	Value         *big.Int
	MakerSpendPub [32]byte
	TakerSpendPub [32]byte
	Timeout1      uint64
	Timeout2      uint64
	Stage         uint8
}, error) {
	return _XmrSwap.Contract.Swaps(&_XmrSwap.CallOpts, swapId)
}

// CancelOffer is a paid mutator transaction binding the contract method 0xf952279e.
//
// Solidity: function cancelOffer(bytes32 offerId) returns()
func (_XmrSwap *XmrSwapTransactor) CancelOffer(opts *bind.TransactOpts, offerId [32]byte) (*types.Transaction, error) {
	return _XmrSwap.contract.Transact(opts, "cancelOffer", offerId)
}

// CancelOffer is a paid mutator transaction binding the contract method 0xf952279e.
//
// Solidity: function cancelOffer(bytes32 offerId) returns()
func (_XmrSwap *XmrSwapSession) CancelOffer(offerId [32]byte) (*types.Transaction, error) {
	return _XmrSwap.Contract.CancelOffer(&_XmrSwap.TransactOpts, offerId)
}

// CancelOffer is a paid mutator transaction binding the contract method 0xf952279e.
//
// Solidity: function cancelOffer(bytes32 offerId) returns()
func (_XmrSwap *XmrSwapTransactorSession) CancelOffer(offerId [32]byte) (*types.Transaction, error) {
	return _XmrSwap.Contract.CancelOffer(&_XmrSwap.TransactOpts, offerId)
}

// Claim is a paid mutator transaction binding the contract method 0x84cc9dfb.
//
// Solidity: function claim(bytes32 swapId, bytes32 secret) returns()
func (_XmrSwap *XmrSwapTransactor) Claim(opts *bind.TransactOpts, swapId [32]byte, secret [32]byte) (*types.Transaction, error) {
	return _XmrSwap.contract.Transact(opts, "claim", swapId, secret)
}

// Claim is a paid mutator transaction binding the contract method 0x84cc9dfb.
//
// Solidity: function claim(bytes32 swapId, bytes32 secret) returns()
func (_XmrSwap *XmrSwapSession) Claim(swapId [32]byte, secret [32]byte) (*types.Transaction, error) {
	return _XmrSwap.Contract.Claim(&_XmrSwap.TransactOpts, swapId, secret)
}

// Claim is a paid mutator transaction binding the contract method 0x84cc9dfb.
//
// Solidity: function claim(bytes32 swapId, bytes32 secret) returns()
func (_XmrSwap *XmrSwapTransactorSession) Claim(swapId [32]byte, secret [32]byte) (*types.Transaction, error) {
	return _XmrSwap.Contract.Claim(&_XmrSwap.TransactOpts, swapId, secret)
}

// PostOffer is a paid mutator transaction binding the contract method 0x3faff491.
//
// Solidity: function postOffer(address asset, uint128 minAmount, uint128 maxAmount, uint256 xmrPerAsset, uint64 expiry, uint64 timeout1Duration, uint64 timeout2Duration, bytes32 makerSpendPub, bytes32 makerViewPriv, address payout) returns(bytes32 offerId)
func (_XmrSwap *XmrSwapTransactor) PostOffer(opts *bind.TransactOpts, asset common.Address, minAmount *big.Int, maxAmount *big.Int, xmrPerAsset *big.Int, expiry uint64, timeout1Duration uint64, timeout2Duration uint64, makerSpendPub [32]byte, makerViewPriv [32]byte, payout common.Address) (*types.Transaction, error) {
	return _XmrSwap.contract.Transact(opts, "postOffer", asset, minAmount, maxAmount, xmrPerAsset, expiry, timeout1Duration, timeout2Duration, makerSpendPub, makerViewPriv, payout)
}

// PostOffer is a paid mutator transaction binding the contract method 0x3faff491.
//
// Solidity: function postOffer(address asset, uint128 minAmount, uint128 maxAmount, uint256 xmrPerAsset, uint64 expiry, uint64 timeout1Duration, uint64 timeout2Duration, bytes32 makerSpendPub, bytes32 makerViewPriv, address payout) returns(bytes32 offerId)
func (_XmrSwap *XmrSwapSession) PostOffer(asset common.Address, minAmount *big.Int, maxAmount *big.Int, xmrPerAsset *big.Int, expiry uint64, timeout1Duration uint64, timeout2Duration uint64, makerSpendPub [32]byte, makerViewPriv [32]byte, payout common.Address) (*types.Transaction, error) {
	return _XmrSwap.Contract.PostOffer(&_XmrSwap.TransactOpts, asset, minAmount, maxAmount, xmrPerAsset, expiry, timeout1Duration, timeout2Duration, makerSpendPub, makerViewPriv, payout)
}

// PostOffer is a paid mutator transaction binding the contract method 0x3faff491.
//
// Solidity: function postOffer(address asset, uint128 minAmount, uint128 maxAmount, uint256 xmrPerAsset, uint64 expiry, uint64 timeout1Duration, uint64 timeout2Duration, bytes32 makerSpendPub, bytes32 makerViewPriv, address payout) returns(bytes32 offerId)
func (_XmrSwap *XmrSwapTransactorSession) PostOffer(asset common.Address, minAmount *big.Int, maxAmount *big.Int, xmrPerAsset *big.Int, expiry uint64, timeout1Duration uint64, timeout2Duration uint64, makerSpendPub [32]byte, makerViewPriv [32]byte, payout common.Address) (*types.Transaction, error) {
	return _XmrSwap.Contract.PostOffer(&_XmrSwap.TransactOpts, asset, minAmount, maxAmount, xmrPerAsset, expiry, timeout1Duration, timeout2Duration, makerSpendPub, makerViewPriv, payout)
}

// Refund is a paid mutator transaction binding the contract method 0xe4683a79.
//
// Solidity: function refund(bytes32 swapId, bytes32 secret) returns()
func (_XmrSwap *XmrSwapTransactor) Refund(opts *bind.TransactOpts, swapId [32]byte, secret [32]byte) (*types.Transaction, error) {
	return _XmrSwap.contract.Transact(opts, "refund", swapId, secret)
}

// Refund is a paid mutator transaction binding the contract method 0xe4683a79.
//
// Solidity: function refund(bytes32 swapId, bytes32 secret) returns()
func (_XmrSwap *XmrSwapSession) Refund(swapId [32]byte, secret [32]byte) (*types.Transaction, error) {
	return _XmrSwap.Contract.Refund(&_XmrSwap.TransactOpts, swapId, secret)
}

// Refund is a paid mutator transaction binding the contract method 0xe4683a79.
//
// Solidity: function refund(bytes32 swapId, bytes32 secret) returns()
func (_XmrSwap *XmrSwapTransactorSession) Refund(swapId [32]byte, secret [32]byte) (*types.Transaction, error) {
	return _XmrSwap.Contract.Refund(&_XmrSwap.TransactOpts, swapId, secret)
}

// SetReady is a paid mutator transaction binding the contract method 0x6bcfee15.
//
// Solidity: function setReady(bytes32 swapId) returns()
func (_XmrSwap *XmrSwapTransactor) SetReady(opts *bind.TransactOpts, swapId [32]byte) (*types.Transaction, error) {
	return _XmrSwap.contract.Transact(opts, "setReady", swapId)
}

// SetReady is a paid mutator transaction binding the contract method 0x6bcfee15.
//
// Solidity: function setReady(bytes32 swapId) returns()
func (_XmrSwap *XmrSwapSession) SetReady(swapId [32]byte) (*types.Transaction, error) {
	return _XmrSwap.Contract.SetReady(&_XmrSwap.TransactOpts, swapId)
}

// SetReady is a paid mutator transaction binding the contract method 0x6bcfee15.
//
// Solidity: function setReady(bytes32 swapId) returns()
func (_XmrSwap *XmrSwapTransactorSession) SetReady(swapId [32]byte) (*types.Transaction, error) {
	return _XmrSwap.Contract.SetReady(&_XmrSwap.TransactOpts, swapId)
}

// TakeOffer is a paid mutator transaction binding the contract method 0xe4cbfed5.
//
// Solidity: function takeOffer(bytes32 offerId, uint256 amount, bytes32 takerSpendPub, bytes32 takerViewPriv) payable returns(bytes32 swapId)
func (_XmrSwap *XmrSwapTransactor) TakeOffer(opts *bind.TransactOpts, offerId [32]byte, amount *big.Int, takerSpendPub [32]byte, takerViewPriv [32]byte) (*types.Transaction, error) {
	return _XmrSwap.contract.Transact(opts, "takeOffer", offerId, amount, takerSpendPub, takerViewPriv)
}

// TakeOffer is a paid mutator transaction binding the contract method 0xe4cbfed5.
//
// Solidity: function takeOffer(bytes32 offerId, uint256 amount, bytes32 takerSpendPub, bytes32 takerViewPriv) payable returns(bytes32 swapId)
func (_XmrSwap *XmrSwapSession) TakeOffer(offerId [32]byte, amount *big.Int, takerSpendPub [32]byte, takerViewPriv [32]byte) (*types.Transaction, error) {
	return _XmrSwap.Contract.TakeOffer(&_XmrSwap.TransactOpts, offerId, amount, takerSpendPub, takerViewPriv)
}

// TakeOffer is a paid mutator transaction binding the contract method 0xe4cbfed5.
//
// Solidity: function takeOffer(bytes32 offerId, uint256 amount, bytes32 takerSpendPub, bytes32 takerViewPriv) payable returns(bytes32 swapId)
func (_XmrSwap *XmrSwapTransactorSession) TakeOffer(offerId [32]byte, amount *big.Int, takerSpendPub [32]byte, takerViewPriv [32]byte) (*types.Transaction, error) {
	return _XmrSwap.Contract.TakeOffer(&_XmrSwap.TransactOpts, offerId, amount, takerSpendPub, takerViewPriv)
}

// XmrSwapClaimedIterator is returned from FilterClaimed and is used to iterate over the raw logs and unpacked data for Claimed events raised by the XmrSwap contract.
type XmrSwapClaimedIterator struct {
	Event *XmrSwapClaimed // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *XmrSwapClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(XmrSwapClaimed)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(XmrSwapClaimed)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *XmrSwapClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *XmrSwapClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// XmrSwapClaimed represents a Claimed event raised by the XmrSwap contract.
type XmrSwapClaimed struct {
	SwapId [32]byte
	Secret [32]byte
	Payout *big.Int
	Fee    *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterClaimed is a free log retrieval operation binding the contract event 0xd9126119fa9c516709ff9f21d9a56ddf245584c2d2ea9261586dc9ca978af618.
//
// Solidity: event Claimed(bytes32 indexed swapId, bytes32 secret, uint256 payout, uint256 fee)
func (_XmrSwap *XmrSwapFilterer) FilterClaimed(opts *bind.FilterOpts, swapId [][32]byte) (*XmrSwapClaimedIterator, error) {

	var swapIdRule []interface{}
	for _, swapIdItem := range swapId {
		swapIdRule = append(swapIdRule, swapIdItem)
	}

	logs, sub, err := _XmrSwap.contract.FilterLogs(opts, "Claimed", swapIdRule)
	if err != nil {
		return nil, err
	}
	return &XmrSwapClaimedIterator{contract: _XmrSwap.contract, event: "Claimed", logs: logs, sub: sub}, nil
}

// WatchClaimed is a free log subscription operation binding the contract event 0xd9126119fa9c516709ff9f21d9a56ddf245584c2d2ea9261586dc9ca978af618.
//
// Solidity: event Claimed(bytes32 indexed swapId, bytes32 secret, uint256 payout, uint256 fee)
func (_XmrSwap *XmrSwapFilterer) WatchClaimed(opts *bind.WatchOpts, sink chan<- *XmrSwapClaimed, swapId [][32]byte) (event.Subscription, error) {

	var swapIdRule []interface{}
	for _, swapIdItem := range swapId {
		swapIdRule = append(swapIdRule, swapIdItem)
	}

	logs, sub, err := _XmrSwap.contract.WatchLogs(opts, "Claimed", swapIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(XmrSwapClaimed)
				if err := _XmrSwap.contract.UnpackLog(event, "Claimed", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseClaimed is a log parse operation binding the contract event 0xd9126119fa9c516709ff9f21d9a56ddf245584c2d2ea9261586dc9ca978af618.
//
// Solidity: event Claimed(bytes32 indexed swapId, bytes32 secret, uint256 payout, uint256 fee)
func (_XmrSwap *XmrSwapFilterer) ParseClaimed(log types.Log) (*XmrSwapClaimed, error) {
	event := new(XmrSwapClaimed)
	if err := _XmrSwap.contract.UnpackLog(event, "Claimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// XmrSwapOfferCancelledIterator is returned from FilterOfferCancelled and is used to iterate over the raw logs and unpacked data for OfferCancelled events raised by the XmrSwap contract.
type XmrSwapOfferCancelledIterator struct {
	Event *XmrSwapOfferCancelled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *XmrSwapOfferCancelledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(XmrSwapOfferCancelled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(XmrSwapOfferCancelled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *XmrSwapOfferCancelledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *XmrSwapOfferCancelledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// XmrSwapOfferCancelled represents a OfferCancelled event raised by the XmrSwap contract.
type XmrSwapOfferCancelled struct {
	OfferId [32]byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterOfferCancelled is a free log retrieval operation binding the contract event 0x3f9cb69d022b6ec319f86f2df848bcce01f2fc51c9f86396779a8081cf6ca2ea.
//
// Solidity: event OfferCancelled(bytes32 indexed offerId)
func (_XmrSwap *XmrSwapFilterer) FilterOfferCancelled(opts *bind.FilterOpts, offerId [][32]byte) (*XmrSwapOfferCancelledIterator, error) {

	var offerIdRule []interface{}
	for _, offerIdItem := range offerId {
		offerIdRule = append(offerIdRule, offerIdItem)
	}

	logs, sub, err := _XmrSwap.contract.FilterLogs(opts, "OfferCancelled", offerIdRule)
	if err != nil {
		return nil, err
	}
	return &XmrSwapOfferCancelledIterator{contract: _XmrSwap.contract, event: "OfferCancelled", logs: logs, sub: sub}, nil
}

// WatchOfferCancelled is a free log subscription operation binding the contract event 0x3f9cb69d022b6ec319f86f2df848bcce01f2fc51c9f86396779a8081cf6ca2ea.
//
// Solidity: event OfferCancelled(bytes32 indexed offerId)
func (_XmrSwap *XmrSwapFilterer) WatchOfferCancelled(opts *bind.WatchOpts, sink chan<- *XmrSwapOfferCancelled, offerId [][32]byte) (event.Subscription, error) {

	var offerIdRule []interface{}
	for _, offerIdItem := range offerId {
		offerIdRule = append(offerIdRule, offerIdItem)
	}

	logs, sub, err := _XmrSwap.contract.WatchLogs(opts, "OfferCancelled", offerIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(XmrSwapOfferCancelled)
				if err := _XmrSwap.contract.UnpackLog(event, "OfferCancelled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOfferCancelled is a log parse operation binding the contract event 0x3f9cb69d022b6ec319f86f2df848bcce01f2fc51c9f86396779a8081cf6ca2ea.
//
// Solidity: event OfferCancelled(bytes32 indexed offerId)
func (_XmrSwap *XmrSwapFilterer) ParseOfferCancelled(log types.Log) (*XmrSwapOfferCancelled, error) {
	event := new(XmrSwapOfferCancelled)
	if err := _XmrSwap.contract.UnpackLog(event, "OfferCancelled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// XmrSwapOfferPostedIterator is returned from FilterOfferPosted and is used to iterate over the raw logs and unpacked data for OfferPosted events raised by the XmrSwap contract.
type XmrSwapOfferPostedIterator struct {
	Event *XmrSwapOfferPosted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *XmrSwapOfferPostedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(XmrSwapOfferPosted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(XmrSwapOfferPosted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *XmrSwapOfferPostedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *XmrSwapOfferPostedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// XmrSwapOfferPosted represents a OfferPosted event raised by the XmrSwap contract.
type XmrSwapOfferPosted struct {
	OfferId     [32]byte
	Maker       common.Address
	Asset       common.Address
	MinAmount   *big.Int
	MaxAmount   *big.Int
	XmrPerAsset *big.Int
	Expiry      uint64
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterOfferPosted is a free log retrieval operation binding the contract event 0x5d98652ad01d0501cc21957f9e428beced1f29456b42459ff3c1b5a4c00fb844.
//
// Solidity: event OfferPosted(bytes32 indexed offerId, address indexed maker, address indexed asset, uint256 minAmount, uint256 maxAmount, uint256 xmrPerAsset, uint64 expiry)
func (_XmrSwap *XmrSwapFilterer) FilterOfferPosted(opts *bind.FilterOpts, offerId [][32]byte, maker []common.Address, asset []common.Address) (*XmrSwapOfferPostedIterator, error) {

	var offerIdRule []interface{}
	for _, offerIdItem := range offerId {
		offerIdRule = append(offerIdRule, offerIdItem)
	}
	var makerRule []interface{}
	for _, makerItem := range maker {
		makerRule = append(makerRule, makerItem)
	}
	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _XmrSwap.contract.FilterLogs(opts, "OfferPosted", offerIdRule, makerRule, assetRule)
	if err != nil {
		return nil, err
	}
	return &XmrSwapOfferPostedIterator{contract: _XmrSwap.contract, event: "OfferPosted", logs: logs, sub: sub}, nil
}

// WatchOfferPosted is a free log subscription operation binding the contract event 0x5d98652ad01d0501cc21957f9e428beced1f29456b42459ff3c1b5a4c00fb844.
//
// Solidity: event OfferPosted(bytes32 indexed offerId, address indexed maker, address indexed asset, uint256 minAmount, uint256 maxAmount, uint256 xmrPerAsset, uint64 expiry)
func (_XmrSwap *XmrSwapFilterer) WatchOfferPosted(opts *bind.WatchOpts, sink chan<- *XmrSwapOfferPosted, offerId [][32]byte, maker []common.Address, asset []common.Address) (event.Subscription, error) {

	var offerIdRule []interface{}
	for _, offerIdItem := range offerId {
		offerIdRule = append(offerIdRule, offerIdItem)
	}
	var makerRule []interface{}
	for _, makerItem := range maker {
		makerRule = append(makerRule, makerItem)
	}
	var assetRule []interface{}
	for _, assetItem := range asset {
		assetRule = append(assetRule, assetItem)
	}

	logs, sub, err := _XmrSwap.contract.WatchLogs(opts, "OfferPosted", offerIdRule, makerRule, assetRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(XmrSwapOfferPosted)
				if err := _XmrSwap.contract.UnpackLog(event, "OfferPosted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOfferPosted is a log parse operation binding the contract event 0x5d98652ad01d0501cc21957f9e428beced1f29456b42459ff3c1b5a4c00fb844.
//
// Solidity: event OfferPosted(bytes32 indexed offerId, address indexed maker, address indexed asset, uint256 minAmount, uint256 maxAmount, uint256 xmrPerAsset, uint64 expiry)
func (_XmrSwap *XmrSwapFilterer) ParseOfferPosted(log types.Log) (*XmrSwapOfferPosted, error) {
	event := new(XmrSwapOfferPosted)
	if err := _XmrSwap.contract.UnpackLog(event, "OfferPosted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// XmrSwapReadyIterator is returned from FilterReady and is used to iterate over the raw logs and unpacked data for Ready events raised by the XmrSwap contract.
type XmrSwapReadyIterator struct {
	Event *XmrSwapReady // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *XmrSwapReadyIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(XmrSwapReady)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(XmrSwapReady)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *XmrSwapReadyIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *XmrSwapReadyIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// XmrSwapReady represents a Ready event raised by the XmrSwap contract.
type XmrSwapReady struct {
	SwapId [32]byte
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterReady is a free log retrieval operation binding the contract event 0x5fc23b25552757626e08b316cc2387ad1bc70ee1594af7204db4ce0c39f5d15f.
//
// Solidity: event Ready(bytes32 indexed swapId)
func (_XmrSwap *XmrSwapFilterer) FilterReady(opts *bind.FilterOpts, swapId [][32]byte) (*XmrSwapReadyIterator, error) {

	var swapIdRule []interface{}
	for _, swapIdItem := range swapId {
		swapIdRule = append(swapIdRule, swapIdItem)
	}

	logs, sub, err := _XmrSwap.contract.FilterLogs(opts, "Ready", swapIdRule)
	if err != nil {
		return nil, err
	}
	return &XmrSwapReadyIterator{contract: _XmrSwap.contract, event: "Ready", logs: logs, sub: sub}, nil
}

// WatchReady is a free log subscription operation binding the contract event 0x5fc23b25552757626e08b316cc2387ad1bc70ee1594af7204db4ce0c39f5d15f.
//
// Solidity: event Ready(bytes32 indexed swapId)
func (_XmrSwap *XmrSwapFilterer) WatchReady(opts *bind.WatchOpts, sink chan<- *XmrSwapReady, swapId [][32]byte) (event.Subscription, error) {

	var swapIdRule []interface{}
	for _, swapIdItem := range swapId {
		swapIdRule = append(swapIdRule, swapIdItem)
	}

	logs, sub, err := _XmrSwap.contract.WatchLogs(opts, "Ready", swapIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(XmrSwapReady)
				if err := _XmrSwap.contract.UnpackLog(event, "Ready", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseReady is a log parse operation binding the contract event 0x5fc23b25552757626e08b316cc2387ad1bc70ee1594af7204db4ce0c39f5d15f.
//
// Solidity: event Ready(bytes32 indexed swapId)
func (_XmrSwap *XmrSwapFilterer) ParseReady(log types.Log) (*XmrSwapReady, error) {
	event := new(XmrSwapReady)
	if err := _XmrSwap.contract.UnpackLog(event, "Ready", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// XmrSwapRefundedIterator is returned from FilterRefunded and is used to iterate over the raw logs and unpacked data for Refunded events raised by the XmrSwap contract.
type XmrSwapRefundedIterator struct {
	Event *XmrSwapRefunded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *XmrSwapRefundedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(XmrSwapRefunded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(XmrSwapRefunded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *XmrSwapRefundedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *XmrSwapRefundedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// XmrSwapRefunded represents a Refunded event raised by the XmrSwap contract.
type XmrSwapRefunded struct {
	SwapId [32]byte
	Secret [32]byte
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRefunded is a free log retrieval operation binding the contract event 0x007c875846b687732a7579c19bb1dade66cd14e9f4f809565e2b2b5e76c72b4f.
//
// Solidity: event Refunded(bytes32 indexed swapId, bytes32 secret)
func (_XmrSwap *XmrSwapFilterer) FilterRefunded(opts *bind.FilterOpts, swapId [][32]byte) (*XmrSwapRefundedIterator, error) {

	var swapIdRule []interface{}
	for _, swapIdItem := range swapId {
		swapIdRule = append(swapIdRule, swapIdItem)
	}

	logs, sub, err := _XmrSwap.contract.FilterLogs(opts, "Refunded", swapIdRule)
	if err != nil {
		return nil, err
	}
	return &XmrSwapRefundedIterator{contract: _XmrSwap.contract, event: "Refunded", logs: logs, sub: sub}, nil
}

// WatchRefunded is a free log subscription operation binding the contract event 0x007c875846b687732a7579c19bb1dade66cd14e9f4f809565e2b2b5e76c72b4f.
//
// Solidity: event Refunded(bytes32 indexed swapId, bytes32 secret)
func (_XmrSwap *XmrSwapFilterer) WatchRefunded(opts *bind.WatchOpts, sink chan<- *XmrSwapRefunded, swapId [][32]byte) (event.Subscription, error) {

	var swapIdRule []interface{}
	for _, swapIdItem := range swapId {
		swapIdRule = append(swapIdRule, swapIdItem)
	}

	logs, sub, err := _XmrSwap.contract.WatchLogs(opts, "Refunded", swapIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(XmrSwapRefunded)
				if err := _XmrSwap.contract.UnpackLog(event, "Refunded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRefunded is a log parse operation binding the contract event 0x007c875846b687732a7579c19bb1dade66cd14e9f4f809565e2b2b5e76c72b4f.
//
// Solidity: event Refunded(bytes32 indexed swapId, bytes32 secret)
func (_XmrSwap *XmrSwapFilterer) ParseRefunded(log types.Log) (*XmrSwapRefunded, error) {
	event := new(XmrSwapRefunded)
	if err := _XmrSwap.contract.UnpackLog(event, "Refunded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// XmrSwapSwapCreatedIterator is returned from FilterSwapCreated and is used to iterate over the raw logs and unpacked data for SwapCreated events raised by the XmrSwap contract.
type XmrSwapSwapCreatedIterator struct {
	Event *XmrSwapSwapCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *XmrSwapSwapCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(XmrSwapSwapCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(XmrSwapSwapCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *XmrSwapSwapCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *XmrSwapSwapCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// XmrSwapSwapCreated represents a SwapCreated event raised by the XmrSwap contract.
type XmrSwapSwapCreated struct {
	SwapId        [32]byte
	OfferId       [32]byte
	Taker         common.Address
	Value         *big.Int
	XmrPiconero   *big.Int
	TakerSpendPub [32]byte
	TakerViewPriv [32]byte
	Timeout1      uint64
	Timeout2      uint64
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterSwapCreated is a free log retrieval operation binding the contract event 0x6ac60c1e1d7b18793066f728ab948222b6ac5790d9af38983264967ab615155b.
//
// Solidity: event SwapCreated(bytes32 indexed swapId, bytes32 indexed offerId, address indexed taker, uint256 value, uint256 xmrPiconero, bytes32 takerSpendPub, bytes32 takerViewPriv, uint64 timeout1, uint64 timeout2)
func (_XmrSwap *XmrSwapFilterer) FilterSwapCreated(opts *bind.FilterOpts, swapId [][32]byte, offerId [][32]byte, taker []common.Address) (*XmrSwapSwapCreatedIterator, error) {

	var swapIdRule []interface{}
	for _, swapIdItem := range swapId {
		swapIdRule = append(swapIdRule, swapIdItem)
	}
	var offerIdRule []interface{}
	for _, offerIdItem := range offerId {
		offerIdRule = append(offerIdRule, offerIdItem)
	}
	var takerRule []interface{}
	for _, takerItem := range taker {
		takerRule = append(takerRule, takerItem)
	}

	logs, sub, err := _XmrSwap.contract.FilterLogs(opts, "SwapCreated", swapIdRule, offerIdRule, takerRule)
	if err != nil {
		return nil, err
	}
	return &XmrSwapSwapCreatedIterator{contract: _XmrSwap.contract, event: "SwapCreated", logs: logs, sub: sub}, nil
}

// WatchSwapCreated is a free log subscription operation binding the contract event 0x6ac60c1e1d7b18793066f728ab948222b6ac5790d9af38983264967ab615155b.
//
// Solidity: event SwapCreated(bytes32 indexed swapId, bytes32 indexed offerId, address indexed taker, uint256 value, uint256 xmrPiconero, bytes32 takerSpendPub, bytes32 takerViewPriv, uint64 timeout1, uint64 timeout2)
func (_XmrSwap *XmrSwapFilterer) WatchSwapCreated(opts *bind.WatchOpts, sink chan<- *XmrSwapSwapCreated, swapId [][32]byte, offerId [][32]byte, taker []common.Address) (event.Subscription, error) {

	var swapIdRule []interface{}
	for _, swapIdItem := range swapId {
		swapIdRule = append(swapIdRule, swapIdItem)
	}
	var offerIdRule []interface{}
	for _, offerIdItem := range offerId {
		offerIdRule = append(offerIdRule, offerIdItem)
	}
	var takerRule []interface{}
	for _, takerItem := range taker {
		takerRule = append(takerRule, takerItem)
	}

	logs, sub, err := _XmrSwap.contract.WatchLogs(opts, "SwapCreated", swapIdRule, offerIdRule, takerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(XmrSwapSwapCreated)
				if err := _XmrSwap.contract.UnpackLog(event, "SwapCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseSwapCreated is a log parse operation binding the contract event 0x6ac60c1e1d7b18793066f728ab948222b6ac5790d9af38983264967ab615155b.
//
// Solidity: event SwapCreated(bytes32 indexed swapId, bytes32 indexed offerId, address indexed taker, uint256 value, uint256 xmrPiconero, bytes32 takerSpendPub, bytes32 takerViewPriv, uint64 timeout1, uint64 timeout2)
func (_XmrSwap *XmrSwapFilterer) ParseSwapCreated(log types.Log) (*XmrSwapSwapCreated, error) {
	event := new(XmrSwapSwapCreated)
	if err := _XmrSwap.contract.UnpackLog(event, "SwapCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
