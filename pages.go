package main

import (
	_ "embed"
	"fmt"
	"html"
	"path"
	"strings"
)

func createHomePageHTML(topics []topicStruct) string {
	topicsListHTMLBuilder := strings.Builder{}
	topicsListHTMLBuilder.WriteString("<ul>")
	for _, topic := range topics {
		topicURL := fmt.Sprintf("/%s", topic.pathId)
		listItemHTMLTemplate := `<li><a href="%s">%s</a></li>`
		listItemHTML := fmt.Sprintf(listItemHTMLTemplate, html.EscapeString(topicURL), html.EscapeString(topic.title))
		topicsListHTMLBuilder.WriteString(listItemHTML)
	}
	topicsListHTMLBuilder.WriteString("</ul>")
	topicsListHTML := topicsListHTMLBuilder.String()

	contentHTMLTemplate := `<h1>Auth book</h1>
<p>This is my personal auth book.
It is a collection of guides, recommendations, and examples for implementing auth in web applications based on my personal opinion.
It is completely free with zero ads.
I hope this is useful for anyone looking to learn more about auth, security, and the web in general.
</p>
<p>
As the name implies, this book focuses heavily on the authentication and login system for your application.
For more general security topics, see the <a href="https://cheatsheetseries.owasp.org/index.html">OWASP Cheat Sheet Series</a>.
</p>
<p>
If you have any questions, feel free to ask them on the <a href="https://discord.gg/zZqCfVUMnX">Discord server</a> or on <a href="https://github.com/pilcrowonpaper/auth.pilcrowonpaper.com/discussions">GitHub Discussions</a>.
</p>
<p>Please also consider supporting my work on <a href="https://github.com/sponsors/pilcrowonpaper">GitHub Sponsors</a>.</p>
<p><i>
Written and maintained by <a href="https://pilcrowonpaper.com">Pilcrow</a>.
Source code available on <a href="https://github.com/pilcrowonpaper/auth.pilcrowonpaper.com">GitHub</a>.
</i></p>
<h2>Topics</h2>
%s
<h2>Examples</h2>
<p>Complete, fully open-source example websites written in Go based on the contents of the book.</p>
<ul>
	<li><a href="https://basic-example.auth.pilcrowonpaper.com">Basic auth example</a>: Password example with email address verification and password reset (<a href="https://github.com/pilcrowonpaper/basic-example.auth.pilcrowonpaper.com">source code</a>).</li>
	<li><a href="https://passwordless-example.auth.pilcrowonpaper.com">Passwordless auth example</a>: Passkey and email code sign-in example with email address verification (<a href="https://github.com/pilcrowonpaper/passwordless-example.auth.pilcrowonpaper.com">source code</a>).</li>
</ul>`
	contentHTML := fmt.Sprintf(contentHTMLTemplate, topicsListHTML)

	pageHTML := createPageHTML("Auth book", "/", contentHTML)

	return pageHTML
}

func createTopicPageHTML(topicTitle string, topicPathId string, contentHTML string) string {
	pageTitle := fmt.Sprintf("%s | Auth book", topicTitle)
	pagePath := fmt.Sprintf("/%s", topicPathId)

	pageHTML := createPageHTML(pageTitle, pagePath, contentHTML)

	return pageHTML
}

type topicStruct struct {
	id     string
	pathId string
	title  string
}

//go:embed assets/stylesheet.css
var stylesheetCSS string

func createPageHTML(pageTitle string, pagePath string, contentHTML string) string {
	pageURL := path.Join("https://auth.pilcrowonpaper.com", pagePath)

	pageHTMLTemplate := `<!DOCTYPE html>
<html lang="en">
<head>
	<title>%s</title>
	<meta name="description" content="Pilcrow's auth book.">

	<meta charset="utf-8">
    <meta name="viewport" content="width=device-width">

	<meta property="og:title" content="%s">
	<meta property="og:type" content="website">
	<meta property="og:locale" content="en_US">
	<meta property="og:site_name" content="Auth book">
	<meta property="og:description" content="Pilcrow's auth book.">
	<meta property="og:url" content="%s">
	<meta property="og:image" content="https://pilcrowonpaper.com/pilcrow.jpeg">

	<meta name="twitter:card" content="summary">
    <meta name="twitter:site" content="@pilcrowonpaper">

    <link rel="icon" type="image/jpeg" href="https://pilcrowonpaper.com/pilcrow.jpeg">

    <link rel="canonical" href="%s">

	<style>%s</style>
</head>

<body>
	<header>
		<a id="home-link" href="/"><img id="pilcrow-icon" src="https://pilcrowonpaper.com/pilcrow.jpeg" alt="Pilcrow"><p>Auth book</p></a>
	</header>
	<main>%s</main>
</body>
</html>`

	pageHTML := fmt.Sprintf(pageHTMLTemplate, html.EscapeString(pageTitle), html.EscapeString(pageTitle), html.EscapeString(pageURL), html.EscapeString(pageURL), stylesheetCSS, contentHTML)

	return pageHTML
}
