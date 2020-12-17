package maker

const (
	totalLevel = 2001
	Buy        = false
	Sell       = true
)

type level struct {
	price float64
	id    string
}

type List struct {
	list      [totalLevel]*level
	sellFront int
	sellRear  int
	buyFront  int
	buyRear   int
}

func initList() *List {
	return &List{list: [totalLevel]*level{}, sellFront: 0, sellRear: 0, buyFront: 0, buyRear: 0}
}

func (q *List) Insert(l *level, i int) {
	q.list[i] = l
}

func (q *List) SetOrderid(id string, i int) {
	q.list[i].id = id
}

func (q *List) UpLevel(levels int) {
	for i := 0; i < levels; i++ {
		q.sellFront++ // 5 -> 6
		q.sellRear++  // 14 -> 15
		q.buyFront++  // -6 -> -5
		q.buyRear++   // -15 -> -14
	}
}

func (q *List) DownLevel(levels int) {
	for i := 0; i < levels; i++ {
		q.sellFront--
		q.sellRear--
		q.buyFront--
		q.buyRear--
	}
}
