package coredto

type UpdateSql struct {
	Versions []Version `xml:"version"`
}

type Version struct {
	Id   string `xml:"id,attr"`
	Sqls []Sql  `xml:"sql"`
}

type Sql struct {
	Sql string `xml:",cdata"`
}
