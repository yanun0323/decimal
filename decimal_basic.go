package decimal

// Predefined basic values for each decimal size.
var (
	Zero128    = Decimal128{}
	One128     = New128FromInt(1)
	Ten128     = New128FromInt(10)
	Hundred128 = New128FromInt(100)

	Zero256    = Decimal256{}
	One256     = New256FromInt(1)
	Ten256     = New256FromInt(10)
	Hundred256 = New256FromInt(100)

	Zero512    = Decimal512{}
	One512     = New512FromInt(1)
	Ten512     = New512FromInt(10)
	Hundred512 = New512FromInt(100)
)
