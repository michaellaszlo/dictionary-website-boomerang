package main

import (
  "os"
  "fmt"
  "strings"
  "unicode/utf8"
  "regexp"
  "net/url"
  "path/filepath"
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
)

func main() {

  root := os.Getenv("DOCUMENT_ROOT")
  dbPath := filepath.Join(root, "db/dictionary.sq3")
  db, err := sql.Open("sqlite3", dbPath)
  if err != nil {
    fmt.Printf("unable to open %s: %s\n", dbPath, err)
    return
  }
  var query, exitMessage string
  //logFile, err := os.OpenFile("log.txt", os.O_APPEND | os.O_WRONLY, 0666)
  //if err != nil {
  //  fmt.Fprintf(os.Stderr, "os.OpenFile error: %s\n", err)
  //}
  //defer logFile.Close()

  request := os.Getenv("REQUEST_URI")
  originalRequest := request
  // If there is no trailing slash, add one.
  if len(request) < 1 || request[len(request)-1] != '/' {
    request += "/"
  }
  // Replace repeated slashes with single slashes.
  request = regexp.MustCompile("/+").ReplaceAllString(request, "/")
  // Count the slashes.
  numSlashes := strings.Count(request, "/")

  getWords := func(query string, parameters ...interface{}) []string {
    rows, _ := db.Query(query, parameters...)
    words := []string{}
    for rows.Next() {
      var word string
      rows.Scan(&word)
      words = append(words, word)
    }
    return words
  }

  var words []string
  var canonicalRequest, initial, title, titleExtension string

  // There are two slashes if the full listing is requested: /listing/
  if numSlashes == 2 {
    canonicalRequest = "/listing/"
    words = getWords("select word from entries order by word")
    title = fmt.Sprintf("%d words in all", len(words))
    titleExtension = " : index of all words"

    // Otherwise, there should be exactly three slashes.
  } else if numSlashes != 3 {

    fmt.Print("Status: 404 Not Found\n")

    titleExtension := " : page not found"
    word := ""

    fmt.Print("Content-Type: text/html; charset=utf-8\n\n")

    fmt.Print(`<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN"
 "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd"> 
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en" lang="en">
<head>
<script type="text/javascript">
  WebFontConfig = {
    google: { families: [ 'Open+Sans::latin', 'Bitter::latin' ] }
  };
  (function() {
    var wf = document.createElement('script');
    wf.src = ('https:' == document.location.protocol ? 'https' : 'http') +
    '://ajax.googleapis.com/ajax/libs/webfont/1/webfont.js';
    wf.type = 'text/javascript';
    wf.async = 'true';
    var s = document.getElementsByTagName('script')[0];
    s.parentNode.insertBefore(wf, s);
  })();
</script>
  <title>The Devil's Dictionary`)
    fmt.Print(titleExtension)
    fmt.Print(`</title>
  <link rel="stylesheet" href="/css/dictionary.css" />
</head>
<body>
<div id="wrapper">

`)
    linkClass := "title"
    if request == "/" {
      linkClass += " home"
    } else {
      linkClass += " notHome"
    }

    fmt.Print(`
<h1><a href="/" class="`)
    fmt.Print(linkClass)
    fmt.Print(`"
     >The Devil's Dictionary<div class="homeLinkIcon"></div></a></h1>

`)

    // Run a database query to get the sorted list of initial letters.

    query = "select distinct substr(word, 1, 1) as initial from entries" +
      " order by initial"
    rows, _ := db.Query(query)
    initials := []string{}
    for rows.Next() {
      var initial string
      rows.Scan(&initial)
      initials = append(initials, initial)
    }

    // Print the initial letters as links to listing pages.

    fmt.Print("<ul class=\"large listing\">")

    for _, initial := range initials {
      path := "/listing/" + initial + "/"
      insert := ""
      if path == request {
        insert = " class=\"currentListing\""
      } else if utf8.RuneCountInString(word) > 0 &&
        string([]rune(word)[:1]) == initial {
        insert = " class=\"relatedListing\""
      }
      display := strings.ToUpper(initial)
      link := "<a href=\"" + path + "\"" + insert + ">" + display + "</a>"
      fmt.Print("<li>" + link + "</li><wbr />")
    }
    first := initials[0]
    last := initials[len(initials)-1]
    label := strings.ToUpper(first + "&#8209;" + last)
    fmt.Print("<li class=\"long\"><a href=\"/listing/\">" +
      label + "</a></li></ul>")

    fmt.Print("</ul>\n")

    fmt.Print(`

<html>
<body>
<h1> The devil you say. Page not found! </h1>

`)

    if originalRequest != "" {

      fmt.Print(`

<p> You requested: </p>

<pre> `)
      fmt.Print(originalRequest)
      fmt.Print(` </pre>

<p> There is no such page. </p>

<p> Perhaps you'd like to visit our <a href="/">home page</a>? </p>

`)
    }

    if exitMessage != "" {

      fmt.Print(`

<p> `)
      fmt.Print(exitMessage)
      fmt.Print(` </p>

`)
    }

    fmt.Print(`
  
</div><!--end wrapper-->
</body>
</html>


`)
    return

  } else {
    // Extract the requested initial.
    parts := strings.Split(request, "/")
    // The request '/entry/abacus/' yields ['', 'entry', 'abacus', ''].
    initial = parts[2]
    // We can ignore Unescape errors and leave the string as is.
    initial, _ = url.QueryUnescape(initial)
    // Render in lower case.
    initial = strings.ToLower(initial)

    // Construct the canonical URL for this page. At this writing, the
    //  token parts[1] resolves to "listing". This may change eventually.
    //  If the virtual directory structure changes, the template will
    //  have to be modified here and elsewhere.
    canonicalRequest = fmt.Sprintf("/%s/%s/",
      parts[1], url.QueryEscape(initial))

    // We should have a single character.
    if utf8.RuneCountInString(initial) != 1 {

      fmt.Print("Status: 404 Not Found\n")

      titleExtension := " : page not found"
      word := ""

      fmt.Print("Content-Type: text/html; charset=utf-8\n\n")

      fmt.Print(`<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN"
 "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd"> 
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en" lang="en">
<head>
<script type="text/javascript">
  WebFontConfig = {
    google: { families: [ 'Open+Sans::latin', 'Bitter::latin' ] }
  };
  (function() {
    var wf = document.createElement('script');
    wf.src = ('https:' == document.location.protocol ? 'https' : 'http') +
    '://ajax.googleapis.com/ajax/libs/webfont/1/webfont.js';
    wf.type = 'text/javascript';
    wf.async = 'true';
    var s = document.getElementsByTagName('script')[0];
    s.parentNode.insertBefore(wf, s);
  })();
</script>
  <title>The Devil's Dictionary`)
      fmt.Print(titleExtension)
      fmt.Print(`</title>
  <link rel="stylesheet" href="/css/dictionary.css" />
</head>
<body>
<div id="wrapper">

`)
      linkClass := "title"
      if request == "/" {
        linkClass += " home"
      } else {
        linkClass += " notHome"
      }

      fmt.Print(`
<h1><a href="/" class="`)
      fmt.Print(linkClass)
      fmt.Print(`"
     >The Devil's Dictionary<div class="homeLinkIcon"></div></a></h1>

`)

      // Run a database query to get the sorted list of initial letters.

      query = "select distinct substr(word, 1, 1) as initial from entries" +
        " order by initial"
      rows, _ := db.Query(query)
      initials := []string{}
      for rows.Next() {
        var initial string
        rows.Scan(&initial)
        initials = append(initials, initial)
      }

      // Print the initial letters as links to listing pages.

      fmt.Print("<ul class=\"large listing\">")

      for _, initial := range initials {
        path := "/listing/" + initial + "/"
        insert := ""
        if path == request {
          insert = " class=\"currentListing\""
        } else if utf8.RuneCountInString(word) > 0 &&
          string([]rune(word)[:1]) == initial {
          insert = " class=\"relatedListing\""
        }
        display := strings.ToUpper(initial)
        link := "<a href=\"" + path + "\"" + insert + ">" + display + "</a>"
        fmt.Print("<li>" + link + "</li><wbr />")
      }
      first := initials[0]
      last := initials[len(initials)-1]
      label := strings.ToUpper(first + "&#8209;" + last)
      fmt.Print("<li class=\"long\"><a href=\"/listing/\">" +
        label + "</a></li></ul>")

      fmt.Print("</ul>\n")

      fmt.Print(`

<html>
<body>
<h1> The devil you say. Page not found! </h1>

`)

      if originalRequest != "" {

        fmt.Print(`

<p> You requested: </p>

<pre> `)
        fmt.Print(originalRequest)
        fmt.Print(` </pre>

<p> There is no such page. </p>

<p> Perhaps you'd like to visit our <a href="/">home page</a>? </p>

`)
      }

      if exitMessage != "" {

        fmt.Print(`

<p> `)
        fmt.Print(exitMessage)
        fmt.Print(` </p>

`)
      }

      fmt.Print(`
  
</div><!--end wrapper-->
</body>
</html>


`)
      return

    }

    // Look up the set of words starting with our initial.
    query = "select word from entries where word like ? order by word"
    words = getWords(query, initial+"%")
    numWords := len(words)

    // If the set is empty, exit with an appropriate message.
    if numWords == 0 {

      fmt.Print("Status: 404 Not Found\n")

      titleExtension := " : page not found"
      word := ""

      fmt.Print("Content-Type: text/html; charset=utf-8\n\n")

      fmt.Print(`<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN"
 "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd"> 
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en" lang="en">
<head>
<script type="text/javascript">
  WebFontConfig = {
    google: { families: [ 'Open+Sans::latin', 'Bitter::latin' ] }
  };
  (function() {
    var wf = document.createElement('script');
    wf.src = ('https:' == document.location.protocol ? 'https' : 'http') +
    '://ajax.googleapis.com/ajax/libs/webfont/1/webfont.js';
    wf.type = 'text/javascript';
    wf.async = 'true';
    var s = document.getElementsByTagName('script')[0];
    s.parentNode.insertBefore(wf, s);
  })();
</script>
  <title>The Devil's Dictionary`)
      fmt.Print(titleExtension)
      fmt.Print(`</title>
  <link rel="stylesheet" href="/css/dictionary.css" />
</head>
<body>
<div id="wrapper">

`)
      linkClass := "title"
      if request == "/" {
        linkClass += " home"
      } else {
        linkClass += " notHome"
      }

      fmt.Print(`
<h1><a href="/" class="`)
      fmt.Print(linkClass)
      fmt.Print(`"
     >The Devil's Dictionary<div class="homeLinkIcon"></div></a></h1>

`)

      // Run a database query to get the sorted list of initial letters.

      query = "select distinct substr(word, 1, 1) as initial from entries" +
        " order by initial"
      rows, _ := db.Query(query)
      initials := []string{}
      for rows.Next() {
        var initial string
        rows.Scan(&initial)
        initials = append(initials, initial)
      }

      // Print the initial letters as links to listing pages.

      fmt.Print("<ul class=\"large listing\">")

      for _, initial := range initials {
        path := "/listing/" + initial + "/"
        insert := ""
        if path == request {
          insert = " class=\"currentListing\""
        } else if utf8.RuneCountInString(word) > 0 &&
          string([]rune(word)[:1]) == initial {
          insert = " class=\"relatedListing\""
        }
        display := strings.ToUpper(initial)
        link := "<a href=\"" + path + "\"" + insert + ">" + display + "</a>"
        fmt.Print("<li>" + link + "</li><wbr />")
      }
      first := initials[0]
      last := initials[len(initials)-1]
      label := strings.ToUpper(first + "&#8209;" + last)
      fmt.Print("<li class=\"long\"><a href=\"/listing/\">" +
        label + "</a></li></ul>")

      fmt.Print("</ul>\n")

      fmt.Print(`

<html>
<body>
<h1> The devil you say. Page not found! </h1>

`)

      if originalRequest != "" {

        fmt.Print(`

<p> You requested: </p>

<pre> `)
        fmt.Print(originalRequest)
        fmt.Print(` </pre>

<p> There is no such page. </p>

<p> Perhaps you'd like to visit our <a href="/">home page</a>? </p>

`)
      }

      if exitMessage != "" {

        fmt.Print(`

<p> `)
        fmt.Print(exitMessage)
        fmt.Print(` </p>

`)
      }

      fmt.Print(`
  
</div><!--end wrapper-->
</body>
</html>


`)
      return

    }

    plural := ""
    if numWords > 1 {
      plural = "s"
    }
    title = fmt.Sprintf("%d word%s starting with %s",
      numWords, plural, strings.ToUpper(initial))
    titleExtension = " : index of words starting with " +
      strings.ToUpper(initial)
  }

  // Redirect to the canonical URL if necessary.
  if originalRequest != canonicalRequest {
    fmt.Print("Status: 301 Moved Permanently\n")
    fmt.Printf("Location: %s\n\n", canonicalRequest)
    return
  }

  word := ""

  fmt.Print("Content-Type: text/html; charset=utf-8\n\n")

  fmt.Print(`<!DOCTYPE html PUBLIC "-//W3C//DTD XHTML 1.0 Strict//EN"
 "http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd"> 
<html xmlns="http://www.w3.org/1999/xhtml" xml:lang="en" lang="en">
<head>
<script type="text/javascript">
  WebFontConfig = {
    google: { families: [ 'Open+Sans::latin', 'Bitter::latin' ] }
  };
  (function() {
    var wf = document.createElement('script');
    wf.src = ('https:' == document.location.protocol ? 'https' : 'http') +
    '://ajax.googleapis.com/ajax/libs/webfont/1/webfont.js';
    wf.type = 'text/javascript';
    wf.async = 'true';
    var s = document.getElementsByTagName('script')[0];
    s.parentNode.insertBefore(wf, s);
  })();
</script>
  <title>The Devil's Dictionary`)
  fmt.Print(titleExtension)
  fmt.Print(`</title>
  <link rel="stylesheet" href="/css/dictionary.css" />
</head>
<body>
<div id="wrapper">

`)
  linkClass := "title"
  if request == "/" {
    linkClass += " home"
  } else {
    linkClass += " notHome"
  }

  fmt.Print(`
<h1><a href="/" class="`)
  fmt.Print(linkClass)
  fmt.Print(`"
     >The Devil's Dictionary<div class="homeLinkIcon"></div></a></h1>

`)

  // Run a database query to get the sorted list of initial letters.

  query = "select distinct substr(word, 1, 1) as initial from entries" +
    " order by initial"
  rows, _ := db.Query(query)
  initials := []string{}
  for rows.Next() {
    var initial string
    rows.Scan(&initial)
    initials = append(initials, initial)
  }

  // Print the initial letters as links to listing pages.

  fmt.Print("<ul class=\"large listing\">")

  for _, initial := range initials {
    path := "/listing/" + initial + "/"
    insert := ""
    if path == request {
      insert = " class=\"currentListing\""
    } else if utf8.RuneCountInString(word) > 0 &&
      string([]rune(word)[:1]) == initial {
      insert = " class=\"relatedListing\""
    }
    display := strings.ToUpper(initial)
    link := "<a href=\"" + path + "\"" + insert + ">" + display + "</a>"
    fmt.Print("<li>" + link + "</li><wbr />")
  }
  first := initials[0]
  last := initials[len(initials)-1]
  label := strings.ToUpper(first + "&#8209;" + last)
  fmt.Print("<li class=\"long\"><a href=\"/listing/\">" +
    label + "</a></li></ul>")

  fmt.Print("</ul>\n")

  fmt.Print(`


<h2> `)
  fmt.Print(title)
  fmt.Print(`: </h2>

<ul class="listing">
`)
  for _, word := range words {
    path := "/entry/" + url.QueryEscape(word) + "/"
    // Replace hyphens with non-breaking hyphens.
    word = strings.Replace(word, "-", "&#8209;", -1)
    // Replace spaces with non-breaking spaces.
    word = strings.Replace(word, " ", "&nbsp;", -1)
    link := fmt.Sprintf("<a href=\"%s\">%s</a>", path, word)

    fmt.Print(`
    <li>`)
    fmt.Print(link)
    fmt.Print(`</li><wbr />
`)
  }

  fmt.Print(`
</ul>
  
</div><!--end wrapper-->
</body>
</html>
`)
}
