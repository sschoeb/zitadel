package command

import (
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/repository/idpconfig"
)

type JWTConfigWriteModel struct {
	eventstore.WriteModel

	IDPConfigID  string
	Issuer       string
	KeysEndpoint string
	State        domain.IDPConfigState
}

func (wm *JWTConfigWriteModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *idpconfig.JWTConfigAddedEvent:
			wm.reduceConfigAddedEvent(e)
		case *idpconfig.JWTConfigChangedEvent:
			wm.reduceConfigChangedEvent(e)
		case *idpconfig.IDPConfigDeactivatedEvent:
			wm.State = domain.IDPConfigStateInactive
		case *idpconfig.IDPConfigReactivatedEvent:
			wm.State = domain.IDPConfigStateActive
		case *idpconfig.IDPConfigRemovedEvent:
			wm.State = domain.IDPConfigStateRemoved
		}
	}

	return wm.WriteModel.Reduce()
}

func (wm *JWTConfigWriteModel) reduceConfigAddedEvent(e *idpconfig.JWTConfigAddedEvent) {
	wm.IDPConfigID = e.IDPConfigID
	wm.Issuer = e.Issuer
	wm.KeysEndpoint = e.KeysEndpoint
	wm.State = domain.IDPConfigStateActive
}

func (wm *JWTConfigWriteModel) reduceConfigChangedEvent(e *idpconfig.JWTConfigChangedEvent) {
	if e.Issuer != nil {
		wm.Issuer = *e.Issuer
	}
	if e.KeysEndpoint != nil {
		wm.KeysEndpoint = *e.KeysEndpoint
	}
}
