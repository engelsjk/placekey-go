package placekey

import (
	"math"
	"testing"
)

func TestToGeo(t *testing.T) {
	ToGeo("@dvt-smp-tvz")
}

func TestToH3(t *testing.T) {
	got := ToH3("@dvt-smp-tvz")
	if got != "8a754e64992ffff" {
		t.Errorf(`ToH3("@dvt-smp-tvz") = "%s"; wanted "8a754e64992ffff"`, got)
	}
}

func TestFromH3(t *testing.T) {
	got := FromH3("8a754e64992ffff")
	if got != "@dvt-smp-tvz" {
		t.Errorf(`FromH3("8a754e64992ffff") = "%s"; wanted "@dvt-smp-tvz"`, got)
	}
}

func TestDistance(t *testing.T) {
	got := Distance("@dvt-smp-tvz", "@5vg-7gq-tjv")
	tolerance := 0.001
	if diff := math.Abs(got - 12795124.895573696); diff > tolerance {
		t.Errorf(`Distance("@dvt-smp-tvz", "@5vg-7gq-tjv") = %f; exceeds %f tolerance`, got, tolerance)
	}
}

func TestFormatIsValid(t *testing.T) {
	got := FormatIsValid("222-227@dvt-smp-tvz")
	if got != true {
		t.Errorf(`FormatIsValid("222-227@dvt-smp-tvz") = %t; wanted true`, got)
	}
	got = FormatIsValid("@123-456-789")
	if got != false {
		t.Errorf(`FormatIsValid("@123-456-789") = %t; wanted false`, got)
	}
}

func TestFromGeo(t *testing.T) {
	got := FromGeo(37.7371, -122.44283)
	if got != "@5vg-82n-kzz" {
		t.Errorf(`FromGeo(37.7371, -122.44283) = "%s"; wanted "@5vg-82n-kzz"`, got)
	}
}
