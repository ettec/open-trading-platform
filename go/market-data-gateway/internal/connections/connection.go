package connections

type Connection interface {

	Subscribe(listingId int)
	Close() error
}

