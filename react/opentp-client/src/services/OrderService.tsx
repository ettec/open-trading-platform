

import Login from '../components/Login';
import { Timestamp } from '../serverapi/modelcommon_pb';
import { Order } from '../serverapi/order_pb';
import { ViewServiceClient } from '../serverapi/ViewserviceServiceClientPb';
import { SubscribeToOrdersWithRootOriginatorIdArgs } from '../serverapi/viewservice_pb';
import Stream from './impl/Stream';




export interface OrderService {
    SubscribeToAllParentOrders(onUpdate : (order : Order)=>void) : void
    GetChildOrders(order : Order) : Array<Order>
}



 export default class OrderServiceImpl implements OrderService {

    viewService = new ViewServiceClient(Login.grpcContext.serviceUrl, null, null)
    orderStream : Stream<Order>
    orders = new Map<string,Order>()
    listeners = new Array<(order : Order)=>void>()

    childOrders = new Map<string, Set<string>>()

    constructor() {
        let after = new Timestamp()

        let startOfLocalDay = new Date()
        startOfLocalDay.setHours(0, 0, 0, 0)
        after.setSeconds(Math.floor(startOfLocalDay.getTime() / 1000))
        let sto = new SubscribeToOrdersWithRootOriginatorIdArgs()
        sto.setAfter(after)
        sto.setRootoriginatorid(Login.desk)
    
    
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
                

                if (this.IsParentOrder(order)) {
                    for( let listener of this.listeners ) {
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

    IsParentOrder(order :Order ) : boolean {
        return order.getOriginatorid() === Login.desk
    }


    SubscribeToAllParentOrders(onUpdate : (order : Order)=>void) {
        this.listeners.push(onUpdate)

        for( var order of this.orders.values()) {
            if (this.IsParentOrder(order)) {
                onUpdate(order)
            }
        }
    }

    GetChildOrders(order : Order) : Array<Order> {

        let result = new Array<Order>()
        let childIds = this.childOrders.get(order.getId())
        if (childIds) {
            for( let id of childIds) {
                let order = this.orders.get(id)
                if(order) {
                    result.push(order)
                }
            }
        }
        
        return result
    }


}