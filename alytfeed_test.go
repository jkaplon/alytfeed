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
	// Read in entire alytfeed.xml file.
	xmlFile, err := os.Open("alytfeed.xml")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer xmlFile.Close()

	XMLdata, _ := ioutil.ReadAll(xmlFile)

	var i RSS
	xml.Unmarshal(XMLdata, &i)

	for _, item := range i.Items.ItemList {

		// Validate RFC 2822 pubdate values according to big, hairy regex.
		var validPubDate = regexp.MustCompile(`^(?:(Sun|Mon|Tue|Wed|Thu|Fri|Sat),\s+)?(0[1-9]|[1-2]?[0-9]|3[01])\s+(Jan|Feb|Mar|Apr|May|Jun|Jul|Aug|Sep|Oct|Nov|Dec)\s+(19[0-9]{2}|[2-9][0-9]{3})\s+(2[0-3]|[0-1][0-9]):([0-5][0-9])(?::(60|[0-5][0-9]))?\s+([-\+][0-9]{2}[0-5][0-9]|(?:UT|GMT|(?:E|C|M|P)(?:ST|DT)|[A-IK-Z]))(\s+|\(([^\(\)]+|\\\(|\\\))*\))*$`)
		matched := validPubDate.MatchString(item.PubDate)
		if matched == false {
			t.Errorf("Bad podcaster/human pubDate value non-conformant w/RFC 2822: %q", item.PubDate)
		}

		// All enclosure links are HTTP (as opposed to HTTPS), pointing to an MP3 file.
		// MP3 format check is brittle if we ever do an OGG feed or video feed.
		// Assume only one enclosure element per item-node and hard-code 0th element.
		var validUrl = regexp.MustCompile(`^http://*`)
		matchedHttp := validUrl.FindString(item.EnclosureList[0].Url)
		if matchedHttp == "" {
			t.Errorf("Bad podcaster/human, enclosure URL's must use http instead of https: %q", item.EnclosureList[0].Url)
		}

		// All enclosure links point to an MP3 file (MP3 format check is brittle if we ever do an OGG feed or video feed).
		// Tried to combine this with HTTP check, but my regex skills were not up to task.
		// Assume only one enclosure element per item-node and hard-code 0th element.
		var validUrlMp3 = regexp.MustCompile(`mp3$`)
		matchedMp3 := validUrlMp3.FindString(item.EnclosureList[0].Url)
		if matchedMp3 == "" {
			t.Errorf("Bad podcaster/human, enclosure URL's must point to an MP3: %q", item.EnclosureList[0].Url)
		}

		// Make sure all enclosure links are valid URL's...I'm declaring this out of scope for now
		// (if it starts with 'http://' and ends with '.mp3', call it good).

		// Check for iframe included in <description> tag
		var iframe = regexp.MustCompile(`iframe`)
		matchedIframe := iframe.MatchString(item.Description)
		if matchedIframe == true {
			t.Errorf("Bad human, naive copy/paste from show notes blog post is discouraged, please remove iframe from description tag: %q", item.Description)
		}
	}
}
