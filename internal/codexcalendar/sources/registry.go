package sources

import "canarias.run/internal/codexcalendar"

func PriorityAdapters() []calendar.Adapter {
	return []calendar.Adapter{
		RunneaAdapter{},
		CarrerasPopularesGCAdapter{},
		CronolineCanariasAdapter{},
		BullTimingAdapter{},
		ConchipCanariasAdapter{},
		ProcronoAdapter{},
		AscensoTimingAdapter{},
	}
}
