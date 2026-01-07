package hateoas

import (
	"fmt"
)

type Link struct {
	Href   string
	Rel    string
	Method string
	Title  string
}

func BuildSelfLink(baseURL string, resourcePath string) Link {
	return Link{
		Href:   fmt.Sprintf("%s%s", baseURL, resourcePath),
		Rel:    "self",
		Method: "GET",
		Title:  "Get this resource",
	}
}

func BuildUpdateLink(baseURL string, resourcePath string) Link {
	return Link{
		Href:   fmt.Sprintf("%s%s", baseURL, resourcePath),
		Rel:    "update",
		Method: "PATCH",
		Title:  "Update this resource",
	}
}

func BuildDeleteLink(baseURL string, resourcePath string) Link {
	return Link{
		Href:   fmt.Sprintf("%s%s", baseURL, resourcePath),
		Rel:    "delete",
		Method: "DELETE",
		Title:  "Delete this resource",
	}
}

func BuildCreateLink(baseURL string, resourcePath string) Link {
	return Link{
		Href:   fmt.Sprintf("%s%s", baseURL, resourcePath),
		Rel:    "create",
		Method: "POST",
		Title:  "Create a new resource",
	}
}

func BuildRelatedLink(url string, rel string, method string, title string) Link {
	return Link{
		Href:   url,
		Rel:    rel,
		Method: method,
		Title:  title,
	}
}

func BuildPaginationLink(baseURL string, resourcePath string, queryParams string, rel string, title string) Link {
	href := fmt.Sprintf("%s%s", baseURL, resourcePath)
	if queryParams != "" {
		href = fmt.Sprintf("%s?%s", href, queryParams)
	}
	return Link{
		Href:   href,
		Rel:    rel,
		Method: "GET",
		Title:  title,
	}
}

func BuildParentLink(baseURL string, resourcePath string) Link {
	return Link{
		Href:   fmt.Sprintf("%s%s", baseURL, resourcePath),
		Rel:    "parent",
		Method: "GET",
		Title:  "Get parent collection",
	}
}
