package keeper

import (
	"context"
	"fmt"

	"github.com/cosmos/cosmos-sdk/crypto"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	AccountKeeper
}

// NewMsgServerImpl returns an implementation of the x/auth MsgServer interface.
func NewMsgServerImpl(ak AccountKeeper) types.MsgServer {
	return &msgServer{
		AccountKeeper: ak,
	}
}

func (ms msgServer) UpdateParams(goCtx context.Context, req *types.MsgUpdateParams) (*types.MsgUpdateParamsResponse, error) {
	if ms.authority != req.Authority {
		return nil, errors.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", ms.authority, req.Authority)
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := ms.SetParams(ctx, req.Params); err != nil {
		return nil, err
	}

	return &types.MsgUpdateParamsResponse{}, nil
}

func (ms msgServer) ChangePubKey(goCtx context.Context, req *types.MsgChangePubKey) (*types.MsgChangePubKeyResponse, error) {

	ctx := sdk.UnwrapSDKContext(goCtx)

	pubKeyBytes, _, err := crypto.UnarmorPubKeyBytes(req.PubKey)
	if err != nil {
		return &types.MsgChangePubKeyResponse{}, fmt.Errorf("invalid public key: %w", err)
	}

	var pubKey cryptotypes.PubKey
	if err := ms.cdc.UnmarshalInterface(pubKeyBytes, &pubKey); err != nil {
		return &types.MsgChangePubKeyResponse{}, fmt.Errorf("cannont unmarshal public key: %w", err)
	}

	acc := ms.AccountKeeper.GetAccount(ctx, sdk.AccAddress(req.Address))
	ms.AccountKeeper.ChangePubKey(ctx, acc, pubKey)

	return &types.MsgChangePubKeyResponse{}, nil
}
