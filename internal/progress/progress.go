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
	if r.callback != nil {
		r.callback(step, percentage, message, chapter)
	}
	if r.channel != nil {
		p := Progress{Step: step, Percentage: percentage, Message: message, Chapter: chapter}
		select {
		case r.channel <- p:
		default:
			select {
			case <-r.channel:
			default:
			}
			select {
			case r.channel <- p:
			default:
			}
		}
	}
}

func (r *Reporter) Close() {
	if r.channel != nil {
		close(r.channel)
		r.channel = nil
	}
}
