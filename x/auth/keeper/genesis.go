package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

// InitGenesis - Init store state from genesis data
//
// CONTRACT: old coins from the FeeCollectionKeeper need to be transferred through
// a genesis port script to the new fee collector account
func (ak AccountKeeper) InitGenesis(ctx sdk.Context, data types.GenesisState) {
	if err := ak.SetParams(ctx, data.Params); err != nil {
		panic(err)
	}

	accounts, err := types.UnpackAccounts(data.Accounts)
	if err != nil {
		panic(err)
	}
	accounts = types.SanitizeGenesisAccounts(accounts)

	for _, a := range accounts {
		acc := ak.NewAccount(ctx, a)
		ak.SetAccount(ctx, acc)
	}

	ak.GetModuleAccount(ctx, types.FeeCollectorName)

	for _, pubKeyMapping := range data.AddressToPubKeys {
		ak.SavePubKeyMapping(ctx, pubKeyMapping)
	}
}

// ExportGenesis returns a GenesisState for a given context and keeper
func (ak AccountKeeper) ExportGenesis(ctx sdk.Context) *types.GenesisState {
	params := ak.GetParams(ctx)

	var genAccounts types.GenesisAccounts
	ak.IterateAccounts(ctx, func(account types.AccountI) bool {
		genAccount := account.(types.GenesisAccount)
		genAccounts = append(genAccounts, genAccount)
		return false
	})

	pubKeyMappings := ak.GetAllPubKeys(ctx)

	return types.NewGenesisState(params, genAccounts, pubKeyMappings)
}
