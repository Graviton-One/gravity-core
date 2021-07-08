package node

import (
	"context"
	"fmt"

	"github.com/Gravity-Tech/gravity-core/common/account"
	"github.com/Gravity-Tech/gravity-core/common/gravity"
	"github.com/Gravity-Tech/gravity-core/common/state"
	"github.com/Gravity-Tech/gravity-core/oracle/extractor"
	"go.uber.org/zap"
)

// pulseId uint64, round state.SubRound, tcHeight uint64, intervalId uint64, roundState *RoundState, ctx context.Context
type roundExecuteProps struct {
	PulseID           uint64
	Round             state.SubRound
	TargetChainHeight uint64
	IntervalID        uint64
	RoundState       *RoundState
	Ctx               context.Context
}

type roundExecutor struct {}


var oracleRoundExecutor = &roundExecutor{}



func (re *roundExecutor) Execute(node *Node, props *roundExecuteProps) error {
	pulseId := props.PulseID
	round := props.Round
	// tcHeight := props.TargetChainHeight
	roundState := props.RoundState
	ctx := props.Ctx
	intervalId := props.IntervalID

	switch round {
	case state.CommitSubRound:
		zap.L().Sugar().Debugf("Commit subround pulseId: %d", pulseId)
		if len(roundState.commitHash) != 0 {
			zap.L().Sugar().Debug("Len(commit hash): ", len(roundState.commitHash), roundState.commitHash)
			return nil
		}
		_, err := node.gravityClient.CommitHash(node.chainType, node.nebulaId, int64(intervalId), int64(pulseId), node.oraclePubKey)
		if err != nil && err != gravity.ErrValueNotFound {
			zap.L().Error(err.Error())
			return err
		} else if err == nil {
			zap.L().Sugar().Debugf("Commit subround pulseId: %d intervalId: %d exists", pulseId, intervalId)
			return nil
		}

		data, err := node.extractor.Extract(ctx)
		if err != nil && err != extractor.NotFoundErr {
			zap.L().Error(err.Error())
			return err
		} else if err == extractor.NotFoundErr {
			return nil
		}

		if data == nil {
			zap.L().Debug("Commit subround Extractor Data is empty")
			return nil
		}
		zap.L().Sugar().Debug("Extracted data ", data)

		commit, err := node.invokeCommitTx(data, intervalId, pulseId)
		if err != nil {
			return err
		}

		roundState.commitHash = commit
		roundState.data = data
		zap.L().Sugar().Debug("Commit round end ", roundState)
	case state.RevealSubRound:
		zap.L().Debug("Reveal subround")
		if len(roundState.commitHash) == 0 || roundState.RevealExist {
			zap.L().Sugar().Debugf("CommitHash is empty: %t, RevealExist: %t", len(roundState.commitHash) == 0, roundState.RevealExist)
			return nil
		}
		_, err := node.gravityClient.Reveal(node.chainType, node.oraclePubKey, node.nebulaId, int64(intervalId), int64(pulseId), roundState.commitHash)
		if err != nil && err != gravity.ErrValueNotFound {
			zap.L().Error(err.Error())
			return err
		} else if err == nil {
			return nil
		}

		err = node.invokeRevealTx(intervalId, pulseId, roundState.data, roundState.commitHash)
		if err != nil {
			zap.L().Error(err.Error())
			return err
		}
		roundState.RevealExist = true
		zap.L().Sugar().Debug("Reveal round end ", roundState)
	case state.ResultSubRound:
		zap.L().Debug("Result subround")
		if roundState.data == nil && !roundState.RevealExist {
			return nil
		}
		if roundState.resultValue != nil {
			zap.L().Debug("Round sign exists")
			return nil
		}
		value, hash, err := node.signRoundResult(intervalId, pulseId, ctx)
		if err != nil {
			zap.L().Error(err.Error())
			return err
		}
		//TODO migrate to err
		if value == nil {
			zap.L().Sugar().Debugf("Value is nil: %t", value == nil)
			return nil
		}

		roundState.resultValue = value
		roundState.resultHash = hash
	case state.SendToTargetChain:
		zap.L().Debug("Send to target chain subround")
		var oracles []account.OraclesPubKey
		var myRound uint64

		if roundState.isSent || roundState.resultValue == nil {
			zap.L().Sugar().Debugf("roundState.isSent: %t, resultValue is nil: %t", roundState.isSent, roundState.resultValue == nil)
			return nil
		}

		oraclesMap, err := node.gravityClient.BftOraclesByNebula(node.chainType, node.nebulaId)
		if err != nil {
			zap.L().Sugar().Debugf("BFT error: %s , \n %s", err, zap.Stack("trace").String)
			return nil
		}
		if _, ok := oraclesMap[node.oraclePubKey.ToString(node.chainType)]; !ok {
			zap.L().Debug("Oracle not found")
			return nil
		}

		var count uint64
		for k, v := range oraclesMap {
			oracle, err := account.StringToOraclePubKey(k, v)
			if err != nil {
				return err
			}
			oracles = append(oracles, oracle)
			if node.oraclePubKey == oracle {
				myRound = count
			}
			count++
		}

		if len(oracles) == 0 {
			zap.L().Debug("Oracles map is empty")
			return nil
		}
		if intervalId%uint64(len(oracles)) != myRound {
			zap.L().Debug("Len oracles != myRound")
			return nil
		}
		zap.L().Sugar().Debugf("Adding pulse id: %d", pulseId)

		txId, err := node.adaptor.AddPulse(node.nebulaId, pulseId, oracles, roundState.resultHash, ctx)

		if err != nil {
			zap.L().Sugar().Debugf("Error: %s", err)
			return err
		}

		if txId != "" {
			err = node.adaptor.WaitTx(txId, ctx)
			if err != nil {
				zap.L().Sugar().Debugf("Error: %s", err)
				return err
			}

			zap.L().Sugar().Infof("Result tx id: %s", txId)

			roundState.isSent = true
			zap.L().Sugar().Debugf("Sending Value to subs, pulse id: %d", pulseId)
			err = node.adaptor.SendValueToSubs(node.nebulaId, pulseId, roundState.resultValue, ctx)
			if err != nil {
				zap.L().Sugar().Debugf("Error: %s", err)
				return err
			}
		} else {
			fmt.Printf("Info: Tx result not sent")
		}
	}

	return nil
}