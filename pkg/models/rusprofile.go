package models

import (
	"bytes"
	"golang.org/x/net/html"
	"io"
)

type Rusprofile struct {
	INN string `json:"inn"`
	KPP string `json:"kpp"`
	CEO string `json:"ceo"`
}

func NewRusprofileFromWEB(node *html.Node) (*Rusprofile, error) {
	re := &Rusprofile{}
	buf := &bytes.Buffer{}
	var f1 func(*html.Node)
	f1 = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for _, a := range n.Attr {
				if a.Val == "clip_inn" {
					w := io.Writer(buf)
					if err := html.Render(w, n.FirstChild); err != nil {
						buf.Reset()
						continue
					}
					re.INN = buf.String()
					buf.Reset()
				}
				if a.Val == "clip_kpp" {
					w := io.Writer(buf)
					if err := html.Render(w, n.FirstChild); err != nil {
						buf.Reset()
						continue
					}
					re.KPP = buf.String()
					buf.Reset()
				}
				if a.Val == "link-arrow gtm_main_fl" {
					w := io.Writer(buf)
					if err := html.Render(w, n.FirstChild.FirstChild); err != nil {
						buf.Reset()
						continue
					}
					re.CEO = buf.String()
					buf.Reset()
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f1(c)
		}
	}
	f1(node)
	return re, nil
}
