package hiero

// SPDX-License-Identifier: Apache-2.0

type FeeAssessmentMethod bool

const (
	FeeAssessmentMethodInclusive FeeAssessmentMethod = false
	FeeAssessmentMethodExclusive FeeAssessmentMethod = true
)

// String returns a string representation of the FeeAssessmentMethod
func (assessment FeeAssessmentMethod) String() string {
	if assessment {
		return "FeeAssessmentMethodExclusive"
	}

	return "FeeAssessmentMethodInclusive"
}
