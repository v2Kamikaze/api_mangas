package scraper

import "net/http"

// Scraper é uma interface para que implementa
// a interface de erro padrão.
type Scraper interface {
	GetPageTitles(int) map[string]string
	GetAllTitles(int) map[string]string
}

// MangaHost é um tipo que implementa a interface Scraper
type MangaHost struct {
	BaseURL string
	Client  *http.Client
}

// NewScraper retorna um ponteiro para MangaHost.
// baseURL é a url base para as páginas: "https://mangahosted.com/mangas/page/",
// client é um http.Client para ser usado nas requisições.
func NewScraper(baseURL string, client http.Client) *MangaHost {
	return &MangaHost{BaseURL: baseURL, Client: &client}
}
