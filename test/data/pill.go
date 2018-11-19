package painkiller

//go:generate astgen -t ../../template/stringer.gogo -p $GOFILE -o pill_stringer.go

type Pill int // +tag stringer:""

const (
	Placebo Pill = iota
	Aspirin
	Ibuprofen
	Paracetamol
	Acetaminophen = Paracetamol
)

const (
	foo = 1
	bar = "string"
)
