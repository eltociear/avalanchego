// Copyright (C) 2019-2021, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package indexes

import (
	"math"
	"sync"
	"time"

	"github.com/ava-labs/avalanchego/database"
	"github.com/ava-labs/avalanchego/ids"
	"github.com/ava-labs/avalanchego/utils"
	"github.com/ava-labs/avalanchego/utils/logging"
	"github.com/ava-labs/avalanchego/utils/units"
	"github.com/ava-labs/avalanchego/vms/proposervm/state"
)

const defaultCommitSizeCap = 1 * units.MiB

var _ HeightIndexer = &heightIndexer{}

type HeightIndexer interface {
	// signals whether index rebuilding it still ongoing
	IsRepaired() bool

	// checks whether index rebuilding is needed and if so, performs it
	RepairHeightIndex() error
}

func NewHeightIndexer(srv BlockServer,
	log logging.Logger,
	indexState state.HeightIndex,
	shutdownChan chan struct{},
	shutdownWg *sync.WaitGroup) HeightIndexer {
	return newHeightIndexer(srv, log, indexState, shutdownChan, shutdownWg)
}

func newHeightIndexer(srv BlockServer,
	log logging.Logger,
	indexState state.HeightIndex,
	shutdownChan chan struct{},
	shutdownWg *sync.WaitGroup) *heightIndexer {
	res := &heightIndexer{
		server:        srv,
		log:           log,
		shutdownChan:  shutdownChan,
		shutdownWg:    shutdownWg,
		indexState:    indexState,
		commitMaxSize: defaultCommitSizeCap,
	}

	return res
}

type heightIndexer struct {
	server       BlockServer
	log          logging.Logger
	shutdownChan chan struct{}
	shutdownWg   *sync.WaitGroup

	jobDone    utils.AtomicBool
	indexState state.HeightIndex

	commitMaxSize int
}

func (hi *heightIndexer) IsRepaired() bool {
	return hi.jobDone.GetValue()
}

// RepairHeightIndex ensures the height -> proBlkID height block index is well formed.
// Starting from last accepted proposerVM block, it will go back to snowman++ activation fork
// or genesis. PreFork blocks will be handled by innerVM height index.
// RepairHeightIndex can take a non-trivial time to complete; hence we make sure
// the process has limited memory footprint, can be resumed from periodic checkpoints
// and works asynchronously without blocking the VM.
func (hi *heightIndexer) RepairHeightIndex() error {
	defer hi.shutdownWg.Done()
	needRepair, startBlkID, err := hi.shouldRepair()
	if err != nil {
		hi.log.Info("Block indexing by height starting: failed. Could not determine if index is complete, error %v", err)
		return err
	}
	if !needRepair {
		forkHeight, err := hi.indexState.GetForkHeight()
		if err != nil {
			return err
		}
		hi.log.Info("Block indexing by height: already complete. Fork height %d", forkHeight)
		return nil
	}
	return hi.doRepair(startBlkID)
}

// shouldRepair checks if height index is complete;
// if not, it returns the checkpoint from which repairing should start.
func (hi *heightIndexer) shouldRepair() (bool, ids.ID, error) {
	switch checkpointID, err := hi.indexState.GetCheckpoint(); err {
	case nil:
		// checkpoint found, repair must be resumed
		hi.log.Info("Block indexing by height starting: success. Retrieved checkpoint %v", checkpointID)
		return true, checkpointID, nil

	case database.ErrNotFound:
		// no checkpoint. Either index is complete or repair was never attempted.
		hi.log.Info("Block indexing by height starting: checkpoint not found. Verifying index is complete...")

	default:
		return true, ids.Empty, err
	}

	// index is complete iff lastAcceptedBlock is indexed
	latestProBlkID, err := hi.server.LastAcceptedWrappingBlkID()
	switch err {
	case nil:
		break

	case database.ErrNotFound:
		// snowman++ has not forked yet; height block index is ok.
		// forkHeight set at math.MaxUint64, aka +infinity
		if err := hi.indexState.SetForkHeight(math.MaxUint64); err != nil {
			return true, ids.Empty, err
		}
		hi.jobDone.SetValue(true)
		hi.log.Info("Block indexing by height starting: Snowman++ fork not reached yet. No need to rebuild index.")
		return false, ids.Empty, nil

	default:
		return true, ids.Empty, err
	}

	lastAcceptedBlk, err := hi.server.GetWrappingBlk(latestProBlkID)
	if err != nil {
		// Could not retrieve last accepted block.
		// We got bigger problems than repairing the index
		return true, ids.Empty, err
	}

	_, err = hi.indexState.GetBlockIDAtHeight(lastAcceptedBlk.Height())
	switch err {
	case nil:
		// index is complete already. Just make sure forkHeight can be read
		if _, err := hi.indexState.GetForkHeight(); err != nil {
			return true, ids.Empty, err
		}
		hi.jobDone.SetValue(true)
		hi.log.Info("Block indexing by height starting: Index already complete, nothing to do.")
		return false, ids.Empty, nil

	case database.ErrNotFound:
		// index needs repairing and it's the first time we do this.
		// Mark the checkpoint so that, in case new blocks are accepted while
		// indexing is ongoing, and the process is terminated before first commit,
		// we do not miss rebuilding the full index.
		if err := hi.indexState.SetCheckpoint(latestProBlkID); err != nil {
			return true, ids.Empty, err
		}
		if err := hi.server.DBCommit(); err != nil {
			return true, ids.Empty, err
		}

		// also duly init forkHeight at math.MaxUint64, aka +infinity
		if err := hi.indexState.SetForkHeight(math.MaxUint64); err != nil {
			return true, ids.Empty, err
		}
		hi.log.Info("Block indexing by height starting: index incomplete. Rebuilding from %v", latestProBlkID)
		return true, latestProBlkID, nil

	default:
		return true, ids.Empty, err
	}
}

// if height index needs repairing, doRepair would do that. It
// iterates back via parents, checking and rebuilding height indexing
func (hi *heightIndexer) doRepair(repairStartBlkID ids.ID) error {
	var (
		currentProBlkID   = repairStartBlkID
		currentInnerBlkID = ids.Empty

		start                     = time.Now()
		lastLogTime               = start
		indexedBlks               = 0
		pendingBytesApproximation = 0 // tracks of the size of uncommitted writes
	)

	for {
		// handle graceful termination
		select {
		case <-hi.shutdownChan:
			// ChainVM.Shutdown is called with ctx.Lock hold. Hence there
			// is not much we can do here (certainly no calls to blockServer). Log and exit
			hi.log.Info("Block indexing by height: shutdown called.")
			return nil
		default:
			// go ahead with index repairing
		}

		currentAcceptedBlk, err := hi.server.GetWrappingBlk(currentProBlkID)
		switch err {
		case nil:

		case database.ErrNotFound:
			// visited all proposerVM blocks. Let's record forkHeight ...
			firstWrappedInnerBlk, err := hi.server.GetInnerBlk(currentInnerBlkID)
			if err != nil {
				return err
			}
			forkHeight := firstWrappedInnerBlk.Height()
			if err := hi.indexState.SetForkHeight(forkHeight); err != nil {
				return err
			}

			// ... delete checkpoint and finally commit
			if err := hi.indexState.DeleteCheckpoint(); err != nil {
				return err
			}
			if err := hi.server.DBCommit(); err != nil {
				return err
			}
			hi.log.Info("Block indexing by height: completed. Indexed %d blocks, duration %v, fork height %d",
				indexedBlks, time.Since(start), forkHeight)
			return nil

		default:
			return err
		}

		currentInnerBlkID = currentAcceptedBlk.GetInnerBlk().ID()

		// Rebuild height block index.
		estimatedByteLen, err := hi.indexState.SetBlockIDAtHeight(currentAcceptedBlk.Height(), currentProBlkID)
		if err != nil {
			return err
		}
		pendingBytesApproximation += estimatedByteLen

		// Let's keep memory footprint under control by committing when a size threshold is reached
		if pendingBytesApproximation > hi.commitMaxSize {
			if err := hi.doCheckpoint(currentAcceptedBlk); err != nil {
				return err
			}
			if err := hi.indexState.SetForkHeight(currentAcceptedBlk.Height()); err != nil {
				return err
			}
			if err := hi.server.DBCommit(); err != nil {
				return err
			}
			hi.log.Info("Block indexing by height: ongoing. Indexed %d blocks, latest committed height %d, committed %d bytes",
				indexedBlks, currentAcceptedBlk.Height()+1, pendingBytesApproximation)
			pendingBytesApproximation = 0
		}

		// Periodically log progress
		indexedBlks++
		if time.Since(lastLogTime) > 15*time.Second {
			lastLogTime = time.Now()
			hi.log.Info("Block indexing by height: ongoing. Indexed %d blocks, latest indexed height %d",
				indexedBlks, currentAcceptedBlk.Height()+1)
		}

		// keep checking the parent
		currentProBlkID = currentAcceptedBlk.Parent()
	}
}

func (hi *heightIndexer) doCheckpoint(currentProBlk WrappingBlock) error {
	// checkpoint is current block's parent, it if exists
	var checkpoint ids.ID
	parentBlkID := currentProBlk.Parent()
	checkpointBlk, err := hi.server.GetWrappingBlk(parentBlkID)
	switch err {
	case nil:
		checkpoint = checkpointBlk.ID()
		if err := hi.indexState.SetCheckpoint(checkpoint); err != nil {
			return err
		}
		hi.log.Info("Block indexing by height. Stored checkpoint %v at height %d",
			currentProBlk.ID(), currentProBlk.Height())
		return nil

	case database.ErrNotFound:
		// parent must be a preFork block. We do not checkpoint here.
		// Process will set forkHeight and terminate
		return nil

	default:
		return err
	}
}
