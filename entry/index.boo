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
  defer runtime.PrintCGI()
?> <?insert /advance.mer ?> <?code

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
?>  <?insert /missing.mer ?>  <?code
  }
  // Extract the requested dictionary headword.
  parts := strings.Split(request, "/")
  // The request '/entry/abacus/' yields ['', 'entry', 'abacus', ''].
  word := parts[2]
  // Apply URL decoding and UTF-8 decoding.
  word, err = url.QueryUnescape(word)
  if err != nil {
?>  <?insert /missing.mer ?>  <?code
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
?>  <?insert /missing.mer ?> <?code
  }

  // Construct the canonical URL for this word. At this writing, the
  //  token parts[1] resolves to "entry". This may change eventually.
  //  If the virtual directory structure changes, the template will
  //  have to be modified here and elsewhere.
  canonicalRequest := "/" + parts[1] + "/" + url.QueryEscape(word) + "/"
  // Redirect to the canonical URL if necessary.
  if originalRequest != canonicalRequest {
    runtime.Print("Status: 301 Moved Permanently\n")
    runtime.Printf("Location: %s\n\n", canonicalRequest)
    return
  }

  titleExtension := " : "+word
?>
  <?insert /header.mer ?>

<h2> <?code runtime.Print(word) ?>: </h2>

  <?insert /definition.mer ?>

  <?insert /footer.mer ?>
<?code
}
?>
