package core

var banks []string

// BankEnum represents an enumeration for accepted bank
type BankEnum int

func ciota(s string) BankEnum {
	banks = append(banks, s)
	return BankEnum(len(banks) - 1)
}

// CA, BNP, HSBC are accepted bank names
var (
	CA   = ciota("CA")
	BNP  = ciota("BNP")
	HSBC = ciota("HSBC")
)

// Bank represents a name of accepted bank
type Bank struct {
	X BankEnum
}

func (b BankEnum) String() string {
	return banks[int(b)]
}

// ExistsBank chech if given bank name is registered in our valid bank names
func ExistsBank(s string) bool {
	for _, v := range banks {
		if v == s {
			return true
		}
	}
	return false
}
