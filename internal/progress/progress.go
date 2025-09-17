package progress

type Progress struct {
	Step       string  `json:"step"`
	Percentage float64 `json:"percentage"`
	Message    string  `json:"message"`
	Chapter    int     `json:"chapter"`
}

type Reporter struct {
	channel  chan Progress
	callback func(step string, percentage float64, message string, chapter int)
}

func NewReporter() *Reporter {
	return &Reporter{}
}

func (r *Reporter) SetChannel(ch chan Progress) {
	r.channel = ch
}

func (r *Reporter) SetCallback(callback func(step string, percentage float64, message string, chapter int)) {
	r.callback = callback
}

func (r *Reporter) Update(step string, percentage float64, message string, chapter int) {
	if r.channel != nil {
		select {
		case r.channel <- Progress{
			Step:       step,
			Percentage: percentage,
			Message:    message,
			Chapter:    chapter,
		}:
		default:
		}
	}
}

func (r *Reporter) Close() {
	if r.channel != nil {
		close(r.channel)
		r.channel = nil
	}
}
