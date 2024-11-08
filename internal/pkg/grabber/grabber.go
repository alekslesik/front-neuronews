package grabber

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/alekslesik/kandinsky"
	"github.com/alekslesik/neuro-news/internal/app/model"
	"github.com/alekslesik/neuro-news/pkg/config"
	"github.com/alekslesik/neuro-news/pkg/logger"

	"golang.org/x/net/html"
)

// Grabber struct
type Grabber struct {
	log  *logger.Logger
	cfg  *config.Config
	home string
}

// Tag struct
type Tag struct {
	Name string
	Key  string
	Val  string
}

// New return new instance of Grabber struct
func New(log *logger.Logger, cfg *config.Config, home string) *Grabber {
	return &Grabber{log: log, cfg: cfg, home: home}
}

// GetGeneratedImage generate, save and return image model
func (g *Grabber) GetGeneratedImage(a *model.Article) (*model.Image, error) {
	const op = "grabber.GetGeneratedImage()"
	var imagePath = "website/static/upload/"

	title := a.Title

	params := kandinsky.Params{
		Width:  1024,
		Height: 680,
		Style:  "UHD",
		GenerateParams: struct {
			Query string "json:\"query\""
		}{title},
	}

	image, err := kandinsky.GetImage(g.cfg.Kand.Key, g.cfg.Kand.Secret, params)
	if err != nil {
		g.log.Error().Msgf("%s: get image from Kandinsky API error > %s", op, err)
		return nil, err
	}

	fImage, err := image.ToFile()
	if err != nil {
		g.log.Error().Msgf("%s: convert image generated from Kandinsky to os.File error > %s", op, err)
		return nil, err
	}

	size, err := getFileSize(*fImage)
	if err != nil {
		g.log.Error().Msgf("%s: get Kandinsky image size error > %s", op, err)
		return nil, err
	}

	imageName := translit(title)
	err = image.SavePNGTo(imageName, imagePath)
	if err != nil {
		g.log.Error().Msgf("%s: save image generated from Kandinsky error > %s", op, err)
		return nil, err
	}

	preparedPath := prepareImagePath(imagePath, imageName)

	model := &model.Image{
		ImagePath: preparedPath,
		Size:      size,
		Name:      imageName,
		Alt:       title,
	}

	return model, nil
}

// GenerateImageFruity generate, save and return image model trough Fruity API
func (g *Grabber) GenerateImageFruity(a *model.Article) (*model.Image, error) {
	const op = "grabber.GenerateImageFruity()"
	var imagePath = "website/static/upload/"

	title := a.Title

	image, err := getFruityImage(title)
	if err != nil {
		g.log.Error().Msgf("%s: get image through FruityBang error > %s", op, err)
		return nil, err
	}

	fImage, err := image.ToFile()
	if err != nil {
		g.log.Error().Msgf("%s: convert image generated from Kandinsky to os.File error > %s", op, err)
		return nil, err
	}

	size, err := getFileSize(*fImage)
	if err != nil {
		g.log.Error().Msgf("%s: get Kandinsky image size error > %s", op, err)
		return nil, err
	}

	imageName := translit(title)
	err = image.SavePNGTo(imageName, imagePath)
	if err != nil {
		g.log.Error().Msgf("%s: save image generated from Kandinsky error > %s", op, err)
		return nil, err
	}

	preparedPath := prepareImagePath(imagePath, imageName)

	model := &model.Image{
		ImagePath: preparedPath,
		Size:      size,
		Name:      imageName,
		Alt:       title,
	}

	return model, nil
}

func getFruityImage(title string) (*kandinsky.Image, error) {
	url := "http://alekslesik1.fvds.ru:8008/v2/images"

	imageName := translit(title)

	type params struct {
		Title string `json:"title"`
		Name  string `json:"name"`
	}

	p := params{
		Title: title,
		Name:  imageName,
	}

	b, err := json.Marshal(&p)
	if err != nil {
		return nil, err
	}

	body := bytes.NewBuffer(b)

	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")

	client := http.Client{}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	result, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	type FruityResponse struct {
		Title           string   `json:"title"`
		Name            string   `json:"name"`
		ID              int      `json:"id"`
		Images          []string `json:"images"`
		ByteImageSizeKB float64  `json:"byte_image_size_kB"`
		Width           int      `json:"width"`
		Height          int      `json:"height"`
		Error           string   `json:"error"`
	}

	var fRes FruityResponse

	err = json.Unmarshal(result, &fRes)
	if err != nil {
		return nil, err
	}

	if fRes.Error != "" {
		return nil, errors.New("error from Fruity API: " + fRes.Error)
	}

	image := &kandinsky.Image{Images: make([]string, 1)}

	err = image.AddBase64(fRes.Images[0])
	if err != nil {
		return nil, err
	}

	return image, nil
}

// GrabArticle grab article from
func (g *Grabber) GrabArticle() (*model.Article, error) {
	const op = "grabber.GrabArticle()"

	// Get last article from list
	last, err := g.LastArticle("parts/news/", Tag{
		Name: "ul",
		Key:  "class",
		Val:  "parts-page__body _parts-news",
	})
	if err != nil {
		g.log.Error().Msgf("%s: get last article error > %s", op, err)
		return nil, err
	}

	// Get article title
	title, err := g.TagInner(last, Tag{Name: "span", Key: "class", Val: "topic-body__title"})
	if err != nil {
		g.log.Error().Msgf("%s: get article title error > %s", op, err)
		return nil, err
	}

	// Get article PreviewText
	previewText, err := g.TagInner(last, Tag{Name: "div", Key: "class", Val: "topic-body__title-yandex"})
	if err != nil {
		g.log.Error().Msgf("%s: get article title error > %s", op, err)
		return nil, err
	}

	// Current article time
	articleTime := time.Now()

	// Get article tag
	tag, err := g.TagInner(last, Tag{Name: "a", Key: "class", Val: "topic-header__item topic-header__rubric"})
	if err != nil {
		g.log.Error().Msgf("%s: get article title error > %s", op, err)
		return nil, err
	}

	// Get detail text
	detailText, err := g.DetailText(last, Tag{Name: "p", Key: "class", Val: "topic-body__content-text"})
	if err != nil {
		g.log.Error().Msgf("%s: get article title error > %s", op, err)
		return nil, err
	}

	// Get article href (translit from title)
	href := translit(title)

	// Get article kind (article)
	kind := "article"

	// TODO extract Comments
	// Get category (translit from tag)
	category := translit(tag)

	// fill article model
	article := &model.Article{
		Title:       title,
		PreviewText: previewText,
		ArticleTime: articleTime,
		Tag:         tag,
		DetailText:  detailText,
		Href:        href,
		Category:    category,
		Kind:        kind,
		VideoID:     0,
	}

	return article, nil
}

// LastArticle return last element url from list
func (g *Grabber) LastArticle(url string, tag Tag) (string, error) {
	const op = "grabber.GetLast()"

	// get node from url
	node, err := getNodePage(g.home + url)
	if err != nil {
		g.log.Error().Msgf("%s: get node page error > %s", op, err)
		return "", err
	}

	// recursive find inner of tag with atr
	var result string
	var f func(*html.Node)
	f = func(n *html.Node) {
		// find ul
		if n.Type == html.ElementNode && n.Data == tag.Name {
			for _, attr := range n.Attr {
				// find class=value
				if attr.Key == tag.Key && strings.Contains(attr.Val, tag.Val) {
					for _, attr := range n.FirstChild.FirstChild.Attr {
						if attr.Key == "href" {
							result = attr.Val
							break
						}
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(node)

	return result, nil
}

// TagInner return inner of tag with attribute from url
func (g *Grabber) TagInner(url string, tag Tag) (string, error) {
	const op = "grabber.TagInner()"

	// get node from url
	node, err := getNodePage(g.home + url)
	if err != nil {
		g.log.Error().Msgf("%s: get node page error > %s", op, err)
		return "", err
	}

	// recursive find result of tag with atr
	var result string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == tag.Name {
			for _, a := range n.Attr {
				if a.Key == tag.Key && a.Val == tag.Val {
					if d := n.FirstChild.Data; d != "" {
						result = d
					}
					break
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(node)

	return result, nil
}

// DetailText return detail article text
func (g *Grabber) DetailText(url string, tag Tag) (string, error) {
	const op = "grabber.DetailText()"

	var result string
	var resArr []string

	// get node from url
	node, err := getNodePage(g.home + url)
	if err != nil {
		g.log.Error().Msgf("%s: get node page error > %s", op, err)
		return "", err
	}

	// recursive find result of tag with atr
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == tag.Name {
			for _, a := range n.Attr {
				if a.Key == tag.Key && a.Val == tag.Val {
					if p := n.FirstChild.Data; p != "" {
						p = "<p>" + p + "</p>"
						resArr = append(resArr, p)
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(node)

	result = strings.Join(resArr, "\n")

	return result, nil
}

// getNodePage return page like a Node
func getNodePage(url string) (*html.Node, error) {
	// response from url
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// read response body
	var body []byte
	if res.StatusCode == http.StatusOK {
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
	}

	// create node from response body
	node, err := html.Parse(strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}

	return node, nil
}

// getTagHref return tag href with attr from url
func getTagHref(url string) (string, error) {
	// get node from url
	node, err := getNodePage(url)
	if err != nil {
		return "", err
	}

	// recursive find inner of tag with atr
	var href string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "ul" {
			for _, a := range n.Attr {
				if a.Key == "class" && strings.Contains(a.Val, "parts-page__body _parts-news") {

					for _, attr := range n.FirstChild.Attr {
						if attr.Key == "href" {
							href = attr.Val
							break
						}
					}
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(node)

	return href, nil
}

// translit
func translit(src string) string {
	if src == "" {
		return ""
	}

	var result string
	var resArr []string

	var dictionary = map[string]string{
		"А": "a", "а": "a", "Б": "b", "б": "b", "В": "v", "в": "v", "Г": "g", "г": "g", "Д": "d",
		"д": "d", "Е": "e", "е": "e", "Ё": "e", "ё": "e", "Ж": "zh", "ж": "zh", "З": "z", "з": "z",
		"И": "i", "и": "i", "Й": "i", "й": "i", "К": "k", "к": "k", "Л": "l", "л": "l", "М": "m",
		"м": "m", "Н": "n", "н": "n", "О": "o", "о": "o", "П": "p", "п": "p", "Р": "r", "р": "r",
		"С": "s", "с": "s", "Т": "t", "т": "t", "У": "u", "у": "u", "Ф": "f", "ф": "f", "Х": "h",
		"х": "h", "Ц": "c", "ц": "c", "Ч": "ch", "ч": "ch", "Ш": "sh", "ш": "sh", "Щ": "sh", "щ": "sh",
		"Ъ": "", "ъ": "", "Ы": "y", "ы": "i", "Ь": "", "ь": "", "Э": "e", "э": "e", "Ю": "iu", "ю": "iu",
		"Я": "ia", "я": "ia", " ": "-", "0": "0", "1": "1", "2": "2", "3": "3", "4": "4", "5": "5", "6": "6",
		"7": "7", "8": "8", "9": "9", "a": "a", "b": "b", "c": "c", "d": "d", "e": "e", "f": "f", "g": "g",
		"h": "h", "i": "i", "j": "j", "k": "k", "l": "l", "m": "m", "n": "n", "o": "o", "p": "p", "q": "q",
		"r": "r", "s": "s", "t": "t", "u": "u", "v": "v", "w": "w", "x": "x", "y": "y", "z": "z", "A": "a",
		"B": "b", "C": "c", "D": "d", "E": "e", "F": "f", "G": "g", "H": "h", "I": "i", "J": "j", "K": "k",
		"L": "l", "M": "m", "N": "n", "O": "o", "P": "p", "Q": "q", "R": "r", "S": "s", "T": "t", "U": "u",
		"V": "v", "W": "w", "X": "x", "Y": "y", "Z": "z", ".": "-", ",": "-", "!": "-", "?": "-", ";": "-", ":": "-",
		"'": "-", "\"": "-", "`": "-", "~": "-", "@": "-", "#": "-", "$": "-", "%": "-", "^": "-", "&": "-", "*": "-",
		"(": "-", ")": "-", "_": "-", "+": "-", "=": "-", "[": "-", "]": "-", "{": "-", "}": "-", "\\": "-", "|": "-",
		"/": "-", "<": "-", ">": "-", "\n": "-", "\t": "-", "-": "-",
	}

	split := strings.Split(src, "")

	for k, s := range split {
		if k != 0 {
			if dictionary[split[k-1]] == "-" && dictionary[s] == "-" {
				continue
			}
		}

		resArr = append(resArr, dictionary[s])
	}

	result = strings.Join(resArr, "")
	// result = strings.Trim(result, " ")
	result = strings.Trim(result, "-")

	return result
}

func getFileSize(f os.File) (int64, error) {
	fInfo, err := f.Stat()
	if err != nil {
		return 0, err
	}

	return fInfo.Size(), nil
}

func prepareImagePath(path, name string) string {
	name = name + ".png"

	path = strings.Replace(path, "website", "", 1)

	return path + name
}
