package gobgg

import "fmt"

const (
	spcPollName = "suggested_numplayers"
)

// ThreeRating is a rating type for Best/Recommended/Not Recommended
type ThreeRating int

const (
	// NotRecommended is not recommended
	NotRecommended ThreeRating = iota
	// Recommended
	Recommended
	// Best
	Best
)

func (t ThreeRating) String() string {
	switch t {
	case NotRecommended:
		return "Not Recommended"
	case Recommended:
		return "Recommended"
	case Best:
		return "Best"
	}

	return fmt.Sprint(int(t))
}

// SuggestedPlayerCount is a structure that shows the suggested player count based on user voting
type SuggestedPlayerCount struct {
	NumPlayers     string
	Best           int
	Recommended    int
	NotRecommended int
}

func percent(i1, i2, i3 int) float32 {
	sum := float32(i1 + i2 + i3)
	if sum <= 0 {
		return 0
	}
	return (float32(i1) / sum) * 100
}

func (sp *SuggestedPlayerCount) Suggestion() (ThreeRating, int, float32) {
	// In case of a tie, the not recommended wins, the recommended, then best
	if sp.Recommended >= sp.Best && sp.Recommended > sp.NotRecommended {
		return Recommended, sp.Recommended, percent(sp.Recommended, sp.Best, sp.NotRecommended)
	}

	if sp.Best > sp.Recommended && sp.Best > sp.NotRecommended {
		return Best, sp.Best, percent(sp.Best, sp.Recommended, sp.NotRecommended)
	}

	return NotRecommended, sp.NotRecommended, percent(sp.NotRecommended, sp.Best, sp.Recommended)
}

func (sp *SuggestedPlayerCount) BestPercentile() float32 {
	return percent(sp.Best, sp.Recommended, sp.NotRecommended)
}

func (sp *SuggestedPlayerCount) RecommendedPercentile() float32 {
	return percent(sp.Recommended, sp.Best, sp.NotRecommended)
}

func (sp *SuggestedPlayerCount) NotRecommendedPercentile() float32 {
	return percent(sp.NotRecommended, sp.Best, sp.Recommended)
}

func getSuggestedPoll(ps []PollStruct) ([]SuggestedPlayerCount, error) {
	var result []SuggestedPlayerCount
	for pi := range ps {
		if ps[pi].Name != spcPollName {
			continue
		}
		for it := range ps[pi].Results {
			item := SuggestedPlayerCount{
				NumPlayers: ps[pi].Results[it].Numplayers,
			}
			for _, single := range ps[pi].Results[it].Result {
				switch single.Value {
				case "Best":
					item.Best = single.Numvotes
				case "Not Recommended":
					item.NotRecommended = single.Numvotes
				case "Recommended":
					item.Recommended = single.Numvotes
				}
			}

			result = append(result, item)
		}
	}
	return result, nil
}
