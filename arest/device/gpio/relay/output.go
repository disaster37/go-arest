package relay

const (
	nc string = "nc"
	no string = "no"
)

// Output represent the relay output type (normaly close or normaly open)
type Output interface {
	// Output return the output type
	Output() (output string)

	// SetOutputNO set the output as Normaly Open
	SetOutputNO()

	// SetOutputNO set the output as Normaly Close
	SetOutputNC()

	// IsNO return true if NO
	IsNO() bool

	// IsNC return true if NC
	IsNC() bool
}

// OutputImp implement output interface
type OutputImp struct {
	output string
}

// NewOutput return new output object
func NewOutput() (output Output) {
	output = &OutputImp{}
	return output
}

// SetOutputNO set the output as Normaly Open
func (o *OutputImp) SetOutputNO() {
	o.output = no
}

// SetOutputNO set the output as Normaly Close
func (o *OutputImp) SetOutputNC() {
	o.output = nc
}

// Output return the output type
func (o *OutputImp) Output() string {
	return o.output
}

// IsNO return true if NO
func (o *OutputImp) IsNO() bool {
	return o.output == no
}

// IsNC return true if NC
func (o *OutputImp) IsNC() bool {
	return o.output == nc
}
