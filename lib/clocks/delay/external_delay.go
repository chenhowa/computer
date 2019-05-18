package delay

/*ExternalDelay represents a delay that occurs due to
an external source of the delay. This external source could be
the network, user input on the command line, user input over the network,
or even just a sleeping goroutine*/
type ExternalDelay struct {
	source delaySource
}

/*MakeExternalDelay is a constructor for an External Delay instance
with a given delaySource*/
func MakeExternalDelay(source delaySource) ExternalDelay {
	delay := ExternalDelay{
		source: source,
	}

	return delay
}

type delaySource interface {
	delay()
}

/*Delay requests an external source to provide a delay. The delay may be
significant, or nonexistent, depending on the source provided during
construction*/
func (d *ExternalDelay) Delay() {
	d.source.delay()
}
