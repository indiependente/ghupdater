package main

type Release struct {
	Assets  []Asset `json:"assets"`
	TagName string  `json:"tag_name"`
}

type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	ContentType        string `json:"content_type"`
	Size               int64  `json:"size"`
}
