package api

func DefaultJobProfile() JobProfile {

	return defaultJobProfile{}
}

type defaultJobProfile struct {
}

func (p defaultJobProfile) CanAdd() bool {
	return true
}
