package domain

type Events []any

func (ev *Events) append(v ...any) {
	new := *ev
	new = append(new, v)
	*ev = new
}

type Allocated struct {
	OrderId  string
	Sku      string
	Qty      int
	BatchRef string
}

type OutOfStock struct {
	Sku string
}
