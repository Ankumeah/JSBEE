package frontend

import (
	"context"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/Ankumeah/JSBEE/backend/internal/provider"
)

var seoSitemapRoutes = []struct {
	path     string
	priority string
	freq     string
}{
	{"/", "1.0", "weekly"},
	{"/volumes", "0.9", "weekly"},
	{"/blog", "0.9", "weekly"},
	{"/paper", "0.8", "weekly"},
	{"/about", "0.8", "monthly"},
	{"/author", "0.7", "weekly"},
	{"/team", "0.8", "monthly"},
	{"/contact", "0.7", "monthly"},
}

func seoRobotsTxt() string {
	return "User-agent: *\n" +
		"Allow: /\n" +
		"Disallow: /api/\n" +
		"Disallow: /profile\n" +
		"Disallow: /review\n" +
		"Disallow: /admin\n" +
		"\n" +
		"Sitemap: " + provider.SiteURL + "/sitemap.xml\n"
}

func seoSitemapXML() string {
	date := time.Now().UTC().Format("2006-01-02")
	var sb strings.Builder
	sb.WriteString(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")
	sb.WriteString(`<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">` + "\n")
	for _, r := range seoSitemapRoutes {
		fmt.Fprintf(&sb,
			"  <url>\n    <loc>%s%s</loc>\n    <lastmod>%s</lastmod>\n    <changefreq>%s</changefreq>\n    <priority>%s</priority>\n  </url>\n",
			provider.SiteURL, r.path, date, r.freq, r.priority,
		)
	}
	sb.WriteString("</urlset>\n")
	return sb.String()
}

func (
	u *ComponentUpdater,
) UpdateSEO(_ context.Context) error {
	if err := os.WriteFile(
		path.Join(u.savePath, "robots.txt"),
		[]byte(seoRobotsTxt()), 0o644,
	); err != nil {
		return err
	}

	return os.WriteFile(
		path.Join(u.savePath, "sitemap.xml"),
		[]byte(seoSitemapXML()), 0o644,
	)
}
