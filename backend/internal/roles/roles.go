package roles

var Owner = newRole("owner", true, true, true)
var Admin = newRole("admin", true, false, false)
var Reviewer = newRole("reviewer", true, false, false)
var Viewer = newRole("viewer", false, false, false)
