<?code
package main

import (
  "os"
  "strings"
  "unicode/utf8"
  "regexp"
  "net/url"
  "path/filepath"
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
)

func main() {
?> <?insert /advance.mer ?> <?code

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
  if numSlashes == 2  {
    canonicalRequest = "/listing/"
    words = getWords("select word from entries order by word")
    title = fmt.Sprintf("%d words in all", len(words))
    titleExtension = " : index of all words"

  // Otherwise, there should be exactly three slashes.
  } else if numSlashes != 3 {
?>  <?insert /missing.mer ?>  <?code

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
?>    <?insert /missing.mer ?>  <?code
    }

    // Look up the set of words starting with our initial.
    query = "select word from entries where word like ? order by word"
    words = getWords(query, initial+"%")
    numWords := len(words)

    // If the set is empty, exit with an appropriate message.
    if (numWords == 0) {
?>    <?insert /missing.mer ?>  <?code
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
  if (originalRequest != canonicalRequest) {
    fmt.Print("Status: 301 Moved Permanently\n")
    fmt.Printf("Location: %s\n\n", canonicalRequest)
    return
  }

  word := ""
?>
  <?insert /header.mer ?>

<h2> <?code fmt.Print(title) ?>: </h2>

<ul class="listing">
<?code
  for _, word := range words {
    path := "/entry/" + url.QueryEscape(word) + "/"
    // Replace hyphens with non-breaking hyphens.
    word = strings.Replace(word, "-", "&#8209;", -1)
    // Replace spaces with non-breaking spaces.
    word = strings.Replace(word, " ", "&nbsp;", -1)
    link := fmt.Sprintf("<a href=\"%s\">%s</a>", path, word)
?>
    <li><?code fmt.Print(link) ?></li><wbr />
<?code
  }
?>
</ul>
  <?insert /footer.mer ?>
<?code
}
?>
