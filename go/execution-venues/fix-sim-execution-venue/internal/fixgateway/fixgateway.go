package fixgateway

import (
	"fmt"
	"github.com/ettec/otp-common/model"
	"github.com/quickfixgo/quickfix"
	"github.com/quickfixgo/quickfix/enum"
	"github.com/quickfixgo/quickfix/field"
	"github.com/quickfixgo/quickfix/fix40/ordercancelreject"
	"github.com/quickfixgo/quickfix/fix42/businessmessagereject"
	"github.com/quickfixgo/quickfix/fix50sp1/ordercancelreplacerequest"
	"github.com/quickfixgo/quickfix/fix50sp1/ordercancelrequest"
	"github.com/quickfixgo/quickfix/fix50sp2/executionreport"
	"github.com/quickfixgo/quickfix/fix50sp2/newordersingle"
	"github.com/shopspring/decimal"
	"log"
	"os"
	"strings"
	"time"
)

type fixOrderGateway struct {
	sessionID    quickfix.SessionID
	orderHandler OrderHandler
}

func NewFixOrderGateway(sessionID quickfix.SessionID) *fixOrderGateway{
	return &fixOrderGateway{
		sessionID: sessionID,
	}
}

func (f *fixOrderGateway) Send(order *model.Order, listing *model.Listing) error {

	side, err := getFixSide(order.Side)
	if err != nil {
		return err
	}

	msg := newordersingle.New(field.NewClOrdID(order.Id), field.NewSide(side),
		field.NewTransactTime(time.Now()), field.NewOrdType(enum.OrdType_LIMIT_OR_BETTER))

	msg.SetOrderQty(toFixDecimal(order.GetQuantity()))
	msg.SetPrice(toFixDecimal(order.GetPrice()))
	msg.SetSymbol(listing.MarketSymbol)

	logSessionMsg(f.sessionID, "sending new order single:"+toReadableString(msg.Message))

	return quickfix.SendToTarget(msg, f.sessionID)
}

func (f *fixOrderGateway) Modify(order *model.Order, listing *model.Listing, quantity *model.Decimal64, price *model.Decimal64) error {
	side, err := getFixSide(order.Side)
	if err != nil {
		return err
	}

	msg := ordercancelreplacerequest.New(field.NewClOrdID(order.Id), field.NewSide(side),
		field.NewTransactTime(time.Now()), field.NewOrdType(enum.OrdType_LIMIT_OR_BETTER))

	msg.SetOrderQty(toFixDecimal(quantity))
	msg.SetPrice(toFixDecimal(price))
	msg.SetSymbol(listing.MarketSymbol)

	logSessionMsg(f.sessionID, "sending order cancel replace request:"+toReadableString(msg.Message))

	return quickfix.SendToTarget(msg, f.sessionID)
}

func toFixDecimal(d *model.Decimal64) (decimal.Decimal, int32) {
	var scale int32 = 0
	if d.Exponent < 0 {
		scale = -d.Exponent
	}

	return decimal.New(d.GetMantissa(), d.GetExponent()), scale
}

func toReadableString(msg *quickfix.Message) string {
	return strings.ReplaceAll(msg.String(), "\001", "|")
}

func (f *fixOrderGateway) Cancel(order *model.Order) error {

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
	SetErrorMsg(orderId string, msg string) error
	AddExecution(orderId string, lastPrice model.Decimal64, lastQty model.Decimal64, execId string) error
}

type fixHandler struct {
	sessionToHandler map[quickfix.SessionID]OrderHandler
	inboundRouter    *quickfix.MessageRouter
	outboundRouter   *quickfix.MessageRouter
	errLog			 *log.Logger
}

func NewFixHandler(sessionID quickfix.SessionID, handler OrderHandler) quickfix.Application {
	f := fixHandler{sessionToHandler: make(map[quickfix.SessionID]OrderHandler),
		errLog: log.New(os.Stderr, "error", log.Ltime | log.Lshortfile),
		}

	f.sessionToHandler[sessionID] = handler
	f.inboundRouter = quickfix.NewMessageRouter()
	f.inboundRouter.AddRoute(executionreport.Route(f.onExecutionReport))
	f.inboundRouter.AddRoute(ordercancelreject.Route(f.onOrderCancelReject))

	return &f
}

func (f *fixHandler) onOrderCancelReject(msg ordercancelreject.OrderCancelReject, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {

	logSessionMsg(sessionID, "received order cancel/modification reject:"+toReadableString(msg.Message))

	handler, exists := f.sessionToHandler[sessionID]
	if !exists {
		logSessionMsg(sessionID, "Error: No handler found for session id")
		return nil
	}

	orderId, err := msg.GetClOrdID()
	if err != nil {
		return err
	}

	errMsg, err := msg.GetText()
	if err != nil {
		return err
	}

	er := handler.SetErrorMsg(orderId, errMsg)
	if er != nil {
		f.errLog.Printf("Failed to set error msg on order %v, message: %v,  error:%v", orderId, errMsg, er)
	}

	return nil
}

func logSessionMsg(sessionID quickfix.SessionID, msg string) {
	log.Print(sessionID.String() + ":" + msg)
}

func logSessionMsgf(sessionID quickfix.SessionID, format string, v ...interface{}) {
	log.Printf(sessionID.String()+":"+format, v)
}

func (f *fixHandler) onOutboundBusinessMessageReject(msg businessmessagereject.BusinessMessageReject, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {

	logSessionMsgf(sessionID, "Sending reject message to target: %v", toReadableString(msg.Message))

	return nil
}

func (f *fixHandler) onExecutionReport(msg executionreport.ExecutionReport, sessionID quickfix.SessionID) (err quickfix.MessageRejectError) {

	logSessionMsg(sessionID, "received execution report:"+toReadableString(msg.Message))

	execType, err := msg.GetExecType()
	if err != nil {
		return err
	}

	handler, exists := f.sessionToHandler[sessionID]
	if !exists {
		logSessionMsg(sessionID, "Error: No handler found for  session id")
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
	case enum.ExecType_REPLACED:
		err := handler.SetOrderStatus(orderId, model.OrderStatus_LIVE)
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

		execId, err := msg.GetExecID()
		if err != nil {
			return err
		}

		er := handler.AddExecution(orderId, *model.ToDecimal64(lastPrice), *model.ToDecimal64(lastQty), execId)
		if er != nil {
			f.errLog.Printf("failed to add execution to order %v, error: %v", orderId, er)
		}
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
