package pin

import (
	"sync"
	"time"
	"yip/src/cryptox"
	"yip/src/slyerrors"
)

type Pin struct {
	AccountId   string    `json:"accountId"`
	Email       string    `json:"email"`
	ECDSAPubKey string    `json:"ecdsaPubKey"`
	Pin         string    `json:"pin,omitempty"`
	Expiration  time.Time `json:"expirationDate"`
}

type PinPool struct {
	pool            map[string]Pin
	mutex           *sync.Mutex
	ExpirationInMin time.Duration
}

func NewPool(expirationDurationInMin time.Duration) PinPool {
	return PinPool{
		pool:            make(map[string]Pin),
		mutex:           &sync.Mutex{},
		ExpirationInMin: expirationDurationInMin,
	}
}

func (p *PinPool) request(accountId string, email string, ecdsaPubKey string) Pin {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	pin := cryptox.GeneratePassword(6, false, false, true)
	for {
		_, ok := p.pool[pin]
		if !ok {
			pin = cryptox.GeneratePassword(6, false, false, true)
			break
		}
	}

	pp := Pin{
		AccountId:   accountId,
		Email:       email,
		ECDSAPubKey: ecdsaPubKey,
		Pin:         pin,
		Expiration:  time.Now().Add(p.ExpirationInMin * time.Minute),
	}

	p.pool[pin] = pp

	p.cleanPool()

	return pp
}

func (p *PinPool) redeem(pin string, pinSignature string) (Pin, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	pp, ok := p.pool[pin]
	if !ok {
		return pp, slyerrors.BadRequest(slyerrors.ErrCodePinNotExistent, "pin not found")
	}

	if pp.Expiration.Before(time.Now()) {
		return pp, slyerrors.Unauthorized(slyerrors.ErrCodePinExpired, "pin expired")
	}

	address, err := cryptox.Recover(pin, pinSignature, cryptox.SignMethodEthereumPrefix, false)
	if err != nil {
		return pp, err
	}

	if address.String() != pp.ECDSAPubKey {
		return pp, slyerrors.Unauthorized(slyerrors.ErrCodeWrongSignature, "invalid signature")
	}

	p.remove(pin)

	return pp, nil
}

func (p *PinPool) remove(pin string) {
	delete(p.pool, pin)
}

func (p *PinPool) list() []Pin {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	pins := make([]Pin, 0)
	for _, v := range p.pool {
		pins = append(pins, v)
	}
	return pins
}

func (p *PinPool) cleanPool() {
	removePin := make([]string, 0)
	now := time.Now()
	for _, v := range p.pool {
		if v.Expiration.Before(now) {
			removePin = append(removePin, v.Pin)
		}
	}

	for _, v := range removePin {
		p.remove(v)
	}
}
