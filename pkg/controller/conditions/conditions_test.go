package conditions_test

import (
	"time"

	"github.com/nmstate/kubernetes-nmstate/pkg/controller/conditions"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	nmstatev1 "github.com/nmstate/kubernetes-nmstate/pkg/apis/nmstate/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var _ = Describe("Conditions", func() {
	Context("when instance is empty", func() {
		It("should return nil", func() {
			condition := conditions.Condition(&nmstatev1.NodeNetworkState{}, nmstatev1.NodeNetworkStateConditionFailing)
			Expect(condition).To(BeNil())
		})
	})

	Context("when composed by one condition with known type and unknown status", func() {
		nodeNetworkState := nmstatev1.NodeNetworkState{
			Status: nmstatev1.NodeNetworkStateStatus{
				Conditions: []nmstatev1.NodeNetworkStateCondition{
					nmstatev1.NodeNetworkStateCondition{
						Type:   nmstatev1.NodeNetworkStateConditionAvailable,
						Status: corev1.ConditionTrue,
					},
					nmstatev1.NodeNetworkStateCondition{
						Type:   nmstatev1.NodeNetworkStateConditionFailing,
						Status: corev1.ConditionUnknown,
					},
				},
			},
		}

		It("should find the correct condition", func() {
			instance := nodeNetworkState.DeepCopy()
			condition := conditions.Condition(instance, nmstatev1.NodeNetworkStateConditionFailing)
			Expect(condition).ToNot(BeNil())
			Expect(condition.Type).To(Equal(nmstatev1.NodeNetworkStateConditionFailing))
			Expect(condition.Status).To(Equal(corev1.ConditionUnknown))
		})
	})

	Context("when is empty", func() {
		nodeNetworkState := &nmstatev1.NodeNetworkState{}

		It("should add new condition", func() {
			instance := nodeNetworkState.DeepCopy()
			conditionType := nmstatev1.NodeNetworkStateConditionFailing

			conditions.SetCondition(instance, conditionType, corev1.ConditionFalse, "foo", "bar")
			condition := conditions.Condition(instance, conditionType)

			Expect(condition.Type).To(Equal(conditionType))
			Expect(condition.Status).To(Equal(corev1.ConditionFalse))
			Expect(condition.Reason).To(Equal("foo"))
			Expect(condition.Message).To(Equal("bar"))
			Expect(condition.LastHeartbeatTime).To(Equal(condition.LastTransitionTime))
		})
	})

	Context("when there is condition", func() {
		nodeNetworkState := nmstatev1.NodeNetworkState{
			Status: nmstatev1.NodeNetworkStateStatus{
				Conditions: []nmstatev1.NodeNetworkStateCondition{
					nmstatev1.NodeNetworkStateCondition{
						Type:    nmstatev1.NodeNetworkStateConditionFailing,
						Status:  corev1.ConditionUnknown,
						Reason:  "foo",
						Message: "bar",
						LastHeartbeatTime: metav1.Time{
							Time: time.Unix(0, 0),
						},
						LastTransitionTime: metav1.Time{
							Time: time.Unix(0, 0),
						},
					},
				},
			},
		}

		Context("and we add different condition", func() {
			instance := nodeNetworkState.DeepCopy()
			conditionType := nmstatev1.NodeNetworkStateConditionAvailable

			It("should add new condition", func() {
				conditions.SetCondition(instance, conditionType, corev1.ConditionTrue, "foo", "bar")
				Expect(len(instance.Status.Conditions)).To(Equal(len(nodeNetworkState.Status.Conditions) + 1))

				condition := conditions.Condition(instance, conditionType)
				Expect(condition.Type).To(Equal(conditionType))
				Expect(condition.Status).To(Equal(corev1.ConditionTrue))
				Expect(condition.Reason).To(Equal("foo"))
				Expect(condition.Message).To(Equal("bar"))
				Expect(condition.LastHeartbeatTime).To(Equal(condition.LastTransitionTime))
			})
		})

		Context("and we update the condition with same values", func() {
			conditionType := nmstatev1.NodeNetworkStateConditionFailing

			It("should change LastHearbeatTime", func() {
				instance := nodeNetworkState.DeepCopy()
				conditions.SetCondition(instance, conditionType, corev1.ConditionFalse, "foo", "bar")

				By("Shouldn't add new condition")
				Expect(len(instance.Status.Conditions)).To(Equal(len(nodeNetworkState.Status.Conditions)))

				condition := conditions.Condition(instance, conditionType)
				originalCondition := conditions.Condition(&nodeNetworkState, conditionType)
				By("Should change LastHeartbeatTime")
				Expect(originalCondition.LastHeartbeatTime.Time.Before(condition.LastHeartbeatTime.Time)).To(BeTrue(), "LastHeartbeatTime of updated condition wasn't updated")
			})
		})

		Context("and we update the condition with different values", func() {
			conditionType := nmstatev1.NodeNetworkStateConditionFailing

			It("should change values and update LastTransitionTime", func() {
				instance := nodeNetworkState.DeepCopy()
				conditions.SetCondition(instance, conditionType, corev1.ConditionTrue, "bar", "foo")

				By("Shouldn't add new condition")
				Expect(len(instance.Status.Conditions)).To(Equal(len(nodeNetworkState.Status.Conditions)))

				condition := conditions.Condition(instance, conditionType)
				originalCondition := conditions.Condition(&nodeNetworkState, conditionType)
				By("Should change different values")
				Expect(condition.Status).To(Equal(corev1.ConditionTrue))
				Expect(condition.Reason).To(Equal("bar"))
				Expect(condition.Message).To(Equal("foo"))

				By("Should change LastTransitionTime")
				Expect(originalCondition.LastTransitionTime.Time.Before(condition.LastTransitionTime.Time)).To(BeTrue(), "LastTransitionTime of updated condition wasn't updated")

				By("Should change LastHeartbeatTime")
				Expect(originalCondition.LastHeartbeatTime.Time.Before(condition.LastHeartbeatTime.Time)).To(BeTrue(), "LastHeartbeatTime of updated condition wasn't updated")
			})
		})
	})
})
