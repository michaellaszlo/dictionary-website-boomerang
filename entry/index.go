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
  // There should now be exactly three slashes.
  if strings.Count(request, "/") != 3 {

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
  // Extract the requested dictionary headword.
  parts := strings.Split(request, "/")
  // The request '/entry/abacus/' yields ['', 'entry', 'abacus', ''].
  word := parts[2]
  // Apply URL decoding and UTF-8 decoding.
  word, err = url.QueryUnescape(word)
  if err != nil {

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
  // Render in lower case.
  word = strings.ToLower(word)

  // Look up the word in the dictionary.
  query = "select definition from entries where word = ?"
  var definition string
  err = db.QueryRow(query, word).Scan(&definition)
  // If there is no such word, exit with an appropriate message.
  if err != nil {
    word = ""

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

  // Construct the canonical URL for this word. At this writing, the
  //  token parts[1] resolves to "entry". This may change eventually.
  //  If the virtual directory structure changes, the template will
  //  have to be modified here and elsewhere.
  canonicalRequest := "/" + parts[1] + "/" + url.QueryEscape(word) + "/"
  // Redirect to the canonical URL if necessary.
  if originalRequest != canonicalRequest {
    fmt.Print("Status: 301 Moved Permanently\n")
    fmt.Printf("Location: %s\n\n", canonicalRequest)
    return
  }

  titleExtension := " : " + word

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
  fmt.Print(word)
  fmt.Print(`: </h2>

  <div class="definition">

`)
  fmt.Print(definition)
  fmt.Print(`


</div><!--end definition-->


  
</div><!--end wrapper-->
</body>
</html>
`)
}
