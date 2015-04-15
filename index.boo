<?code
package main

import (
  "os"
  "strings"
  "unicode/utf8"
  "regexp"
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
    request = request + "/"
  }
  // Replace repeated slashes with single slashes.
  request = regexp.MustCompile("/+").ReplaceAllString(request, "/")
  // There should now be exactly one slash.
  if strings.Count(request, "/") != 1 {
?>  <?insert /missing.mer ?>  <?code
  }

  // Construct the canonical URL for this page.
  canonicalRequest := "/"
  // Redirect to the canonical URL if necessary.
  //fmt.Fprintf(logFile, "%s  originalRequest = \"%s\", request = \"%s\"" +
  //    ", canonicalRequest = \"%s\"\n",
  //    time.Now().Format(time.Stamp), originalRequest, request, canonicalRequest)
  if originalRequest != canonicalRequest {
    fmt.Print("Status: 301 Moved Permanently\n")
    fmt.Printf("Location: %s\n\n", canonicalRequest)
    return
  }

  titleExtension := " by Ambrose Bierce"
  word := ""
?>
  <?insert /header.mer ?>

<h2> Random entry: </h2>

<?code

  // Get a random word from the database.

  query = "select word from entries order by random() limit 1"
  db.QueryRow(query).Scan(&word)

  // Look up the definition of that word.

  query = "select definition from entries where word = ?"
  var definition string
  db.QueryRow(query, word).Scan(&definition)
?>
  <?insert /definition.mer ?>

<h2> About this dictionary </h2>

<p> <i>The Devil's Dictionary</i> is the work of the American writer <a
href="http://en.wikipedia.org/wiki/Ambrose_Bierce">Ambrose Bierce</a>,
whose oeuvre includes the heavily anthologized short story "An Occurrence
at Owl Creek Bridge". Bierce began to write satirical dictionary
definitions in 1875. For several decades, he wrote them sporadically
in newspaper columns. The first book-length collection of his devilish
definitions appeared in 1906. While some of Bierce's allusions are
outmoded, his observations of human frailty remain sharp and true. </p>


  <?insert /footer.mer ?>
<?code
}
?>
