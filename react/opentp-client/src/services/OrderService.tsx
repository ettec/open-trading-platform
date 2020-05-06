

import { Timestamp } from '../serverapi/modelcommon_pb';
import { Order } from '../serverapi/order_pb';

import Stream from './impl/Stream';
import Login from '../components/Login';
import { ViewServiceClient } from '../serverapi/ViewserviceServiceClientPb';
import { SubscribeToOrdersWithRootOriginatorIdArgs } from '../serverapi/viewservice_pb';



export interface OrderService {
    SubscribeToAllParentOrders(onUpdate : (order : Order)=>void) : void
    GetChildOrders(order : Order) : Array<Order>
}



 export default class OrderServiceImpl implements OrderService {

    viewService = new ViewServiceClient(Login.grpcContext.serviceUrl, null, null)
    orderStream : Stream<Order>
    orders = new Map<string,Order>()
    listeners = new Array<(order : Order)=>void>()


    constructor() {
        let after = new Timestamp()

        let startOfLocalDay = new Date()
        startOfLocalDay.setHours(0, 0, 0, 0)
        after.setSeconds(Math.floor(startOfLocalDay.getTime() / 1000))
        let sto = new SubscribeToOrdersWithRootOriginatorIdArgs()
        sto.setAfter(after)
        sto.setRootoriginatorid("desk")
    
        
        this.orderStream = new Stream(() => {
            return  this.viewService.subscribeToOrdersWithRootOriginatorId(sto, Login.grpcContext.grpcMetaData)
        }, (order : Order)=> {

            let updateOrders = false
                
            let currentOrder =  this.orders.get(order.getId())
            if (currentOrder)  {
                if (order.getVersion() > currentOrder.getVersion()) {
                    updateOrders = true
                }
            } else {
                updateOrders = true
            }

            if( updateOrders ) {
                this.orders.set(order.getId(), order)
                for( let listener of this.listeners ) {
                    listener(order)
                }
            }


        }, "order updates stream")

    }

    SubscribeToAllParentOrders(onUpdate : (order : Order)=>void) {
        this.listeners.push(onUpdate)

        for( var order of this.orders.values()) {
            onUpdate(order)
        }
    }

    GetChildOrders(order : Order) : Array<Order> {
        let result = new Array<Order>()
        result.push(order)
        return result
    }


}