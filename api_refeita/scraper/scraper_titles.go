package scraper

import (
	"fmt"
	"log"
	"math"
	"progressbar"
	"strconv"
	"time"

	gq "github.com/PuerkitoBio/goquery"
)

// GetPageTitles irá recuperar os títulos e links dos mangás de uma página do MangáHost.
func (mh *MangaHost) GetPageTitles(pageNumber int) map[string]string {
	pageTitles := make(map[string]string)
	res, err := mh.Client.Get(mh.BaseURL + strconv.Itoa(pageNumber))
	if err != nil {
		log.Fatalf("erro ao requisitar página: %+v", err.Error())
	}
	if res.StatusCode != 200 {
		log.Fatalf("página não encontrada. código: %s", res.Status)
	}
	defer res.Body.Close()
	doc, err := gq.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalf("os dados da resposta não podem ser convertidos para html. erro: %+v", err)
	}
	doc.Find("a.manga-block-title-link").Each(func(index int, s *gq.Selection) {
		link, linkExists := s.Attr("href")
		title, titleExists := s.Attr("title")
		if linkExists && titleExists {
			if _, contains := pageTitles[title]; title != "" && !contains {
				pageTitles[title] = link
			} else {
				title = s.Text()
				pageTitles[title] = link
			}
		}
	})
	return pageTitles
}

// GetAllTitles irá pegar os títulos de todos os mangás.
// Usa GetPageTitles para pegar os títulos em cada página.
// Implementa uma workerPool para requisitar concorrentemente cada página.
// Retorna um map contendo todos os títulos presentes em MangáHost.
func (mh *MangaHost) GetAllTitles(numPages int) map[string]string {
	// Iniciando a barra de progresso.
	bar := progressbar.NewBar(float64(numPages), "█")
	fmt.Println("Requisitando páginas...")
	bar.Init()
	// numWorkers irá definir o total de goroutines rodando baseado no número de páginas passadas.
	numWorkers := int(math.Ceil(float64(numPages/10 + 1)))
	allTitles := make(map[string]string)
	titlesChannel := make(chan map[string]string, numPages)
	sendPageChannel := make(chan int)
	start := time.Now()
	// Iniciando os workers
	for w := 0; w < numWorkers; w++ {
		go mh.workerPageTitles(sendPageChannel, titlesChannel)
	}

	// Enviando as páginas para os workers
	for p := 1; p <= numPages; p++ {
		sendPageChannel <- p
	}
	close(sendPageChannel)

	// Recebendo os maps com os títulos de cada página do MangáHost.
	// Os valores de cada map serão repassados ao map allTitles.
	for t := 0; t < numPages; t++ {
		bar.Increment()
		receivedTitles := <-titlesChannel
		for title, url := range receivedTitles {
			allTitles[title] = url
		}
	}
	close(titlesChannel)

	fmt.Printf("\nTodos os títulos requisitados\nDuração: %.2fs.\n", time.Since(start).Seconds())
	return allTitles
}

// workerPageTitles recebe o número das páginas pelo channel sendPageChannel, passando os resultados
// de GetPageTitles para o channel de recebimento titlesChannel.
// É usado para implementar uma workerPool.
func (mh *MangaHost) workerPageTitles(sendPageChannel <-chan int, titlesChannel chan<- map[string]string) {
	for pageNumber := range sendPageChannel {
		titlesChannel <- mh.GetPageTitles(pageNumber)
	}
}
