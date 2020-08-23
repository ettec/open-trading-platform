import Login from '../components/Login';
import { Order } from '../serverapi/order_pb';
import { ViewServiceClient } from '../serverapi/ViewserviceServiceClientPb';
import { SubscribeToOrdersWithRootOriginatorIdArgs, OrderHistory, GetOrderHistoryArgs } from '../serverapi/viewservice_pb';
import Stream from './impl/Stream';
import * as grpcWeb from 'grpc-web';

export interface OrderService {
    SubscribeToAllParentOrders(onUpdate: (order: Order) => void): void
    GetChildOrders(order: Order): Array<Order>
    GetOrderHistory(order: Order, callback: (err: grpcWeb.Error, response: OrderHistory) => void): void

}



export default class OrderServiceImpl implements OrderService {

    viewService = new ViewServiceClient(Login.grpcContext.serviceUrl, null, null)
    orderStream: Stream<Order>
    orders = new Map<string, Order>()
    listeners = new Array<(order: Order) => void>()

    childOrders = new Map<string, Set<string>>()

    constructor() {

        let sto = new SubscribeToOrdersWithRootOriginatorIdArgs()

        sto.setRootoriginatorid(Login.desk)


        this.orderStream = new Stream<Order>((): grpcWeb.ClientReadableStream<any> => {
            return this.viewService.subscribeToOrdersWithRootOriginatorId(sto, Login.grpcContext.grpcMetaData)
        }, (order: Order) => {

            let updateOrders = false

            let currentOrder = this.orders.get(order.getId())
            if (currentOrder) {
                if (order.getVersion() > currentOrder.getVersion()) {
                    updateOrders = true
                }
            } else {
                updateOrders = true
            }

            if (updateOrders) {
                this.orders.set(order.getId(), order)


                if (this.IsParentOrder(order)) {
                    for (let listener of this.listeners) {
                        listener(order)
                    }
                } else {
                    let childIds = this.childOrders.get(order.getOriginatorref())

                    if (!childIds) {
                        childIds = new Set<string>()
                        this.childOrders.set(order.getOriginatorref(), childIds)
                    }

                    childIds.add(order.getId())
                }
            }


        }, "order updates stream")

    }

    IsParentOrder(order: Order): boolean {
        return order.getOriginatorid() === Login.desk
    }


    SubscribeToAllParentOrders(onUpdate: (order: Order) => void) {
        this.listeners.push(onUpdate)

        for (var order of this.orders.values()) {
            if (this.IsParentOrder(order)) {
                onUpdate(order)
            }
        }
    }

    GetChildOrders(order: Order): Array<Order> {

        let result = new Array<Order>()
        let childIds = this.childOrders.get(order.getId())
        if (childIds) {
            for (let id of childIds) {
                let order = this.orders.get(id)
                if (order) {
                    result.push(order)
                }
            }
        }

        return result
    }

    GetOrderHistory(order: Order, callback: (err: grpcWeb.Error,
        response: OrderHistory) => void): void {
        let args = new GetOrderHistoryArgs()
        args.setOrderid(order.getId())
        args.setToversion(order.getVersion())

        this.viewService.getOrderHistory(args, Login.grpcContext.grpcMetaData, callback)
    }


}