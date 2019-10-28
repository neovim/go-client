package msgpack

import (
	"reflect"
	"testing"
	"time"
)

type taggedStruct struct {
	Int      int         `msgpack:"int"`
	Int8     int8        `msgpack:"int8"`
	Int16    int16       `msgpack:"int16"`
	Int32    int32       `msgpack:"int32"`
	Int64    int64       `msgpack:"int64"`
	String   string      `msgpack:"string"`
	Bool     bool        `msgpack:"bool"`
	SString  []string    `msgpack:"sstring"`
	SInt     []int       `msgpack:"sint"`
	SInt8    []int8      `msgpack:"sint8"`
	SInt16   []int16     `msgpack:"sint16"`
	SInt32   []int32     `msgpack:"sint32"`
	SInt64   []int64     `msgpack:"sint64"`
	SFloat32 []float32   `msgpack:"sfloat32"`
	SFloat64 []float64   `msgpack:"sfloat64"`
	SBool    []bool      `msgpack:"sbool"`
	Time     time.Time   `msgpack:"time"`
	Stime    []time.Time `msgpack:"stime"`
	Struct   struct {
		Number        int64
		Height        int64
		AnotherStruct struct {
			Image string
		} `msgpack:",array"`
	} `msgpack:"struct"`
	Latitude           float32 `msgpack:"lat"`
	Longitude          float32 `msgpack:"long"`
	CreditCardNumber   string  `msgpack:"cc_number"`
	CreditCardType     string  `msgpack:"cc_type"`
	Email              string  `msgpack:"email"`
	DomainName         string  `msgpack:"domain_name"`
	IPV4               string  `msgpack:"ipv4"`
	IPV6               string  `msgpack:"ipv6"`
	Password           string  `msgpack:"password"`
	PhoneNumber        string  `msgpack:"phone_number"`
	MacAddress         string  `msgpack:"mac_address"`
	URL                string  `msgpack:"url"`
	UserName           string  `msgpack:"username"`
	TollFreeNumber     string  `msgpack:"toll_free_number"`
	E164PhoneNumber    string  `msgpack:"e_164_phone_number"`
	TitleMale          string  `msgpack:"title_male"`
	TitleFemale        string  `msgpack:"title_female"`
	FirstName          string  `msgpack:"first_name"`
	FirstNameMale      string  `msgpack:"first_name_male"`
	FirstNameFemale    string  `msgpack:"first_name_female"`
	LastName           string  `msgpack:"last_name"`
	Name               string  `msgpack:"name"`
	UnixTime           int64   `msgpack:"unix_time"`
	Date               string  `msgpack:"date"`
	MonthName          string  `msgpack:"month_name"`
	Year               string  `msgpack:"year"`
	DayOfWeek          string  `msgpack:"day_of_week"`
	DayOfMonth         string  `msgpack:"day_of_month"`
	Timestamp          string  `msgpack:"timestamp"`
	Century            string  `msgpack:"century"`
	TimeZone           string  `msgpack:"timezone"`
	TimePeriod         string  `msgpack:"time_period"`
	Word               string  `msgpack:"word"`
	Sentence           string  `msgpack:"sentence"`
	Paragraph          string  `msgpack:"paragraph"`
	Currency           string  `msgpack:"currency"`
	Amount             float32 `msgpack:"amount"`
	AmountWithCurrency string  `msgpack:"amount_with_currency"`
	ID                 string  `msgpack:"uuid_digit"`
	HyphenatedID       string  `msgpack:"uuid_hyphenated"`
}

type omitemptyStruct struct {
	Int      int         `msgpack:"int,omitempty"`
	Int8     int8        `msgpack:"int8,omitempty"`
	Int16    int16       `msgpack:"int16,omitempty"`
	Int32    int32       `msgpack:"int32,omitempty"`
	Int64    int64       `msgpack:"int64,omitempty"`
	String   string      `msgpack:"string"`
	Bool     bool        `msgpack:"bool,omitempty"`
	SString  []string    `msgpack:"sstring,omitempty"`
	SInt     []int       `msgpack:"sint,omitempty"`
	SInt8    []int8      `msgpack:"sint8,omitempty"`
	SInt16   []int16     `msgpack:"sint16,omitempty"`
	SInt32   []int32     `msgpack:"sint32,omitempty"`
	SInt64   []int64     `msgpack:"sint64,omitempty"`
	SFloat32 []float32   `msgpack:"sfloat32,omitempty"`
	SFloat64 []float64   `msgpack:"sfloat64,omitempty"`
	SBool    []bool      `msgpack:"sbool,omitempty"`
	Time     time.Time   `msgpack:"time,omitempty"`
	Stime    []time.Time `msgpack:"stime,omitempty"`
	Struct   struct {
		Number        int64
		Height        int64
		AnotherStruct struct {
			Image string
		} `msgpack:",array,omitempty"`
	} `msgpack:"struct,omitempty"`
	Latitude           float32 `msgpack:"lat,omitempty"`
	Longitude          float32 `msgpack:"long,omitempty"`
	CreditCardNumber   string  `msgpack:"cc_number,omitempty"`
	CreditCardType     string  `msgpack:"cc_type,omitempty"`
	Email              string  `msgpack:"email,omitempty"`
	DomainName         string  `msgpack:"domain_name,omitempty"`
	IPV4               string  `msgpack:"ipv4,omitempty"`
	IPV6               string  `msgpack:"ipv6,omitempty"`
	Password           string  `msgpack:"password,omitempty"`
	PhoneNumber        string  `msgpack:"phone_number,omitempty"`
	MacAddress         string  `msgpack:"mac_address,omitempty"`
	URL                string  `msgpack:"url,omitempty"`
	UserName           string  `msgpack:"username,omitempty"`
	TollFreeNumber     string  `msgpack:"toll_free_number,omitempty"`
	E164PhoneNumber    string  `msgpack:"e_164_phone_number,omitempty"`
	TitleMale          string  `msgpack:"title_male,omitempty"`
	TitleFemale        string  `msgpack:"title_female,omitempty"`
	FirstName          string  `msgpack:"first_name,omitempty"`
	FirstNameMale      string  `msgpack:"first_name_male,omitempty"`
	FirstNameFemale    string  `msgpack:"first_name_female,omitempty"`
	LastName           string  `msgpack:"last_name,omitempty"`
	Name               string  `msgpack:"name,omitempty"`
	UnixTime           int64   `msgpack:"unix_time,omitempty"`
	Date               string  `msgpack:"date,omitempty"`
	MonthName          string  `msgpack:"month_name,omitempty"`
	Year               string  `msgpack:"year,omitempty"`
	DayOfWeek          string  `msgpack:"day_of_week,omitempty"`
	DayOfMonth         string  `msgpack:"day_of_month,omitempty"`
	Timestamp          string  `msgpack:"timestamp,omitempty"`
	Century            string  `msgpack:"century,omitempty"`
	TimeZone           string  `msgpack:"timezone,omitempty"`
	TimePeriod         string  `msgpack:"time_period,omitempty"`
	Word               string  `msgpack:"word,omitempty"`
	Sentence           string  `msgpack:"sentence,omitempty"`
	Paragraph          string  `msgpack:"paragraph,omitempty"`
	Currency           string  `msgpack:"currency,omitempty"`
	Amount             float32 `msgpack:"amount,omitempty"`
	AmountWithCurrency string  `msgpack:"amount_with_currency,omitempty"`
	ID                 string  `msgpack:"uuid_digit,omitempty"`
	HyphenatedID       string  `msgpack:"uuid_hyphenated,omitempty"`
}

func Benchmark_correctFilels(b *testing.B) {
	b.ReportAllocs()

	b.Run("named", func(b *testing.B) {
		t := reflect.ValueOf(taggedStruct{}).Type()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = collectFields(nil, t, make(map[reflect.Type]bool), make(map[string]int), nil)
		}
	})

	b.Run("omitempty", func(b *testing.B) {
		t := reflect.ValueOf(omitemptyStruct{}).Type()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_ = collectFields(nil, t, make(map[reflect.Type]bool), make(map[string]int), nil)
		}
	})
}
