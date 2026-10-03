package components

import (
	"context"
	"strings"
	"testing"

	"github.com/a-h/templ"
)

// Smoke test: every public page must render a complete document.
func TestPagesRender(t *testing.T) {
	cfg := map[string]any{"apiKey": "x"}
	pages := map[string]templ.Component{
		"index":   IndexPage(cfg),
		"volumes": VolumesPage(cfg),
		"paper":   PaperPage(cfg),
		"blog":    BlogPage(cfg),
		"about":   AboutPage(cfg, "https://example.com/about.md"),
		"author":  AuthorPage(cfg),
		"team":    TeamPage(cfg),
		"contact": ContactPage(cfg),
		"profile": ProfilePage(cfg),
		"review":  ReviewDashboardPage(cfg),
		"admin":   AdminDashboardPage(cfg, "https://example.com/about.md"),
		"404":     NotFoundPage(cfg),
	}
	for name, p := range pages {
		var sb strings.Builder
		if err := p.Render(context.Background(), &sb); err != nil {
			t.Fatalf("%s render failed: %v", name, err)
		}
		html := strings.ToLower(sb.String())
		if !strings.Contains(html, "<!doctype html>") || !strings.Contains(html, "</html>") {
			t.Errorf("%s: incomplete document", name)
		}
	}

	var idx strings.Builder
	if err := IndexPage(cfg).Render(context.Background(), &idx); err != nil {
		t.Fatalf("index render failed: %v", err)
	}
	for _, id := range []string{"stats", "featured", "latest", "researchers", "publish", "stat-papers"} {
		if !strings.Contains(idx.String(), `id="`+id+`"`) {
			t.Errorf("index: missing section %s", id)
		}
	}
}
