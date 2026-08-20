package database

type dbDriver interface {
}
type DBController struct { db dbDriver }
