// Copyright (C) 2019-2023, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package faultinjection

import (
	"fmt"
	"math/rand"
	"strings"

	ginkgo "github.com/onsi/ginkgo/v2"
	"github.com/stretchr/testify/require"

	runner_sdk "github.com/ava-labs/avalanche-network-runner-sdk"

	"github.com/ava-labs/avalanchego/staking"
	"github.com/ava-labs/avalanchego/tests/e2e"
)

var _ = ginkgo.Describe("Duplicate node handling", func() {
	require := require.New(ginkgo.GinkgoT())

	ginkgo.It("should ensure that a given Node ID (i.e. staking keypair) can be used at most once on a network", func() {
		// TODO(marun) Ensure reliable node removal on teardown

		// Minimize the potential for node name collision to allow for
		// iterating with a persistent network. The use of non-secure random
		// value is not security critical so the lint is ignored.
		//nolint:gosec
		baseNodeName := fmt.Sprintf("e2e-duplicate-node-%08x-", rand.Uint32())
		node1Name := baseNodeName + "1"

		var stakingCert, stakingKey []byte
		ginkgo.By("generating a staking keypair", func() {
			var err error
			stakingCert, stakingKey, err = staking.NewCertAndKeyBytes()
			require.NoError(err)
		})

		ginkgo.By("creating the first node using the staking keypair", func() {
			err := e2e.Env.AddNode(node1Name, runner_sdk.WithGlobalNodeConfig(fmt.Sprintf(
				`{""--staking-tls-cert-contents": "%s", "--staking-tls-key-contents": "%s"}`,
				// Escape newlines in the PEM-formatted keys to
				// ensure compatibility with grpc.
				strings.ReplaceAll(string(stakingCert), "\n", "\\n"),
				strings.ReplaceAll(string(stakingKey), "\n", "\\n"),
			)))
			require.NoError(err)
			err = e2e.Env.CheckHealth()
			require.NoError(err)
		})

		// TODO(marun)
		// - Boot badNode2 with the same Certs
		// - Check the network peers are connected amongst the network with badNode1 (the first to go up)
		// - Check that badNode2 can't bootstrap (invalid peer alias)
		// - Remove badNode1
		// - Check that badNode2 can now bootstrap
		// - Check both nodes have/had the same nodeID
		// - Check the network peers are connected amongst the network with badNode2
	})
})
