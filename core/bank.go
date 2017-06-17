package core

var banks []string

// BankEnum represents an enumeration for accepted bank
type BankEnum int

func ciota(s string) BankEnum {
	banks = append(banks, s)
	return BankEnum(len(banks) - 1)
}

var (
	CA   = ciota("CA")
	BNP  = ciota("BNP")
	HSBC = ciota("HSBC")
)

type Bank struct {
	X BankEnum
}

func (b BankEnum) String() string {
	return banks[int(b)]
}

func ExistsBank(s string) bool {
	for _, v := range banks {
		if v == s {
			return true
		}
	}
	return false
}
