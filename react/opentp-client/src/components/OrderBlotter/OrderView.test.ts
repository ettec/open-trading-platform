import { Listing } from "../../serverapi/listing_pb"
import { Timestamp } from "../../serverapi/modelcommon_pb"
import { Order, OrderStatus } from "../../serverapi/order_pb"
import { ListingService } from "../../services/ListingService"
import { OrdersView } from "./OrderView"


test("order view", () => {
    let view = new OrdersView(new TestListingService(), ()=>{})


    view.addOrUpdateOrder(createTestOrder("1", 1, OrderStatus.LIVE))
    view.addOrUpdateOrder(createTestOrder("2", 2, OrderStatus.LIVE))
    view.addOrUpdateOrder(createTestOrder("3", 3, OrderStatus.CANCELLED))
    view.addOrUpdateOrder(createTestOrder("5", 5, OrderStatus.LIVE))
    view.addOrUpdateOrder(createTestOrder("4", 4, OrderStatus.LIVE))
    
    view.setFilter((order: Order) : boolean => {
        return order.getStatus() === OrderStatus.LIVE
    })

    view.setSort((a:Order, b:Order) : number => {
       
        let aCreated = a.getCreated()
        let bCreated = b.getCreated() 
        if( aCreated  && bCreated) {
            return aCreated.getSeconds() - bCreated.getSeconds()
        } else {
            return 0
        }
    })

    let orders = view.getOrders()

    expect(orders.length).toEqual(4)
    expect(orders[0].id).toEqual("1")
    expect(orders[1].id).toEqual("2")
    expect(orders[2].id).toEqual("4")
    expect(orders[3].id).toEqual("5")



})

function createTestOrder(id : string, createdTime: number, status: OrderStatus) : Order {
    let o1 = new Order()
    o1.setId(id)
    o1.setStatus(status)
    let t1 = new Timestamp()
    t1.setSeconds(createdTime)
    o1.setCreated(t1)
    return o1
}

class TestListingService implements ListingService {

    GetListingImmediate(listingId: number): Listing | undefined {
        return undefined
    }

    GetListing(listingId: number, listener: (
      response: Listing) => void): void {

      }
    
}
