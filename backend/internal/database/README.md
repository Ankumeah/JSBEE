# database

This package contains the DBContoller struct
which is responsible for executing any DB queries

By default this project only uses the "modernc.org/sqlite"
driver but more can be added, to do so create file called
`<driver>/<driver>_options.go` and do the following within it:
- build flag for the driver
- register the driver with "database/sql"
- `DriverName string` (name to use when calling `sql.Open`)
- `func isUniqueViolation(err error) bool` (checks for unique violation)
- `func isForeignKeyViolation(err error) bool` (check for foreign key violation)

Then create a sturct that impliments `DBContoller`.
DB drivers are sqitched at compile time with build flags

Happy Coding and Good Luck!
