import { Listing } from "../../serverapi/listing_pb";
import { Order } from "../../serverapi/order_pb";


export class ListingContext {

    selectedListing?: Listing;

    private listeners: Array<(listing: Listing) => void>;

    constructor() {
        this.listeners = new Array<(listing: Listing) => void>();

    }

    setSelectedListing(listing: Listing) {
        this.selectedListing = listing;
        this.listeners.forEach(l => l(listing));
    }

    addListener(listener: (listing: Listing) => void) {
        if (this.selectedListing) {
            listener(this.selectedListing);
        }

        this.listeners.push(listener);
    }

}

export class OrderContext {

    selectedOrder?: Order;
    private listeners: Array<(order: Order) => void>;

    constructor() {
        this.listeners = new Array<(order: Order) => void>();
    }

    setSelectedOrder(order: Order) {
        this.selectedOrder = order;
        this.listeners.forEach(l => l(order));
    }

    addListener(listener: (order: Order) => void) {
        if (this.selectedOrder) {
            listener(this.selectedOrder);
        }
        this.listeners.push(listener);
    }

}
