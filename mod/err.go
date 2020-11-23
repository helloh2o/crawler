package mod

type Err struct {
	code     int
	reason   string
	callback string
}

func NewStringErr(err error) *Err {
	return &Err{
		code:     -1,
		reason:   err.Error(),
		callback: "",
	}
}
