package fixgateway

import (
	"fmt"
	"github.com/ettec/open-trading-platform/go/execution-venue/internal/model"
	"github.com/ettec/open-trading-platform/go/execution-venue/internal/ordergateway"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/enum"
	"github.com/quickfixgo/quickfix/field"
	"github.com/quickfixgo/quickfix/fix42/businessmessagereject"
	"github.com/quickfixgo/quickfix/fix50sp1/ordercancelrequest"
	"github.com/quickfixgo/quickfix/fix50sp2/executionreport"
	"github.com/quickfixgo/quickfix/fix50sp2/newordersingle"
	"github.com/shopspring/decimal"
	"log"
	"math/big"
	"time"
)

type FixOrderGateway struct {
	sessionID    quickfix.SessionID
	orderHandler OrderHandler
}

func NewFixOrderGateway(sessionID quickfix.SessionID) ordergateway.OrderGateway {
	return &FixOrderGateway{
		sessionID: sessionID,
	}
}

func (f *FixOrderGateway) Send(order *model.Order, listing *model.Listing) error {

	side, err := getFixSide(order.Side)
	if err != nil {
		return err
	}

	msg := newordersingle.New(field.NewClOrdID(order.Id), field.NewSide(side),
		field.NewTransactTime(time.Now()), field.NewOrdType(enum.OrdType_LIMIT_OR_BETTER))

	msg.SetOrderQty(decimal.NewFromBigInt(big.NewInt(order.GetQuantity().Mantissa), order.GetQuantity().GetExponent()), 0)
	msg.SetPrice(decimal.NewFromBigInt(big.NewInt(order.GetPrice().Mantissa), order.GetPrice().GetExponent()), 0)
	msg.SetSymbol(listing.MarketSymbol)

	return quickfix.SendToTarget(msg, f.sessionID)

}

func (f *FixOrderGateway) Cancel(order *model.Order) error {

	side, err := getFixSide(order.Side)
	if err != nil {
		return err
	}

	msg := ordercancelrequest.New(field.NewClOrdID(order.Id), field.NewSide(side),
		field.NewTransactTime(time.Now()))

	return quickfix.SendToTarget(msg, f.sessionID)
}


func getFixSide(side model.Side) (enum.Side, error) {

	switch side {
	case model.Side_BUY:
		return enum.Side_BUY, nil
	case model.Side_SELL:
		return enum.Side_SELL, nil
	default:
		return "", fmt.Errorf("side %v not supported", side)
	}

}

type OrderHandler interface {
	SetOrderStatus(orderId string, status model.OrderStatus) error
	UpdateTradedQuantity(orderId string, lastPrice model.Decimal64, lastQty model.Decimal64) error
}

type fixHandler struct {
	sessionToHandler map[quickfix.SessionID]OrderHandler
	inboundRouter    *quickfix.MessageRouter
	outboundRouter   *quickfix.MessageRouter
}

func NewFixHandler(sessionID quickfix.SessionID, handler OrderHandler) quickfix.Application {
	f := fixHandler{sessionToHandler: make(map[quickfix.SessionID]OrderHandler)}

	f.sessionToHandler[sessionID] = handler
	f.inboundRouter = quickfix.NewMessageRouter()
	f.inboundRouter.AddRoute(executionreport.Route(f.onExecutionReport))


	return &f
}

func logSessionMsg(sessionID quickfix.SessionID, msg string) {
	log.Print(sessionID.String() + ":" + msg)
}

func logSessionMsgf(sessionID quickfix.SessionID, format string, v ...interface{}) {
	log.Printf(sessionID.String()+":"+format, v)
}

func (f *fixHandler) onOutboundBusinessMessageReject(msg businessmessagereject.BusinessMessageReject, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {

	logSessionMsgf(sessionID, "Sending reject message to target: %v", msg)

	return nil
}

func (f *fixHandler) onExecutionReport(msg executionreport.ExecutionReport, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {

	logSessionMsg(sessionID, "received execution report:"+msg.Message.String())

	execType, err := msg.GetExecType()
	if err != nil {
		return err
	}

	handler, exists := f.sessionToHandler[sessionID]
	if !exists {
		logSessionMsg(sessionID, "Error: No handler found for session id")
		return nil
	}

	orderId, err := msg.GetClOrdID()
	if err != nil {
		return err
	}

	switch execType {
	case enum.ExecType_NEW:
		err := handler.SetOrderStatus(orderId, model.OrderStatus_LIVE)
		if err != nil {
			logSessionMsgf(sessionID, "Error: Failed to set order status: %v", err)
			return nil
		}
	case enum.ExecType_CANCELED:
		err := handler.SetOrderStatus(orderId, model.OrderStatus_CANCELLED)
		if err != nil {
			logSessionMsgf(sessionID, "Error: Failed to set order status: %v", err)
			return nil
		}
	case enum.ExecType_TRADE:
		lastQty, err := msg.GetLastQty()

		if err != nil {
			return err
		}

		lastPrice, err := msg.GetLastPx()
		if err != nil {
			return err
		}

		handler.UpdateTradedQuantity(orderId, *model.ToDecimal64(lastPrice), *model.ToDecimal64(lastQty))
	}

	return nil
}

//Notification of a session begin created.
func (f *fixHandler) OnCreate(sessionID quickfix.SessionID) {
	logSessionMsg(sessionID, "created")
}

//Notification of a session successfully logging on.
func (f *fixHandler) OnLogon(sessionID quickfix.SessionID) {
	logSessionMsg(sessionID, "logon received")
}

//Notification of a session logging off or disconnecting.
func (f *fixHandler) OnLogout(sessionID quickfix.SessionID) {
	logSessionMsg(sessionID, "logout received")
}

//Notification of admin message being sent to target.
func (f *fixHandler) ToAdmin(message *quickfix.Message, sessionID quickfix.SessionID) {
}

//Notification of app message being sent to target.
func (f *fixHandler) ToApp(message *quickfix.Message, sessionID quickfix.SessionID) error {

	return nil
}

//Notification of admin message being received from target.
func (f *fixHandler) FromAdmin(message *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {

	return nil
}

//Notification of app message being received from target.
func (f *fixHandler) FromApp(message *quickfix.Message, sessionID quickfix.SessionID) quickfix.MessageRejectError {

	return f.inboundRouter.Route(message, sessionID)
}
