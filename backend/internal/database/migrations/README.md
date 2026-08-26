# database/migrations

This package contains all DB migrations
Note that the migrations have no revert, rather you are
encouraged to create a new migration to revert that migration
to help keep the timeline lenier and avoide branching migrations

Any new migration should have its own file with its verison number
be very careful to not give any two migrations the same number.
Migrations must add themselves to the `migrations` local variable within
their `init()` function

Happy Coding and Good Luck!
