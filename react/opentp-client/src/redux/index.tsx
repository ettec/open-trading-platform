import { Order } from "../model/Model";


export interface RootState {
    selectedOrder? : Order;
}


export const SET_SELECTED_ORDER = 'SET_SELECTED_ORDER';
export type SET_SELECTED_ORDER = typeof SET_SELECTED_ORDER;

//define action interfaces
export interface SetSelectedOrderAction {
    type: SET_SELECTED_ORDER;
    order: Order;
}


export type OrderActionTypes = SetSelectedOrderAction

//define actions
export function setSelectedOrder(theOrder: Order): SetSelectedOrderAction {
    return {
        type: SET_SELECTED_ORDER,
        order: theOrder
    }
    
};


export function omsReducer(state = {},
     action: OrderActionTypes) : RootState {

    switch( action.type ) {
        case SET_SELECTED_ORDER:
            console.log("Setting order id to:" + action.order.id)
            return { ...state,  selectedOrder: action.order }
    }

    return state;
}




