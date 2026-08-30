# frontend

This package contains the files to generate the frontend's static files.
Its is reccamended to call `*ComponentUpdater.UpdateAll` and `assets.SaveAssets` dueing the apps
initalisation to make sure that all assets exist at all times

This project uses "github.com/a-h/templ" to generate its frontend files
so make sure to regularly run `templ generate`

The pattern the this project follow for its static files is as follows,
static files are generated upon startup and then served from teh reverse proxy's
or CDN's cache. Volumes page cache is perged and the page is rebuilt everytime a new
paper is approved to refrect and cache the newly generated page

For every new file to be generated keep the file name in `values.go`

Any new static asset such as images are to be kept in the `assets` package
Take a look at `assets/README.md` for more info

Happy Coding and Good Luck!
