package urlshort

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParseYAML(t *testing.T) {
	testGoodYaml := `
    - path: /test123
      url: http://test123
    - path: /test456
      url: http://test456
    `
	testBadYaml := `
    path: /test123
	      url: http://test123
    path: /test456
      url: http://test456
    `
	expected := []pathUrl{
		{
			URL:  "http://test123",
			Path: "/test123",
		},
		{
			URL:  "http://test456",
			Path: "/test456",
		},
	}
	Convey("Given yaml content", t, func() {
		Convey("When parsing good yaml", func() {
			parsedYAML, err := parseYAML([]byte(testGoodYaml))
			Convey("It should return slice config and nil err", func() {
				So(err, ShouldBeNil)
				So(parsedYAML, ShouldResemble, expected)
			})
		})

		Convey("When parsing bad yaml", func() {
			_, err := parseYAML([]byte(testBadYaml))
			Convey("It should return slice config and nil err", func() {
				So(err, ShouldNotBeNil)
			})
		})
	})
}
