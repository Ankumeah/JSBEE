# database

This package contains the DBContoller struct
which is responsible for executing any DB queries

By default this project only uses the "modernc.org/sqlite"
driver but more can be added, to do so create file called
`<driver>_options.go` and impliments the same values and functions
as the sqlite driver. Then create a sturct that impliments `DBContoller`.
DB drivers are sqitched at compile time with build flags

Happy Coding and Good Luck!
