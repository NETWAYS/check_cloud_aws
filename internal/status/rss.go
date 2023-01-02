package status

type Rss struct {
	Channel struct {
		Title string `xml:"title"`
		Link  string `xml:"link"`
		Desc  string `xml:"description"`
		Items []Item `xml:"item"`
	} `xml:"channel"`
}

type Item struct {
	Title string `xml:"title"`
}
