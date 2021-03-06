package e2e

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/tidwall/gjson"

	nmstatev1alpha1 "github.com/nmstate/kubernetes-nmstate/pkg/apis/nmstate/v1alpha1"
)

func hasVlans(result gjson.Result, maxVlan int) {
	vlans := result.Array()
	ExpectWithOffset(1, vlans).To(HaveLen(maxVlan))
	for i, vlan := range vlans {
		Expect(vlan.Get("vlan").Int()).To(Equal(int64(i + 1)))
	}
}

var _ = Describe("NodeNetworkState", func() {
	var (
		bond1Up = nmstatev1alpha1.State(`interfaces:
  - name: eth1
    type: ethernet
    state: up
  - name: bond1
    type: bond
    state: up
    link-aggregation:
      mode: active-backup
      slaves:
        - eth1
      options:
        miimon: '120'
`)
		br1Up = nmstatev1alpha1.State(`interfaces:
  - name: eth1
    type: ethernet
    state: up
  - name: br1
    type: linux-bridge
    state: up
    bridge:
      options:
        stp:
          enabled: false
      port:
        - name: eth1
          stp-hairpin-mode: false
          stp-path-cost: 100
          stp-priority: 32
`)
		br1WithBond1Up = nmstatev1alpha1.State(`interfaces:
  - name: eth1
    type: ethernet
    state: up
  - name: bond1
    type: bond
    state: up
    link-aggregation:
      mode: active-backup
      slaves:
        - eth1
      options:
        miimon: '120'
  - name: br1
    type: linux-bridge
    state: up
    bridge:
      options:
        stp:
          enabled: false
      port:
        - name: bond1
          stp-hairpin-mode: false
          stp-path-cost: 100
          stp-priority: 32
`)

		br1Absent = nmstatev1alpha1.State(`interfaces:
  - name: br1
    type: linux-bridge
    state: absent
  - name: eth1
    type: ethernet
    state: absent
`)
	)
	Context("when desiredState is configured", func() {
		Context("with a linux bridge up", func() {
			BeforeEach(func() {
				updateDesiredState(br1Up)
			})
			AfterEach(func() {

				// First we clean desired state if we
				// don't do that nmstate recreates the bridge
				resetDesiredStateForNodes()

				// TODO: Add status conditions to ensure that
				//       it has being really reset so we can
				//       remove this ugly sleep
				time.Sleep(1 * time.Second)

				// Let's clean the bridge directly in the node
				// bypassing nmstate
				deleteConnectionAtNodes("eth1")
				deleteConnectionAtNodes("br1")
			})
			It("should have the linux bridge at currentState", func() {
				for _, node := range nodes {
					interfacesNameForNode(node).Should(ContainElement("br1"))
				}
			})
		})
		Context("with a linux bridge absent", func() {
			BeforeEach(func() {
				createBridgeAtNodes("br1")
				updateDesiredState(br1Absent)
			})
			AfterEach(func() {
				// If not br1 is going to be removed if created manually
				resetDesiredStateForNodes()
			})
			It("should have the linux bridge at currentState", func() {
				for _, node := range nodes {
					interfacesNameForNode(node).ShouldNot(ContainElement("br1"))
				}
			})
		})
		Context("with a active-backup miimon 100 bond interface up", func() {
			BeforeEach(func() {
				updateDesiredState(bond1Up)
			})
			AfterEach(func() {

				resetDesiredStateForNodes()

				// TODO: Add status conditions to ensure that
				//       it has being really reset so we can
				//       remove this ugly sleep
				time.Sleep(1 * time.Second)

				deleteConnectionAtNodes("bond1")
				deleteConnectionAtNodes("eth1")
			})
			It("should have the bond interface at currentState", func() {
				var (
					expectedBond = interfaceByName(interfaces(bond1Up), "bond1")
				)

				for _, node := range nodes {
					interfacesForNode(node).Should(ContainElement(SatisfyAll(
						HaveKeyWithValue("name", expectedBond["name"]),
						HaveKeyWithValue("type", expectedBond["type"]),
						HaveKeyWithValue("state", expectedBond["state"]),
						HaveKeyWithValue("link-aggregation", expectedBond["link-aggregation"]),
					)))
				}
			})
		})
		Context("with the bond interface as linux bridge port", func() {
			BeforeEach(func() {
				createBridgeAtNodes("br2", "eth2")
				updateDesiredState(br1WithBond1Up)
			})
			AfterEach(func() {

				resetDesiredStateForNodes()

				// TODO: Add status conditions to ensure that
				//       it has being really reset so we can
				//       remove this ugly sleep
				time.Sleep(1 * time.Second)

				deleteConnectionAtNodes("eth1")
				deleteConnectionAtNodes("eth2")
				deleteConnectionAtNodes("br1")
				deleteConnectionAtNodes("br2")
				deleteConnectionAtNodes("bond1")
			})
			It("should have the bond in the linux bridge as port at currentState", func() {
				var (
					expectedInterfaces  = interfaces(br1WithBond1Up)
					expectedBond        = interfaceByName(expectedInterfaces, "bond1")
					expectedBridge      = interfaceByName(expectedInterfaces, "br1")
					expectedBridgePorts = expectedBridge["bridge"].(map[string]interface{})["port"]
				)
				for _, node := range nodes {
					interfacesForNode(node).Should(SatisfyAll(
						ContainElement(SatisfyAll(
							HaveKeyWithValue("name", expectedBond["name"]),
							HaveKeyWithValue("type", expectedBond["type"]),
							HaveKeyWithValue("state", expectedBond["state"]),
							HaveKeyWithValue("link-aggregation", expectedBond["link-aggregation"]),
						)),
						ContainElement(SatisfyAll(
							HaveKeyWithValue("name", expectedBridge["name"]),
							HaveKeyWithValue("type", expectedBridge["type"]),
							HaveKeyWithValue("state", expectedBridge["state"]),
							HaveKeyWithValue("bridge", HaveKeyWithValue("port", expectedBridgePorts)),
						))))
				}
				for _, bridgeVlans := range bridgeVlansAtNodes() {
					parsedVlans := gjson.Parse(bridgeVlans)

					hasVlans(parsedVlans.Get("bond1"), 4094)
					hasVlans(parsedVlans.Get("br1"), 4094)

					br2 := parsedVlans.Get("br2")
					hasVlans(br2, 1)

					eth1Vlans := parsedVlans.Get("eth1").Array()
					Expect(eth1Vlans).To(BeEmpty())

					eth2Vlans := parsedVlans.Get("eth2").Array()
					Expect(eth2Vlans).To(BeEmpty())
				}
			})
		})
	})
})
