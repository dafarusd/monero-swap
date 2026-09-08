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
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"feeBps_\",\"type\":\"uint16\",\"internalType\":\"uint16\"},{\"name\":\"feeRecipient_\",\"type\":\"address\",\"internalType\":\"addresspayable\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"BPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MAX_FEE_BPS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MIN_TIMEOUT\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint64\",\"internalType\":\"uint64\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"cancelOffer\",\"inputs\":[{\"name\":\"offerId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"claim\",\"inputs\":[{\"name\":\"swapId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"secret\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"feeBps\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint16\",\"internalType\":\"uint16\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"feeRecipient\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"addresspayable\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"keyUsed\",\"inputs\":[{\"name\":\"spendPub\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"offers\",\"inputs\":[{\"name\":\"offerId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"maker\",\"type\":\"address\",\"internalType\":\"addresspayable\"},{\"name\":\"asset\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"minAmount\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"maxAmount\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"xmrPerAsset\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expiry\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"timeout1Duration\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"timeout2Duration\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"makerSpendPub\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"makerViewPriv\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"active\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"postOffer\",\"inputs\":[{\"name\":\"asset\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"minAmount\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"maxAmount\",\"type\":\"uint128\",\"internalType\":\"uint128\"},{\"name\":\"xmrPerAsset\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"expiry\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"timeout1Duration\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"timeout2Duration\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"makerSpendPub\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"makerViewPriv\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"offerId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"refund\",\"inputs\":[{\"name\":\"swapId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"secret\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"secretMatches\",\"inputs\":[{\"name\":\"secret\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"pubKey\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"setReady\",\"inputs\":[{\"name\":\"swapId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"swaps\",\"inputs\":[{\"name\":\"swapId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"taker\",\"type\":\"address\",\"internalType\":\"addresspayable\"},{\"name\":\"maker\",\"type\":\"address\",\"internalType\":\"addresspayable\"},{\"name\":\"asset\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"makerSpendPub\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"takerSpendPub\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"timeout1\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"timeout2\",\"type\":\"uint64\",\"internalType\":\"uint64\"},{\"name\":\"stage\",\"type\":\"uint8\",\"internalType\":\"enumXmrSwap.Stage\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"takeOffer\",\"inputs\":[{\"name\":\"offerId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"takerSpendPub\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"},{\"name\":\"takerViewPriv\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"outputs\":[{\"name\":\"swapId\",\"type\":\"bytes32\",\"internalType\":\"bytes32\"}],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"Claimed\",\"inputs\":[{\"name\":\"swapId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"secret\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"payout\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"fee\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OfferCancelled\",\"inputs\":[{\"name\":\"offerId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"OfferPosted\",\"inputs\":[{\"name\":\"offerId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"maker\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"asset\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"minAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"maxAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"xmrPerAsset\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"expiry\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"timeout1Duration\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"timeout2Duration\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"makerSpendPub\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"makerViewPriv\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Ready\",\"inputs\":[{\"name\":\"swapId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Refunded\",\"inputs\":[{\"name\":\"swapId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"secret\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"SwapCreated\",\"inputs\":[{\"name\":\"swapId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"offerId\",\"type\":\"bytes32\",\"indexed\":true,\"internalType\":\"bytes32\"},{\"name\":\"taker\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"value\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"xmrPiconero\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"takerSpendPub\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"takerViewPriv\",\"type\":\"bytes32\",\"indexed\":false,\"internalType\":\"bytes32\"},{\"name\":\"timeout1\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"},{\"name\":\"timeout2\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"error\",\"name\":\"AlreadyReady\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"AmountOutOfRange\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"BadAmounts\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"FeeTooHigh\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSecret\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"InvalidSwap\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"KeyAlreadyUsed\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotMaker\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotTaker\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"NotTimeToRefund\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OfferExpired\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"OfferNotActive\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ReentrancyGuardReentrantCall\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"SafeERC20FailedOperation\",\"inputs\":[{\"name\":\"token\",\"type\":\"address\",\"internalType\":\"address\"}]},{\"type\":\"error\",\"name\":\"SwapCompleted\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TimeoutTooShort\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TooEarlyToClaim\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"TooLateToClaim\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"WrongValue\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroAddress\",\"inputs\":[]},{\"type\":\"error\",\"name\":\"ZeroKey\",\"inputs\":[]}]",
	Bin: "0x60c060405234801561000f575f5ffd5b5060405161223e38038061223e83398101604081905261002e916100ca565b60017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055606461ffff831611156100785760405163cd4e616760e01b815260040160405180910390fd5b61ffff82161580159061009257506001600160a01b038116155b156100b05760405163d92e233d60e01b815260040160405180910390fd5b61ffff9091166080526001600160a01b031660a052610113565b5f5f604083850312156100db575f5ffd5b825161ffff811681146100ec575f5ffd5b60208401519092506001600160a01b0381168114610108575f5ffd5b809150509250929050565b60805160a0516120fc6101425f395f8181610159015261083c01525f818161012601526107bb01526120fc5ff3fe6080604052600436106100e4575f3560e01c806384cc9dfb11610087578063e4cbfed511610057578063e4cbfed514610347578063eb84e7f21461035a578063f952279e146103fd578063fb0e72b81461041c575f5ffd5b806384cc9dfb146102c8578063a6f23215146102e7578063d55be8c614610314578063e4683a7914610328575f5ffd5b8063474d3ff0116100c2578063474d3ff0146101935780634977149d1461024a578063543ad1df146102795780636bcfee15146102a7575f5ffd5b8063249d39e9146100e857806324a9d853146101155780634690484014610148575b5f5ffd5b3480156100f3575f5ffd5b506100fd61271081565b60405161ffff90911681526020015b60405180910390f35b348015610120575f5ffd5b506100fd7f000000000000000000000000000000000000000000000000000000000000000081565b348015610153575f5ffd5b5061017b7f000000000000000000000000000000000000000000000000000000000000000081565b6040516001600160a01b03909116815260200161010c565b34801561019e575f5ffd5b506102336101ad366004611dba565b5f60208190529081526040902080546001820154600283015460038401546004850154600586015460068701546007909701546001600160a01b039687169796909516956001600160801b0380861696600160801b968790049091169567ffffffffffffffff808616956801000000000000000081048216959290041692909160ff168b565b60405161010c9b9a99989796959493929190611dd1565b348015610255575f5ffd5b50610269610264366004611e55565b61044a565b604051901515815260200161010c565b348015610284575f5ffd5b5061028e610e1081565b60405167ffffffffffffffff909116815260200161010c565b3480156102b2575f5ffd5b506102c66102c1366004611dba565b61045e565b005b3480156102d3575f5ffd5b506102c66102e2366004611e55565b6105eb565b3480156102f2575f5ffd5b50610306610301366004611ea7565b6108d4565b60405190815260200161010c565b34801561031f575f5ffd5b506100fd606481565b348015610333575f5ffd5b506102c6610342366004611e55565b610d6e565b610306610355366004611f44565b610f96565b348015610365575f5ffd5b506103e8610374366004611dba565b600160208190525f918252604090912080549181015460028201546003830154600484015460058501546006909501546001600160a01b03968716969485169593909416939192909167ffffffffffffffff808216916801000000000000000081049091169060ff600160801b9091041689565b60405161010c99989796959493929190611f87565b348015610408575f5ffd5b506102c6610417366004611dba565b611477565b348015610427575f5ffd5b50610269610436366004611dba565b60026020525f908152604090205460ff1681565b5f610455838361150d565b90505b92915050565b5f818152600160205260408120906006820154600160801b900460ff16600381111561048c5761048c611f73565b036104aa57604051631115766760e01b815260040160405180910390fd5b60036006820154600160801b900460ff1660038111156104cc576104cc611f73565b036104ea5760405163066916a960e01b815260040160405180910390fd5b60026006820154600160801b900460ff16600381111561050c5761050c611f73565b03610543576040517f2000ee6c00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b80546001600160a01b0316331461056d57604051635211a07960e01b815260040160405180910390fd5b600681015467ffffffffffffffff16421061059b5760405163497df9d160e01b815260040160405180910390fd5b60068101805460ff60801b191670020000000000000000000000000000000017905560405182907f5fc23b25552757626e08b316cc2387ad1bc70ee1594af7204db4ce0c39f5d15f905f90a25050565b6105f361154f565b5f828152600160205260408120906006820154600160801b900460ff16600381111561062157610621611f73565b0361063f57604051631115766760e01b815260040160405180910390fd5b60036006820154600160801b900460ff16600381111561066157610661611f73565b0361067f5760405163066916a960e01b815260040160405180910390fd5b60018101546001600160a01b031633146106ac5760405163b331e42160e01b815260040160405180910390fd5b600681015467ffffffffffffffff16421080156106e9575060026006820154600160801b900460ff1660038111156106e6576106e6611f73565b14155b15610720576040517fd71d60b500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b600681015468010000000000000000900467ffffffffffffffff16421061075a5760405163497df9d160e01b815260040160405180910390fd5b61076882826004015461150d565b6107855760405163abab6bd760e01b815260040160405180910390fd5b60068101805460ff60801b191670030000000000000000000000000000000017905560038101545f90612710906107e19061ffff7f00000000000000000000000000000000000000000000000000000000000000001690612030565b6107eb9190612047565b90505f8183600301546107fe9190612066565b60028401546001850154919250610822916001600160a01b0391821691168361157d565b8115610861576002830154610861906001600160a01b03167f00000000000000000000000000000000000000000000000000000000000000008461157d565b604080518581526020810183905290810183905285907fd9126119fa9c516709ff9f21d9a56ddf245584c2d2ea9261586dc9ca978af6189060600160405180910390a25050506108d060017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055565b5050565b5f8215806108e0575081155b156108fe57604051637dec24e360e11b815260040160405180910390fd5b6001600160801b03891615806109255750886001600160801b0316886001600160801b0316105b8061092e575086155b15610965576040517fc6fd446400000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610e1067ffffffffffffffff8616108061098a5750610e1067ffffffffffffffff8516105b156109c1576040517fe143262500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b428667ffffffffffffffff16116109eb57604051639cb1308760e01b815260040160405180910390fd5b5f8381526002602052604090205460ff1615610a1a576040516312b3b8ff60e01b815260040160405180910390fd5b5f838152600260205260408120805460ff191660011790556003805446923092339291610a4683612079565b909155506040805160208101959095526bffffffffffffffffffffffff19606094851b8116918601919091529190921b1660548301526068820152608801604051602081830303815290604052805190602001209050604051806101600160405280336001600160a01b031681526020018b6001600160a01b031681526020018a6001600160801b03168152602001896001600160801b031681526020018881526020018767ffffffffffffffff1681526020018667ffffffffffffffff1681526020018567ffffffffffffffff168152602001848152602001838152602001600115158152505f5f8381526020019081526020015f205f820151815f015f6101000a8154816001600160a01b0302191690836001600160a01b031602179055506020820151816001015f6101000a8154816001600160a01b0302191690836001600160a01b031602179055506040820151816002015f6101000a8154816001600160801b0302191690836001600160801b0316021790555060608201518160020160106101000a8154816001600160801b0302191690836001600160801b031602179055506080820151816003015560a0820151816004015f6101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555060c08201518160040160086101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555060e08201518160040160106101000a81548167ffffffffffffffff021916908367ffffffffffffffff16021790555061010082015181600501556101208201518160060155610140820151816007015f6101000a81548160ff021916908315150217905550905050896001600160a01b0316336001600160a01b0316827f860b380beec098396bb1504e1f061d000184df09ed95118654bd1eb30f394d3f8c8c8c8c8c8c8c8c604051610d599897969594939291906001600160801b039889168152969097166020870152604086019490945267ffffffffffffffff928316606086015290821660808501521660a083015260c082015260e08101919091526101000190565b60405180910390a49998505050505050505050565b610d7661154f565b5f828152600160205260408120906006820154600160801b900460ff166003811115610da457610da4611f73565b03610dc257604051631115766760e01b815260040160405180910390fd5b60036006820154600160801b900460ff166003811115610de457610de4611f73565b03610e025760405163066916a960e01b815260040160405180910390fd5b80546001600160a01b03163314610e2c57604051635211a07960e01b815260040160405180910390fd5b600681015468010000000000000000900467ffffffffffffffff1642108015610e8c5750600681015467ffffffffffffffff1642101580610e8c575060026006820154600160801b900460ff166003811115610e8a57610e8a611f73565b145b15610ec3576040517f65430c1e00000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b610ed182826005015461150d565b610eee5760405163abab6bd760e01b815260040160405180910390fd5b60068101805460ff60801b1916700300000000000000000000000000000000179055600281015481546003830154610f33926001600160a01b0390811692169061157d565b827e7c875846b687732a7579c19bb1dade66cd14e9f4f809565e2b2b5e76c72b4f83604051610f6491815260200190565b60405180910390a2506108d060017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055565b5f610f9f61154f565b5f858152602081905260409020600781015460ff16610fd157604051634d80826b60e11b815260040160405180910390fd5b600481015467ffffffffffffffff164210610fff57604051639cb1308760e01b815260040160405180910390fd5b83158061100a575082155b1561102857604051637dec24e360e11b815260040160405180910390fd5b60028101546001600160801b031685108061105657506002810154600160801b90046001600160801b031685115b1561108d576040517fc64200e900000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b5f8481526002602052604090205460ff16156110bc576040516312b3b8ff60e01b815260040160405180910390fd5b5f848152600260205260409020805460ff191660019081179091558101546001600160a01b031661110c5784341461110757604051632635240760e21b815260040160405180910390fd5b611145565b341561112b57604051632635240760e21b815260040160405180910390fd5b6001810154611145906001600160a01b0316333088611668565b60078101805460ff1916905560048101545f906111789068010000000000000000900467ffffffffffffffff1642612091565b60048301549091505f9061119d90600160801b900467ffffffffffffffff1683612091565b60408051602081018b90526bffffffffffffffffffffffff193360601b16918101919091526054810188905290915060740160408051601f19818403018152919052805160209091012093505f5f85815260016020526040902060060154600160801b900460ff16600381111561121657611216611f73565b1461123457604051631115766760e01b815260040160405180910390fd5b604080516101208101825233815284546001600160a01b03908116602083015260018087015490911692820192909252606081018990526005850154608082015260a0810188905267ffffffffffffffff80851660c0830152831660e0820152906101008201525f85815260016020818152604092839020845181546001600160a01b039182167fffffffffffffffffffffffff0000000000000000000000000000000000000000918216178355928601519382018054948216948416949094179093559284015160028401805491909316911617905560608201516003808301919091556080830151600483015560a0830151600583015560c083015160068301805460e086015167ffffffffffffffff90811668010000000000000000027fffffffffffffffffffffffffffffffff00000000000000000000000000000000909216931692909217919091178082556101008501519260ff60801b1990911690600160801b9084908111156113ad576113ad611f73565b02179055509050505f670de0b6b3a76400008460030154896113cf9190612030565b6113d99190612047565b604080518a8152602081018390529081018990526060810188905267ffffffffffffffff8086166080830152841660a082015290915033908a9087907f6ac60c1e1d7b18793066f728ab948222b6ac5790d9af38983264967ab615155b9060c00160405180910390a45050505061146f60017f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055565b949350505050565b5f818152602081905260409020600781015460ff166114a957604051634d80826b60e11b815260040160405180910390fd5b80546001600160a01b031633146114d35760405163b331e42160e01b815260040160405180910390fd5b60078101805460ff1916905560405182907f3f9cb69d022b6ec319f86f2df848bcce01f2fc51c9f86396779a8081cf6ca2ea905f90a25050565b5f806115188461169e565b9050805f0361152a575f915050610458565b5f5f611535836116d3565b915091508461154483836117f4565b149695505050505050565b611557611827565b60027f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0055565b6001600160a01b03831661164f575f826001600160a01b0316826040515f6040518083038185875af1925050503d805f81146115d4576040519150601f19603f3d011682016040523d82523d5f602084013e6115d9565b606091505b5050905080611649576040517f08c379a000000000000000000000000000000000000000000000000000000000815260206004820152601360248201527f455448207472616e73666572206661696c65640000000000000000000000000060448201526064015b60405180910390fd5b50505050565b6116636001600160a01b0384168383611884565b505050565b6116768484848460016118b9565b61164957604051635274afe760e01b81526001600160a01b0385166004820152602401611640565b5f8060205b80156116cc57600884811c9460ff1692901b9190911790806116c4816120b1565b9150506116a3565b5092915050565b5f5f6116f660405180606001604052805f81526020015f81526020015f81525090565b61171760405180606001604052805f81526020015f81526020015f81525090565b7f216936d3cd6e53fec0a4e231fdd6dc5c692cc7609525a7b2c9562d608f25d51a82527f6666666666666666666666666666666666666666666666666666666666666658602080840191909152600160408085018290525f8452918301819052908201525b84156117b1578460011660010361179a57611797818361193f565b90505b600185901c94506117aa82611b20565b915061177c565b5f6117bf8260400151611cdc565b90506013600160ff1b03825182900982526013600160ff1b038183602001510960208301819052915196919550909350505050565b5f7f800000000000000000000000000000000000000000000000000000000000000060ff84901b16821761146f8161169e565b7f9b779b17422d0df92223018b32b4d1fa46e071723d6817e2486d003becc55f0054600203611882576040517f3ee5aeb500000000000000000000000000000000000000000000000000000000815260040160405180910390fd5b565b6118918383836001611d3f565b61166357604051635274afe760e01b81526001600160a01b0384166004820152602401611640565b6040517f23b872dd000000000000000000000000000000000000000000000000000000005f8181526001600160a01b038781166004528616602452604485905291602083606481808c5af1925060015f5114831661192e578383151615611922573d5f823e3d81fd5b5f883b113d1516831692505b604052505f60605295945050505050565b61196060405180606001604052805f81526020015f81526020015f81525090565b6119a06040518061010001604052805f81526020015f81526020015f81526020015f81526020015f81526020015f81526020015f81526020015f81525090565b6013600160ff1b03836040015185604001510981526013600160ff1b038151800960208201526013600160ff1b03835185510960408201526013600160ff1b03836020015185602001510960608201526013600160ff1b038082606001518360400151097f52036cee2b6ffe738cc740797779e89800700a4d4141d8ab75eb4dca135978a309608082018190526013600160ff1b0390611a409082612066565b82602001510860a08201526013600160ff1b03816080015182602001510860c08201526013600160ff1b03806060830151611a82906013600160ff1b03612066565b6013600160ff1b036040850151611aa0906013600160ff1b03612066565b6013600160ff1b038060208a01518a51086013600160ff1b0360208c01518c51080908086013600160ff1b0360a08401518451090982526013600160ff1b038082604001518360600151086013600160ff1b0360c08401518451090960208301526013600160ff1b038160c001518260a001510960408301525092915050565b611b4160405180606001604052805f81526020015f81526020015f81525090565b611b816040518061010001604052805f81526020015f81526020015f81526020015f81526020015f81526020015f81526020015f81526020015f81525090565b6013600160ff1b03602084015184510881526013600160ff1b038151800960208201526013600160ff1b038351800960408201526013600160ff1b036020840151800960608201526040810151611bdf906013600160ff1b03612066565b6080820181905260608201516013600160ff1b03910860a08201526013600160ff1b036040840151800960e08201526013600160ff1b03808260e00151600209611c30906013600160ff1b03612066565b8260a001510860c08201526013600160ff1b0360c08201516013600160ff1b036060840151611c66906013600160ff1b03612066565b6013600160ff1b036040860151611c84906013600160ff1b03612066565b866020015108080982526013600160ff1b03806060830151611cad906013600160ff1b03612066565b8360800151088260a001510960208301526013600160ff1b038160c001518260a0015109604083015250919050565b5f80611cf060026013600160ff1b03612066565b90505f6013600160ff1b03905060405160208152602080820152602060408201528460608201528260808201528160a082015260208160c08360055f19fa611d36575f5ffd5b51949350505050565b6040517fa9059cbb000000000000000000000000000000000000000000000000000000005f8181526001600160a01b038616600452602485905291602083604481808b5af1925060015f51148316611dae578383151615611da2573d5f823e3d81fd5b5f873b113d1516831692505b60405250949350505050565b5f60208284031215611dca575f5ffd5b5035919050565b6001600160a01b038c811682528b1660208201526001600160801b038a81166040830152891660608201526080810188905267ffffffffffffffff87811660a083015286811660c0830152851660e082015261016081018461010083015283610120830152611e4561014083018415159052565b9c9b505050505050505050505050565b5f5f60408385031215611e66575f5ffd5b50508035926020909101359150565b80356001600160801b0381168114611e8b575f5ffd5b919050565b803567ffffffffffffffff81168114611e8b575f5ffd5b5f5f5f5f5f5f5f5f5f6101208a8c031215611ec0575f5ffd5b89356001600160a01b0381168114611ed6575f5ffd5b9850611ee460208b01611e75565b9750611ef260408b01611e75565b965060608a01359550611f0760808b01611e90565b9450611f1560a08b01611e90565b9350611f2360c08b01611e90565b989b979a50959894979396929550929360e081013593506101000135919050565b5f5f5f5f60808587031215611f57575f5ffd5b5050823594602084013594506040840135936060013592509050565b634e487b7160e01b5f52602160045260245ffd5b5f610120820190506001600160a01b038b1682526001600160a01b038a1660208301526001600160a01b03891660408301528760608301528660808301528560a083015267ffffffffffffffff851660c083015267ffffffffffffffff841660e08301526004831061200757634e487b7160e01b5f52602160045260245ffd5b826101008301529a9950505050505050505050565b634e487b7160e01b5f52601160045260245ffd5b80820281158282048414176104585761045861201c565b5f8261206157634e487b7160e01b5f52601260045260245ffd5b500490565b818103818111156104585761045861201c565b5f6001820161208a5761208a61201c565b5060010190565b67ffffffffffffffff81811683821601908111156104585761045861201c565b5f816120bf576120bf61201c565b505f19019056fea2646970667358221220936eece96f381e34f6ec0a43e4361f06bb93348ff601fd523baa4ae4d9d70cc964736f6c634300081c0033",
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
// Solidity: function offers(bytes32 offerId) view returns(address maker, address asset, uint128 minAmount, uint128 maxAmount, uint256 xmrPerAsset, uint64 expiry, uint64 timeout1Duration, uint64 timeout2Duration, bytes32 makerSpendPub, bytes32 makerViewPriv, bool active)
func (_XmrSwap *XmrSwapCaller) Offers(opts *bind.CallOpts, offerId [32]byte) (struct {
	Maker            common.Address
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
	outstruct.Asset = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.MinAmount = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.MaxAmount = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.XmrPerAsset = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.Expiry = *abi.ConvertType(out[5], new(uint64)).(*uint64)
	outstruct.Timeout1Duration = *abi.ConvertType(out[6], new(uint64)).(*uint64)
	outstruct.Timeout2Duration = *abi.ConvertType(out[7], new(uint64)).(*uint64)
	outstruct.MakerSpendPub = *abi.ConvertType(out[8], new([32]byte)).(*[32]byte)
	outstruct.MakerViewPriv = *abi.ConvertType(out[9], new([32]byte)).(*[32]byte)
	outstruct.Active = *abi.ConvertType(out[10], new(bool)).(*bool)

	return *outstruct, err

}

// Offers is a free data retrieval call binding the contract method 0x474d3ff0.
//
// Solidity: function offers(bytes32 offerId) view returns(address maker, address asset, uint128 minAmount, uint128 maxAmount, uint256 xmrPerAsset, uint64 expiry, uint64 timeout1Duration, uint64 timeout2Duration, bytes32 makerSpendPub, bytes32 makerViewPriv, bool active)
func (_XmrSwap *XmrSwapSession) Offers(offerId [32]byte) (struct {
	Maker            common.Address
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
// Solidity: function offers(bytes32 offerId) view returns(address maker, address asset, uint128 minAmount, uint128 maxAmount, uint256 xmrPerAsset, uint64 expiry, uint64 timeout1Duration, uint64 timeout2Duration, bytes32 makerSpendPub, bytes32 makerViewPriv, bool active)
func (_XmrSwap *XmrSwapCallerSession) Offers(offerId [32]byte) (struct {
	Maker            common.Address
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
// Solidity: function swaps(bytes32 swapId) view returns(address taker, address maker, address asset, uint256 value, bytes32 makerSpendPub, bytes32 takerSpendPub, uint64 timeout1, uint64 timeout2, uint8 stage)
func (_XmrSwap *XmrSwapCaller) Swaps(opts *bind.CallOpts, swapId [32]byte) (struct {
	Taker         common.Address
	Maker         common.Address
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
	outstruct.Asset = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.Value = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.MakerSpendPub = *abi.ConvertType(out[4], new([32]byte)).(*[32]byte)
	outstruct.TakerSpendPub = *abi.ConvertType(out[5], new([32]byte)).(*[32]byte)
	outstruct.Timeout1 = *abi.ConvertType(out[6], new(uint64)).(*uint64)
	outstruct.Timeout2 = *abi.ConvertType(out[7], new(uint64)).(*uint64)
	outstruct.Stage = *abi.ConvertType(out[8], new(uint8)).(*uint8)

	return *outstruct, err

}

// Swaps is a free data retrieval call binding the contract method 0xeb84e7f2.
//
// Solidity: function swaps(bytes32 swapId) view returns(address taker, address maker, address asset, uint256 value, bytes32 makerSpendPub, bytes32 takerSpendPub, uint64 timeout1, uint64 timeout2, uint8 stage)
func (_XmrSwap *XmrSwapSession) Swaps(swapId [32]byte) (struct {
	Taker         common.Address
	Maker         common.Address
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
// Solidity: function swaps(bytes32 swapId) view returns(address taker, address maker, address asset, uint256 value, bytes32 makerSpendPub, bytes32 takerSpendPub, uint64 timeout1, uint64 timeout2, uint8 stage)
func (_XmrSwap *XmrSwapCallerSession) Swaps(swapId [32]byte) (struct {
	Taker         common.Address
	Maker         common.Address
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

// PostOffer is a paid mutator transaction binding the contract method 0xa6f23215.
//
// Solidity: function postOffer(address asset, uint128 minAmount, uint128 maxAmount, uint256 xmrPerAsset, uint64 expiry, uint64 timeout1Duration, uint64 timeout2Duration, bytes32 makerSpendPub, bytes32 makerViewPriv) returns(bytes32 offerId)
func (_XmrSwap *XmrSwapTransactor) PostOffer(opts *bind.TransactOpts, asset common.Address, minAmount *big.Int, maxAmount *big.Int, xmrPerAsset *big.Int, expiry uint64, timeout1Duration uint64, timeout2Duration uint64, makerSpendPub [32]byte, makerViewPriv [32]byte) (*types.Transaction, error) {
	return _XmrSwap.contract.Transact(opts, "postOffer", asset, minAmount, maxAmount, xmrPerAsset, expiry, timeout1Duration, timeout2Duration, makerSpendPub, makerViewPriv)
}

// PostOffer is a paid mutator transaction binding the contract method 0xa6f23215.
//
// Solidity: function postOffer(address asset, uint128 minAmount, uint128 maxAmount, uint256 xmrPerAsset, uint64 expiry, uint64 timeout1Duration, uint64 timeout2Duration, bytes32 makerSpendPub, bytes32 makerViewPriv) returns(bytes32 offerId)
func (_XmrSwap *XmrSwapSession) PostOffer(asset common.Address, minAmount *big.Int, maxAmount *big.Int, xmrPerAsset *big.Int, expiry uint64, timeout1Duration uint64, timeout2Duration uint64, makerSpendPub [32]byte, makerViewPriv [32]byte) (*types.Transaction, error) {
	return _XmrSwap.Contract.PostOffer(&_XmrSwap.TransactOpts, asset, minAmount, maxAmount, xmrPerAsset, expiry, timeout1Duration, timeout2Duration, makerSpendPub, makerViewPriv)
}

// PostOffer is a paid mutator transaction binding the contract method 0xa6f23215.
//
// Solidity: function postOffer(address asset, uint128 minAmount, uint128 maxAmount, uint256 xmrPerAsset, uint64 expiry, uint64 timeout1Duration, uint64 timeout2Duration, bytes32 makerSpendPub, bytes32 makerViewPriv) returns(bytes32 offerId)
func (_XmrSwap *XmrSwapTransactorSession) PostOffer(asset common.Address, minAmount *big.Int, maxAmount *big.Int, xmrPerAsset *big.Int, expiry uint64, timeout1Duration uint64, timeout2Duration uint64, makerSpendPub [32]byte, makerViewPriv [32]byte) (*types.Transaction, error) {
	return _XmrSwap.Contract.PostOffer(&_XmrSwap.TransactOpts, asset, minAmount, maxAmount, xmrPerAsset, expiry, timeout1Duration, timeout2Duration, makerSpendPub, makerViewPriv)
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
	OfferId          [32]byte
	Maker            common.Address
	Asset            common.Address
	MinAmount        *big.Int
	MaxAmount        *big.Int
	XmrPerAsset      *big.Int
	Expiry           uint64
	Timeout1Duration uint64
	Timeout2Duration uint64
	MakerSpendPub    [32]byte
	MakerViewPriv    [32]byte
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterOfferPosted is a free log retrieval operation binding the contract event 0x860b380beec098396bb1504e1f061d000184df09ed95118654bd1eb30f394d3f.
//
// Solidity: event OfferPosted(bytes32 indexed offerId, address indexed maker, address indexed asset, uint256 minAmount, uint256 maxAmount, uint256 xmrPerAsset, uint64 expiry, uint64 timeout1Duration, uint64 timeout2Duration, bytes32 makerSpendPub, bytes32 makerViewPriv)
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

// WatchOfferPosted is a free log subscription operation binding the contract event 0x860b380beec098396bb1504e1f061d000184df09ed95118654bd1eb30f394d3f.
//
// Solidity: event OfferPosted(bytes32 indexed offerId, address indexed maker, address indexed asset, uint256 minAmount, uint256 maxAmount, uint256 xmrPerAsset, uint64 expiry, uint64 timeout1Duration, uint64 timeout2Duration, bytes32 makerSpendPub, bytes32 makerViewPriv)
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

// ParseOfferPosted is a log parse operation binding the contract event 0x860b380beec098396bb1504e1f061d000184df09ed95118654bd1eb30f394d3f.
//
// Solidity: event OfferPosted(bytes32 indexed offerId, address indexed maker, address indexed asset, uint256 minAmount, uint256 maxAmount, uint256 xmrPerAsset, uint64 expiry, uint64 timeout1Duration, uint64 timeout2Duration, bytes32 makerSpendPub, bytes32 makerViewPriv)
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
