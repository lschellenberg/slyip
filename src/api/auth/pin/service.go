package pin

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
	"yip/src/api/auth/verifier"
	"yip/src/config"
	"yip/src/providers"
	"yip/src/repositories"
	"yip/src/slyerrors"
)

type Service struct {
	config   *config.Config
	verifier *verifier.Verifier
	pool     PinPool
	database repositories.Database
	ep       *providers.EmailProvider
}

func NewService(
	config *config.Config,
	verifier *verifier.Verifier,
	useDB repositories.Database,
	ep *providers.EmailProvider,
) Service {
	return Service{
		config:   config,
		verifier: verifier,
		pool:     NewPool(60 * time.Minute),
		database: useDB,
		ep:       ep,
	}
}

func (s *Service) RequestPin(ctx context.Context, body *PinRequestDTO) (*PinRequestResponse, error) {
	account, err := s.database.GetAccountByEmail(ctx, body.Email)
	if err != nil {
		fmt.Println("trying to get account by email but failed: ", err.Error())
		account, err = s.database.RegisterAccountWithEmail(ctx, body.Email)
		if err != nil {
			fmt.Println("trying to register account by email but failed: ", err.Error())
			return nil, err
		}
	}

	pin := s.pool.request(account.ID, account.Email, body.ECDSAPubKey)
	response := &PinRequestResponse{
		AccountId:   account.ID,
		Email:       pin.Email,
		ECDSAPubKey: pin.ECDSAPubKey,
		Expiration:  pin.Expiration,
	}

	// test
	if !s.config.Test.On {
		err = s.ep.SendPinMail(pin.Email, pin.Pin)
		if err != nil {
			return nil, err
		}
	}

	if s.config.Test.On {
		response.Pin = pin.Pin
	}

	return response, nil
}

func (s *Service) Redeem(ctx context.Context, body *PinRedeemDTO) (*verifier.Token, error) {
	if !s.config.VerifyAudiencesExist(body.Audiences) {
		return nil, slyerrors.BadRequest(slyerrors.ErrCodeAudienceDoesntExist, "audience(s) dont exist")
	}
	pin, err := s.pool.redeem(body.Pin, body.PinSignature)
	if err != nil {
		return nil, err
	}

	uu, err := uuid.Parse(pin.AccountId)
	if err != nil {
		return nil, slyerrors.Unexpected(slyerrors.ErrCodeUnknown, err.Error())
	}

	devices, err := s.database.GetDevices(ctx, uu)
	if err != nil {
		return nil, err
	}

	deviceIsRegistered := false

	for _, d := range devices {
		if d.Address == pin.ECDSAPubKey {
			deviceIsRegistered = true
			break
		}
	}

	if !deviceIsRegistered {
		_, err = s.database.AddDevice(ctx, uu, pin.ECDSAPubKey)
		if err != nil {
			return nil, err
		}
	}

	account, err := s.database.GetAccountById(ctx, uu)
	if err != nil {
		return nil, err
	}

	return s.verifier.CreateToken(body.Audiences, account.ID, pin.ECDSAPubKey, account.LastUsedSLYWallet, verifier.RoleBasic)
}

func (s *Service) ListPins(ctx context.Context) ([]Pin, error) {
	return s.pool.list(), nil
}
