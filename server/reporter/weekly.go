package reporter

// WeekReport contains details about elements to construct a report.
type WeekReport struct {
	StartDate         string
	EndDate           string
	Watch             int
	Star              int
	Fork              int
	ContributorsCount int
	PullReuqestCount  int
}

func (r *Reporter) weeklyReport() error {
	if _, err := r.construcWeekReport(); err != nil {
		return err
	}

	return nil
}

func (r *Reporter) construcWeekReport() (WeekReport, error) {
	var wr WeekReport
	repo, err := r.client.GetRepository()
	if err != nil {
		return wr, err
	}

	wr.Watch = *(repo.WatchersCount)
	wr.Star = *(repo.StargazersCount)
	wr.Fork = *(repo.ForksCount)

	return wr, nil
}
