package server

import (
	"context"
	"github.com/a1sarpi/QuietPlace/currency/data"
	protos "github.com/a1sarpi/QuietPlace/currency/protos/currency"
	"github.com/hashicorp/go-hclog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"time"
)

// Currency is a gRPC server it implements the methods defined by the CurrencyServer interface
type Currency struct {
	protos.UnimplementedCurrencyServer // Встраиваем структуру UnimplementedCurrencyServer
	rates                              *data.ExchangeRates
	log                                hclog.Logger
	subscriptions                      map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest
}

// NewCurrency creates a new Currency server
func NewCurrency(r *data.ExchangeRates, l hclog.Logger) *Currency {
	c := &Currency{rates: r, log: l, subscriptions: make(map[protos.Currency_SubscribeRatesServer][]*protos.RateRequest)}
	go c.handleUpdates()

	return c
}

func (c *Currency) handleUpdates() {
	ru := c.rates.MonitorRates(5 * time.Second)
	for range ru {
		c.log.Info("Got Update Rates")

		// loop over subscribed clients
		for k, v := range c.subscriptions {
			// loop over subscribed rates
			for _, rr := range v {
				r, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
				if err != nil {
					c.log.Error("Unable to get updated rate", "base", rr.GetBase().String(), "destination", rr.GetDestination().String(), "destination", rr.GetDestination().String())
				}

				err = k.Send(&protos.RateResponse{Base: rr.Base, Destination: rr.Destination, Rate: r})
				if err != nil {
					c.log.Error("Unable to send updated rate", "base", rr.GetBase().String(), "destination", rr.GetDestination().String(), "destination", rr.GetDestination().String())
				}
			}
		}
	}
}

// GetRate implements the CurrencyServer GetRate method and returns the currency exchange rate
// for the two given currencies
func (c *Currency) GetRate(ctx context.Context, rr *protos.RateRequest) (*protos.RateResponse, error) {
	c.log.Info("Handle GetRate", "base", rr.GetBase(), "destination", rr.GetDestination())

	if rr.Base == rr.Destination {
		err := status.Newf(
			codes.InvalidArgument,
			"Base currency %s and Destination currency %s cannot be the same",
			rr.Base.String(),
			rr.Destination.String(),
		)

		err, wde := err.WithDetails(rr)
		if wde != nil {
			return nil, wde
		}

		return nil, err.Err()
	}

	rate, err := c.rates.GetRate(rr.GetBase().String(), rr.GetDestination().String())
	if err != nil {
		return nil, err
	}

	return &protos.RateResponse{Base: rr.Base, Destination: rr.Destination, Rate: rate}, nil
}

// SubscribeRates implements the gRPC bidirectional streaming method for the server
func (c *Currency) SubscribeRates(src protos.Currency_SubscribeRatesServer) error {

	// handle client messages
	for {
		rr, err := src.Recv() // Recv is a blocking method which returns on client data

		// io.EOF signals that the client has closed the connection
		if err == io.EOF {
			c.log.Info("Client has closed the connection")
			return err
		}

		// any other error means the transport between the server and client is unavailable
		if err != nil {
			c.log.Error("Unable to read from client", "error", err)
			return err
		}

		c.log.Info("Handle client request", "request", rr)
		rrs, ok := c.subscriptions[src]
		if !ok {
			rrs = []*protos.RateRequest{}
		}

		// check that subscription does not exist
		for _, v := range rrs {
			if v.Base == rr.Base && v.Destination == rr.Destination {
				// subscription exists return errors
			}
		}

		rrs = append(rrs, rr)
		c.subscriptions[src] = rrs
	}

	return nil
}
