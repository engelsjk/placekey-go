package placekey

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/uber/h3-go"
)

var (
	resolution       int    = 10
	baseResolution   int    = 12
	alphabet         string = "23456789BCDFGHJKMNPQRSTVWXYZ"
	codeLength       int    = 9
	tupleLength      int    = 3
	paddingChar      string = "a"
	replacementChars string = "eu"
	replacementMap          = map[string]string{
		"prn":   "pre",
		"f4nny": "f4nne",
		"tw4t":  "tw4e",
		"ngr":   "ngu",
		"dck":   "dce",
		"vjn":   "vju",
		"fck":   "fce",
		"pns":   "pne",
		"sht":   "she",
		"kkk":   "kke",
		"fgt":   "fgu",
		"dyk":   "dye",
		"bch":   "bce",
	}
	alphabetLength         int
	headerBits             string
	headerInt              uint64
	baseCellShift          uint64
	unusedResolutionFiller uint64
	firstTupleRegex        string
	tupleRegex             string
	whereRegex             *regexp.Regexp
	whatRegex              *regexp.Regexp
)

func init() {
	alphabet = strings.ToLower(alphabet)
	alphabetLength = len(alphabet)
	headerBits = fmt.Sprintf("%064s", strconv.FormatUint(uint64(h3.FromGeo(h3.GeoCoord{0.0, 0.0}, resolution)), 2))[:12]
	baseCellShift = 1 << (3 * 15)
	unusedResolutionFiller = 1<<(3*(15-baseResolution)) - 1
	firstTupleRegex = "[" + alphabet + replacementChars + paddingChar + "]{3}"
	tupleRegex = "[" + alphabet + replacementChars + "]{3}"
	whereRegex = regexp.MustCompile("^" + strings.Join([]string{firstTupleRegex, tupleRegex, tupleRegex}, "-") + "$")
	whatRegex = regexp.MustCompile("^[" + alphabet + "]{3}(-[" + alphabet + "]{3})?$")
}

func getHeaderInt() uint64 {
	i, err := strconv.ParseInt(headerBits, 2, 64)
	if err != nil {
		panic(err)
	}
	headerInt := uint64(i) * 1 << 52
	return headerInt
}

func init() {
	headerInt = getHeaderInt()
}

// FromGeo converts a (latitude, longitude) into a Placekey.
func FromGeo(lat, lon float64) string {
	return encodeH3Int(uint64(h3.FromGeo(h3.GeoCoord{lat, lon}, resolution)))
}

// ToGeo converts a Placekey into a (latitude, longitude).
func ToGeo(placekey string) (float64, float64) {
	geo := h3.ToGeo(h3.H3Index(ToH3Int(placekey)))
	return geo.Latitude, geo.Longitude
}

// ToH3 converts a Placekey string into an H3 string.
func ToH3(placekey string) string {
	_, where := parsePlacekey(placekey)
	return h3.ToString(h3.H3Index(decodeToH3Int(where)))
}

// FromH3 converts an H3 hexadecimal string into a Placekey string.
func FromH3(h3String string) string {
	return encodeH3Int(uint64(h3.FromString(h3String)))
}

// GetPrefixDistanceMap returns a map of the length of a shared Placekey prefix to the
// maximal distance in meters between two Placekeys sharing a prefix of that length.
func GetPrefixDistanceMap() map[int]float64 {
	return map[int]float64{
		0: 2.004e7,
		1: 2.004e7,
		2: 2.777e6,
		3: 1.065e6,
		4: 1.524e5,
		5: 2.177e4,
		6: 8227.0,
		7: 1176.0,
		8: 444.3,
		9: 63.47,
	}
}

// H3IntToPlacekey converts an H3 integer into a Placekey.
func FromH3Int(h3Int uint64) string {
	return encodeH3Int(h3Int)
}

// ToH3Int converts a Placekey to an H3 integer.
func ToH3Int(placekey string) uint64 {
	_, where := parsePlacekey(placekey)
	return decodeToH3Int(where)
}

// todo: ToHexBoundary
// todo: ToPolygon
// todo: ToWKT
// todo: ToGeojson
// todo: FromPolygon
// todo: FromWKT
// todo: FromGeoJSON

// FormatIsValid returns a oolean for whether or not the format of a Placekey is valid, including
// checks for valid encoding of location.
func FormatIsValid(placekey string) bool {
	what, where := parsePlacekey(placekey)
	if what != "" {
		return wherePartIsValid(where) && whatRegex.Match([]byte(what))
	}
	return wherePartIsValid(where)
}

// Distance returns the distance in meters between the centers of two Placekeys.
func Distance(placekey1, placekey2 string) float64 {
	geo1 := h3.ToGeo(h3.H3Index(ToH3Int(placekey1)))
	geo2 := h3.ToGeo(h3.H3Index(ToH3Int(placekey2)))
	return geoDistance(geo1, geo2)
}

///////////////////////////////////////////////////
///////////////////////////////////////////////////

// split a Placekey in to what and where parts.
func parsePlacekey(placekey string) (string, string) {
	if strings.Contains(placekey, "@") {
		whatwhere := strings.Split(placekey, "@")
		return whatwhere[0], whatwhere[1]
	}
	return "", placekey
}

func wherePartIsValid(where string) bool {
	return whereRegex.Match([]byte(where)) && h3.IsValid(h3.H3Index(ToH3Int(where)))
}

func geoDistance(geo1, geo2 h3.GeoCoord) float64 {
	earthRadius := 6371.0 // km

	lat1 := rad(geo1.Latitude)
	lon1 := rad(geo1.Longitude)
	lat2 := rad(geo2.Latitude)
	lon2 := rad(geo2.Longitude)

	havLat := 0.5 * (1 - math.Cos(lat1-lat2))
	havLon := 0.5 * (1 - math.Cos(lon1-lon2))
	radical := math.Sqrt(havLat + math.Cos(lat1)*math.Cos(lat2)*havLon)
	return 2.0 * earthRadius * math.Asin(radical) * 1000
}

// shorten an H3 integer to only include location data up to the base resolution
func encodeH3Int(h3Int uint64) string {
	shortH3Int := shortenH3Int(h3Int)
	encodedShortH3 := encodeShortInt(int(shortH3Int))

	cleanEncodedShortH3 := cleanString(encodedShortH3)
	if len(cleanEncodedShortH3) <= codeLength {
		cleanEncodedShortH3 = strings.Repeat(paddingChar, codeLength-len(cleanEncodedShortH3)) + cleanEncodedShortH3
	}

	tuples := []string{}
	for i := 0; i < len(cleanEncodedShortH3); i += tupleLength {
		tuples = append(tuples, cleanEncodedShortH3[i:i+tupleLength])
	}

	return "@" + strings.Join(tuples, "-")
}

func encodeShortInt(x int) string {
	if x == 0 {
		return string(alphabet[0])
	}
	res := ``
	for x > 0 {
		remainder := x % alphabetLength
		res = string(alphabet[remainder]) + res
		x = x / alphabetLength
	}
	return res
}

func decodeToH3Int(wherePart string) uint64 {
	code := stripEncoding(wherePart)
	dirtyEncoding := dirtyString(code)
	shortH3Int := decodeString(dirtyEncoding)
	return unshortenH3Int(shortH3Int)
}

func decodeString(s string) uint64 {
	var val uint64
	for i := len(s) - 1; i >= 0; i-- {
		val += power64(uint64(alphabetLength), len(s)-1-i) * uint64(strings.Index(alphabet, string(s[i])))
	}
	return val
}

func stripEncoding(s string) string {
	s = strings.ReplaceAll(s, "@", "")
	s = strings.ReplaceAll(s, "-", "")
	s = strings.ReplaceAll(s, paddingChar, "")
	return s
}

func shortenH3Int(h3Int uint64) uint64 {
	// Cuts off the 12 left-most bits that don't code location
	out := (h3Int + baseCellShift) % (1 << 52)
	// Cuts off the rightmost bits corresponding to resolutions greater than the base resolution
	out = out >> (3 * (15 - baseResolution))
	return out
}

func unshortenH3Int(shortH3Int uint64) uint64 {
	unshiftedInt := shortH3Int << (3 * (15 - baseResolution))
	rebuiltInt := headerInt + unusedResolutionFiller - baseCellShift + unshiftedInt
	return rebuiltInt
}

func cleanString(s string) string {
	// Replacement should be in order
	for k, v := range replacementMap {
		if strings.ContainsAny(s, k) {
			s = strings.ReplaceAll(s, k, v)
		}
	}
	return s
}

func dirtyString(s string) string {
	// Replacement should be in (reversed) order
	for k, v := range replacementMap {
		if strings.ContainsAny(s, v) {
			s = strings.ReplaceAll(s, v, k)
		}
	}
	return s
}

func power64(base uint64, exponent int) uint64 {
	if exponent == 0 {
		return 1
	}
	return (base * power64(base, exponent-1))
}

///////////////////////////////////////////////////
///////////////////////////////////////////////////

func rad(d float64) float64 {
	return d * math.Pi / 180
}

func deg(r float64) float64 {
	return r / math.Pi * 180
}
