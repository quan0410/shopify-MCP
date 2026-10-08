package tools

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/datumbridge/shopify-mcp/internal/mcp"
)

func registerContentTools(add toolAdder) {
	// Pages
	add(
		"shopify_list_pages",
		"List static store pages (e.g. About Us, Contact, FAQ, Terms) with optional filtering.",
		baseProps(map[string]interface{}{
			"limit":    map[string]interface{}{"type": "integer", "description": "Number of pages to return (default 50, max 250)"},
			"since_id": map[string]interface{}{"type": "string", "description": "Restrict results to after the specified ID"},
			"title":    map[string]interface{}{"type": "string", "description": "Filter pages by title"},
		}),
		nil,
		handleShopifyListPages,
	)

	add(
		"shopify_get_page",
		"Retrieve details and HTML content of a specific static page by ID.",
		baseProps(map[string]interface{}{
			"page_id": map[string]interface{}{"type": "string", "description": "The numeric ID of the page"},
		}),
		[]string{"page_id"},
		handleShopifyGetPage,
	)

	add(
		"shopify_create_page",
		"Create a new static store page with title and body HTML.",
		baseProps(map[string]interface{}{
			"title":        map[string]interface{}{"type": "string", "description": "Title of the page"},
			"body_html":    map[string]interface{}{"type": "string", "description": "Page content in HTML or plain text"},
			"author":       map[string]interface{}{"type": "string", "description": "Name of the page author (optional)"},
			"is_published": map[string]interface{}{"type": "boolean", "description": "Whether the page is published immediately (default true)"},
		}),
		[]string{"title"},
		handleShopifyCreatePage,
	)

	add(
		"shopify_update_page",
		"Update an existing static page's title, body HTML, or publication status.",
		baseProps(map[string]interface{}{
			"page_id":      map[string]interface{}{"type": "string", "description": "The numeric ID of the page to update"},
			"title":        map[string]interface{}{"type": "string", "description": "New title for the page"},
			"body_html":    map[string]interface{}{"type": "string", "description": "New HTML content for the page"},
			"is_published": map[string]interface{}{"type": "boolean", "description": "Whether the page is published"},
		}),
		[]string{"page_id"},
		handleShopifyUpdatePage,
	)

	add(
		"shopify_delete_page",
		"Permanently delete a static page by ID.",
		baseProps(map[string]interface{}{
			"page_id": map[string]interface{}{"type": "string", "description": "The numeric ID of the page to delete"},
		}),
		[]string{"page_id"},
		handleShopifyDeletePage,
	)

	// Blogs & Articles
	add(
		"shopify_list_blogs",
		"List all blogs created on the Shopify store.",
		baseProps(nil),
		nil,
		handleShopifyListBlogs,
	)

	add(
		"shopify_list_articles",
		"List articles/blog posts published in a specific blog.",
		baseProps(map[string]interface{}{
			"blog_id":  map[string]interface{}{"type": "string", "description": "The numeric ID of the blog"},
			"limit":    map[string]interface{}{"type": "integer", "description": "Number of articles to return (default 50)"},
			"since_id": map[string]interface{}{"type": "string", "description": "Restrict results to after the specified ID"},
		}),
		[]string{"blog_id"},
		handleShopifyListArticles,
	)

	add(
		"shopify_get_article",
		"Retrieve full details of a specific blog article by blog ID and article ID.",
		baseProps(map[string]interface{}{
			"blog_id":    map[string]interface{}{"type": "string", "description": "The numeric ID of the blog"},
			"article_id": map[string]interface{}{"type": "string", "description": "The numeric ID of the article"},
		}),
		[]string{"blog_id", "article_id"},
		handleShopifyGetArticle,
	)

	add(
		"shopify_create_article",
		"Create a new article/blog post under a specific blog.",
		baseProps(map[string]interface{}{
			"blog_id":      map[string]interface{}{"type": "string", "description": "The numeric ID of the blog"},
			"title":        map[string]interface{}{"type": "string", "description": "Article title"},
			"author":       map[string]interface{}{"type": "string", "description": "Author name"},
			"body_html":    map[string]interface{}{"type": "string", "description": "Article HTML or text body"},
			"tags":         map[string]interface{}{"type": "string", "description": "Comma-separated list of tags"},
			"summary_html": map[string]interface{}{"type": "string", "description": "Article summary/excerpt (optional)"},
			"is_published": map[string]interface{}{"type": "boolean", "description": "Whether to publish the article immediately (default true)"},
		}),
		[]string{"blog_id", "title", "body_html"},
		handleShopifyCreateArticle,
	)

	add(
		"shopify_update_article",
		"Update an existing blog article's title, body, author, tags, summary, or publication status (is_published).",
		baseProps(map[string]interface{}{
			"article_id":   map[string]interface{}{"type": "string", "description": "The numeric ID of the article to update (e.g. '622314520859' or 'articles/622314520859')"},
			"blog_id":      map[string]interface{}{"type": "string", "description": "The numeric ID of the blog (optional)"},
			"title":        map[string]interface{}{"type": "string", "description": "New title for the article"},
			"author":       map[string]interface{}{"type": "string", "description": "Author name"},
			"body_html":    map[string]interface{}{"type": "string", "description": "New HTML or text content"},
			"tags":         map[string]interface{}{"type": "string", "description": "Comma-separated list of tags"},
			"summary_html": map[string]interface{}{"type": "string", "description": "Article summary/excerpt"},
			"is_published": map[string]interface{}{"type": "boolean", "description": "Whether the article is published (true to publish, false to unpublish/draft)"},
		}),
		[]string{"article_id"},
		handleShopifyUpdateArticle,
	)
}

func handleShopifyListPages(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	q := url.Values{}
	if l := intArg(m, "limit", 50); l > 0 {
		q.Set("limit", strconv.Itoa(l))
	}
	if s := strArg(m, "since_id"); s != "" {
		q.Set("since_id", s)
	}
	if t := strArg(m, "title"); t != "" {
		q.Set("title", t)
	}

	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodGet, "/pages.json", q, nil)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

func handleShopifyGetPage(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	pageID := strArg(m, "page_id")
	if pageID == "" {
		return mcp.ToolResultError("page_id is required")
	}

	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodGet, fmt.Sprintf("/pages/%s.json", pageID), url.Values{}, nil)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

func handleShopifyCreatePage(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	title := strArg(m, "title")
	if title == "" {
		return mcp.ToolResultError("title is required")
	}

	page := map[string]interface{}{
		"title":     title,
		"published": boolArg(m, "is_published", true),
	}
	if b := strArg(m, "body_html"); b != "" {
		page["body_html"] = b
	}
	if a := strArg(m, "author"); a != "" {
		page["author"] = a
	}

	payload := map[string]interface{}{"page": page}
	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodPost, "/pages.json", url.Values{}, payload)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

func handleShopifyUpdatePage(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	pageID := strArg(m, "page_id")
	if pageID == "" {
		return mcp.ToolResultError("page_id is required")
	}

	page := map[string]interface{}{"id": pageID}
	if t := strArg(m, "title"); t != "" {
		page["title"] = t
	}
	if b := strArg(m, "body_html"); b != "" {
		page["body_html"] = b
	}
	if _, ok := m["is_published"]; ok {
		page["published"] = boolArg(m, "is_published", true)
	}

	payload := map[string]interface{}{"page": page}
	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodPut, fmt.Sprintf("/pages/%s.json", pageID), url.Values{}, payload)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

func handleShopifyDeletePage(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	pageID := strArg(m, "page_id")
	if pageID == "" {
		return mcp.ToolResultError("page_id is required")
	}

	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodDelete, fmt.Sprintf("/pages/%s.json", pageID), url.Values{}, nil)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

func handleShopifyListBlogs(args json.RawMessage) map[string]interface{} {
	client, _, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}

	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodGet, "/blogs.json", url.Values{}, nil)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

func handleShopifyListArticles(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	blogID := strArg(m, "blog_id")
	if blogID == "" {
		return mcp.ToolResultError("blog_id is required")
	}
	q := url.Values{}
	if l := intArg(m, "limit", 50); l > 0 {
		q.Set("limit", strconv.Itoa(l))
	}
	if s := strArg(m, "since_id"); s != "" {
		q.Set("since_id", s)
	}

	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodGet, fmt.Sprintf("/blogs/%s/articles.json", blogID), q, nil)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

func handleShopifyGetArticle(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	blogID := strArg(m, "blog_id")
	articleID := strArg(m, "article_id")
	if blogID == "" || articleID == "" {
		return mcp.ToolResultError("blog_id and article_id are required")
	}

	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodGet, fmt.Sprintf("/blogs/%s/articles/%s.json", blogID, articleID), url.Values{}, nil)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

func handleShopifyCreateArticle(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	blogID := strArg(m, "blog_id")
	title := strArg(m, "title")
	body := strArg(m, "body_html")
	if blogID == "" || title == "" || body == "" {
		return mcp.ToolResultError("blog_id, title, and body_html are required")
	}

	article := map[string]interface{}{
		"title":     title,
		"body_html": body,
		"published": boolArg(m, "is_published", true),
	}
	if a := strArg(m, "author"); a != "" {
		article["author"] = a
	}
	if t := strArg(m, "tags"); t != "" {
		article["tags"] = t
	}
	if s := strArg(m, "summary_html"); s != "" {
		article["summary_html"] = s
	}

	payload := map[string]interface{}{"article": article}
	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodPost, fmt.Sprintf("/blogs/%s/articles.json", blogID), url.Values{}, payload)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

func handleShopifyUpdateArticle(args json.RawMessage) map[string]interface{} {
	client, m, err := clientFrom(args)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	blogID := strArg(m, "blog_id")
	articleID := strArg(m, "article_id")
	articleID = strings.TrimPrefix(articleID, "articles/")
	articleID = strings.TrimPrefix(articleID, "gid://shopify/Article/")
	if articleID == "" {
		return mcp.ToolResultError("article_id is required")
	}

	if blogID == "" {
		gid := fmt.Sprintf("gid://shopify/Article/%s", articleID)
		artInput := map[string]interface{}{}
		if t := strArg(m, "title"); t != "" {
			artInput["title"] = t
		}
		if b := strArg(m, "body_html"); b != "" {
			artInput["body"] = b
		}
		if a := strArg(m, "author"); a != "" {
			artInput["author"] = map[string]interface{}{"name": a}
		}
		if tg := strArg(m, "tags"); tg != "" {
			var tags []string
			for _, tag := range strings.Split(tg, ",") {
				if trimmed := strings.TrimSpace(tag); trimmed != "" {
					tags = append(tags, trimmed)
				}
			}
			artInput["tags"] = tags
		}
		if s := strArg(m, "summary_html"); s != "" {
			artInput["summary"] = s
		}
		if _, ok := m["is_published"]; ok {
			artInput["isPublished"] = boolArg(m, "is_published", true)
		}

		mutation := `mutation articleUpdate($id: ID!, $article: ArticleUpdateInput!) {
			articleUpdate(id: $id, article: $article) {
				article {
					id
					title
					isPublished
					publishedAt
				}
				userErrors {
					field
					message
				}
			}
		}`
		vars := map[string]interface{}{
			"id":      gid,
			"article": artInput,
		}
		c, cancel := ctx()
		defer cancel()
		data, _, err := client.GraphQL(c, mutation, vars)
		if err != nil {
			return mcp.ToolResultError(err.Error())
		}

		var resp struct {
			Data struct {
				ArticleUpdate struct {
					Article    map[string]interface{} `json:"article"`
					UserErrors []struct {
						Field   []string `json:"field"`
						Message string   `json:"message"`
					} `json:"userErrors"`
				} `json:"articleUpdate"`
			} `json:"data"`
		}
		if err := json.Unmarshal(data, &resp); err == nil && len(resp.Data.ArticleUpdate.UserErrors) > 0 {
			var errMsgs []string
			for _, ue := range resp.Data.ArticleUpdate.UserErrors {
				errMsgs = append(errMsgs, fmt.Sprintf("%s: %s", strings.Join(ue.Field, "."), ue.Message))
			}
			return mcp.ToolResultError(fmt.Sprintf("Shopify error: %s", strings.Join(errMsgs, ", ")))
		}

		return rawResult(data)
	}

	article := map[string]interface{}{"id": articleID}
	if t := strArg(m, "title"); t != "" {
		article["title"] = t
	}
	if b := strArg(m, "body_html"); b != "" {
		article["body_html"] = b
	}
	if a := strArg(m, "author"); a != "" {
		article["author"] = a
	}
	if tg := strArg(m, "tags"); tg != "" {
		article["tags"] = tg
	}
	if s := strArg(m, "summary_html"); s != "" {
		article["summary_html"] = s
	}
	if _, ok := m["is_published"]; ok {
		article["published"] = boolArg(m, "is_published", true)
	}

	payload := map[string]interface{}{"article": article}
	c, cancel := ctx()
	defer cancel()
	data, _, err := client.Request(c, http.MethodPut, fmt.Sprintf("/blogs/%s/articles/%s.json", blogID, articleID), url.Values{}, payload)
	if err != nil {
		return mcp.ToolResultError(err.Error())
	}
	return rawResult(data)
}

