import React from 'react';
import query from 'query-string';


interface OrderParameters {
    quantity: number,
    price: number,
    side: string,
    instrumentId: string
}


export default class OrderTicket extends React.Component<{}, OrderParameters> {

    constructor() {
        super({});

        this.state = {
          quantity : 0,
          price : 0,
          side: '',
          instrumentId: ''
        };
    }      

   

    public render() {


        return ( 
            
          <div>
            <div>
            <label>Side </label>
            <input
              type="text"
              name="side"
              value={this.state.side}
              onChange={
                (e: any) =>{
                    this.setState({side: e.target.value})
                }
                    
              }
            />
            </div>
            <div>
            <label>Quantity </label>
            <input
              type="text"
              name="quantity"
              value={this.state.quantity}
              onChange={
                (e: any) =>{
                    this.setState({quantity: e.target.value})
                }
                    
              }
            />
            </div>
            <div>
            <label>Instrument </label>
            <input
              type="text"
              name="instrument"
              value={this.state.instrumentId}
              onChange={
                (e: any) =>{
                    this.setState({instrumentId: e.target.value})
                }

              }
            />
            </div>
            <div>
            <label>Price </label>
            <input
              type="text"
              name="price"
              value={this.state.price}
              onChange={
                (e: any) =>{
                    this.setState({price: e.target.value})
                }

              }
            />
            </div>
            <div>
              <button onClick={e=>this.sendOrder(this.state) }>Send</button>
            </div>



           
          </div>
        );

        }



        sendOrder(params: OrderParameters ) {
          
     let queryParams:string  = query.stringify(params);

          fetch('http://192.168.1.102:32413/order-management/create-and-route-order' + '?' + queryParams, {
              method: 'POST',
              mode: 'no-cors'
            })
              .then(

                response =>{ console.log(response.statusText) } 
                )

                
              
              .catch(error => {
                throw new Error(error);
              });



        }

        
       




}