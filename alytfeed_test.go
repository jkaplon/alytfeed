package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"testing"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Items   Items    `xml:"channel"`
}
type Items struct {
	XMLName  xml.Name `xml:"channel"`
	ItemList []Item   `xml:"item"`
}
type Item struct {
	Title         string      `xml:"title"`
	Link          string      `xml:"link"`
	Description   string      `xml:"description"`
	PubDate       string      `xml:"pubDate"`
	EnclosureList []Enclosure `xml:"enclosure"`
}
type Enclosure struct {
	Url string `xml:"url,attr"`
}

func TestAlytfeed(t *testing.T) {
	// read in entire alytfeed.xml file...may need another package import.
	xmlFile, err := os.Open("alytfeed.xml")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer xmlFile.Close()

	XMLdata, _ := ioutil.ReadAll(xmlFile)
	//fmt.Printf(XMLdata)

	var i RSS
	xml.Unmarshal(XMLdata, &i)
	//fmt.Printf("%#v", i)

	for _, item := range i.Items.ItemList {
		//fmt.Printf("\t%d: %s\n", c, item.Title)
		//fmt.Printf("\t%d: %s\n", c, item.PubDate)

		// Validate RFC 2822 pubdate values according to big, hairy regex.
		var validPubDate = regexp.MustCompile(`^(?:(Sun|Mon|Tue|Wed|Thu|Fri|Sat),\s+)?(0[1-9]|[1-2]?[0-9]|3[01])\s+(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)\s+(19[0-9]{2}|[2-9][0-9]{3})\s+(2[0-3]|[0-1][0-9]):([0-5][0-9])(?::(60|[0-5][0-9]))?\s+([-\+][0-9]{2}[0-5][0-9]|(?:UT|GMT|(?:E|C|M|P)(?:ST|DT)|[A-IK-Z]))(\s+|\(([^\(\)]+|\\\(|\\\))*\))*$`)
		matched := validPubDate.MatchString(item.PubDate)
		if matched == false {
			t.Errorf("Bad podcaster/human pubDate value non-conformant w/RFC 2822: %q", item.PubDate)
		}

		// All enclosure links are HTTP (as opposed to HTTPS).
		// TODO, fix this...item.Url coming back empty.
		//var validUrl = regexp.MustCompile(`^http://*`)
		//matchedHttp := validUrl.FindString(item.Url)
		//fmt.Printf("%q\n", item.Url)
		//if matchedHttp == "" {
		//fmt.Printf("\t%d: %s\n", c, item.Url)
		//t.Errorf("Bad Jody, use http instead of https: %q", item.Url)
		//}

		// TODO, add these also:
		// Make sure all enclosure links are valid URL's
		// All enclosure links are to MP3 files (brittle if we ever do an OGG feed or video feed).
	}

}
