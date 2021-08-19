package hedera

type FeeAssessmentMethod bool

const (
	FeeAssessmentMethodInclusive FeeAssessmentMethod = false
	FeeAssessmentMethodExclusive FeeAssessmentMethod = true
)

func (assessment FeeAssessmentMethod) String() string {
	if assessment == true {
		return "FeeAssessmentMethodExclusive"
	}

	return "FeeAssessmentMethodInclusive"
}
