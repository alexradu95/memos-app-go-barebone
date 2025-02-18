package posts

type QueryParams struct {
	AccountId  int64  `query:"accountId"`
	SearchText string `query:"searchText"`
	DateFrom   string `query:"dateFrom"`
	DateTo     string `query:"dateTo"`
	PageNumber int64  `query:"pageNumber"`
	PageSize   int64  `query:"pageSize"`
}
