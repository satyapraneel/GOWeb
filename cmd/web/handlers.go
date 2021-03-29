package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/satyapraneel/goweb/pkg/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// panic("oops! something went wrong") // Deliberate panic

	s, err := app.snippets.Latest()

	// if err != nil {
	// 	app.serverError(w, err)
	// 	return
	// }
	// for _, snippet := range s {
	// 	fmt.Fprintf(w, "%v\n", snippet)
	// }
	// Initialize a slice containing the paths to the two files. Note that the
	// home.page.tmpl file must be the *first* file in the slice.
	// files := []string{
	// 	"./ui/html/home.page.tmpl",
	// 	"./ui/html/base.layout.tmpl",
	// 	"./ui/html/footer.partial.tmpl",
	// }
	// Use the template.ParseFiles() function to read the files and store the
	// templates in a template set. Notice that we can pass the slice of file paths
	// as a variadic parameter?
	// ts, err := template.ParseFiles(files...)
	// data := &templateData{Snippets: s}
	if err != nil {
		// Because the home handler function is now a method against application
		// it can access its fields, including the error logger. We'll write the log
		// message to this instead of the standard logger.
		app.errorLog.Println(err.Error())
		app.serverError(w, err)
		return
	}
	app.render(w, r, "home.page.tmpl", &templateData{
		Snippets: s,
	})
	/* err = ts.Execute(w, data)
	if err != nil {

		// Also update the code here to use the error logger from the application
		// struct.
		app.errorLog.Println(err.Error())
		// http.Error(w, "Internal Server Error", 500)
		app.serverError(w, err)
	} */
}
func (app *application) showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		// http.NotFound(w, r)
		app.notFound(w)
		return
	}
	s, err := app.snippets.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, err)
		}
		return
	}

	// Use the new render helper.
	app.render(w, r, "show.page.tmpl", &templateData{
		Snippet: s,
	})

	// Initialize a slice containing the paths to the show.page.tmpl file,
	// plus the base layout and footer partial that we made earlier.
	/* files := []string{
		"./ui/html/show.page.tmpl",
		"./ui/html/base.layout.tmpl",
		"./ui/html/footer.partial.tmpl",
	}
	// Parse the template files...
	ts, err := template.ParseFiles(files...)
	if err != nil {
		app.serverError(w, err)
		return
	}
	// Create an instance of a templateData struct holding the snippet data.
	data := &templateData{Snippet: s}

	// And then execute them. Notice how we are passing in the snippet
	// data (a models.Snippet struct) as the final parameter.
	// Pass in the templateData struct when executing the template.
	err = ts.Execute(w, data)
	if err != nil {
		app.serverError(w, err)
	} */
}
func (app *application) createSnippet(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		// http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}
	// Create some variables holding dummy data. We'll remove these later on
	// during the build.
	title := "O snail"
	content := "O snail\nClimb Mount Fuji,\nBut slowly, slowly!\n\n– Kobayashi Issa"
	expires := "7"
	// Pass the data to the SnippetModel.Insert() method, receiving the
	// ID of the new record back.
	id, err := app.snippets.Insert(title, content, expires)
	if err != nil {
		app.serverError(w, err)
		return
	}

	// Redirect the user to the relevant page for the snippet
	http.Redirect(w, r, fmt.Sprintf("/snippet?id=%d", id), http.StatusSeeOther)
}
