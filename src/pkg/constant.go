package pkg

const (
	H_LANG            = "z-lang"
	H_USERID          = "z-userid"
	H_CURRENCY        = "z-currency"
	H_VISIBILITY      = "z-visibility-level"
	H_XTIMEZONE       = "x-timezone"
	H_XEARN_MAINTANCE = "X-Earn-Maintenance"
	H_FEATURE_LEVEL   = "z-feature-level"
	CURRENCY_USD      = "USD"
	KEY               = "75e9a69672164793937ceba23aa5f5e1"
)

const (
	SECRETMACAROON = "abcd"
)

var AllowIp []string = []string{"192.168.1.100", "192.168.1.101", "192.168.1.102"}

var WhitelistCIDRs = []string{
	"127.0.0.1/32",   // localhost
	"192.168.1.0/24", // LAN
	"203.0.113.5/32", // IP publik tunggal
	"::1/128",
}
