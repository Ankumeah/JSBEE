# apis
This package contains all the http routes to be handled by this program
Any new route group is to follow this format

```golang
// Comment explaining the route
// Note that the function name is to match the route path
func routePath(r *gin.RouterGroup, app *a.App) {
	group := r.Group("/routepath", middlewares...)

	group.POST(route, func(c *gin.Context) {
    ...
	}
}
```

After this you are to call this new created func from within `Apis()` (`apis.go`)

Happy Coding and Good Luck!
