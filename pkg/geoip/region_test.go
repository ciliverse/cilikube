package geoip

import "testing"

func TestParseRegion(t *testing.T) {
	legacy := parseRegion("中国|0|广东省|深圳市|电信")
	if legacy.Country != "中国" || legacy.Province != "广东省" || legacy.City != "深圳市" || legacy.ISP != "电信" {
		t.Fatalf("legacy fields: %+v", legacy)
	}
	if legacy.Label != "中国 广东省 深圳市" {
		t.Fatalf("legacy label=%q", legacy.Label)
	}

	v3 := parseRegion("中国|江苏省|南京市|0|CN")
	if v3.Country != "中国" || v3.Province != "江苏省" || v3.City != "南京市" {
		t.Fatalf("v3 fields: %+v", v3)
	}
	if v3.Label != "中国 江苏省 南京市" {
		t.Fatalf("v3 label=%q", v3.Label)
	}

	us := parseRegion("United States|California|0|Google LLC|US")
	if us.Country != "United States" || us.Province != "California" || us.ISP != "Google LLC" {
		t.Fatalf("us fields: %+v", us)
	}
	if us.Label != "United States California" {
		t.Fatalf("us label=%q", us.Label)
	}
}

func TestClassifyNonPublic(t *testing.T) {
	if loc := classifyNonPublic("127.0.0.1"); loc == nil || loc.Label != "本机回环" {
		t.Fatalf("loopback: %+v", loc)
	}
	if loc := classifyNonPublic("192.168.1.1"); loc == nil || loc.Label != "内网" {
		t.Fatalf("private: %+v", loc)
	}
}
