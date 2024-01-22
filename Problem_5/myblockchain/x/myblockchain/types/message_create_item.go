package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateItem{}

func NewMsgCreateItem(creator string, id string, name string, details string) *MsgCreateItem {
	return &MsgCreateItem{
		Creator: creator,
		Id:      id,
		Name:    name,
		Details: details,
	}
}

func (msg *MsgCreateItem) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
