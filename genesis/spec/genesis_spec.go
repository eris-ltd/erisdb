package spec

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	acm "github.com/hyperledger/burrow/account"
	"github.com/hyperledger/burrow/genesis"
	"github.com/hyperledger/burrow/keys"
	"github.com/hyperledger/burrow/permission"
	ptypes "github.com/hyperledger/burrow/permission/types"
)

const DefaultAmount uint64 = 1000000
const DefaultAmountBonded uint64 = 10000

// A GenesisSpec is schematic representation of a genesis state, that is it is a template
// for a GenesisDoc excluding that which needs to be instantiated at the point of genesis
// so it describes the type and number of accounts, the genesis salt, but not the
// account keys or addresses, or the GenesisTime. It is responsible for generating keys
// by interacting with the KeysClient it is passed and other information not known at
// specification time
type GenesisSpec struct {
	GenesisTime       *time.Time        `json:",omitempty" toml:",omitempty"`
	ChainName         string            `json:",omitempty" toml:",omitempty"`
	Salt              []byte            `json:",omitempty" toml:",omitempty"`
	GlobalPermissions []string          `json:",omitempty" toml:",omitempty"`
	Accounts          []TemplateAccount `json:",omitempty" toml:",omitempty"`
}

type TemplateAccount struct {
	// Template accounts sharing a name will be merged when merging genesis specs
	Name string `json:",omitempty" toml:",omitempty"`
	// Address  is convenient to have in file for reference, but otherwise ignored since derived from PublicKey
	Address   *acm.Address   `json:",omitempty" toml:",omitempty"`
	PublicKey *acm.PublicKey `json:",omitempty" toml:",omitempty"`
	Amount    *uint64        `json:",omitempty" toml:",omitempty"`
	// If any bonded amount then this account is also a Validator
	AmountBonded *uint64  `json:",omitempty" toml:",omitempty"`
	Permissions  []string `json:",omitempty" toml:",omitempty"`
	Roles        []string `json:",omitempty" toml:",omitempty"`
}

func (ta TemplateAccount) Validator(keyClient keys.KeyClient, index int) (acm.Validator, error) {

	publicKey, _, err := ta.RealisePubKeyAndAddress(keyClient)
	if err != nil {
		return nil, err
	}
	validator := acm.NewValidator(publicKey, 100, 1)
	return validator, nil
}

func (ta TemplateAccount) AccountPermissions() (ptypes.AccountPermissions, error) {
	basePerms, err := permission.BasePermissionsFromStringList(ta.Permissions)
	if err != nil {
		return permission.ZeroAccountPermissions, nil
	}
	return ptypes.AccountPermissions{
		Base:  basePerms,
		Roles: ta.Roles,
	}, nil
}

func (ta TemplateAccount) Account(keyClient keys.KeyClient, index int) (acm.Account, error) {
	publicKey, _, err := ta.RealisePubKeyAndAddress(keyClient)
	if err != nil {
		return nil, err
	}
	account := acm.NewConcreteAccount(publicKey)

	if ta.Amount == nil {
		account.Balance = DefaultAmount
	} else {
		account.Balance = *ta.Amount
	}

	if ta.Permissions == nil {
		account.Permissions = permission.DefaultAccountPermissions.Clone()
	} else {
		account.Permissions, err = ta.AccountPermissions()
		if err != nil {
			return nil, err
		}
	}
	return account.Account(), nil
}

// Adds a public key and address to the template. If PublicKey will try to fetch it by Address.
// If both PublicKey and Address are not set will use the keyClient to generate a new keypair
func (ta TemplateAccount) RealisePubKeyAndAddress(keyClient keys.KeyClient) (pubKey acm.PublicKey, address acm.Address, err error) {
	if ta.PublicKey == nil {
		if ta.Address == nil {
			// If neither PublicKey or Address set then generate a new one
			address, err = keyClient.Generate(ta.Name, keys.KeyTypeEd25519Ripemd160)
			if err != nil {
				return
			}
		} else {
			address = *ta.Address
		}
		// Get the (possibly existing) key
		pubKey, err = keyClient.PublicKey(address)
		if err != nil {
			return
		}
	} else {
		address = ta.PublicKey.Address()
		if ta.Address != nil && *ta.Address != address {
			err = fmt.Errorf("template address %s does not match public key derived address %s", ta.Address,
				ta.PublicKey)
		}
		pubKey = *ta.PublicKey
	}
	return
}

func (gs *GenesisSpec) RealiseKeys(keyClient keys.KeyClient) error {
	for _, templateAccount := range gs.Accounts {
		_, _, err := templateAccount.RealisePubKeyAndAddress(keyClient)
		if err != nil {
			return err
		}
	}
	return nil
}

// Produce a fully realised GenesisDoc from a template GenesisDoc that may omit values
func (gs *GenesisSpec) GenesisDoc(keyClient keys.KeyClient) (*genesis.GenesisDoc, error) {
	var genesisTime time.Time
	if gs.GenesisTime == nil {
		genesisTime = time.Now()
	} else {
		genesisTime = *gs.GenesisTime
	}

	var chainName string
	if gs.ChainName == "" {
		chainName = fmt.Sprintf("BurrowChain_%X", gs.ShortHash())
	} else {
		chainName = gs.ChainName
	}

	var globalPermissions ptypes.AccountPermissions
	if len(gs.GlobalPermissions) == 0 {
		globalPermissions = permission.DefaultAccountPermissions.Clone()
	} else {
		basePerms, err := permission.BasePermissionsFromStringList(gs.GlobalPermissions)
		if err != nil {
			return nil, err
		}
		globalPermissions = ptypes.AccountPermissions{
			Base: basePerms,
		}
	}

	templateAccounts := gs.Accounts
	if len(gs.Accounts) == 0 {
		amountBonded := DefaultAmountBonded
		templateAccounts = append(templateAccounts, TemplateAccount{
			AmountBonded: &amountBonded,
		})
	}

	var accounts []acm.Account
	var validators []acm.Validator

	for i, templateAccount := range templateAccounts {
		account, err := templateAccount.Account(keyClient, i)
		if err != nil {
			return nil, fmt.Errorf("could not create Account from template: %v", err)
		}
		accounts = append(accounts, account)
		// Create a corresponding validator
		if templateAccount.AmountBonded != nil {
			// Note this does not modify the input template
			addr := account.Address()
			templateAccount.Address = &addr
			validator, err := templateAccount.Validator(keyClient, i)
			if err != nil {
				return nil, fmt.Errorf("could not create Validator from template: %v", err)
			}
			validators = append(validators, validator)
		}
	}

	genesisDoc := genesis.MakeGenesisDocFromAccounts(chainName, nil, genesisTime, globalPermissions, accounts, validators)

	return &genesisDoc, nil
}

func (gs *GenesisSpec) JSONBytes() ([]byte, error) {
	bs, err := json.Marshal(gs)
	if err != nil {
		return nil, err
	}
	// rewrite buffer with indentation
	indentedBuffer := new(bytes.Buffer)
	if err := json.Indent(indentedBuffer, bs, "", "\t"); err != nil {
		return nil, err
	}
	return indentedBuffer.Bytes(), nil
}

func (gs *GenesisSpec) Hash() []byte {
	gsBytes, err := gs.JSONBytes()
	if err != nil {
		panic(fmt.Errorf("could not create hash of GenesisDoc: %v", err))
	}
	hasher := sha256.New()
	hasher.Write(gsBytes)
	return hasher.Sum(nil)
}

func (gs *GenesisSpec) ShortHash() []byte {
	return gs.Hash()[:genesis.ShortHashSuffixBytes]
}

func GenesisSpecFromJSON(jsonBlob []byte) (*GenesisSpec, error) {
	genDoc := new(GenesisSpec)
	err := json.Unmarshal(jsonBlob, genDoc)
	if err != nil {
		return nil, fmt.Errorf("couldn't read GenesisSpec: %v", err)
	}
	return genDoc, nil
}

func accountNameFromIndex(index int) string {
	return fmt.Sprintf("Account_%v", index)
}
