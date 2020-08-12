package mint

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/mint/internal/types"
)

// Joe.He
func beginPerBlockMintBlocker(ctx sdk.Context, k Keeper, minter *types.Minter, params types.Params,
	totalStakingSupply sdk.Int) (mintCoin sdk.Coins) {

	if len(params.PerBlockMint) == 0 {
		return sdk.NewCoins(sdk.NewCoin(params.MintDenom, sdk.ZeroInt()))
	}
	mintCoin = k.GetCurrPerBlockMint(ctx)

	if mintCoin.AmountOf(params.MintDenom).IsZero() {
		return sdk.NewCoins(sdk.NewCoin(params.MintDenom, sdk.ZeroInt()))
	}
	minter.AnnualProvisions = sdk.NewDec(int64(params.BlocksPerYear)).MulInt(mintCoin.AmountOf(params.MintDenom))
	minter.Inflation = minter.AnnualProvisions.QuoInt(totalStakingSupply)
	return
}

// Joe.He
func mintCoins(ctx sdk.Context, k Keeper, minter *types.Minter, params types.Params, mintedCoins sdk.Coins) {
	if params.MinterSupplyAddress != nil {
		_, err := k.GetBankKeeper().SubtractCoins(ctx, params.MinterSupplyAddress, mintedCoins)
		if err == nil {
			minter.AnnualProvisions = sdk.ZeroDec()
			minter.Inflation = sdk.ZeroDec()
			if err := k.MintSelfCoins(ctx, mintedCoins); err != nil {
				panic(err)
			}
			return
		}
	}

	if err := k.MintCoins(ctx, mintedCoins); err != nil {
		panic(err)
	}
}

// BeginBlocker mints new tokens for the previous block.
func BeginBlocker(ctx sdk.Context, k Keeper) {
	// fetch stored minter & params
	minter := k.GetMinter(ctx)
	params := k.GetParams(ctx)

	// recalculate inflation rate
	totalStakingSupply := k.StakingTokenSupply(ctx)
	bondedRatio := k.BondedRatio(ctx)
	// Joe.He
	//minter.Inflation = minter.NextInflationRate(params, bondedRatio)
	//minter.AnnualProvisions = minter.NextAnnualProvisions(params, totalStakingSupply)

	minter.Inflation = sdk.ZeroDec()
	minter.AnnualProvisions = sdk.ZeroDec()

	// Joe.He
	mintedCoins := beginPerBlockMintBlocker(ctx, k, &minter, params, totalStakingSupply)

	// Joe.He
	if !mintedCoins.AmountOf(params.MintDenom).IsZero() {
		mintCoins(ctx, k, &minter, params, mintedCoins)
		if err := k.AddCollectedFees(ctx, mintedCoins); err != nil {
			panic(err)
		}
	}
	k.SetMinter(ctx, minter)

	//// mint coins, update supply
	//mintedCoin := minter.BlockProvision(params)
	//mintedCoins := sdk.NewCoins(mintedCoin)
	//
	//err := k.MintCoins(ctx, mintedCoins)
	//if err != nil {
	//	panic(err)
	//}
	//
	// send the minted coins to the fee collector account
	//err := k.AddCollectedFees(ctx, mintedCoins)
	//if err != nil {
	//	panic(err)
	//}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(types.AttributeKeyBondedRatio, bondedRatio.String()),
			sdk.NewAttribute(types.AttributeKeyInflation, minter.Inflation.String()),
			sdk.NewAttribute(types.AttributeKeyAnnualProvisions, minter.AnnualProvisions.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, mintedCoins.String()),
		),
	)
}
