package parsing

import "flag"

type Params struct {
	Command      *string
	ID           *int64
	RecipientID  *int64
	WeightGrams  *int64
	PriceKopecks *int64
	Cover        *string
	ExpireString *string
	NotGiven     *bool
	Args         []string
}

func Parse(f *Params) {
	f.Command = flag.String("command", "", "name of command")
	f.ID = flag.Int64("id", 0, "order id")
	f.RecipientID = flag.Int64("recipient", 0, "recipient id")
	f.WeightGrams = flag.Int64("weight", 0, "order weight")
	f.PriceKopecks = flag.Int64("price", 0, "order price")
	f.Cover = flag.String("cover", "", "order cover")
	f.ExpireString = flag.String("expire", "", "expire date")
	f.NotGiven = flag.Bool("t", false, "return only not given orders")

	flag.Parse()
	f.Args = flag.Args()
}
